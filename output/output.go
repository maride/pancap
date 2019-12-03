package output

import (
	"flag"
	"github.com/fatih/color"
)

var (
	fullOutput *bool
)

func RegisterFlags() {
	fullOutput = flag.Bool("full-output", false, "Show full output instead of limiting submodule output")
}

// Called at the very end, before terminating pancap
func Finalize() {
	// Check if we snipped, to add a notice how to show the whole block
	if DidSnip {
		// We snipped - inform user about this process
		printer := color.New(color.Bold, color.BgBlack)
		printer.Print("\nOutput is snipped at one or more positions. Add --full-output to avoid snipping.")
	}
}
