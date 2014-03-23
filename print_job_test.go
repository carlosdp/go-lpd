package lpd

import (
	"bytes"
	"errors"
	"io"
	"testing"
	"time"
)

func TestNewPrintJob(t *testing.T) {
	file := bytes.NewReader([]byte{0x1})

	job, err := NewPrintJob("test", file)

	if err != nil {
		t.Fatal(err)
	}

	if job.QueueName != "test" {
		t.Fatal("Queue name not set properly")
	}

	buf := make([]byte, 1)

	n, err := job.DataFile.ReadAt(buf, 0)

	if err != nil && err != io.EOF {
		t.Fatal(err)
	}

	if buf[0] != 0x1 {
		t.Fatal("Data file not copied correctly, bytes read: ", n)
	}
}

func TestReceiveJob(t *testing.T) {
	dataFile := []byte{0x1}
	// TODO: Replace this with proper control file
	controlFile := []byte{0x1}

	dataCmd := &subCommand{
		Code:     ReceiveDataFile,
		NumBytes: 1,
		FileName: "test.pdf",
	}

	rawDataCmd := append(marshalSubCommand(dataCmd), 0x10)

	controlCmd := &subCommand{
		Code:     ReceiveControlFile,
		NumBytes: 1,
	}

	rawControlCmd := append(marshalSubCommand(controlCmd), 0x10)

	reader, receiveWrite := io.Pipe()
	receiveRead, writer := io.Pipe()

	jobChan := make(chan *PrintJob, 1)
	errChan := make(chan error, 1)

	go func(reader io.Reader, writer io.Writer, jobChan chan *PrintJob, errChan chan error) {
		job, err := receiveJob(reader, writer)

		if err != nil {
			errChan <- err
			return
		}

		jobChan <- job
	}(receiveRead, receiveWrite, jobChan, errChan)

	_, err := writer.Write(rawControlCmd)

	if err != nil {
		t.Fatal(err)
	}

	_, err = writer.Write(controlFile)

	if err != nil {
		t.Fatal(err)
	}

	err = waitForAck(reader, 2)

	if err != nil {
		t.Fatal(err)
	}

	_, err = writer.Write(rawDataCmd)

	if err != nil {
		t.Fatal(err)
	}

	_, err = writer.Write(dataFile)

	if err != nil {
		t.Fatal(err)
	}

	err = waitForAck(reader, 2)

	if err != nil {
		t.Fatal(err)
	}

	writer.Close()

	timeoutChan := time.After(time.Duration(500) * time.Millisecond)

	for {
		select {
		case <-timeoutChan:
			t.Fatal("receiveJob timed out")
		case err := <-errChan:
			t.Fatal(err)
		case job := <-jobChan:
			buf := make([]byte, 1)

			_, err = job.DataFile.ReadAt(buf, 0)

			if err != nil {
				t.Fatal("Error reading data file: ", err)
			}

			if buf[0] != 0x1 {
				t.Fatal("Data file corrupted")
			}

			// TODO: Test control file integrity

			return
		}
	}
}

func waitForAck(reader io.Reader, timeout int) error {
	readChan := make(chan []byte, 1)
	errChan := make(chan error, 1)

	go func(readChan chan []byte, errChan chan error) {
		buf := make([]byte, 1)
		_, err := reader.Read(buf)

		if err != nil {
			errChan <- err
			return
		}

		readChan <- buf
	}(readChan, errChan)

	timeoutChan := time.After(time.Duration(timeout) * time.Second)

	for {
		select {
		case <-timeoutChan:
			return errors.New("Ack timed out")
		case err := <-errChan:
			return err
		case buf := <-readChan:
			if buf[0] == 0x0 {
				return nil
			} else {
				return errors.New("Error acknowledgment received")
			}
		}
	}
}
