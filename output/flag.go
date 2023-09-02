package output

import "flag"

var (
	fullOutput       *bool
	printEmptyBlocks *bool
	targetFiles      *string
	targetAllFiles   *bool
	targetOutput     *string
	graphOutput      *string
)

func RegisterFlags() {
	fullOutput = flag.Bool("full-output", false, "Show full output instead of limiting submodule output")
	printEmptyBlocks = flag.Bool("print-empty-blocks", false, "Prints blocks (submodule output) even if the submodule doesn't have any content to print.")
	targetFiles = flag.String("extract-these", "", "Comma-separated list of files to extract.")
	targetAllFiles = flag.Bool("extract-all", false, "Extract all files found.")
	targetOutput = flag.String("extract-to", "./extracted", "Directory to store extracted files in.")
	graphOutput = flag.String("create-graph", "", "Create a Graphviz graph out of collected communication")
}
