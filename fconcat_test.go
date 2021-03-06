package fconcat_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/ciaranarcher/fconcat"
)

var bucket = os.Getenv("BUCKET")

var urlsTwo = []string{
	fmt.Sprintf("%s/one.mp3", bucket),
	fmt.Sprintf("%s/two.mp3", bucket),
}

var urlsThree = []string{
	fmt.Sprintf("%s/one.mp3", bucket),
	fmt.Sprintf("%s/two.mp3", bucket),
	fmt.Sprintf("%s/three.mp3", bucket),
}

func TestConcat(t *testing.T) {
	fName, err := fconcat.Concat(urlsThree)

	if err != nil {
		t.Error("unable to concatenate files: ", err)
	}

	f, err := os.Open(fName)
	if err != nil {
		t.Error("unable to open concatenated file.")
	}

	fi, err := f.Stat()
	if err != nil {
		t.Error("unable to get info about concatenated file.")
	}

	expectedSize := int64(870192)

	if fi.Size() != expectedSize {
		t.Error(fmt.Sprintf("byte size incorrect: expected %d and got %d", expectedSize, fi.Size()))
	}
}

func BenchmarkConcat2(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_, _ = fconcat.Concat(urlsTwo)
	}
}

func BenchmarkConcat3(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_, _ = fconcat.Concat(urlsThree)
	}
}
