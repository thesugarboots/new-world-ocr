package main

import (
	"os"

	"github.com/thesugarboots/new-world-ocr/warresults"
)

func main() {
	var inDir, outFile string
	args := os.Args[1:]
	if len(args) > 0 {
		inDir = args[0]
	} else {
		inDir = "."
	}

	if len(args) > 1 {
		outFile = args[1]
	} else {
		outFile = "./war_results.csv"
	}

	warresults.ProcessWarResults(inDir, outFile)

}
