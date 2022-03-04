package socks5server

import (
	"context"

	"github.com/Felamande/go-socks5"
	"github.com/therecipe/qt/core"
)

type Socks5Server struct {
	core.QObject
	lAddr    string
	cancelFn func()
	_        func(bool)        `signal:"runStateChange"`
	_        func(interface{}) `signal:"receiveRunningError"`
}

func (s *Socks5Server) Init() *Socks5Server {
	return s
}

func (s *Socks5Server) StopServer() {
	if s.cancelFn != nil {
		s.cancelFn()
		s.RunStateChange(false)
		return
	}
}

func (s *Socks5Server) StartServer(port string) {

	conf := &socks5.Config{}
	server, _ := socks5.New(conf)
	lAddr := ":" + port
	s.lAddr = lAddr
	ctx, cancelFn := context.WithCancel(context.Background())
	s.cancelFn = cancelFn

	s.RunStateChange(true)
	go func(ctxCancel context.Context) {
		defer s.RunStateChange(false)
		if err := server.ListenAndServeWithCtx("tcp", lAddr, ctxCancel); err != nil {
			s.ReceiveRunningError(err)
			return
		}
	}(ctx)
}
