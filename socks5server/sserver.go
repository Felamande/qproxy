package socks5server

import (
	"context"
	"time"

	"github.com/Felamande/go-socks5"
	"github.com/Workiva/go-datastructures/queue"
)

type Socks5Server struct {
	lAddr    string
	errQ     *queue.Queue
	cancelFn func()
}

func NewSocks5Server() *Socks5Server {
	return &Socks5Server{
		errQ: queue.New(16),
	}
}

func (c *Socks5Server) Stop() bool {
	if c.cancelFn != nil {
		c.cancelFn()
		return true
	}
	return false
}

func (c *Socks5Server) Start(port string) error {
	conf := &socks5.Config{}
	server, _ := socks5.New(conf)

	lAddr := ":" + port
	c.lAddr = lAddr
	ctx, cancelFn := context.WithCancel(context.Background())
	c.cancelFn = cancelFn

	go func(ctxCancel context.Context) {
		if err := server.ListenAndServeWithCtx("tcp", lAddr, ctxCancel); err != nil {
			c.errQ.Put(err)
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
