package server

import (
	"bufio"
	"context"
	"fmt"
	"github.com/shybily/beansproxy/commands"
	"github.com/shybily/beansproxy/resources"
	"io"
	"net"
)

type Proxy struct {
	conn net.Conn
	wr   *bufio.ReadWriter
}

func NewProxy(con net.Conn) *Proxy {
	return &Proxy{
		conn: con,
		wr:   bufio.NewReadWriter(bufio.NewReader(con), bufio.NewWriter(con)),
	}
}

func (p *Proxy) run() error {
	for {
		cmd, err := commands.CommandParse(p.conn)
		if err == io.EOF {
			return err
		} else if err != nil {
			if err := p.sendResp(err); err != nil {
				return err
			}
		}
		if resp, err := p.sendCmd(cmd); err == nil {
			if err := p.sendResp(resp); err != nil {
				return err
			}
		}
	}
}

func (p *Proxy) sendResp(i interface{}) error {
	var resp []byte
	switch i.(type) {
	case error:
		resp = []byte(fmt.Sprintf("%s\r\n", i.(error).Error()))
	case []byte:
		resp = i.([]byte)
	default:
		resp = []byte("UNKNOWN_ERROR\r\n")
	}
	_, err := p.conn.Write(resp)
	return err
}

func (p *Proxy) sendCmd(cmd commands.Command) ([]byte, error) {
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
		return cmd.Send()
	}
}
