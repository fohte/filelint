package main

import (
	"fmt"
	"os"
	"strings"

	"srcd.works/go-git.v4"
	"srcd.works/go-git.v4/plumbing/object"

	. "srcd.works/go-git.v4/_examples"
)

func main() {
	CheckArgs("<url> <path>")
	url := os.Args[1]
	path := os.Args[2]

	// Clone the given repository, creating the remote, the local branches
	// and fetching the objects, exactly as:
	Info("git clone %s %s", url, path)

	r, err := git.PlainClone(path, false, &git.CloneOptions{URL: url})
	CheckIfError(err)

	// Getting the latest commit on the current branch
	Info("git log -1")

	// ... retrieving the branch being pointed by HEAD
	ref, err := r.Head()
	CheckIfError(err)

	// ... retrieving the commit object
	commit, err := r.Commit(ref.Hash())
	CheckIfError(err)
	fmt.Println(commit)

	// List the tree from HEAD
	Info("git ls-tree -r HEAD")

	// ... retrieve the tree from the commit
	tree, err := commit.Tree()
	CheckIfError(err)

	// ... get the files iterator and print the file
	tree.Files().ForEach(func(f *object.File) error {
		fmt.Printf("100644 blob %s    %s\n", f.Hash, f.Name)
		return nil
	})

	// List the history of the repository
	Info("git log --oneline")

	commits, err := commit.History()
	CheckIfError(err)

	for _, c := range commits {
		hash := c.Hash.String()
		line := strings.Split(c.Message, "\n")
		fmt.Println(hash[:7], line[0])
	}
}
