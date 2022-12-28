package main

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"

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

	playerScores := make(map[string]PlayerScore)

	initI := 0
	numberOfFiles := 6

	var wg sync.WaitGroup
	wg.Add(numberOfFiles)

	for i := initI; i < initI+numberOfFiles; i++ {
		filenum := strconv.Itoa(i)
		go processWarResultsFile("./test_images/WarResults/score-"+filenum+".jpg", playerScores, &wg)
		println("Finished: " + filenum)
		//println(playerScores)

	}

	wg.Wait()

	// processWarResultsFile("./test_images/WarResults/score-0.jpg", playerScores, nil)

	fmt.Println(len(playerScores))
	for _, results := range playerScores {
		fmt.Println(results)
	}

}

func processWarResultsFile(filepath string, playerScores map[string]PlayerScore, wg *sync.WaitGroup) {
	fmt.Println("Processing file:", filepath)
	// Read image from file that already exists
	warResults := loadImage(filepath)

	//old
	//yStart,yDiff := 431,662
	yStart, yDiff := 431, 662
	//names
	nameText := processImage(warResults, 807, yStart, 227, yDiff, .5, "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ.  ")
	names := strings.Split(nameText, "\n")
	fmt.Println("Names(", len(names), "):", names)

	// fmt.Println(len(playerNames))
	//score
	scoresText := processImage(warResults, 1044, yStart, 138, yDiff, .50, "0123456789 ")
	scores := strings.Split(scoresText, "\n")
	fmt.Println("Scores(", len(scores), "):", scores)

	//kills
	killsText := processImage(warResults, 1172, yStart, 138, yDiff, .50, "0123456789 ")
	kills := strings.Split(killsText, "\n")
	fmt.Println("Kills(", len(kills), "):", kills)

	//deaths
	deathsText := processImage(warResults, 1300, yStart, 138, yDiff, .50, "0123456789 ")
	deaths := strings.Split(deathsText, "\n")
	fmt.Println("Deaths(", len(deaths), "):", deaths)

	//assists
	assistsText := processImage(warResults, 1435, yStart, 138, yDiff, .50, "0123456789 ")
	assists := strings.Split(assistsText, "\n")
	fmt.Println("Assists(", len(assists), "):", assists)

	//healing
	healingText := processImage(warResults, 1560, yStart, 138, yDiff, .50, "0123456789 ")
	healing := strings.Split(healingText, "\n")
	fmt.Println("Healing(", len(healing), "):", healing)

	//damage
	damageText := processImage(warResults, 1717, yStart, 138, yDiff, .50, "0123456789 ")
	damage := strings.Split(damageText, "\n")
	fmt.Println("Damage(", len(damage), "):", damage)

	namesI, scoresI, killsI, deathsI, assistsI, healingI, damageI := 0, 0, 0, 0, 0, 0, 0

	for namesI < len(names) {

		name, namesNewI, err := nextNonEmptyElement(names, namesI)
		if err != nil {
			fmt.Println(err)
			break
		} else {
			namesI = namesNewI
		}

		score, scoresNewI, err := nextNonEmptyElement(scores, scoresI)
		if err != nil {
			fmt.Println(err)
			break
		} else {
			scoresI = scoresNewI
		}

		kill, killsNewI, err := nextNonEmptyElement(kills, killsI)
		if err != nil {
			fmt.Println(err)
			break
		} else {
			killsI = killsNewI
		}

		death, deathsNewI, err := nextNonEmptyElement(deaths, deathsI)
		if err != nil {
			fmt.Println(err)
			break
		} else {
			deathsI = deathsNewI
		}

		assist, assistsNewI, err := nextNonEmptyElement(assists, assistsI)
		if err != nil {
			fmt.Println(err)
			break
		} else {
			assistsI = assistsNewI
		}

		heal, healingNewI, err := nextNonEmptyElement(healing, healingI)
		if err != nil {
			fmt.Println(err)
			break
		} else {
			healingI = healingNewI
		}

		dmg, damageNewI, err := nextNonEmptyElement(damage, damageI)
		if err != nil {
			fmt.Println(err)
			break
		} else {
			damageI = damageNewI
		}

		addPlayerScore(playerScores, name, score, kill, death, assist, heal, dmg)
	}

	//println(playerScores)

	fmt.Println("File processed:", filepath)
	if wg != nil {
		wg.Done()
	}

}

func nextNonEmptyElement(array []string, i int) (string, int, error) {
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

func addPlayerScore(playerScores map[string]PlayerScore, name string, scoreText string, killsText string, deathsText string, assistsText string, healingText string, damageText string) {
	score, err := strconv.Atoi(scoreText)
	if err != nil {
		fmt.Println("bad score", err)
		fmt.Println(scoreText)
	}

	kills, err := strconv.Atoi(killsText)
	if err != nil {
		fmt.Println("bad kills", err)
		fmt.Println(killsText)
	}

	deaths, err := strconv.Atoi(deathsText)
	if err != nil {
		fmt.Println("bad deaths", err)
		fmt.Println(deathsText)
	}

	assists, err := strconv.Atoi(assistsText)
	if err != nil {
		fmt.Println("bad assists", err)
		fmt.Println(assistsText)
	}

	healing, err := strconv.Atoi(healingText)
	if err != nil {
		fmt.Println("bad healing", err)
		fmt.Println(healingText)
	}

	damage, err := strconv.Atoi(damageText)
	if err != nil {
		fmt.Println("bad damage", err)
		fmt.Println(damageText)
	}

	playerScores[name] = PlayerScore{
		name,
		score,
		kills,
		deaths,
		assists,
		healing,
		damage,
	}
}

func loadImage(filepath string) image.Image {
	warResultsFile, err := os.Open(filepath)
	if err != nil {
		fmt.Println(err)
	}

	defer warResultsFile.Close()

	warResultsFile.Seek(0, 0)
	warResults, _, err := image.Decode(warResultsFile)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println(imageInfo)
	return warResults
}

func processImage(img image.Image, x0 int, y0 int, xdelta int, ydelta int, greyBoundaryMod float32, whitelist string) string {
	imgBoundaries := image.Rect(x0, y0, x0+xdelta, y0+ydelta)

	croppedImg := cropReadOnlyImage(img, imgBoundaries)
	fileNum := strconv.Itoa(int(rand.Uint32()))
	saveImage(croppedImg, "./test_images/WarResults/results/score-"+fileNum+"0-cropped.jpg")
	greyScaleImg := imageToGrayScale(croppedImg, imgBoundaries)
	saveImage(greyScaleImg, "./test_images/WarResults/results/score-"+fileNum+"1-grey.jpg")
	blackOrWhiteImg := imageToBlackOrWhite(greyScaleImg, greyBoundaryMod, imgBoundaries)
	saveImage(blackOrWhiteImg, "./test_images/WarResults/results/score-"+fileNum+"2-bw.jpg")

	text := toText(blackOrWhiteImg, whitelist)
	return text
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

	//fmt.Println(text)

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
