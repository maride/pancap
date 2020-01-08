package output

import (
	"fmt"
	"github.com/fatih/color"
	"strings"
)

const (
	MaxContentLines = 50
	SnipMark = "----- cut at 50 entries -----"
)

var (
	DidSnip bool
	DidAvoidEmptyBlock bool
)

// Prints a block of information with the given headline
// If content is empty, printing the headline is omitted.
// If the content is longer than MaxContentLines, content is cut.
func PrintBlock(headline string, content string) {
	// Avoid printing empty blocks - at least if user didn't specify it otherwise
	if len(content) == 0 && !*printEmptyBlocks {
		// No content and we are not forced to print empty blocks, return
		DidAvoidEmptyBlock = true
		return
	}

	// Print a newline to add some space between blocks
	fmt.Println("")

	// Check if we need to print a headline
	if len(content) > 0 {
		// We have content, we can print the headline
		headlineColor := color.New(color.FgRed, color.Bold)
		headlineColor.Println(headline)
	}

	// Cut to MaxContentLines if required
	if !(*fullOutput) {
		// User states that they don't want to see the whole output - cut content.
		content = cutContent(content)
	}

	// Print our content.
	fmt.Print(content)

	// Check if we need to add a newline character at the end
	if !strings.HasSuffix(content, "\n") {
		fmt.Print("\n")
	}
}

// Cut content after MaxContentLines lines
func cutContent(content string) string {
	numNewlines := 0

	// iterate over every character
	for i, c := range content {
		// check if character is newline
		if c == '\n' {
			// it is, count occurrence
			numNewlines++

			// Check if we already hit our limit yet
			if numNewlines == MaxContentLines {
				// Found nth newline, return content up to this position and add a notice about it.
				DidSnip = true
				return addSnipMark(content[:i])
			}
		}
	}

	// We are done before reaching the cut limit; return the whole content
	return content
}

// Adds a notice about the snipping process
func addSnipMark(content string) string {
	printer := color.New(color.Bold, color.BgBlack)
	return content + "\n" + printer.Sprint(SnipMark) + "\n"
}
