package output

import (
	"crypto/sha256"
	"fmt"
)

type File struct {
	name    string
	content []byte
	origin  string
	hash    string
}

// Creates a new file object and calculates the hash of the given content
func NewFile(name string, content []byte, origin string) File {
	hash := fmt.Sprintf("%x", sha256.Sum256(content))[:6]
	return File{
		name:    name,
		content: content,
		origin:  origin,
		hash:    hash,
	}
}
