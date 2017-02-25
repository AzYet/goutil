package go_utils

import (
	"github.com/AzYet/logrus-prefixed-formatter"
	"github.com/getsentry/raven-go"
	"github.com/evalphobia/logrus_sentry"
	"github.com/Sirupsen/logrus"
)

func NewLogrusWithSentryHook(DSL, release string) *logrus.Logger {
	l := logrus.New()
	fmtr := &prefixed.TextFormatter{}
	fmtr.TimestampFormat = "2006/01/02 15:04:05"
	l.Formatter = fmtr
	l.Level = logrus.DebugLevel

	if DSL != "" {
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
				Enable :true,
				// the level at which to start capturing stacktraces
				Level : logrus.ErrorLevel,
			}
			l.Hooks.Add(h)
			l.Debug("sentry hook added.")
		}
	}
	return l
}

func LogBuilder(l *logrus.Logger, okLvl, errLvl int) func(err error, okStr, failStr string, fileds... string) func(values... interface{}) bool {
	return func(err error, okStr, failStr string, fileds... string) func(values... interface{}) bool {
		return func(values... interface{}) bool {
			var logBody *logrus.Entry
			if err == nil {
				logBody = l.WithFields(logrus.Fields{})
			} else {
				logBody = l.WithError(err)
			}
			vLen := len(values)
			for k, v := range fileds {
				if k == vLen {
					break
				}
				logBody = logBody.WithField(v, values[k])
			}
			if err == nil {
				switch okLvl {
				case 4:logBody.Info(okStr)
				case 5 :logBody.Debug(okStr)
				}
			} else {
				switch errLvl {
				case 0:logBody.Panic(failStr)
				case 1:logBody.Fatal(failStr)
				case 2:logBody.Error(failStr)
				case 3:logBody.Warn(failStr)
				}
			}
			return err == nil
		}
	}
}
