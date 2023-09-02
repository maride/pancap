package output

import "flag"

var (
	fullOutput       bool
	printEmptyBlocks bool
	targetFiles      string
	targetAllFiles   bool
	targetOutput     string
	graphOutput      string
)

func RegisterFlags() {
	flag.BoolVar(&fullOutput, "full-output", false, "Show full output instead of limiting submodule output")
	flag.BoolVar(&printEmptyBlocks, "print-empty-blocks", false, "Prints blocks (submodule output) even if the submodule doesn't have any content to print.")
	flag.StringVar(&targetFiles, "extract-these", "", "Comma-separated list of files to extract.")
	flag.BoolVar(&targetAllFiles, "extract-all", false, "Extract all files found.")
	flag.StringVar(&targetOutput, "extract-to", "./extracted", "Directory to store extracted files in.")
	flag.StringVar(&graphOutput, "create-graph", "", "Create a Graphviz graph out of collected communication")
}
