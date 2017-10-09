package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/tthanh/godup"
	"github.com/tthanh/godup/pkg/helper"
)

var (
	logger = log.New(os.Stderr, "", 0)
)

func main() {
	groups := make(map[int64][]*godup.FileInfo)
	paths := os.Args[1:]

	if len(paths) < 1 {
		logger.Fatal("path not found")
		return
	}

	for _, path := range paths {
		src, err := os.Stat(path)
		if err != nil {
			logger.Fatal(err)
		}
		if !src.IsDir() {
			logger.Fatalf("%s is not directory", path)
		}
	}

	for _, path := range paths {
		path, err := filepath.EvalSymlinks(path)
		if err != nil {
			logger.Fatal(err)
		}

		err = filepath.Walk(path, func(path string, fi os.FileInfo, err error) error {
			if !fi.IsDir() {
				groups[fi.Size()] = append(groups[fi.Size()], &godup.FileInfo{
					Name: fi.Name(),
					Size: fi.Size(),
					Path: path,
				})
			}

			return nil
		})

		if err != nil {
			logger.Fatal(err)
		}
	}

	fmt.Printf("Found %d group of candidate\n", len(groups))

	for _, group := range groups {
		result, err := validateGroup(group)
		if err != nil {
			log.Fatal(err)
		}

		if len(result) > 1 {
			finalResult, sErr := compareBytes(result)
			if sErr != nil {
				log.Fatal(err)
			}

			if len(finalResult) > 1 {
				fmt.Println("")
				for _, r := range result {
					fmt.Printf("%s\t%x\t%s\n", r.Name, r.Hash, r.Path)
				}
			}
		}
	}
}

func validateGroup(files []*godup.FileInfo) (result []*godup.FileInfo, err error) {
	for _, file := range files {
		file.Hash, err = helper.Hash(file.Path)
		if err != nil {
			return result, err
		}
	}

	for i := 0; i < len(files); i++ {
		for j := i + 1; j < len(files); j++ {
			if bytes.Equal(files[i].Hash, files[j].Hash) {
				result = helper.AppendNotExistFile(result, files[i])
				result = helper.AppendNotExistFile(result, files[j])
			}
		}
	}

	return result, nil
}

func compareBytes(files []*godup.FileInfo) (result []*godup.FileInfo, err error) {
	if len(files) < 2 {
		return result, nil
	}

	for i := 0; i < len(files); i++ {
		for j := i + 1; j < len(files); j++ {
			f1, err := ioutil.ReadFile(files[i].Path)
			if err != nil {
				panic(err)
			}

			f2, err := ioutil.ReadFile(files[j].Path)
			if err != nil {
				panic(err)
			}

			if bytes.Equal(f1, f2) {
				result = helper.AppendNotExistFile(result, files[i])
				result = helper.AppendNotExistFile(result, files[j])
			}
		}
	}

	return result, nil
}
