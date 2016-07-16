package main

import (
	"bytes"
	"crypto/sha1"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"sync"
)

var (
	allFile map[int64][]*File
)

// File struct
type File struct {
	Name string
	Size int64
	Path string
	Hash []byte
}

func main() {
	var cPath string
	flag.StringVar(&cPath, "p", "", "check path")
	flag.Parse()

	if len(cPath) == 0 {
		fmt.Println("check path not found")
		os.Exit(1)
	}

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

	wg := &sync.WaitGroup{}
	for _, files := range allFile {
		if len(files) > 1 {
			wg.Add(1)
			go func(files []*File) {
				defer wg.Done()
				sameHash := compareHash(files)
				sameByte := compareByte(sameHash)
				if len(sameByte) > 1 {
					fmt.Printf("SHA1: %x\n", sameByte[0].Hash)
					for _, file := range files {
						fmt.Println(file.Path)
					}
				}
			}(files)
		}
	}
	wg.Wait()
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

func compareHash(files []*File) []*File {
	sameHash := []*File{}

	for _, f := range files {
		hash, err := computeHash(f.Path)
		if err != nil {
			panic(err)
		}
		f.Hash = hash
	}

	for _, i := range files {
		for _, j := range files {
			if !reflect.DeepEqual(i, j) {
				if bytes.Equal(i.Hash, j.Hash) {
					if !checkFilesContain(sameHash, i) {
						sameHash = append(sameHash, i)
					}
					if !checkFilesContain(sameHash, j) {

						sameHash = append(sameHash, j)
					}
				}
			}
		}
	}
	return sameHash
}

func checkFilesContain(files []*File, file *File) bool {
	for _, f := range files {
		if reflect.DeepEqual(f, file) {
			return true
		}
	}
	return false
}

func computeHash(path string) (result []byte, err error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}

	hash := sha1.Sum(data)
	result = hash[:]
	return
}

func compareByte(files []*File) []*File {
	sameByte := []*File{}
	for _, i := range files {
		for _, j := range files {
			if reflect.DeepEqual(i, j) {
				break
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
				if !checkFilesContain(sameByte, i) {
					sameByte = append(sameByte, i)
				}
				if !checkFilesContain(sameByte, j) {
					sameByte = append(sameByte, j)
				}
			}
		}
	}
	return sameByte
}
