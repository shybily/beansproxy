package commands

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/shybily/beansproxy/resources"
	"github.com/shybily/go-utils"
	"io"
	"net"
	"strings"
)

var (
	ErrUnknownCommand = errors.New("UNKNOWN_COMMAND")
	NotSupportCommand = errors.New("COMMAND_NOT_SUPPORT")
)

type Command interface {
	ToBytes() []byte
	Send() ([]byte, error)
}

type baseCommand struct {
	cmdLine []byte
}

func (b *baseCommand) ToBytes() []byte {
	return b.cmdLine
}

func (b *baseCommand) Send() ([]byte, error) {
	return sendAndRead(b.cmdLine)
}

func CommandParse(con net.Conn) (Command, error) {
	data, err := readResp(con)
	if err != nil {
		return nil, err
	}
	tmp := bytes.Split(data, []byte(" "))
	if len(tmp) <= 0 {
		return nil, ErrUnknownCommand
	}
	switch strings.ToLower(strings.TrimRight(utils.ByteToString(tmp[0]), "\r\n")) {
	case "put":
		return NewPutCommand(data, con)
	case "use":
		return NewUseCommand(data)
	case "watch":
		return NewWatchCommand(data)
	case "ignore":
		return NewIgnoreCommand(data)
	case "stats":
		return NewStatsCommand(data)
	case "reserve": //暂不支持
		return nil, NotSupportCommand
	case "quit":
		return nil, io.EOF
	default:
		return nil, NotSupportCommand
	}
}

func readResp(conn net.Conn) ([]byte, error) {
	r := bufio.NewReader(conn)
	resp, err := r.ReadBytes('\n')
	if err != nil {
		return nil, err
	}
	var (
		status string
		length int
	)
	_, _ = fmt.Sscanf(utils.ByteToString(resp), "%s %d", &status, &length)
	if status != "OK" {
		return resp, nil
	}
	if length > 0 {
		body := make([]byte, length)
		if _, err := r.Read(body); err != nil {
			return nil, err
		}
		resp = append(resp, body...)
	}
	return resp, nil
}

func sendAndRead(b []byte) ([]byte, error) {
	ins, err := resources.GetInstance()
	if err != nil {
		return nil, err
	}
	if con, err := ins.GetConn(context.TODO()); err != nil {
		return nil, err
	} else {
		defer func() {
			ins.Put(con)
		}()
		if _, err := con.Write(b); err != nil {
			return nil, err
		}
		return readResp(con)
	}
}
