package vmproxy

import (
	"context"

	"gitea.com/lunny/tango"
	"gitea.com/tango/binding"
	"github.com/therecipe/qt/core"
)

type VmProxyServer struct {
	core.QObject

	logChan   chan error
	t         *tango.Tango
	logCancel func()

	_ func() `constructor:"init"`

	_ func(string) `signal:"receiveLog"`
	_ func()       `signal:"startFailed"`
	_ func(bool)   `signal:"runStatusChange"`

	_ func(string) `slot:"Start"`
}

func (s *VmProxyServer) init() {
	s.logChan = make(chan error, 8)

	s.t = tango.Classic(NewTangoChanLogger(s.logChan))
	route := NewVmProxyRoute(s.logChan)
	InitDefault(s.logChan)
	s.t.Use(binding.Bind())
	s.t.Route("GET", "/getvms", route)
	s.t.Route("POST:ReqStart", "/reqstart", route)

	s.ConnectStart(s.start)
}

func (s *VmProxyServer) start(port string) {
	s.RunStatusChange(true)
	ctx, cancel := context.WithCancel(context.Background())
	s.logCancel = cancel
	go s.receiveLogging(ctx)

	go func(ctx context.Context) {
		defer s.logCancel()
		defer s.RunStatusChange(false)

		s.t.Run(":" + port)
	}(ctx)
}

func (s *VmProxyServer) receiveLogging(ctx context.Context) {
	for {
		select {
		case e := <-s.logChan:
			s.ReceiveLog(e.Error())
		case <-ctx.Done():
			return
		}
	}
}
