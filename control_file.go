package lpd

import (
	"io"
)

type ControlFile struct {
	File []byte
}

func NewControlFile() *ControlFile {
	cFile := new(ControlFile)

	return cFile
}

func readControlFile(reader io.Reader, numBytes uint64) (*ControlFile, error) {
	file := make([]byte, numBytes)

	_, err := io.ReadFull(reader, file)

	if err != nil {
		return nil, err
	}

	// For now, absorb control file so we can play with data first
	cFile := new(ControlFile)
	cFile.File = file

	return cFile, nil
}
