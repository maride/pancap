package main

import (
	"flag"
	"fmt"
	"log"
)

func main() {
	// important things first
	printMOTD()

	// register flags
	registerFileFlags()
	flag.Parse()

	// Open the given PCAP
	packetSource, linkType, fileErr := openPCAP()
	if fileErr != nil {
		// Encountered problems with the PCAP - permission and/or existance error
		log.Fatalf("Error occured while opeining specified file: %s", fileErr.Error())
	}

	// Start analyzing
	analyzeErr := analyzePCAP(packetSource, linkType)
	if analyzeErr != nil {
		// Mh, encountered some problems while analyzing file
		log.Fatalf("Error occurred while analyzing: %s", analyzeErr.Error())
	}
}

// Prints a simple figlet-style ASCII art
func printMOTD() {
	fmt.Println(" _ __   __ _ _ __   ___ __ _ _ __")
	fmt.Println("| '_ \\ / _` | '_ \\ / __/ _` | '_ \\")
	fmt.Println("| |_) | (_| | | | | (_| (_| | |_) |")
	fmt.Println("| .__/ \\__,_|_| |_|\\___\\__,_| .__/")
	fmt.Println("|_|                         |_|")
	fmt.Println("PanCAP: Analyzer for capture files\n")
}
