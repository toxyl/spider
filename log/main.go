package log

import (
	"github.com/toxyl/glog"
)

var (
	logger = glog.NewLoggerSimple("spider")
)

func Error(format string, a ...interface{}) {
	for i, ai := range a {
		switch t := ai.(type) {
		case error:
			a[i] = glog.Error(t)
		}
	}
	logger.Error(format, a...)
}

func Failed(format string, a ...interface{}) {
	for i, ai := range a {
		switch t := ai.(type) {
		case error:
			a[i] = glog.Error(t)
		}
	}
	logger.NotOK(format, a...)
}

func Warning(format string, a ...interface{}) {
	for i, ai := range a {
		switch t := ai.(type) {
		case error:
			a[i] = glog.Error(t)
		}
	}
	logger.Warning(format, a...)
}

func Normal(format string, a ...interface{}) {
	logger.Default(format, a...)
}

func Info(format string, a ...interface{}) {
	logger.Info(format, a...)
}

func Table(column ...*glog.TableColumn) {
	logger.Table(column...)
}
