package imageprocessing

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"os"
	"sync"

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

func ProcessImage(img image.Image, x0 int, y0 int, xdelta int, ydelta int, greyBoundaryMod float32, whitelist string, text *string, wg *sync.WaitGroup) {
	imgBoundaries := image.Rect(x0, y0, x0+xdelta, y0+ydelta)

	croppedImg := cropReadOnlyImage(img, imgBoundaries)
	//fileNum := strconv.Itoa(int(rand.Uint32()))
	//convert to be debuggable
	//saveImage(croppedImg, "./test_images/WarResults/results/score-"+fileNum+"0-cropped.jpg")
	greyScaleImg := imageToGrayScale(croppedImg, imgBoundaries)
	//convert to be debuggable
	//saveImage(greyScaleImg, "./test_images/WarResults/results/score-"+fileNum+"1-grey.jpg")
	blackOrWhiteImg := imageToBlackOrWhite(greyScaleImg, greyBoundaryMod, imgBoundaries)
	//convert to be debuggable
	//saveImage(blackOrWhiteImg, "./test_images/WarResults/results/score-"+fileNum+"2-bw.jpg")

	*text = toText(blackOrWhiteImg, whitelist)

	if wg != nil {
		wg.Done()
	}
}

func cropReadOnlyImage(img image.Image, rect image.Rectangle) image.Image {
	size := img.Bounds().Size()
	croppedImg := image.NewRGBA(rect)

	for x := 0; x < size.X; x++ {
		for y := 0; y < size.Y; y++ {
			croppedImg.Set(x, y, img.At(x, y))
		}
	}

	return croppedImg
}

func imageToGrayScale(img image.Image, rect image.Rectangle) image.Image {
	greyScaleImg := image.NewRGBA(rect)

	for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
		for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
			pixel := img.At(x, y)
			originalColor := color.RGBAModel.Convert(pixel).(color.RGBA)
			//fmt.Println(originalColor)
			// Offset colors a little, adjust it to your taste
			r := float64(originalColor.R) * 0.92126
			g := float64(originalColor.G) * 0.97152
			b := float64(originalColor.B) * 0.90722
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

func imageToBlackOrWhite(img image.Image, greyBoundaryMod float32, rect image.Rectangle) image.Image {
	blackOrWhiteImg := image.NewRGBA(rect)

	//Get grey boundary between white and black
	r, _, _, _ := color.White.RGBA()
	modifier := greyBoundaryMod * float32(r)
	greyBoundary := uint32(modifier)

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

func saveImage(img image.Image, filepath string) {

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

func toText(blackOrWhiteScore image.Image, whitelist string) string {
	client := gosseract.NewClient()
	defer client.Close()

	client.SetWhitelist(whitelist)
	client.SetImageFromBytes(toBytes(blackOrWhiteScore))
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
