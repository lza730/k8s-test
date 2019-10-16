package errorHelper

import (
	"github.com/golang/glog"
	"os"
)

func PanicOnError(err error, msg string) {
	if err != nil {
		glog.Error(msg)
		panic(err.Error())
	}
}

func InfoOnError(err error, msg string) {
	if err != nil {
		glog.Info(msg)
	}
}

func ExitOnError(err error, msg string, code int) {
	if err != nil {
		glog.Error(msg)
		os.Exit(code)
	}
}