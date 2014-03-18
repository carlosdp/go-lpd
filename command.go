package lpd

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
)

const (
	PrintWaitingJobs    = 0x1
	ReceivePrintJob     = 0x2
	SendQueueStateShort = 0x3
	SendQueueStateLong  = 0x4
	RemoveJobs          = 0x5

	AbortJob           = 0x1
	ReceiveControlFile = 0x2
	ReceiveDataFile    = 0x3
)

type command struct {
	Code     byte
	Queue    string
	Username string
	Other    [][]byte
}

type subCommand struct {
	Code     byte
	NumBytes uint64
	FileName string
}

func newPrintWaitingJobsCommand(queue string) *command {
	cmd := new(command)
	cmd.Code = PrintWaitingJobs
	cmd.Queue = queue

	return cmd
}

func marshalCommand(cmd *command) []byte {
	output := []byte{cmd.Code}

	bQueue := []byte(cmd.Queue)

	for _, b := range bQueue {
		output = append(output, b)
	}

	output = append(output, 0x32)

	bUsername := []byte(cmd.Username)

	for _, b := range bUsername {
		output = append(output, b)
	}

	for _, other := range cmd.Other {
		output = append(output, 0x32)
		for _, b := range other {
			output = append(output, b)
		}
	}

	return output
}

func unmarshalCommand(rawCommand []byte) (*command, error) {
	cmd := new(command)

	byteReader := bytes.NewReader(rawCommand)

	reader := bufio.NewReader(byteReader)

	code, err := reader.ReadByte()

	if code == 0x00 {
		return nil, errors.New("Command missing code")
	}

	cmd.Code = code

	if err != nil {
		return cmd, nil
	}

	queueName, err := reader.ReadBytes(0x32)

	if len(queueName) < 1 {
		return nil, errors.New("Command missing queue name")
	}

	if queueName[len(queueName)-1] == 0x32 {
		queueName = queueName[:len(queueName)-1]
	}

	cmd.Queue = string(queueName)

	if err != nil {
		return cmd, nil
	}

	if cmd.Code == RemoveJobs {
		username, err := reader.ReadBytes(0x32)

		if len(username) < 1 {
			return nil, errors.New("Command missing username")
		}

		if username[len(username)-1] == 0x32 {
			username = username[:len(username)-1]
		}

		cmd.Username = string(username)

		if err != nil {
			return cmd, nil
		}
	}

	cmd.Other = make([][]byte, 0)

	for {
		other, err := reader.ReadBytes(0x32)

		if len(other) > 0 {
			if other[len(other)-1] == 0x32 {
				other = other[:len(other)-1]
			}

			cmd.Other = append(cmd.Other, other)
		}

		if err != nil {
			// Reached end of command
			return cmd, nil
		}
	}
}

func marshalSubCommand(sbCmd *subCommand) []byte {
	output := []byte{sbCmd.Code}

	bNumBytes := make([]byte, 8)

	binary.LittleEndian.PutUint64(bNumBytes, sbCmd.NumBytes)

	for _, b := range bNumBytes {
		output = append(output, b)
	}

	output = append(output, 0x32)

	bFileName := []byte(sbCmd.FileName)

	for _, b := range bFileName {
		output = append(output, b)
	}

	return output
}

func unmarshalSubCommand(rawSubCommand []byte) (*subCommand, error) {
	subCmd := new(subCommand)

	byteReader := bytes.NewReader(rawSubCommand)

	reader := bufio.NewReader(byteReader)

	code, err := reader.ReadByte()

	if code == 0x0 {
		return nil, errors.New("Command missing code")
	}

	subCmd.Code = code

	if err != nil {
		return subCmd, nil
	}

	bNumBytes, err := reader.ReadBytes(0x32)

	if len(bNumBytes) < 1 {
		return nil, errors.New("Command missing number of bytes")
	}

	if bNumBytes[len(bNumBytes)-1] == 0x32 {
		bNumBytes = bNumBytes[:len(bNumBytes)-1]
	}

	subCmd.NumBytes = binary.LittleEndian.Uint64(bNumBytes)

	fileName, _ := reader.ReadBytes(0x32)

	if len(fileName) < 1 {
		return nil, errors.New("Command missing filename")
	}

	if fileName[len(fileName)-1] == 0x32 {
		fileName = fileName[:len(fileName)-1]
	}

	subCmd.FileName = string(fileName)

	return subCmd, nil
}
