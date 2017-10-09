package helper

import (
	"testing"

	"github.com/tthanh/godup"
)

func TestCheckFilesContain(t *testing.T) {
	f1 := &godup.FileInfo{
		Name: "f1",
		Size: 123456,
		Path: "p1",
	}
	f2 := &godup.FileInfo{
		Name: "f2",
		Size: 123457,
		Path: "p2",
	}
	f3 := &godup.FileInfo{
		Name: "f3",
		Size: 123458,
		Path: "p4",
	}

	files := []*godup.FileInfo{f1, f2, f3}

	result := checkFilesContain(files, f1)
	if !result {
		t.Fatalf("CheckFilesContain failed")
	}
}
