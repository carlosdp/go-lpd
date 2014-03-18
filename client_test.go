package lpd

import (
	"net"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	client, err := NewClient("localhost:5155")

	defer client.Close()

	if err != nil {
		t.Fatal(err)
	}
}

func TestSendCommand(t *testing.T) {
	addr, _ := net.ResolveTCPAddr("tcp", "localhost:5155")

	server, err := net.ListenTCP("tcp", addr)

	if err != nil {
		t.Fatal("could not start test listener: ", err)
	}

	defer server.Close()

	receiveChan := make(chan []byte, 1)

	go func(s *net.TCPListener, rChan chan []byte) {
		conn, err := s.Accept()

		if err != nil {
			return
		}

		defer conn.Close()

		var buf [1024]byte

		bytes, err := conn.Read(buf[:])

		if err != nil {
			return
		}

		rChan <- buf[:bytes]
	}(server, receiveChan)

	client, err := NewClient("localhost:5155")

	if err != nil {
		t.Fatal("client error: ", err)
	}

	cmd := newPrintWaitingJobsCommand("test")

	err = client.sendCommand(cmd)

	if err != nil {
		t.Fatal(err)
	}

	timeoutChan := time.After(time.Duration(1) * time.Second)

	for {
		select {
		case buf := <-receiveChan:
			tCmd, err := unmarshalCommand(buf)

			if err != nil {
				t.Fatal(err)
			}

			if tCmd.Code != cmd.Code || tCmd.Queue != cmd.Queue {
				t.Fatal("command corrupted")
			}

			return
		case <-timeoutChan:
			t.Fatal("test timed out")
		}
	}
}
