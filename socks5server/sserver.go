package socks5server

import (
	"context"
	"sync"
	"time"

	"github.com/Felamande/go-socks5"
	"github.com/Workiva/go-datastructures/queue"
)

type Socks5Server struct {
	lAddr    string
	errQ     *queue.Queue
	cancelFn func()
	lock     sync.Mutex
	state    bool
}

func NewSocks5Server() *Socks5Server {
	return &Socks5Server{
		errQ:  queue.New(16),
		state: false,
	}
}

func (s *Socks5Server) Stop() bool {
	if s.cancelFn != nil {
		s.cancelFn()
		return true
	}
	return false
}

func (s *Socks5Server) Start(port string) error {
	conf := &socks5.Config{}
	server, _ := socks5.New(conf)

	lAddr := ":" + port
	s.lAddr = lAddr
	ctx, cancelFn := context.WithCancel(context.Background())
	s.cancelFn = cancelFn

	s.changeState(true)
	go func(ctxCancel context.Context) {
		defer s.changeState(false)
		if err := server.ListenAndServeWithCtx("tcp", lAddr, ctxCancel); err != nil {
			s.errQ.Put(err)
			return
		}
	}(ctx)

	return nil
}

func (s *Socks5Server) GetAllError(sec int64) []interface{} {
	items, _ := s.errQ.Poll(16, time.Second*time.Duration(sec))
	return items
}

func (s *Socks5Server) PeekError() interface{} {
	items, _ := s.errQ.Peek()
	return items
}

func (s *Socks5Server) changeState(isRunning bool) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.state = isRunning
}

func (s *Socks5Server) GetRunState() bool {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.state
}
