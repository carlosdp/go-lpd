package lpd

import (
	"encoding/binary"
	"testing"
)

func TestMarshalCommand(t *testing.T) {
	cmd := &command{
		Code:     0x5,
		Queue:    "Test",
		Username: "User",
	}

	rawCommand := marshalCommand(cmd)

	if rawCommand[0] != 0x5 {
		t.Fatal("Code not encoded correctly")
	}

	i := 1

	testQueue := []byte("Test")

	for _, b := range testQueue {
		if rawCommand[i] != b {
			t.Fatal("Queue not encoded correctly")
		}

		i++
	}

	if rawCommand[i] != 0x32 {
		t.Fatal("Space not included after queue name")
	}
	i++

	testUsername := []byte("User")

	for _, b := range testUsername {
		if rawCommand[i] != b {
			t.Fatal("Username not encoded correctly")
		}

		i++
	}
}

func TestUnmarshalCommand(t *testing.T) {
	rawCommand := []byte{0x5}

	bQueue := []byte("Test")

	for _, b := range bQueue {
		rawCommand = append(rawCommand, b)
	}

	rawCommand = append(rawCommand, 0x32)

	bUsername := []byte("User")

	for _, b := range bUsername {
		rawCommand = append(rawCommand, b)
	}

	cmd, err := unmarshalCommand(rawCommand)

	if err != nil {
		t.Fatal(err)
	}

	if cmd.Code != 0x5 {
		t.Fatal("Code not decoded correctly: ", cmd.Code)
	}

	if cmd.Queue != "Test" {
		t.Fatal("Queue not decoded correctly: ", cmd.Queue)
	}

	if cmd.Username != "User" {
		t.Fatal("User not decoded correctly: ", cmd.Username)
	}
}

func TestMashalSubCommand(t *testing.T) {
	subCmd := &subCommand{
		Code:     0x2,
		NumBytes: 10,
		FileName: "Test.file",
	}

	rawSubCommand := marshalSubCommand(subCmd)

	if rawSubCommand[0] != 0x2 {
		t.Fatal("Code not encoded correctly")
	}

	i := 1

	bNumBytes := make([]byte, 8)

	binary.LittleEndian.PutUint64(bNumBytes, 10)

	for _, b := range bNumBytes {
		if rawSubCommand[i] != b {
			t.Fatal("Number of bytes not encoded correctly")
		}

		i++
	}

	if rawSubCommand[i] != 0x32 {
		t.Fatal("Space not included after number of bytes")
	}

	i++

	bFileName := []byte("Test.file")

	for _, b := range bFileName {
		if rawSubCommand[i] != b {
			t.Fatal("File name not encoded correctly")
		}

		i++
	}
}
