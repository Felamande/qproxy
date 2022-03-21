package vmproxy

import "fmt"

type tangoChanLogger struct {
	logChan chan error
}

func NewTangoChanLogger(logChan chan error) *tangoChanLogger {
	return &tangoChanLogger{
		logChan: logChan,
	}
}

func (t *tangoChanLogger) Debugf(format string, v ...interface{}) {
	go func() {
		t.logChan <- fmt.Errorf("[debug]"+format, v...)
	}()
}
func (t *tangoChanLogger) Debug(v ...interface{}) {
	go func() {
		t.logChan <- fmt.Errorf("[debug]" + fmt.Sprint(v...))
	}()
}
func (t *tangoChanLogger) Infof(format string, v ...interface{}) {
	go func() {
		t.logChan <- fmt.Errorf("[info]"+format, v...)
	}()
}
func (t *tangoChanLogger) Info(v ...interface{}) {
	go func() {
		t.logChan <- fmt.Errorf("[info]" + fmt.Sprint(v...))
	}()
}
func (t *tangoChanLogger) Warnf(format string, v ...interface{}) {
	go func() {
		t.logChan <- fmt.Errorf("[warn]"+format, v...)
	}()
}
func (t *tangoChanLogger) Warn(v ...interface{}) {
	go func() {
		t.logChan <- fmt.Errorf("[warn]" + fmt.Sprint(v...))
	}()
}
func (t *tangoChanLogger) Errorf(format string, v ...interface{}) {
	go func() {
		t.logChan <- fmt.Errorf("[error]"+format, v...)
	}()
}
func (t *tangoChanLogger) Error(v ...interface{}) {
	go func() {
		t.logChan <- fmt.Errorf("[error]" + fmt.Sprint(v...))
	}()
}
