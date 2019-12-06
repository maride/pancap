package output

import "flag"

var (
	fullOutput *bool
	printEmptyBlocks *bool
)

func RegisterFlags() {
	fullOutput = flag.Bool("full-output", false, "Show full output instead of limiting submodule output")
	printEmptyBlocks = flag.Bool("print-empty-blocks", false, "Prints blocks (submodule output) even if the submodule doesn't have any content to print.")
}


