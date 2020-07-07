package goutil

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/evalphobia/logrus_sentry"
	"github.com/getsentry/raven-go"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

//NewLogger creates a new *logrus.Logger with sentry hook if DSN and Version provided
//current lvlNameSrcEtc are level(int), app name(string), enable source (bool)
func NewLogger(logPath, DSN, release string, color bool, lvlNameSrcEtc ...interface{}) *logrus.Logger {
	l := &logrus.Logger{
		Out:   os.Stdout,
		Hooks: make(logrus.LevelHooks),
		Formatter: &logrus.TextFormatter{
			CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
				dir, file := path.Split(frame.File)
				return "", fmt.Sprintf("%s:%d",
					path.Join(path.Base(strings.TrimSuffix(dir, "/")), file), frame.Line)
			},
			TimestampFormat: "2006/01/02 15:04:05",
			FullTimestamp:   true,
			ForceColors:     color,
			DisableColors:   !color,
			FieldMap:        logrus.FieldMap{logrus.FieldKeyFile: "~src"},
		},
		ReportCaller: len(lvlNameSrcEtc) >= 3 && lvlNameSrcEtc[2] == true,
		Level:        logrus.InfoLevel,
		ExitFunc:     os.Exit,
	}
	if logPath != "" {
		l.Out = &lumberjack.Logger{
			Filename:   logPath,
			MaxSize:    150, // megabytes
			MaxBackups: 3,
			MaxAge:     28, //days
		}
	}
	if len(lvlNameSrcEtc) > 0 {
		if lvl, ok := lvlNameSrcEtc[0].(int); ok && lvl >= 0 && lvl < 6 {
			l.WithField("level", logrus.Level(lvl)).Info("log level set ok.")
			l.SetLevel(logrus.Level(lvl))
		}
	}
	if DSN == "" {
		return l
	}
	client, err := raven.New(DSN)
	if err != nil {
		panic(err)
	}
	client.SetRelease(release)
	if len(lvlNameSrcEtc) > 1 {
		if name, ok := lvlNameSrcEtc[1].(string); ok {
			client.Tags = map[string]string{"name": name}
		}
	}
	h, err := logrus_sentry.NewWithClientSentryHook(client, []logrus.Level{logrus.PanicLevel, logrus.FatalLevel, logrus.ErrorLevel})
	if err != nil {
		panic(err)
	}
	h.StacktraceConfiguration = logrus_sentry.StackTraceConfiguration{
		Enable: true,
		// the level at which to start capturing stacktraces
		Level: logrus.ErrorLevel,
	}
	l.Hooks.Add(h)
	return l
}

//NewLogrusWithSentryHook creates a new logrus.Logger with sentry hook if DSL and Version provided
func NewLogrusWithSentryHook(color bool, DSL, release string) *logrus.Logger {
	l := logrus.New()
	fmtr := &logrus.TextFormatter{}
	fmtr.TimestampFormat = "2006/01/02 15:04:05"
	if color {
		fmtr.ForceColors = true
	}
	l.Formatter = fmtr
	l.Level = logrus.DebugLevel

	if DSL == "" {
		return l
	}
	client, err := raven.New(DSL)
	client.SetRelease(release)
	if err != nil {
		l.Errorf("failed to create raven agent %v.", err)
	} else if h, err := logrus_sentry.NewWithClientSentryHook(client, []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
	}); err != nil {
		l.Errorf("failed to create hook %v.", err)
	} else {
		h.StacktraceConfiguration = logrus_sentry.StackTraceConfiguration{
			Enable: true,
			// the level at which to start capturing stacktraces
			Level: logrus.ErrorLevel,
		}
		l.Hooks.Add(h)
	}
	return l
}

//LogBuilder deprecated: turned out to be slow in benchmark
// if you want no out put when no error, either set okLvl to -1 or okStr to ""
func LogBuilder(l *logrus.Logger, okLvl, errLvl int) func(err error, okStr, failStr string, fileds ...string) func(values ...interface{}) bool {

	return func(err error, okStr, failStr string, fileds ...string) func(values ...interface{}) bool {
		var (
			logBody *logrus.Entry
			ok      = err == nil
			vLen    int
		)

		if !ok {
			logBody = l.WithError(err)
			return func(values ...interface{}) bool {
				vLen = len(values)
				for k, v := range fileds {
					if k < vLen {
						logBody.WithField(v, values[k])
					}
				}
				switch errLvl {
				case 0:
					logBody.Panic(failStr)
				case 1:
					logBody.Fatal(failStr)
				case 2:
					logBody.Error(failStr)
				case 3:
					logBody.Warn(failStr)
				}
				return ok
			}
		} else if okLvl > 0 && okStr != "" {
			logBody = l.WithFields(logrus.Fields{})
			return func(values ...interface{}) bool {
				vLen = len(values)
				for k, v := range fileds {
					if k < vLen {
						logBody.WithField(v, values[k])
					}
				}
				switch errLvl {
				case 3:
					logBody.Warn(okStr)
				case 4:
					logBody.Info(okStr)
				case 5:
					logBody.Debug(okStr)
				}
				return ok
			}
		} else {
			return func(values ...interface{}) bool {
				return ok
			}
		}
	}
}
