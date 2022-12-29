package main

import (
	"os"

	"github.com/thesugarboots/new-world-ocr/wargroups"
	"github.com/thesugarboots/new-world-ocr/warresults"
)

func main() {
	args := os.Args[1:]

	switch args[0] {
	case "results":
		var inDir, outFile string

		if len(args) > 1 {
			inDir = args[1]
		} else {
			inDir = "."
		}

		if len(args) > 2 {
			outFile = args[2]
		} else {
			outFile = "./war_results.csv"
		}

		warresults.ProcessWarResults(inDir, outFile)
	case "groups":
		var inDir, outFile string

		if len(args) > 1 {
			inDir = args[1]
		} else {
			inDir = "."
		}

		if len(args) > 2 {
			outFile = args[2]
		} else {
			outFile = "./war_groups.csv"
		}
		wargroups.ProcessWarGroups(inDir, outFile)
	}

}
