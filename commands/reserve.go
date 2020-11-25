package commands

import (
	"context"
	"github.com/shybily/beansproxy/resources"
)

type ReserveCommand struct {
	baseCommand
	ch    chan []byte
	close chan interface{}
}

func NewReserveCommand(cmdLine []byte) (*ReserveCommand, error) {
	return &ReserveCommand{
		baseCommand: baseCommand{cmdLine: cmdLine},
		ch:          make(chan []byte),
		close:       make(chan interface{}),
	}, nil
}

func (r *ReserveCommand) Send() ([]byte, error) {
	resources.Range(func(ins *resources.Instance) {
		con, err := ins.GetConn(context.TODO())
		if err != nil {
			return
		}
		defer func() {
			ins.Put(con)
		}()
		if _, err := con.Write(r.ToBytes()); err != nil {
			return
		}
		select {
		case <-r.close:
			return
		}
	})

	for {
		select {
		case resp := <-r.ch:
			return resp, nil
		}
	}
}
