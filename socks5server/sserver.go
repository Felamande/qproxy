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
	logChan  chan error
	server   *socks5.Server

	_ func()       `constructor:"init"`
	_ func(bool)   `signal:"runStateChange"`
	_ func(string) `signal:"receiveRunningError"`
	_ func(string) `signal:"receiveServingError"`
	_ func(string) `slot:"StartServer"`
	_ func()       `slot:"StopServer"`
}

func (s *Socks5Server) init() {
	s.ConnectStartServer(s.startServer)
	s.ConnectStopServer(s.stopServer)
}

func (s *Socks5Server) stopServer() {
	if s.cancelFn != nil {
		s.cancelFn()
		s.RunStateChange(false)
		return
	}
}

func (s *Socks5Server) startServer(port string) {

	s.logChan = make(chan error, 6)
	conf := &socks5.Config{
		LogChan: s.logChan,
	}
	server, _ := socks5.New(conf)
	s.server = server
	lAddr := ":" + port
	s.lAddr = lAddr
	ctx, cancelFn := context.WithCancel(context.Background())
	s.cancelFn = cancelFn

	s.RunStateChange(true)
	go func(ctxCancel context.Context) {

		defer s.RunStateChange(false)
		if err := server.ListenAndServeWithCtx("tcp", lAddr, ctxCancel); err != nil {
			s.ReceiveRunningError(err.Error())
		}
	}(ctx)

	go func() {
		for e := range s.logChan {
			s.ReceiveServingError(e.Error())
		}
	}()
}
