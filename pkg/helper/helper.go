package helper

import (
	"crypto/sha256"
	"os"
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
