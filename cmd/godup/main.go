package main

import (
	"bytes"
	"crypto/md5"
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

func validateGroup(files []*godup.FileInfo) (results map[string]*godup.FileInfo, err error) {
	results = make(map[string]*godup.FileInfo)

	for _, file := range files {
		file.Hash, err = helper.Hash(file.Path)
		if err != nil {
			return results, err
		}
	}

	for i := 0; i < len(files); i++ {
		for j := i + 1; j < len(files); j++ {
			if bytes.Equal(files[i].Hash, files[j].Hash) {
				_x := md5.Sum([]byte(files[i].Name + files[i].Path))
				fi := _x[:]
				if _, ok := results[string(fi)]; !ok {
					results[string(fi)] = files[i]
				}

				_x = md5.Sum([]byte(files[j].Name + files[j].Path))
				fj := _x[:]
				if _, ok := results[string(fj)]; !ok {
					results[string(fj)] = files[j]
				}
			}
		}
	}

	return results, nil
}

func compareBytes(files map[string]*godup.FileInfo) (results map[string]*godup.FileInfo, err error) {
	results = make(map[string]*godup.FileInfo)

	if len(files) < 2 {
		return results, nil
	}

	for ki, i := range files {
		for kj, j := range files {
			if ki == kj {
				continue
			}

			f1, err := ioutil.ReadFile(i.Path)
			if err != nil {
				panic(err)
			}

			f2, err := ioutil.ReadFile(j.Path)
			if err != nil {
				panic(err)
			}

			if bytes.Equal(f1, f2) {
				if _, ok := results[ki]; !ok {
					results[ki] = i
				}
				if _, ok := results[kj]; !ok {
					results[kj] = j
				}
			}
		}
	}

	return results, nil
}
