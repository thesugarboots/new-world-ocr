package wargroups

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/thesugarboots/new-world-ocr/imageprocessing"
)

func ProcessWarGroups(inFile string, outFile string) {
	img, err := imageprocessing.LoadImage(inFile)
	if err != nil {
		fmt.Println(err)
	}

	groupedPlayers := make(map[string]string)

	groups := make([]string, 10)
	var wgPI sync.WaitGroup
	wgPI.Add(10)

	yStart1, yStart2, yDiff, xDiff := 370, 780, 278, 140
	go imageprocessing.ProcessImage(img, 680, yStart1, xDiff, yDiff, .3, "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ.  '", &groups[0], &wgPI)
	go imageprocessing.ProcessImage(img, 979, yStart1, xDiff, yDiff, .3, "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ.  '", &groups[1], &wgPI)
	go imageprocessing.ProcessImage(img, 1278, yStart1, xDiff, yDiff, .3, "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ.  '", &groups[2], &wgPI)
	go imageprocessing.ProcessImage(img, 1576, yStart1, xDiff, yDiff, .3, "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ.  '", &groups[3], &wgPI)
	go imageprocessing.ProcessImage(img, 1875, yStart1, xDiff, yDiff, .3, "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ.  '", &groups[4], &wgPI)

	go imageprocessing.ProcessImage(img, 680, yStart2, xDiff, yDiff, .3, "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ.  '", &groups[5], &wgPI)
	go imageprocessing.ProcessImage(img, 979, yStart2, xDiff, yDiff, .3, "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ.  '", &groups[6], &wgPI)
	go imageprocessing.ProcessImage(img, 1278, yStart2, xDiff, yDiff, .3, "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ.  '", &groups[7], &wgPI)
	go imageprocessing.ProcessImage(img, 1576, yStart2, xDiff, yDiff, .3, "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ.  '", &groups[8], &wgPI)
	go imageprocessing.ProcessImage(img, 1875, yStart2, xDiff, yDiff, .3, "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ.  '", &groups[9], &wgPI)
	wgPI.Wait()

	for i, group := range groups {
		players := strings.Split(group, "\n")
		for j := 0; j < len(players); {
			player, playersNewJ, err := imageprocessing.NextNonEmptyElement(players, j)
			if err != nil {
				fmt.Println(err)
				break
			} else {
				j = playersNewJ
			}
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
