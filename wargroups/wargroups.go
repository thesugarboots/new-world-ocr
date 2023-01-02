package wargroups

import (
	"encoding/csv"
	"fmt"
	"image"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

	imgproc "github.com/thesugarboots/new-world-ocr/imageprocessing"
)

func ProcessWarGroups(inFile string, outFile string) {
	img, err := imgproc.LoadImage(inFile)
	if err != nil {
		fmt.Println(err)
	}

	groupedPlayers := make(map[string]string)

	groups := make([]string, 10)
	var wgPI sync.WaitGroup
	wgPI.Add(10)

	charNameWhitelist := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ. '"
	yStart1, yStart2, yDiff, xDiff := 370, 780, 278, 145
	xStart1, xStart2, xStart3, xStart4, xStart5 := 675, 974, 1273, 1571, 1870

	greyBoundary := float32(.395)
	go processImage(img, xStart1, yStart1, xDiff, yDiff, greyBoundary, charNameWhitelist, &groups[0], &wgPI)
	go processImage(img, xStart2, yStart1, xDiff, yDiff, greyBoundary, charNameWhitelist, &groups[1], &wgPI)
	go processImage(img, xStart3, yStart1, xDiff, yDiff, greyBoundary, charNameWhitelist, &groups[2], &wgPI)
	go processImage(img, xStart4, yStart1, xDiff, yDiff, greyBoundary, charNameWhitelist, &groups[3], &wgPI)
	go processImage(img, xStart5, yStart1, xDiff, yDiff, greyBoundary, charNameWhitelist, &groups[4], &wgPI)

	go processImage(img, xStart1, yStart2, xDiff, yDiff, greyBoundary, charNameWhitelist, &groups[5], &wgPI)
	go processImage(img, xStart2, yStart2, xDiff, yDiff, greyBoundary, charNameWhitelist, &groups[6], &wgPI)
	go processImage(img, xStart3, yStart2, xDiff, yDiff, greyBoundary, charNameWhitelist, &groups[7], &wgPI)
	go processImage(img, xStart4, yStart2, xDiff, yDiff, greyBoundary, charNameWhitelist, &groups[8], &wgPI)
	go processImage(img, xStart5, yStart2, xDiff, yDiff, greyBoundary, charNameWhitelist, &groups[9], &wgPI)
	wgPI.Wait()

	for i, group := range groups {
		players := strings.Split(group, "\n")
		for j := 0; j < len(players); {
			player, playersNewJ, err := imgproc.NextNonEmptyElement(players, j)
			if err != nil {
				fmt.Println(err)
				break
			} else {
				j = playersNewJ
			}

			//Empty seem to be parsing as "ty", so skipping in this case.
			if player == "ty" {
				continue
			}

			groupedPlayers[player] = strconv.Itoa(i + 1)

		}
	}

	fmt.Println(groupedPlayers)
	fmt.Println("Player extracted:", len(groupedPlayers))

	f, err := os.Create(outFile)

	if err != nil {

		fmt.Println("failed to open file", err)
	}

	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	headers := []string{"Name", "Group"}

	if err := w.Write(headers); err != nil {
		log.Fatalln("error writing record to file", err)
	}

	for name, group := range groupedPlayers {
		playerEntry := []string{name, group}
		if err := w.Write(playerEntry); err != nil {
			log.Fatalln("error writing record to file", err)
		}
	}

}

func processImage(img image.Image, x0 int, y0 int, xdelta int, ydelta int, greyBoundaryMod float32, whitelist string, text *string, wg *sync.WaitGroup) {
	imgBoundaries := image.Rect(x0, y0, x0+xdelta, y0+ydelta)
	//whitenMin, whitenMax := float32(.00), float32(1.0)

	croppedImg := imgproc.CropReadOnlyImage(img, imgBoundaries)
	rootPath := "./wargroups/test_images/" + strconv.Itoa(x0) + "-" + strconv.Itoa(y0)
	imgproc.SaveImage(croppedImg, rootPath+"0-cropped.jpg") //convert to be debuggable

	// whitenedImg := imgproc.AdjustToWhite(croppedImg, imgBoundaries, whitenMin, whitenMax)
	// imgproc.SaveImage(whitenedImg, rootPath+"1-whitened.jpg") //convert to be debuggable

	greyScaleImg := imgproc.GrayScaleImage(croppedImg, imgBoundaries)
	imgproc.SaveImage(greyScaleImg, rootPath+"2-grey.jpg") //convert to be debuggable

	blackOrWhiteImg := imgproc.BlackOrWhiteImage(greyScaleImg, greyBoundaryMod, imgBoundaries)
	imgproc.SaveImage(blackOrWhiteImg, rootPath+"3-bw.jpg") //convert to be debuggable

	*text = imgproc.Text(blackOrWhiteImg, whitelist)

	if wg != nil {
		wg.Done()
	}
}
