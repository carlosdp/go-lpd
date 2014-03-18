package lpd

import (
	"net"
)

// A Client sends commands to a server for status
// updates and creation of print jobs.
type Client struct {
	addr *net.TCPAddr
}

// NewClient creates and returns an instance of Client that
// will be connecting to a server at address. The standard port
// for an LPD server is 515. NewClient returns and error if and
// only if the address is invalid.
func NewClient(address string) (*Client, error) {
	addr, err := net.ResolveTCPAddr("tcp", address)

	if err != nil {
		return nil, err
	}

	c := new(Client)
	c.addr = addr

	return c, nil
}

// PrintWaitingJobs sends a command to the LPD server that tells
// it to print any waiting jobs in the specified queue.
func (c *Client) PrintWaitingJobs(queue string) error {
	cmd := newPrintWaitingJobsCommand(queue)

	err := c.sendCommand(cmd)

	if err != nil {
		return err
	}

	return nil
}

// Close closes the Client and blocks until all pending
// commands finish sending.
func (c *Client) Close() {
	// Do nothing, for now
}

func (c *Client) sendCommand(cmd *command) error {
	conn, err := net.DialTCP("tcp", nil, c.addr)

	if err != nil {
		return err
	}

	defer conn.Close()

	rawCmd := marshalCommand(cmd)

	_, err = conn.Write(rawCmd)

	if err != nil {
		return err
	}

	return nil
}
