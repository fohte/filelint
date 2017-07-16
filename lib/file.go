package lib

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	"github.com/andrew-d/isbinary"
)

var (
	ErrNotGitRepository = errors.New("Not a git repository")
)

func FindGitRootPath(path string) (string, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	if IsExist(filepath.Join(path, ".git")) {
		return path, nil
	}

	parentPath := filepath.Join(path, "..")
	if path == parentPath { // if root directory
		return "", ErrNotGitRepository
	}

	return FindGitRootPath(parentPath)
}

func FindGitIgnore() (string, error) {
	gitRoot, err := FindGitRootPath(".")
	if err != nil {
		return "", err
	}
	return filepath.Join(gitRoot, ".gitignore"), nil
}

func IsExist(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func IsDir(path string) bool {
	s, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return s.IsDir()
}

func IsBinary(path string) bool {
	fp, err := os.Open(path)
	if err != nil {
		return false
	}
	defer fp.Close()

	buf, err := ioutil.ReadAll(fp)
	if err != nil {
		return false
	}

	return isbinary.Test(buf)
}

func FindTextFiles(paths []string) []string {
	new := make([]string, 0, len(paths))

	for _, path := range paths {
		if IsDir(path) || IsBinary(path) {
			continue
		}
		new = append(new, path)
	}

	return new
}

func GetHomeDir() string {
	var dir string

	switch runtime.GOOS {
	default:
		dir = filepath.Join(os.Getenv("HOME"), ".config")
	case "windows":
		dir = os.Getenv("APPDATA")
		if dir == "" {
			dir = filepath.Join(os.Getenv("USERPROFILE"), "Application Data")
		}
	}

	return dir
}
