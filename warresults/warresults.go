package warresults

import (
	"encoding/csv"
	"fmt"
	"image"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	imgproc "github.com/thesugarboots/new-world-ocr/imageprocessing"
)

type PlayerScore struct {
	Name    string
	Score   int
	Kills   int
	Deaths  int
	Assists int
	Healing int
	Damage  int
}

func ProcessWarResults(inDir string, outFile string) {
	startTime := time.Now().UnixNano()
	playerScores := make(map[string]PlayerScore)

	warResultImgs := imgproc.LoadImages(inDir)

	var wg sync.WaitGroup
	wg.Add(len(warResultImgs))

	for _, warResultImg := range warResultImgs {
		go processWarResultsFile(warResultImg, playerScores, &wg)
	}

	wg.Wait()

	// processWarResultsFile("./test_images/WarResults/score-0.jpg", playerScores, nil)
	f, err := os.Create(outFile)
	defer f.Close()

	if err != nil {

		fmt.Println("failed to open file", err)
	}

	w := csv.NewWriter(f)
	defer w.Flush()

	headers := []string{"Name", "Score", "Kills", "Deaths", "Assists", "Healing", "Damage"}

	if err := w.Write(headers); err != nil {
		log.Fatalln("error writing record to file", err)
	}

	for _, playerScore := range playerScores {
		if err := w.Write(playerScoreToArray(playerScore)); err != nil {
			log.Fatalln("error writing record to file", err)
		}
	}

	fmt.Println("Elapsed time(ns):", time.Now().UnixNano()-startTime)
	fmt.Println("Scores extracted:", len(playerScores))
}

func playerScoreToArray(playerScore PlayerScore) []string {
	result := make([]string, 7)
	result[0] = playerScore.Name
	result[1] = strconv.Itoa(playerScore.Score)
	result[2] = strconv.Itoa(playerScore.Kills)
	result[3] = strconv.Itoa(playerScore.Deaths)
	result[4] = strconv.Itoa(playerScore.Assists)
	result[5] = strconv.Itoa(playerScore.Healing)
	result[6] = strconv.Itoa(playerScore.Damage)

	return result
}

func processWarResultsFile(warResults image.Image, playerScores map[string]PlayerScore, wg *sync.WaitGroup) {

	//old
	//yStart,yDiff := 431,662
	yStart, yDiff := 431, 662
	var nameText, scoresText, killsText, deathsText, assistsText, healingText, damageText string
	var wgPI sync.WaitGroup
	wgPI.Add(7)
	greyBoundary := float32(.4)
	//names
	go processImage(warResults, 807, yStart, 227, yDiff, greyBoundary, "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ. '", &nameText, &wgPI)
	//score
	go processImage(warResults, 1044, yStart, 138, yDiff, greyBoundary, "0123456789 ", &scoresText, &wgPI)
	//kills
	go processImage(warResults, 1172, yStart, 138, yDiff, greyBoundary, "0123456789 ", &killsText, &wgPI)
	//deaths
	go processImage(warResults, 1300, yStart, 138, yDiff, greyBoundary, "0123456789 ", &deathsText, &wgPI)
	//assists
	go processImage(warResults, 1435, yStart, 138, yDiff, greyBoundary, "0123456789 ", &assistsText, &wgPI)
	//healing
	go processImage(warResults, 1560, yStart, 138, yDiff, greyBoundary, "0123456789 ", &healingText, &wgPI)
	//damage
	go processImage(warResults, 1717, yStart, 138, yDiff, greyBoundary, "0123456789 ", &damageText, &wgPI)
	wgPI.Wait()

	names := strings.Split(nameText, "\n")
	// fmt.Println("Names(", len(names), "):", names)
	scores := strings.Split(scoresText, "\n")
	// fmt.Println("Scores(", len(scores), "):", scores)
	kills := strings.Split(killsText, "\n")
	// fmt.Println("Kills(", len(kills), "):", kills)
	deaths := strings.Split(deathsText, "\n")
	// fmt.Println("Deaths(", len(deaths), "):", deaths)
	assists := strings.Split(assistsText, "\n")
	// fmt.Println("Assists(", len(assists), "):", assists)
	healing := strings.Split(healingText, "\n")
	// fmt.Println("Healing(", len(healing), "):", healing)
	damage := strings.Split(damageText, "\n")
	// fmt.Println("Damage(", len(damage), "):", damage)

	namesI, scoresI, killsI, deathsI, assistsI, healingI, damageI := 0, 0, 0, 0, 0, 0, 0

	for namesI < len(names) {

		name, namesNewI, err := imgproc.NextNonEmptyElement(names, namesI)
		if err != nil {
			fmt.Println(err)
			break
		} else {
			namesI = namesNewI
		}

		score, scoresNewI, err := imgproc.NextNonEmptyElement(scores, scoresI)
		if err != nil {
			fmt.Println(err)
			break
		} else {
			scoresI = scoresNewI
		}

		kill, killsNewI, err := imgproc.NextNonEmptyElement(kills, killsI)
		if err != nil {
			fmt.Println(err)
			break
		} else {
			killsI = killsNewI
		}

		death, deathsNewI, err := imgproc.NextNonEmptyElement(deaths, deathsI)
		if err != nil {
			fmt.Println(err)
			break
		} else {
			deathsI = deathsNewI
		}

		assist, assistsNewI, err := imgproc.NextNonEmptyElement(assists, assistsI)
		if err != nil {
			fmt.Println(err)
			break
		} else {
			assistsI = assistsNewI
		}

		heal, healingNewI, err := imgproc.NextNonEmptyElement(healing, healingI)
		if err != nil {
			fmt.Println(err)
			break
		} else {
			healingI = healingNewI
		}

		dmg, damageNewI, err := imgproc.NextNonEmptyElement(damage, damageI)
		if err != nil {
			fmt.Println(err)
			break
		} else {
			damageI = damageNewI
		}

		addPlayerScore(playerScores, name, score, kill, death, assist, heal, dmg)
	}

	if wg != nil {
		wg.Done()
	}

}

func processImage(img image.Image, x0 int, y0 int, xdelta int, ydelta int, greyBoundaryMod float32, whitelist string, text *string, wg *sync.WaitGroup) {
	imgBoundaries := image.Rect(x0, y0, x0+xdelta, y0+ydelta)

	croppedImg := imgproc.CropReadOnlyImage(img, imgBoundaries)
	fileNum := strconv.Itoa(int(rand.Uint32()))
	//convert to be debuggable
	imgproc.SaveImage(croppedImg, "./warresults/test_images/nyewar/results/score-"+fileNum+"0-cropped.jpg")
	greyScaleImg := imgproc.GrayScaleImage(croppedImg, imgBoundaries)
	//convert to be debuggable
	imgproc.SaveImage(greyScaleImg, "./warresults/test_images/nyewar/results/score-"+fileNum+"1-grey.jpg")
	blackOrWhiteImg := imgproc.BlackOrWhiteImage(greyScaleImg, greyBoundaryMod, imgBoundaries)
	//convert to be debuggable
	imgproc.SaveImage(blackOrWhiteImg, "./warresults/test_images/nyewar/results/score-"+fileNum+"2-bw.jpg")

	*text = imgproc.Text(blackOrWhiteImg, whitelist)

	if wg != nil {
		wg.Done()
	}
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
