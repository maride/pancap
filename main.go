package main

import (
	"flag"
	"fmt"
	"git.darknebu.la/maride/pancap/analyze"
	"git.darknebu.la/maride/pancap/output"
	"log"
	"math/rand"
	"time"
)

func main() {
	// important things first
	printMOTD()

	// register flags
	registerFileFlags()
	output.RegisterFlags()
	flag.Parse()

	// Open the given PCAP
	packetSource, _, fileErr := openPCAP()
	if fileErr != nil {
		// Encountered problems with the PCAP - permission and/or existance error
		log.Fatalf("Error occured while opeining specified file: %s", fileErr.Error())
	}

	// Start analyzing
	analyzeErr := analyze.Analyze(packetSource)
	if analyzeErr != nil {
		// Mh, encountered some problems while analyzing file
		log.Fatalf("Error occurred while analyzing: %s", analyzeErr.Error())
	}

	// Extract found and requested files
	output.StoreFiles()

	// Show user analysis
	analyze.PrintSummary()

	// Print filemanager summary
	output.PrintSummary()

	// Finalize output
	output.Finalize()
}

// Prints a simple figlet-style ASCII art and a random quote
func printMOTD() {
	randomQuotes := []string{
		"PanCAP: Analyzer for capture files",
		"PanCAP: Analyzer for pancake files",
		"You want some syrup with these packets?",
		"Check out CONTRIBUTORS.md!",
		"Push your commits to git.darknebu.la/maride/pancap",
		"Don't let the white noise traffic confuse you.",
		"Grab a Club Mate if you don't have one yet.",
		"In Soviet Russia, traffic analyzes you.",
		"Who captures the captors?",
		"Respect other's privacy. Always.",
		"Make public data available, protect private data.", // https://www.ccc.de/en/hackerethik
	}

	// Maybe switch to urand? Possibly a security issue... ;)
	rand.Seed(time.Now().Unix())

	fmt.Println(" _ __   __ _ _ __   ___ __ _ _ __")
	fmt.Println("| '_ \\ / _` | '_ \\ / __/ _` | '_ \\")
	fmt.Println("| |_) | (_| | | | | (_| (_| | |_) |")
	fmt.Println("| .__/ \\__,_|_| |_|\\___\\__,_| .__/")
	fmt.Println("|_|                         |_|")
	fmt.Println(randomQuotes[rand.Intn(len(randomQuotes))])
	fmt.Println("")
}
