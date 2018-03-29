package main

import "fmt"

// Logger describes basic logging functions.
type Logger interface {
	Info(...interface{})
	Infof(string, ...interface{})
	Error(...interface{})
	Errorf(string, ...interface{})
}

type voidLog struct{}

func (l *voidLog) Info(...interface{})           {}
func (l *voidLog) Infof(string, ...interface{})  {}
func (l *voidLog) Error(...interface{})          {}
func (l *voidLog) Errorf(string, ...interface{}) {}

type trackingLog struct {
	Logger
	prefix string
}

func makeTrackingLog(log Logger, prefix string, id, ct int) trackingLog {
	if log == nil {
		log = &voidLog{}
	}

	return trackingLog{
		Logger: log,
		prefix: fmt.Sprintf("%02d-%s-%03d: ", id, prefix, ct),
	}
}

func (l *trackingLog) logf(format string, args ...interface{}) {
	l.Infof(l.prefix+format, args...)
}

func (l *trackingLog) logerr(err error) {
	l.Error(l.prefix + err.Error())
}
