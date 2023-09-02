package output

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/maride/pancap/common"
)

var (
	registeredFiles []File
	notFound        []string
	extractedFiles  int
)

// Registers a file with the given name and content.
// This function takes care of filesystem I/O handling and flag parsing.
// This means that a module should _always_ call this function when a file is encountered.
// origin is a descriptive string where the file comes from, e.g. the module name.
func RegisterFile(filename string, content []byte, origin string) {
	// Check if there even is anything to register
	if len(content) == 0 {
		// File is empty, won't register the void
		log.Printf("Avoided registering file from %s because it is empty.", origin)
		return
	}
	thisFile := NewFile(filename, content, origin)
	// To avoid doubles, we need to check if that hash is already present
	for _, f := range registeredFiles {
		if f.hash == thisFile.hash {
			// Found - stop here
			log.Printf("Avoided registering file from %s because it has the same content as an already registered file ", origin)
			return
		}
	}

	// None found, add to list
	registeredFiles = append(registeredFiles, thisFile)
}

// Iterates over all registered files and checks if they should be extracted and stored, and does exactly that.
func StoreFiles() {
	var filesToExtract []File

	// Check different flag scenarios
	if *targetAllFiles {
		// We should extract all files.
		filesToExtract = registeredFiles
	} else {
		// We should extract only a given set of files
		fileList := strings.Split(*targetFiles, ",")
		for _, f := range fileList {
			// Iterate over desired files
			found := false
			for _, a := range registeredFiles {
				// Iterate over available (registered) files
				if f == a.hash {
					// Found the file
					found = true
					filesToExtract = append(filesToExtract, a)
					break
				}
			}

			if !found {
				// No file found, notify user
				notFound = append(notFound, fmt.Sprintf("File with hash %s requested but not found.", f))
			}
		}
	}

	// Iterate over all target files and write it them out
	for _, f := range filesToExtract {
		writeOut(f)
	}
}

// Writes the given file object to disk, along with a stats file placed next to it.
func writeOut(f File) {
	targetName := fmt.Sprintf("%s%c%s", *targetOutput, os.PathSeparator, f.hash)
	targetDescName := fmt.Sprintf("%s.info", targetName)
	targetDescription := fmt.Sprintf("Filename: %s\nHash: %s\nOrigin: %s\nSize: %d", f.name, f.hash, f.origin, len(f.content))

	// Write target file
	targetWriteErr := ioutil.WriteFile(targetName, f.content, 0644)
	if targetWriteErr != nil {
		log.Printf("Unable to write file %s: %s", targetName, targetWriteErr.Error())
		return
	}

	// Write stats file
	statsWriteErr := ioutil.WriteFile(targetDescName, []byte(targetDescription), 0644)
	if statsWriteErr != nil {
		log.Printf("Unable to write file %s: %s", targetName, targetWriteErr.Error())
		return
	}

	// Raise stats
	extractedFiles++
}

// Prints a brief summary about the extracted files
func PrintSummary() {
	summary := fmt.Sprintf("%d files found in stream.\n%d files extracted from stream.", len(registeredFiles), extractedFiles)

	// Generate list of found files
	var strFileList []string
	for _, f := range registeredFiles {
		name := f.name
		if name == "" {
			name = "(no name found)"
		}

		strFileList = append(strFileList, fmt.Sprintf("%s: %s (%s), %d bytes", f.hash, name, f.origin, len(f.content)))
	}

	// Print list of files as a tree
	if len(strFileList) > 0 {
		summary += "\nFound files:"
		summary += "\n" + common.GenerateTree(strFileList)
	}

	// Check if we left a few requested files unanswered
	if len(notFound) > 0 {
		summary += "\nUnable to find requested file(s) " + strings.Join(notFound, ", ")
	}

	// Print constructed summary
	PrintBlock("Files", summary)
}
