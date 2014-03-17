package lpd

import (
	"bufio"
	"bytes"
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
	NumBytes int
	FileName string
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

	if err != nil {
		return nil, errors.New("Command missing code")
	}

	cmd.Code = code

	queueName, err := reader.ReadBytes(0x32)

	if err != nil && len(queueName) < 1 {
		return nil, errors.New("Command missing queue name")
	}

	if queueName[len(queueName)-1] == 0x32 {
		queueName = queueName[:len(queueName)-1]
	}

	cmd.Queue = string(queueName)

	if cmd.Code == RemoveJobs {
		username, err := reader.ReadBytes(0x32)

		if err != nil && len(username) < 1 {
			return nil, errors.New("Command missing username")
		}

		if username[len(username)-1] == 0x32 {
			username = username[:len(username)-1]
		}

		cmd.Username = string(username)
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
