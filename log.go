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
			h.StacktraceConfiguration.Enable = true
			l.Hooks.Add(h)
			l.Infoln("sentry hook added.")
		}
	}
	return l
}

