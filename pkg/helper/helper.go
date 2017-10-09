package helper

import (
	"crypto/sha256"
	"os"
	"reflect"

	"github.com/tthanh/godup"
)

// Hash ...
func Hash(path string) (result []byte, err error) {
	buf := make([]byte, 16)

	f, err := os.Open(path)
	if err != nil {
		return result, err
	}
	defer func() {
		err := f.Close()
		if err != nil {
			panic(err)
		}
	}()

	_, err = f.Read(buf)
	if err != nil {
		return result, err
	}

	r := sha256.Sum256(buf)
	result = r[:]
	return result, nil
}

func checkFilesContain(files []*godup.FileInfo, file *godup.FileInfo) bool {
	for _, f := range files {
		if reflect.DeepEqual(f, file) {
			return true
		}
	}
	return false
}

// AppendNotExistFile ...
func AppendNotExistFile(files []*godup.FileInfo, file *godup.FileInfo) []*godup.FileInfo {
	if !checkFilesContain(files, file) {
		files = append(files, file)
	}
	return files
}
