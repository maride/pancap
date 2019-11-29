package common

import "fmt"

// Appends the appendee to the array if it does not contain appendee yet
func AppendIfUnique(appendee string, array []string) []string {
	// Iterate over all elements and check values
	for _, elem := range array {
		if elem == appendee {
			// ... found. Stop here
			return array
		}
	}

	// None found, append
	return append(array, appendee)
}

// Prints each element, along with a small ASCII tree
func PrintTree(strarr []string) {
	// iterate over each element
	for iter, elem := range strarr {
		// check if we got the last element
		if iter < len(strarr) - 1 {
			fmt.Printf("|- %s\n", elem)
		} else {
			fmt.Printf("'- %s\n\n", elem)
		}
	}
}
