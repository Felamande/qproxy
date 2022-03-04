package main

import (
	"fmt"
	"time"
)

func SprintfTimeln(format string, v ...interface{}) string {
	args := append([]interface{}{time.Now().Format("01-02T15:04:05.999")}, v...)
	return fmt.Sprintf("[time=%v]"+format+"\n", args...)
}
