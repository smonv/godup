package godup

import "testing"

func TestCheckFilesContain(t *testing.T) {
	f1 := &File{
		Name: "f1",
		Size: 123456,
		Path: "p1",
	}
	f2 := &File{
		Name: "f2",
		Size: 123457,
		Path: "p2",
	}
	f3 := &File{
		Name: "f3",
		Size: 123458,
		Path: "p4",
	}

	files := []*File{f1, f2, f3}

	result := checkFilesContain(files, f1)
	if !result {
		t.Fatalf("CheckFilesContain failed")
	}
}
