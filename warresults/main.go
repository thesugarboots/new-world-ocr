package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"os"
	"strconv"
	"strings"

	"github.com/otiai10/gosseract/v2"
)

type SubImager interface {
	SubImage(r image.Rectangle) image.Image
}

type PlayerScore struct {
	Name    string
	Score   int
	Kills   int
	Deaths  int
	Assists int
	Healing int
	Damage  int
}

func main() {

	filepath := "./test_images/WarResults/score-0.jpg"

	// Read image from file that already exists
	warResultsFile, err := os.Open(filepath)
	if err != nil {
		fmt.Println(err)
	}

	defer warResultsFile.Close()

	warResultsFile.Seek(0, 0)
	warResults, imageInfo, err := image.Decode(warResultsFile)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(imageInfo)

	playerScores := make(map[string]PlayerScore)

	//names
	nameText := processImage(warResults, 807, 431, 227, 662, .6, "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ.  ")
	fmt.Println(nameText)
	playerNames := strings.Split(nameText, "\n")
	fmt.Println(len(playerNames))

	//scores
	scoreText := processImage(warResults, 1043, 431, 840, 662, .5, "0123456789 ")
	scores := strings.Split(scoreText, "\n")

	if len(playerNames) == len(scores) {
		for i := 0; i < len(playerNames); i++ {
			playerRanking := strings.Split(scores[i], " ")
			score, _ := strconv.Atoi(playerRanking[0])
			kills, _ := strconv.Atoi(playerRanking[1])
			deaths, _ := strconv.Atoi(playerRanking[2])
			assists, _ := strconv.Atoi(playerRanking[3])
			healing, _ := strconv.Atoi(playerRanking[4])
			damage, _ := strconv.Atoi(playerRanking[5])
			playerScores[playerNames[i]] = PlayerScore{
				playerNames[i],
				score,
				kills,
				deaths,
				assists,
				healing,
				damage,
			}
		}
	}

	fmt.Println(playerScores)
}

func processImage(warResults image.Image, x0 int, y0 int, xdelta int, ydelta int, greyBoundaryMod float32, whitelist string) string {
	scoreBoundaries := image.Rect(x0, y0, x0+xdelta, y0+ydelta)

	blackOrWhiteScore := imageToBlackOrWhite(warResults, greyBoundaryMod, scoreBoundaries)
	//saveImage(blackOrWhiteScore, "./test_images/WarResults/results/score-0 BW.jpg")

	text := toText(blackOrWhiteScore, whitelist)
	return text
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

func toBytes(score image.Image) []byte {
	buf := new(bytes.Buffer)
	err1 := jpeg.Encode(buf, score, nil)
	if err1 != nil {
		fmt.Println(err1)
	}
	return buf.Bytes()
}

func imageToBlackOrWhite(img image.Image, greyBoundaryMod float32, rect image.Rectangle) image.Image {
	size := img.Bounds().Size()
	wImg := image.NewRGBA(rect)

	r, _, _, _ := color.White.RGBA()
	modifier := greyBoundaryMod * float32(r)
	greyBoundary := uint32(modifier)

	for x := 0; x < size.X; x++ {
		for y := 0; y < size.Y; y++ {
			r, g, b, _ := img.At(x, y).RGBA()
			c := color.White

			if r > greyBoundary && g > greyBoundary && b > greyBoundary {
				c = color.Black
			}

			wImg.Set(x, y, c)
		}
	}

	return wImg
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

	err = jpeg.Encode(out, img, &opt) // put quality to 80%
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
