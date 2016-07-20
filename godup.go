package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

var (
	allFile map[int64][]*File
	mutex   sync.Mutex
	wg      sync.WaitGroup
)

func main() {
	paths := os.Args[1:]

	allFile = make(map[int64][]*File)

	if len(paths) < 1 {
		fmt.Println("path not found")
		return
	}

	for _, path := range paths {
		err := check(path)
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
	cic := make(chan []*File) // compare input channel
	coc := make(chan []*File) // compare output channel

	workers := runtime.NumCPU()
	wg.Add(workers)

	for i := 0; i < workers; i++ {
		go func() {
			compareWorker(done, cic, coc)
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(coc)
	}()

	go func() {
		defer close(cic)
		for _, files := range allFile {
			if len(files) > 1 {
				cic <- files
			}
		}
	}()

	for files := range coc {
		if len(files) > 1 {
			fmt.Printf("\n")
			fmt.Printf("Size: %d. HASH: %x\n", files[0].Size, files[0].Hash)
			for _, file := range files {
				fmt.Printf("Path: %s\n", file.Path)
			}
		}
	}

	defer close(done)
}

func check(path string) error {
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
