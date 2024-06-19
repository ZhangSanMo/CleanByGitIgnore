package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/monochromegane/go-gitignore"
)

var SEP = string(os.PathSeparator)

func main() {
	// check args
	if len(os.Args) < 2 {
		fmt.Println("Usage: clean.exe <dir>")
		return
	}
	dir, err := filepath.Abs(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}
	dirInfo, err := os.Stat(dir)
	if err != nil {
		fmt.Println("No such file or directory")
		return
	}
	if !dirInfo.IsDir() {
		fmt.Println("Must be directory")
		return
	}
	ignoreFilePath := dir + SEP + ".gitignore"
	matcher, err := gitignore.NewGitIgnore(ignoreFilePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	xx := walkDir(dir, []string{}, []string{}, matcher)
	deleteWaitlist := []string{}
	for _, file := range xx {
		fileinfo, err := os.Stat(file)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if matcher.Match(file, fileinfo.IsDir()) {
			fmt.Println(file)
			deleteWaitlist = append(deleteWaitlist, file)
		}
	}
	if len(deleteWaitlist) == 0 {
		fmt.Println("No files to delete")
		return
	}
	// ask usr if confirm to delete
	fmt.Println("Do you want to delete these files? (y/n)")
	var input string
	fmt.Scanln(&input)
	if input == "y" {
		for _, file := range deleteWaitlist {
			err := os.RemoveAll(file)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

// walk dir, if match gitignore, skip ignored sub folder
func walkDir(dir string, domain []string, collections []string, mathcer gitignore.IgnoreMatcher) []string {
	// copy domain
	domainCopy := make([]string, len(domain)+1)
	copy(domainCopy, domain)
	domainCopy[len(domainCopy)-1] = dir
	path := strings.Join(domainCopy, SEP) + SEP
	files, err := os.ReadDir(path)
	if err != nil {
		fmt.Println(err)
		return collections
	}
	for _, file := range files {
		// domain + sub
		domainCopyCopy := make([]string, len(domainCopy)+1)
		copy(domainCopyCopy, domainCopy)
		domainCopyCopy[len(domainCopyCopy)-1] = file.Name()
		filePath := strings.Join(domainCopyCopy, SEP)
		collections = append(collections, filePath)
		// if exculde, skip sub folder
		isDir := file.IsDir()
		if isDir && !mathcer.Match(filePath, isDir) {
			collections = walkDir(file.Name(), domainCopy, collections, mathcer)
		}
	}
	return collections
}
