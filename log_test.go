package goutil

import (
	"testing"
	"fmt"
	"github.com/pkg/errors"
	"github.com/Sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

func TestLogBuilder1(t *testing.T) {
	defer func() {
		fmt.Println(recover())
	}()
	l := NewLogger("", "", "", true)
	//test normal func
	fmt.Println("_________________should info success")
	l.Info("success")
	fmt.Println("_________________should warn error")
	l.WI(errors.New("i am err"), "success", "fail")()
	fmt.Println("_________________should panic")
	//l.PI(errors.New("i am panic"), "success", "fail")()
}
func TestNewLog(t *testing.T) {
	l := NewLogger("/tmp/test.log", "", "", true, true)


	// print err only
	l.Infoln("_________________should be no print")
	l.EI(nil, "", "fail")()
	l.Infoln("_________________should print err")
	l.EI(errors.New("i am error"), "ok", "fail")()

	// info info when no error, and error when error
	l.Infoln("_________________should info ok")
	l.EI(nil, "ok", "fail")()
	l.Infoln("_________________should print error")
	l.EI(errors.New("i am error"), "ok", "fail")()
	l.Infoln("_________________should be no print")
	l.EI(nil, "", "fail")()

	// info info when no error, and error when error
	l.Infoln("_________________should debug ok")
	l.ED(nil, "ok", "fail")()
	l.Infoln("_________________should print warn")
	l.ED(errors.New("i am error"), "ok", "fail")()
	l.Infoln("_________________should be no print")
	l.ED(nil, "", "fail")()

}

func BenchmarkNewCustLogger(b *testing.B) {
	l := NewLogger("/tmp/test.log", "", "", true, true)
	//l := NewLogger("", "", "", true, false)

	for i := 0; i < b.N; i++ {
		// print err only
		l.Infoln("_________________should be no print")
		l.EI(nil, "", "fail")()
		l.Infoln("_________________should print err")
		l.EI(errors.New("i am error"), "ok", "fail")()

		// info info when no error, and error when error
		l.Infoln("_________________should info ok")
		l.EI(nil, "ok", "fail")()
		l.Infoln("_________________should print error")
		l.EI(errors.New("i am error"), "ok", "fail")()
		l.Infoln("_________________should be no print")
		l.EI(nil, "", "fail")()

		// info info when no error, and error when error
		l.Infoln("_________________should debug ok")
		l.ED(nil, "ok", "fail")()
		l.Infoln("_________________should print warn")
		l.ED(errors.New("i am error"), "ok", "fail")()
		l.Infoln("_________________should be no print")
		l.ED(nil, "", "fail")()
	}
}

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