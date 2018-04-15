package main

import (
	"fmt"
)

// Logger describes basic logging functions.
type Logger interface {
	Info(...interface{})
	Infof(string, ...interface{})
	Error(...interface{})
	Errorf(string, ...interface{})
	Println(...interface{})
	Printf(string, ...interface{})
}

type voidLog struct{}

func (l *voidLog) Info(...interface{})           {}
func (l *voidLog) Infof(string, ...interface{})  {}
func (l *voidLog) Error(...interface{})          {}
func (l *voidLog) Errorf(string, ...interface{}) {}
func (l *voidLog) Println(...interface{})        {}
func (l *voidLog) Printf(string, ...interface{}) {}

type replScopedLog struct {
	Logger
	prefix string
}

func newReplScopedLog(log Logger, prefix string, id, ct int) *replScopedLog {
	if log == nil {
		log = &voidLog{}
	}

	return &replScopedLog{
		Logger: log,
		prefix: fmt.Sprintf("%02d-%s-%03d: ", id, prefix, ct),
	}
}

func (l *replScopedLog) Info(args ...interface{}) {
	l.Logger.Info(l.prefix + fmt.Sprint(args...))
}

func (l *replScopedLog) Infof(format string, args ...interface{}) {
	l.Logger.Infof(l.prefix+format, args...)
}

func (l *replScopedLog) Error(args ...interface{}) {
	l.Logger.Error(l.prefix + fmt.Sprint(args...))
}

func (l *replScopedLog) Errorf(format string, args ...interface{}) {
	l.Logger.Errorf(l.prefix+format, args...)
}

func (l *replScopedLog) Println(args ...interface{}) {
	l.Error(args...)
}

func (l *replScopedLog) Printf(format string, args ...interface{}) {
	l.Errorf(format, args...)
}
