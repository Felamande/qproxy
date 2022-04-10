package portfwd

import (
	"context"
	"fmt"
	"io"
	"net"
	"sync"

	"github.com/Felamande/go-socks5"
)

type PortForwarder struct {
	listenPort int
	fwdToIp    string
	fwdToPort  string
	updateLock sync.Mutex
	status     bool
	cancelFn   func()
	logChan    chan error
}

func (f *PortForwarder) Info() (listenPort int, fwdToIp, fwdToPort string) {
	return f.listenPort, f.fwdToIp, f.fwdToPort
}

func NewPortForwarder(logChan chan error) *PortForwarder {
	return &PortForwarder{
		logChan: logChan,
		status:  false,
	}
}

func (f *PortForwarder) Start(fwdToIp, fwdToPort string) error {
	f.updateStatus(true)

	ctx, cancel := context.WithCancel(context.Background())
	f.cancelFn = cancel

	l, err := socks5.NewChanlisten(ctx, "tcp", "0.0.0.0:0")
	if err != nil {
		return err
	}

	localPort := l.Port()
	f.listenPort = localPort
	f.fwdToIp = fwdToIp
	f.fwdToPort = fwdToPort

	go func() {
		f.logChan <- fmt.Errorf("portforwarder start forward from 0.0.0.0:%v to %s:%s", localPort, fwdToIp, fwdToPort)
	}()

	go func(f *PortForwarder, ctx context.Context) {
		defer f.updateStatus(false)

		cchan, echan := l.Accept()
		for {
			select {
			case conn := <-cchan:
				fwdConn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", fwdToIp, fwdToPort))
				if err != nil {
					go func(e error) {
						f.logChan <- err
					}(err)
					continue
				}

				go func(fConn, lConn net.Conn) {
					_, err := io.Copy(fConn, lConn)
					if err != nil {
						f.logChan <- err
					}
				}(fwdConn, conn)

				go func(fConn, lConn net.Conn) {
					_, err := io.Copy(lConn, fConn)
					if err != nil {
						f.logChan <- err
					}
				}(fwdConn, conn)

			case err := <-echan:
				go func(e error) {
					if err != nil {
						f.logChan <- err
					}
				}(err)
			case <-ctx.Done():
				err := l.Close()
				go func(e error) {
					if err != nil {
						f.logChan <- err
					}
				}(err)
			}

		}
	}(f, ctx)
	return nil
}

func (f *PortForwarder) updateStatus(isrunning bool) {
	f.updateLock.Lock()
	defer f.updateLock.Unlock()
	f.status = isrunning
}

func (f *PortForwarder) IsRunning() bool {
	f.updateLock.Lock()
	defer f.updateLock.Unlock()
	return f.status
}

func (f *PortForwarder) Stop() error {
	f.updateStatus(false)

	f.cancelFn()

	return nil
}
