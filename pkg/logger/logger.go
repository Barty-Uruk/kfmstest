package logger

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"sync"

	kitlog "github.com/go-kit/kit/log"
	kitloglevel "github.com/go-kit/kit/log/level"
)

var (
	ErrLoggerLevel = errors.New("can`t find level in (emerg, alert, crit, err, warn, notice, info, debug)")
	Log            *Logger
	// onceInit guarantee initialize logger only once
	onceInit sync.Once
)

type Logger struct {
	kitlog.Logger
}

func NewLogger(level, timeFormat string) *Logger {
	lvl, err := getLevel(level)
	if err != nil {
		fmt.Fprintf(os.Stderr, "logger init: %s", err)
		os.Exit(1)
	}
	format := "plain"
	var klog kitlog.Logger

	if format == "json" {
		klog = kitlog.NewJSONLogger(kitlog.NewSyncWriter(os.Stdout))
	} else {
		klog = kitlog.NewLogfmtLogger(kitlog.NewSyncWriter(os.Stdout))
	}
	klog = kitloglevel.NewFilter(klog, lvl)
	klog = kitlog.With(klog, "ts", kitlog.DefaultTimestampUTC)

	onceInit.Do(func() {
		Log = &Logger{klog}
	})

	return Log
}

/*func NewService(cfg configs.Logger) error {
	level, err := getLevel(cfg.Level)
	if err != nil {
		return err
	}

	once.Do(func() {
		logger := kitlog.NewJSONLogger(kitlog.NewSyncWriter(os.Stdout))
		logger = kitloglevel.NewFilter(logger, level)
		logger = kitlog.With(logger, "ts", kitlog.DefaultTimestampUTC, "caller", kitlog.DefaultCaller)

		Log = Logger{logger}
	})

	return nil
}*/

//Logger object

func (s *Logger) With(keyvals ...interface{}) *Logger {
	return &Logger{kitlog.With(s, keyvals...)}
}

func (s *Logger) Fatal(keyvals ...interface{}) {
	keyvals = append(keyvals, "caller")
	keyvals = append(keyvals, caller())
	kitloglevel.Error(s).Log(keyvals...)
}

func (s *Logger) Error(keyvals ...interface{}) {
	keyvals = append(keyvals, "caller")
	keyvals = append(keyvals, caller())
	kitloglevel.Error(s).Log(keyvals...)
}

func (s *Logger) Warn(keyvals ...interface{}) {
	keyvals = append(keyvals, "caller")
	keyvals = append(keyvals, caller())
	kitloglevel.Warn(s).Log(keyvals...)
}

func (s *Logger) Info(keyvals ...interface{}) {
	keyvals = append(keyvals, "caller")
	keyvals = append(keyvals, caller())
	kitloglevel.Info(s).Log(keyvals...)
}

func (s *Logger) Print(keyvals ...interface{}) {
	keyvals = append(keyvals, "caller")
	keyvals = append(keyvals, caller())
	kitloglevel.Info(s).Log(keyvals...)
}

func (s *Logger) Debug(keyvals ...interface{}) {
	keyvals = append(keyvals, "caller")
	keyvals = append(keyvals, caller())
	kitloglevel.Debug(s).Log(keyvals...)
}

// default logger
func With(keyvals ...interface{}) *Logger {
	return &Logger{kitlog.With(Log, keyvals...)}
}

func Fatal(keyvals ...interface{}) {
	kitloglevel.Error(Log).Log(keyvals...)
}

func Error(keyvals ...interface{}) {
	kitloglevel.Error(Log).Log(keyvals...)
}

func Warn(keyvals ...interface{}) {
	kitloglevel.Warn(Log).Log(keyvals...)
}

func Info(keyvals ...interface{}) {
	kitloglevel.Info(Log).Log(keyvals...)
}

func Print(keyvals ...interface{}) {
	kitloglevel.Info(Log).Log(keyvals...)
}

func Debug(keyvals ...interface{}) {
	kitloglevel.Debug(Log).Log(keyvals...)
}

func getLevel(lvl string) (kitloglevel.Option, error) {
	switch lvl {
	case "emerg":
		return kitloglevel.AllowError(), nil
	case "alert":
		return kitloglevel.AllowError(), nil
	case "crit":
		return kitloglevel.AllowError(), nil
	case "err":
		return kitloglevel.AllowError(), nil
	case "warning":
		return kitloglevel.AllowWarn(), nil
	case "notice":
		return kitloglevel.AllowInfo(), nil
	case "info":
		return kitloglevel.AllowInfo(), nil
	case "debug":
		return kitloglevel.AllowDebug(), nil
	}
	return nil, fmt.Errorf("level %s is incorrect. Level can be (emerg, alert, crit, err, warn, notice, info, debug)", lvl)
}

func caller() string {
	_, file, no, ok := runtime.Caller(2)
	if ok {
		/*short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		file = short*/
		return fmt.Sprintf("%v:%v ", file, no)
	}
	return "???:0 "
}
