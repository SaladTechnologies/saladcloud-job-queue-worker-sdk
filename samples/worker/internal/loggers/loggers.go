package loggers

import (
	"fmt"
	"path"
	"runtime"

	"github.com/sirupsen/logrus"
)

// Logger ...
var Logger *logrus.Logger

func init() {
	Logger = logrus.StandardLogger()
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&logrus.TextFormatter{
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {

			location := fmt.Sprintf("[%s:%d]", path.Base(frame.File), frame.Line)
			return "", location
		},
	})
}
