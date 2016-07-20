package main

import (
	"bytes"
	"io/ioutil"
	"reflect"
)

func compareWorker(done chan struct{}, cic <-chan []*File, coc chan<- []*File) {
	for files := range cic {
		select {
		case coc <- compare(files):
		case <-done:
			return
		}
	}
}

func compare(files []*File) (result []*File) {
	if len(files) < 2 {
		result = files
		return
	}
	result = compareHash(files)

	if len(result) < 2 {
		return
	}
	result = compareByte(files)

	return
}

func compareHash(files []*File) (result []*File) {
	for _, f := range files {
		r, err := hash(f.Path)
		if err != nil {
			panic(err)
		}
		f.Hash = r
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
