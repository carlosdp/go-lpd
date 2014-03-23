package lpd

import (
	"bufio"
	"errors"
	"io"
	"io/ioutil"
	"os"
)

// PrintJob describes a print job for a queue in an LPD printer with a
// control file and data file.
type PrintJob struct {
	QueueName   string
	ControlFile *ControlFile
	DataFile    *os.File
}

// NewPrintJob returns a PrintJob configured for a queue with a data file.
// NewPrintJob returns an error if and only if the data file fails to copy
// to a temporary file.
func NewPrintJob(queue string, dataFile io.Reader) (*PrintJob, error) {
	job := new(PrintJob)
	job.QueueName = queue
	job.ControlFile = NewControlFile()

	tempFile, err := ioutil.TempFile(os.TempDir(), "go-lpd")

	if err != nil {
		return nil, err
	}

	_, err = io.Copy(tempFile, dataFile)

	if err != nil {
		return nil, err
	}

	job.DataFile = tempFile

	return job, nil
}

// Handle series of subcommands that describe a print job
func receiveJob(reader io.Reader, writer io.Writer) (*PrintJob, error) {
	job := new(PrintJob)

	bufReader := bufio.NewReader(reader)

	for {
		select {
		default:
			rawSubCommand, err := bufReader.ReadBytes(0x10)

			if err != nil {
				if job.ControlFile != nil && job.DataFile != nil {
					return job, nil
				}

				return nil, err
			}

			subCmd, err := unmarshalSubCommand(rawSubCommand)

			if err != nil {
				return nil, err
			}

			switch subCmd.Code {
			case AbortJob:
				ackSubCommand(writer)
				return nil, errors.New("job aborted")
			case ReceiveControlFile:
				cFile, err := readControlFile(reader, subCmd.NumBytes)

				if err != nil {
					return nil, err
				}

				ackSubCommand(writer)

				job.ControlFile = cFile
			case ReceiveDataFile:
				dataFile, err := readDataFile(reader, subCmd.NumBytes)

				if err != nil {
					return nil, err
				}

				ackSubCommand(writer)

				job.DataFile = dataFile
			}
		}
	}
}

// Send an octect of 0 to acknowledge a subcommand
func ackSubCommand(writer io.Writer) error {
	_, err := writer.Write([]byte{0x0})

	if err != nil {
		return err
	}

	return nil
}
