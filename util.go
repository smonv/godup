package main

import (
	"crypto/sha256"
	"io/ioutil"
	"reflect"
)

func hash(path string) (result []byte, err error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}

	r := sha256.Sum256(data)
	result = r[:]
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
