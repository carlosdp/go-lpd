package lpd

import (
	"io"
	"io/ioutil"
	"os"
)

func readDataFile(reader io.Reader, numBytes uint64) (*os.File, error) {
	tempFile, err := ioutil.TempFile(os.TempDir(), "go-lpd")

	if err != nil {
		tempFile.Close()
		return nil, err
	}

	_, err = io.CopyN(tempFile, reader, int64(numBytes))

	if err != nil {
		tempFile.Close()
		return nil, err
	}

	return tempFile, nil
}
