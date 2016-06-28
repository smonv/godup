package main

import (
	"fmt"
	"os"
	"path/filepath"
)

var (
	cPath   = "/run/media/tthanh/media/images"
	allFile map[int64][]*File
)

// File struct
type File struct {
	Name string
	Size int64
	Path string
}

func main() {
	allFile = make(map[int64][]*File)

	src, err := os.Stat(cPath)
	if err != nil {
		panic(err)
	}

	if src.IsDir() {
		fullPath, err := filepath.Abs(cPath)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Checking %s\n", fullPath)
		if err = filepath.Walk(fullPath, walker); err != nil {
			fmt.Println(err)
		}
	}

	for _, v := range allFile {
		for _, f := range v {
			fmt.Printf("Size: %d. Name: %s. Path: %s.\n", f.Size, f.Name, f.Path)
		}
	}
}

func walker(path string, fi os.FileInfo, err error) error {
	if !fi.IsDir() {
		file := &File{
			Name: fi.Name(),
			Size: fi.Size(),
			Path: path,
		}
		files := allFile[file.Size]
		files = append(files, file)
		allFile[file.Size] = files
	}
	return nil
}
