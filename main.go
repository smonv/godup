package main

import (
	"bytes"
	"crypto/md5"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
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

	for _, files := range allFile {
		if len(files) > 1 {
			files = compareHash(files)
			files = compareByte(files)
			if len(files) > 1 {
				fmt.Println("+++o")
				for _, file := range files {
					fmt.Println(file.Path)
				}
			}
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

func compareHash(files []*File) []*File {
	sameHash := []*File{}

	for _, f := range files {
		hash, err := computeMd5(f.Path)
		if err != nil {
			panic(err)
		}
		f.Hash = hash
	}

	for _, i := range files {
		for _, j := range files {
			if !reflect.DeepEqual(i, j) {
				if bytes.Compare(i.Hash, j.Hash) == 0 && !checkFilesContain(sameHash, j) && !checkFilesContain(sameHash, i) {
					sameHash = append(sameHash, i)
					sameHash = append(sameHash, j)
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

func computeMd5(path string) ([]byte, error) {
	var result []byte
	file, err := os.Open(path)
	if err != nil {
		return result, err
	}
	digest := md5.New()
	if _, err := io.Copy(digest, file); err != nil {
		return result, err
	}
	return digest.Sum(result), nil
}

func compareByte(files []*File) []*File {
	sameByte := []*File{}
	for _, i := range files {
		for _, j := range files {
			if reflect.DeepEqual(i, j) {
				break
			}
			if !checkFilesContain(sameByte, i) && !checkFilesContain(sameByte, j) {
				f1, err := ioutil.ReadFile(i.Path)
				if err != nil {
					panic(err)
				}
				f2, err := ioutil.ReadFile(j.Path)
				if err != nil {
					panic(err)
				}
				if bytes.Compare(f1, f2) == 0 {
					sameByte = append(sameByte, i)
					sameByte = append(sameByte, j)
				}
			}
		}
	}
	return sameByte
}
