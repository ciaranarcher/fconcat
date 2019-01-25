package fconcat

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"

	"github.com/pkg/errors"
)

// Concat will take an array of URLs, fetch each in turn
// and return a handle to a file.
func Concat(urls []string) (string, error) {

	tmpFiles, err := fetchFiles(urls)

	if err != nil {
		fmt.Println("error fetching files: ", err)
		return "", err
	}

	outputFile, err := concatenateFiles(tmpFiles)

	if err != nil {
		fmt.Println("error fetching files: ", err)
		return "", err
	}

	fmt.Println("outputFile:", outputFile)

	return outputFile, nil
}

// fetchFiles will take a string of URLs and attempt to download and
// save each as a local tempfile. The array of files created is returned.
func fetchFiles(urls []string) ([]string, error) {
	saved := make([]string, len(urls))
	var wg sync.WaitGroup

	// Fetch files in parallel and wait for all to finish using a waitGroup.
	// We save the f to a TempFile, and store the file path in the correct
	// order in the saved array.
	for i, f := range urls {
		wg.Add(1)

		go func(f string, i int) {
			defer wg.Done()

			resp, err := http.Get(f)
			if err != nil {
				fmt.Println(fmt.Sprintf("fetching file %s (GET) failed", f))
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != 200 {
				fmt.Println(fmt.Sprintf("unable to fetch %s: %d %s", f, resp.StatusCode, resp.Status))
				return
			}

			// Create temp output file
			fname, err := parseFilenameFromURL(f)
			if err != nil {
				fmt.Println(fmt.Sprintf("unable to parse file %s: %o ", fname, err))
				return
			}

			out, err := ioutil.TempFile("", fname)
			defer out.Close()
			if err != nil {
				fmt.Println(fmt.Sprintf("unable to create temp file %s", fname))
				return
			}

			_, err = io.Copy(out, resp.Body)
			if err != nil {
				fmt.Println("error copying response to temp file")
				return
			}

			saved[i] = out.Name()
		}(f, i)
	}

	wg.Wait()
	return saved, nil
}

// parseFilenameFromURL will take a URL like http://some/url/file.mp3 and return just
// the file part, e.g. file.mp3.
func parseFilenameFromURL(urlStr string) (string, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return "", err
	}
	parsedPath := strings.Split(u.Path, "/")
	return parsedPath[len(parsedPath)-1], nil
}

func concatenateFiles(files []string) (string, error) {
	// Open output / destination file.
	outputFile, err := ioutil.TempFile("", "output")
	if err != nil {
		return "", errors.Wrap(err, "error opening the temp output file")
	}

	// Close the output file when we are done.
	defer func() {
		if err := outputFile.Close(); err != nil {
			fmt.Println("error closing temp file")
		}
	}()

	// Create a new buffered writer.
	writeBuffer := bufio.NewWriter(outputFile)

	// Range over our files.
	for _, f := range files {
		// Open each.
		fHandle, err := os.Open(f)
		if err != nil {
			return "", err
		}

		// Remember to close it later!
		defer func() {
			if err := fHandle.Close(); err != nil {
				fmt.Println("error closing temp file")
			}
		}()

		// Create a new buffered reader
		readBuffer := bufio.NewReader(fHandle)

		// Create a buffer to keep chunks that are read.
		buf := make([]byte, 1024)
		for {
			// Read a chunk
			n, err := readBuffer.Read(buf)
			if err != nil && err != io.EOF {
				return "", err
			}
			if n == 0 {
				break
			}

			// Write a chunk
			if _, err := writeBuffer.Write(buf[:n]); err != nil {
				return "", err
			}
		}

		// We're done with this recording, so we can delete it.
		err = os.Remove(f)
		if err != nil {
			return "", err
		}
	}

	// Close the write buffer.
	if err = writeBuffer.Flush(); err != nil {
		return "", err
	}

	// Return the output file name.
	return outputFile.Name(), nil
}
