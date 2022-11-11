package log

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestDefualt(t *testing.T) {
	Infof("this is %s", "infof")
	Debugf("this is %s", "debugf")
	Warnf("this is %s", "warnf")
	Fatalf("this is %s", "fatalf")
	Info("this is", "infof")
	Debug("this is", "debugf")
	Warn("this is", "warnf")
	Fatal("this is", "fatalf")
}

func TestLogger(t *testing.T) {
	conf := &Conf{
		LogFile: "./test.log-*-*",
		Expire:  1,
		Trace:   -1,
	}

	var log *Logger
	var msg string

	conf.Level = 0x0
	msg = "ALL"
	log = New(context.TODO(), conf)
	log.Debug(msg, conf)
	log.Info(msg, conf)
	log.Warn(msg, conf)
	log.Fatal(msg, conf)
	defer os.Remove(log.Fname)

	conf.Level = 0x1
	msg = "DEBUG"
	log = New(context.TODO(), conf)
	log.Debug(msg, conf)
	log.Info(msg, conf)
	log.Warn(msg, conf)
	log.Fatal(msg, conf)
	defer os.Remove(log.Fname)

	conf.Level = 0x2
	msg = "INFO"
	log = New(context.TODO(), conf)
	log.Debug(msg, conf)
	log.Info(msg, conf)
	log.Warn(msg, conf)
	log.Fatal(msg, conf)
	defer os.Remove(log.Fname)

	conf.Level = 0x4
	msg = "WARN"
	log = New(context.TODO(), conf)
	log.Debug(msg, conf)
	log.Info(msg, conf)
	log.Warn(msg, conf)
	log.Fatal(msg, conf)
	defer os.Remove(log.Fname)

	conf.Level = 0x8
	msg = "FATAL"
	log = New(context.TODO(), conf)
	log.Debug(msg, conf)
	log.Info(msg, conf)
	log.Warn(msg, conf)
	log.Fatal(msg, conf)
	defer os.Remove(log.Fname)

	conf.Level = 0xf
	msg = "OFF"
	log = New(context.TODO(), conf)
	log.Debug(msg, conf)
	log.Info(msg, conf)
	log.Warn(msg, conf)
	log.Fatal(msg, conf)
	defer os.Remove(log.Fname)

	if content, err := ioutil.ReadFile(log.Fname); err != nil {
		t.Fatal(err)
	} else {
		fmt.Println(string(content))
	}
}

func Benchmark_LogCloseTrace(b *testing.B) {
	conf := &Conf{
		LogFile: "./test-no-trace.log-*-*",
		Level:   0x0,
		Expire:  1,
		Trace:   -1,
	}

	log := New(context.TODO(), conf)
	defer os.Remove(log.Fname)
	for n := 0; n < b.N; n++ {
		log.Fatalf("111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111")
	}
}

func Benchmark_LogOpenTrace(b *testing.B) {
	conf := &Conf{
		LogFile: "./test-trace.log-*-*",
		Level:   0x0,
		Expire:  1,
		Trace:   1,
	}

	log := New(context.TODO(), conf)
	defer os.Remove(log.Fname)
	for n := 0; n < b.N; n++ {
		log.Fatalf("111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111")
	}
}

func TestLoggers(t *testing.T) {
	var confMap = map[string]*Conf{
		"access": &Conf{"./access.log-*-*-*", 0x1, 1, -1},
		"run":    &Conf{"./run.log-*-*-*", 0x2, 1, -1},
		"state":  &Conf{"./state.log-*-*-*", 0x4, 1, -1},
		"error":  &Conf{"./error.log-*-*-*", 0x8, 1, -1},
	}

	Init(context.TODO(), confMap)

	Get("access").Infof("access")
	Get("run").Infof("access")
	Get("state").Infof("access")
	Get("error").Infof("access")
	Get("noset").Infof("access")

	os.Remove(Get("access").Fname)
	os.Remove(Get("run").Fname)
	os.Remove(Get("state").Fname)
	os.Remove(Get("error").Fname)
}
