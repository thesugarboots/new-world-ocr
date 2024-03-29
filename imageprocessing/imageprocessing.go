package imageprocessing

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"os"

	"github.com/otiai10/gosseract/v2"
)

type SubImager interface {
	SubImage(r image.Rectangle) image.Image
}

func LoadImages(dir string) []image.Image {
	imgs := make([]image.Image, 0)

	files, err := os.ReadDir(dir)
	if err != nil {
		fmt.Println("failed to read directory", err)
	}

	for _, file := range files {
		if !file.IsDir() {
			img, err := LoadImage(dir + "/" + file.Name())
			if err != nil {
				fmt.Println("failed to read directory", err)
			} else {
				imgs = append(imgs, img)
			}

		}
	}

	return imgs

}

func LoadImage(filepath string) (image.Image, error) {
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Println(err)
	}

	defer file.Close()

	file.Seek(0, 0)
	img, _, err := image.Decode(file)

	return img, err
}

func CropReadOnlyImage(img image.Image, rect image.Rectangle) image.Image {
	size := img.Bounds().Size()
	croppedImg := image.NewRGBA(rect)

	for x := 0; x < size.X; x++ {
		for y := 0; y < size.Y; y++ {
			croppedImg.Set(x, y, img.At(x, y))
		}
	}

	return croppedImg
}

func AdjustToWhite(img image.Image, rect image.Rectangle, greyBoundaryModMin float32, greyBoundaryModMax float32) image.Image {
	whitenedImg := image.NewRGBA(rect)

	greyBoundaryMin, greyBoundaryMax := greyBoundary(greyBoundaryModMin), greyBoundary(greyBoundaryModMax)

	for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
		for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
			r, g, b, a := img.At(x, y).RGBA()

			// average
			if greyBoundaryMin < r && r < greyBoundaryMax && greyBoundaryMin < g && g < greyBoundaryMax && greyBoundaryMin < b && b < greyBoundaryMax {
				// fmt.Println(r, g, b)
				rgbAry := []uint32{r, g, b}
				newWhite := maxInt(rgbAry)
				// fmt.Println(newWhite)
				r, g, b = newWhite, newWhite, newWhite
			}

			c := color.RGBA{
				R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a),
			}
			whitenedImg.Set(x, y, c)
		}
	}

	return whitenedImg
}

func maxInt(numbers []uint32) uint32 {
	maxNum := uint32(0)

	for i := 0; i < len(numbers); i++ {
		if i == 0 {
			maxNum = numbers[i]
		}
		if maxNum <= numbers[i] {
			maxNum = numbers[i]
		}
	}

	return maxNum
}

func GrayScaleImage(img image.Image, rect image.Rectangle) image.Image {
	greyScaleImg := image.NewRGBA(rect)

	rWgt, gWgt, bWgt := 0.92126, 0.97152, 0.90722

	for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
		for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
			pixel := img.At(x, y)
			originalColor := color.RGBAModel.Convert(pixel).(color.RGBA)
			//fmt.Println(originalColor)
			// Offset colors a little, adjust it to your taste
			r := float64(originalColor.R) * rWgt
			g := float64(originalColor.G) * gWgt
			b := float64(originalColor.B) * bWgt
			// average
			grey := uint8((r + g + b) / 3)
			c := color.RGBA{
				R: grey, G: grey, B: grey, A: originalColor.A,
			}
			greyScaleImg.Set(x, y, c)
		}
	}

	return greyScaleImg
}

func BlackOrWhiteImage(img image.Image, greyBoundaryMod float32, rect image.Rectangle) image.Image {
	blackOrWhiteImg := image.NewRGBA(rect)

	//Get grey boundary between white and black
	greyBoundary := greyBoundary(greyBoundaryMod)

	for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
		for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
			r, g, b, _ := img.At(x, y).RGBA()
			c := color.White

			if greyBoundary < r && greyBoundary < g && greyBoundary < b {
				c = color.Black
			}

			blackOrWhiteImg.Set(x, y, c)
		}
	}

	return blackOrWhiteImg
}

func greyBoundary(greyBoundaryMod float32) uint32 {
	r, _, _, _ := color.White.RGBA()
	modifier := greyBoundaryMod * float32(r)
	return uint32(modifier)
}

func SaveImage(img image.Image, filepath string) {

	out, err := os.Create(filepath)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var opt jpeg.Options

	opt.Quality = 100
	// ok, write out the data into the new JPEG file

	err = jpeg.Encode(out, img, &opt)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}

func Text(img image.Image, whitelist string) string {
	client := gosseract.NewClient()
	defer client.Close()

	client.SetWhitelist(whitelist)
	client.SetImageFromBytes(toBytes(img))
	text, err := client.Text()

	if err != nil {
		fmt.Println(err)
	}

	return text
}

func NextNonEmptyElement(array []string, i int) (string, int, error) {
	var element string
	var err error

	for ; i < len(array); i++ {
		if array[i] != "" {
			element = array[i]
			i++
			break
		}
	}

	if element == "" {
		err = errors.New("no more elements")
	}

	return element, i, err
}

func toBytes(score image.Image) []byte {
	buf := new(bytes.Buffer)
	err1 := jpeg.Encode(buf, score, nil)
	if err1 != nil {
		fmt.Println(err1)
	}
	return buf.Bytes()
}
