package output

import "github.com/fatih/color"

// Called at the very end, before terminating pancap
func Finalize() {
	printer := color.New(color.Bold, color.BgBlack)

	// Check if we snipped, to add a notice how to show the whole block
	if DidSnip {
		// We snipped - inform user about this process
		printer.Println("Output is snipped at one or more positions. Add --full-output to avoid snipping.")
	}

	// Check if we skipped printing an empty block
	if DidAvoidEmptyBlock {
		// We did - inform user about this
		printer.Println("Some submodule output was hidden. Add --print-empty-blocks to show it.")
	}
}
