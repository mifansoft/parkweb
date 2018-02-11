package logger_test

import (
	"mifanpark/utilities/logger"
	"testing"
)

func TestNewLogger(t *testing.T) {
	logger.Info("hello world abcd")
	logger.Info("my name mifan park")
	conf := logger.NewConfig()
	conf.SetName("newLogName.log")
	lg := logger.NewLogger(conf)
	lg.Error("hello world this is new logger")

	conf2 := logger.NewConfig("log.conf")
	ll := logger.NewLogger(conf2)
	ll.Info("hello world")
}
