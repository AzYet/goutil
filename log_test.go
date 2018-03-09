package goutil

import (
	"reflect"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	"gopkg.in/natefinch/lumberjack.v2"
)

func BenchmarkOrignLogger(b *testing.B) {
	l := logrus.New()
	fmtr := &logrus.TextFormatter{}
	fmtr.ForceColors = true
	l.Formatter = fmtr
	//l := NewLogger("", "", "", true, false)
	rotator := &lumberjack.Logger{
		Filename:   "/tmp/test.log",
		MaxSize:    150, // megabytes
		MaxBackups: 4,
		MaxAge:     28, //days
	}
	l.Out = (rotator)

	for i := 0; i < b.N; i++ {
		// print err only
		var err error
		l.Infoln("_________________should be no print")
		if err == nil {
			l.Info("ok")
		}
		l.Infoln("_________________should print err")
		if e := errors.New("err"); e != nil {
			l.WithError(e).Error("fail")
		}
		// info info when no error, and error when error
		l.Infoln("_________________should info ok")
		if err == nil {
			l.Info("ok")
		}
		l.Infoln("_________________should print error")
		if e := errors.New("err"); e != nil {
			l.WithError(e).Error("fail")
		}
		l.Infoln("_________________should be no print")
		if err != nil {
			l.Info("ok")
		}
		// info info when no error, and error when error
		l.Infoln("_________________should debug ok")
		if err == nil {
			l.Debug("ok")
		}
		l.Infoln("_________________should print warn")
		if e := errors.New("err"); e != nil {
			l.WithError(e).Warn("fail")
		}
		l.Infoln("_________________should be no print")
		if err != nil {
			l.Info("ok")
		}
	}
}

func TestNewLogger(t *testing.T) {
	type args struct {
		logPath string
		DSN     string
		release string
		color   bool
	}
	tests := []struct {
		name string
		args args
		want *logrus.Logger
	}{
		{
			"empty args",
			args{"", "", "", true},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewLogger(tt.args.logPath, tt.args.DSN, tt.args.release, tt.args.color); !reflect.DeepEqual(got, tt.want) {
				if got.Out == nil {
					t.Errorf("NewLogger() = %+v, want %+v", got, tt.want)
				}
			}
		})
	}
}
