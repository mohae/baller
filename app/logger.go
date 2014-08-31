// Contains log related stuff.
package clitpl

import (
	seelog "github.com/cihub/seelog"
)

// constants for loglevels
const (
	Trace    = "trace"
	Debug    = "debug"
	Info     = "info"
	Warn     = "warn"
	Error    = "error"
	Critical = "critical"
	Off      = "off"
)

var logger seelog.LoggerInterface

func init() {
	//Disable logger by default
	DisableLog()
}

// DisableLog disables all package output
func DisableLog() {
	logger = seelog.Disabled
}

// UseLoggers uses a specified seelog.LoggerInterface to output package to log.
func UseLogger(newLogger seelog.LoggerInterface) {
	logger = newLogger
}

// FlushLog, call before app shutdown.
// With the way this application is structured, logging needs to be set and
// flushed at the command/name.go level. If you aren't going to support
func FlushLog() {
	logger.Flush()
}
