// Copyright 2019 The LFX. All rights reserved.
// Based on standard library "log",you can use any method in the standard library

package log

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"
)

const (
	LEVEL_ALL   = 0x0
	LEVEL_DEBUG = 0x1
	LEVEL_INFO  = 0x2
	LEVEL_WARN  = 0x4
	LEVEL_FATAL = 0x8
	LEVEL_OFF   = 0xf
)

var levelTpl = map[int]string{LEVEL_DEBUG: "DEBUG", LEVEL_INFO: "INFO", LEVEL_WARN: "WARN", LEVEL_FATAL: "FATAL"}

type Conf struct {
	LogFile string //The log file path, the file name is in this format demo.log-*-*-*
	Level   int    //Contents : ALL > DEBUG > INFO > WARN > FATAL > OFF.
	Expire  int    //The log file saved days. -1 save always
	Trace   int    //Print file call addr detail. if < 0, close all trace. if >= 0, open trace with set value.
}

type Logger struct {
	*log.Logger
	conf  *Conf
	Fptr  *os.File
	Fname string
}

func New(ctx context.Context, conf *Conf) *Logger {
	l := &Logger{conf: conf}

	if conf.Trace >= 0 {
		l.conf.Trace = l.conf.Trace + 3
		l.Logger = log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
	} else {
		l.Logger = log.New(os.Stdout, "", log.LstdFlags)
	}

	l.auto(ctx)

	return l
}

func (l *Logger) auto(ctx context.Context) {
	if len(l.conf.LogFile) <= 0 {
		return
	}

	fname, deadline := parse_log_fname(l.conf.LogFile)

	fp, err := open_log_file(fname)
	if err != nil {
		fp = os.Stderr
	}

	l.Fptr = fp
	l.Fname = fname
	l.SetOutput(fp)

	if l.conf.LogFile == "/dev/null" {
		return
	}

	go func() {
		select {
		case <-ctx.Done():
			return
		case <-time.After(deadline.Sub(time.Now())):
			l.auto(ctx)
		}
	}()
}

func (l *Logger) pf(level int, f string, v ...any) {
	if level < l.conf.Level {
		return
	}
	l.Output(l.conf.Trace, levelTpl[level]+" "+fmt.Sprintf(f, v...))
}

func (l *Logger) pln(level int, v ...any) {
	if level < l.conf.Level {
		return
	}
	l.Output(l.conf.Trace, levelTpl[level]+" "+fmt.Sprintln(v...))
}

func (l *Logger) Debugf(f string, v ...any) {
	l.pf(LEVEL_DEBUG, f, v...)
}

func (l *Logger) Debug(v ...any) {
	l.pln(LEVEL_DEBUG, v...)
}

func (l *Logger) Infof(f string, v ...any) {
	l.pf(LEVEL_INFO, f, v...)
}

func (l *Logger) Info(v ...any) {
	l.pln(LEVEL_INFO, v...)
}

func (l *Logger) Warnf(f string, v ...any) {
	l.pf(LEVEL_WARN, f, v...)
}

func (l *Logger) Warn(v ...any) {
	l.pln(LEVEL_WARN, v...)
}

func (l *Logger) Fatalf(f string, v ...any) {
	l.pf(LEVEL_FATAL, f, v...)
}

func (l *Logger) Fatal(v ...any) {
	l.pln(LEVEL_FATAL, v...)
}

var (
	defaultLogger, DevNull *Logger
	loggers                map[string]*Logger
)

func init() {
	defaultLogger = New(context.TODO(), &Conf{"", 0x0, -1, 1})
	DevNull = New(context.TODO(), &Conf{"/dev/null", 0x0, -1, 1})
	loggers = make(map[string]*Logger)
}

func Init(ctx context.Context, configs map[string]*Conf) {
	for sn, conf := range configs {
		loggers[sn] = New(ctx, conf)
	}
}

func Get(sn string) *Logger {
	if l, ok := loggers[sn]; ok {
		return l
	}
	return defaultLogger
}

func Debugf(f string, v ...any) {
	defaultLogger.Debugf(f, v...)
}

func Debug(v ...any) {
	defaultLogger.Debug(v...)
}

func Infof(f string, v ...any) {
	defaultLogger.Infof(f, v...)
}

func Info(v ...any) {
	defaultLogger.Info(v...)
}

func Warnf(f string, v ...any) {
	defaultLogger.Warnf(f, v...)
}

func Warn(v ...any) {
	defaultLogger.Warn(v...)
}

func Fatalf(f string, v ...any) {
	defaultLogger.Fatalf(f, v...)
}

func Fatal(v ...any) {
	defaultLogger.Fatal(v...)
}
