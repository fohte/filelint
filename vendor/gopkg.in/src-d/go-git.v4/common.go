package git

import (
	"strings"

	"srcd.works/go-git.v4/config"
	"srcd.works/go-git.v4/plumbing/storer"
)

// Storer is a generic storage of objects, references and any information
// related to a particular repository. The package srcd.works/go-git.v4/storage
// contains two implementation a filesystem base implementation (such as `.git`)
// and a memory implementations being ephemeral
type Storer interface {
	storer.EncodedObjectStorer
	storer.ReferenceStorer
	storer.ShallowStorer
	storer.IndexStorer
	config.ConfigStorer
}

// countLines returns the number of lines in a string à la git, this is
// The newline character is assumed to be '\n'.  The empty string
// contains 0 lines.  If the last line of the string doesn't end with a
// newline, it will still be considered a line.
func countLines(s string) int {
	if s == "" {
		return 0
	}

	nEOL := strings.Count(s, "\n")
	if strings.HasSuffix(s, "\n") {
		return nEOL
	}

	return nEOL + 1
}
