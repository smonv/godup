package walker

import "github.com/tthanh/godup"

type WalkFunc func(path string, fi godup.FileInfo, err error) error

func Walk(root string, walkFn WalkFunc) error {

	return nil
}
