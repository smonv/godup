package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"sync"
)

var (
	allFile map[int64][]*File
	mutex   sync.Mutex
	wg      sync.WaitGroup
)

func main() {
	defer os.Exit(1)

	paths := os.Args[1:]

	allFile = make(map[int64][]*File)

	if len(paths) < 1 {
		fmt.Println("path not found")
		return
	}

	for _, path := range paths {
		err := checkPath(path)
		if err != nil {
			fmt.Println(err)
		}
	}

	if len(allFile) == 0 {
		fmt.Println("cannot find any file")
		return
	}

	fmt.Printf("found %d files\n", len(allFile))

	done := make(chan struct{})
	hic := make(chan []*File) // hash input channel
	hoc := make(chan []*File) // hash output channel

	workers := runtime.NumCPU()
	wg.Add(workers)

	for i := 0; i < workers; i++ {
		go func() {
			hashWorker(done, hic, hoc)
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(hoc)
	}()

	go func() {
		defer close(hic)
		for _, files := range allFile {
			hic <- files
		}
	}()

	for files := range hoc {
		if len(files) > 1 {
			fmt.Printf("Size: %d\n", files[0].Size)
			for _, file := range files {
				fmt.Printf("Path: %s\n", file.Name)
			}
		}
	}

	defer close(done)

	// for _, files := range allFile {
	// 	if len(files) > 1 {
	// 		wg.Add(1)
	// 		go func(files []*File) {
	// 			defer wg.Done()
	// 			sameHash := compareHash(files)
	// 			sameByte := compareByte(sameHash)
	// 			if len(sameByte) > 1 {
	// 				fmt.Printf("\n")
	// 				fmt.Printf("SHA256: %x\n", sameByte[0].Hash)
	// 				for _, file := range files {
	// 					fmt.Println(file.Path)
	// 				}
	// 			}
	// 		}(files)
	// 	}
	// }
}

func hashWorker(done chan struct{}, hashc <-chan []*File, c chan<- []*File) {
	for files := range hashc {
		select {
		case c <- hash(files):
		case <-done:
			return
		}
	}
}

func hash(files []*File) []*File {
	if len(files) < 2 {
		return files
	}

	for _, file := range files {
		data, _ := ioutil.ReadFile(file.Path)
		// if err != nil {
		// 	return files
		// }

		hash := sha256.Sum256(data)
		file.Hash = hash[:]
	}
	return files
}

func checkPath(path string) error {
	src, err := os.Stat(path)
	if err != nil {
		return err
	}

	if src.IsDir() {
		fullPath, err := filepath.Abs(path)
		if err != nil {
			return err
		}
		fmt.Printf("Checking %s\n", fullPath)
		if err = filepath.Walk(fullPath, walker); err != nil {
			return err
		}
	}
	return nil
}

func walker(path string, fi os.FileInfo, err error) error {
	if !fi.IsDir() {
		mutex.Lock()
		defer mutex.Unlock()

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

func compareHash(files []*File) (result []*File) {
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
					if !checkFilesContain(result, i) {
						result = append(result, i)
					}
					if !checkFilesContain(result, j) {
						result = append(result, j)
					}
				}
			}
		}
	}
	return
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

	hash := sha256.Sum256(data)
	result = hash[:]
	return
}

func compareByte(files []*File) (result []*File) {
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
				if !checkFilesContain(result, i) {
					result = append(result, i)
				}
				if !checkFilesContain(result, j) {
					result = append(result, j)
				}
			}
		}
	}
	return
}
