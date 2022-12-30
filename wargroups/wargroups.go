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

	yStart1, yStart2, yDiff, xDiff := 370, 780, 278, 152
	greyBoundary := float32(.4)
	go processImage(img, 680, yStart1, xDiff, yDiff, greyBoundary, "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ.  '", &groups[0], &wgPI)
	go processImage(img, 979, yStart1, xDiff, yDiff, greyBoundary, "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ.  '", &groups[1], &wgPI)
	go processImage(img, 1278, yStart1, xDiff, yDiff, greyBoundary, "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ.  '", &groups[2], &wgPI)
	go processImage(img, 1576, yStart1, xDiff, yDiff, greyBoundary, "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ.  '", &groups[3], &wgPI)
	go processImage(img, 1875, yStart1, xDiff, yDiff, greyBoundary, "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ.  '", &groups[4], &wgPI)

	go processImage(img, 680, yStart2, xDiff, yDiff, greyBoundary, "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ.  '", &groups[5], &wgPI)
	go processImage(img, 979, yStart2, xDiff, yDiff, greyBoundary, "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ.  '", &groups[6], &wgPI)
	go processImage(img, 1278, yStart2, xDiff, yDiff, greyBoundary, "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ.  '", &groups[7], &wgPI)
	go processImage(img, 1576, yStart2, xDiff, yDiff, greyBoundary, "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ.  '", &groups[8], &wgPI)
	go processImage(img, 1875, yStart2, xDiff, yDiff, greyBoundary, "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ.  '", &groups[9], &wgPI)
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

			//Crowns seem to be parsing as " u", so removing that at the end.
			//if strings.

			groupedPlayers[player] = strconv.Itoa(i + 1)

		}
	}

	fmt.Println(groupedPlayers)

	f, err := os.Create(outFile)
	defer f.Close()

	if err != nil {

		fmt.Println("failed to open file", err)
	}

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

	croppedImg := imgproc.CropReadOnlyImage(img, imgBoundaries)
	rootPath := "./wargroups/test_images/" + strconv.Itoa(x0) + "-" + strconv.Itoa(y0)
	//convert to be debuggable
	imgproc.SaveImage(croppedImg, rootPath+"0-cropped.jpg")
	greyScaleImg := imgproc.GrayScaleImage(croppedImg, imgBoundaries)
	//convert to be debuggable
	imgproc.SaveImage(greyScaleImg, rootPath+"1-grey.jpg")
	blackOrWhiteImg := imgproc.BlackOrWhiteImage(greyScaleImg, greyBoundaryMod, imgBoundaries)
	//convert to be debuggable
	imgproc.SaveImage(blackOrWhiteImg, rootPath+"2-bw.jpg")

	*text = imgproc.Text(blackOrWhiteImg, whitelist)

	if wg != nil {
		wg.Done()
	}
}
