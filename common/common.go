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

// Generates a small ASCII tree for the given string array
func GenerateTree(strarr []string) string {
	tmpstr := ""

	// iterate over each element
	for iter, elem := range strarr {
		// check if we got the last element
		if iter < len(strarr) - 1 {
			tmpstr = fmt.Sprintf("%s|- %s\n", tmpstr, elem)
		} else {
			tmpstr = fmt.Sprintf( "%s'- %s\n", tmpstr, elem)
		}
	}

	// Return constructed (grown?) tree
	return tmpstr
}
