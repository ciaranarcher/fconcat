package fconcat_test

import (
	"os"
	"testing"

	"github.com/ciaranarcher/fconcat"
)

var urlsTwo = []string{
	"https://s3-us-west-2.amazonaws.com/labday-eocarroll-nagius/123/call_recordings/123/leg_recordings/one.mp3",
	"https://s3-us-west-2.amazonaws.com/labday-eocarroll-nagius/123/call_recordings/123/leg_recordings/two.mp3",
}

var urlsThree = []string{
	"https://s3-us-west-2.amazonaws.com/labday-eocarroll-nagius/123/call_recordings/123/leg_recordings/one.mp3",
	"https://s3-us-west-2.amazonaws.com/labday-eocarroll-nagius/123/call_recordings/123/leg_recordings/two.mp3",
}

func TestConcat(t *testing.T) {
	fName, err := fconcat.Concat(urlsTwo)

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

	if fi.Size() != 625372 {
		t.Error("size incorrect")
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
