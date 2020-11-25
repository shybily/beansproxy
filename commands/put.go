package commands

import (
	"bufio"
	"io"
	"net"
)

type PutCommand struct {
	baseCommand
	Data []byte
}

func NewPutCommand(cmdLine []byte, conn net.Conn) (*PutCommand, error) {
	c := &PutCommand{
		baseCommand: baseCommand{cmdLine: cmdLine},
	}
	if err := c.readBody(conn); err != nil {
		return nil, err
	}
	return c, nil
}

func (p *PutCommand) ToBytes() []byte {
	return append(p.baseCommand.cmdLine, p.Data...)
}

func (p *PutCommand) Send() ([]byte, error) {
	return sendAndRead(p.ToBytes())
}

func (p *PutCommand) readBody(rd io.Reader) error {
	reader := bufio.NewReader(rd)
	var err error
	p.Data, err = reader.ReadBytes('\n')
	if err != nil {
		return err
	}
	return nil
}
