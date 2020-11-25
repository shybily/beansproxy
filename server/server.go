package server

import (
	"go.uber.org/zap"
	"io"
	"net"
)

type ProxyServer struct {
	opt    *ProxyOptions
	logger *zap.SugaredLogger
	cons   int
}

func NewProxyServer(opt *ProxyOptions) *ProxyServer {
	logger, _ := zap.NewProduction()

	return &ProxyServer{opt: opt, logger: logger.Sugar()}
}

func (p *ProxyServer) Listen() error {
	listen, err := p.run()
	if err != nil {
		return err
	}
	for {
		conn, err := listen.Accept()
		if p.cons >= p.opt.MaxConnections {
			_, _ = conn.Write([]byte("TOO MANY CONNECTIONS\r\n"))
			p.logger.With(zap.Int("max", p.opt.MaxConnections), zap.String("client", conn.RemoteAddr().String())).
				Errorf("too many connections")
			_ = conn.Close()
			continue
		}
		if err != nil {
			p.logger.With(zap.String("client", conn.RemoteAddr().String())).Error(err.Error())
			continue
		}
		p.cons++
		p.logger.With(zap.String("client", conn.RemoteAddr().String())).Info("connected")
		go func() {
			defer func() { _ = conn.Close(); p.cons-- }()
			err := (NewProxy(conn)).run()
			if err == io.EOF {
				p.logger.With(zap.Any("client", conn.RemoteAddr())).Info("client closed")
			} else {
				p.logger.With(zap.String("client", conn.RemoteAddr().String())).Error(err)
			}
		}()
	}
}

func (p *ProxyServer) run() (net.Listener, error) {
	return net.Listen("tcp", p.opt.Listen)
}

func (p *ProxyServer) proxy(conn net.Conn) {

}
