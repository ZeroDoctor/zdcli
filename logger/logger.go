package logger

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/zerodoctor/go-logging"
)

var (
	// log default logger used for termainal output
	log *logging.Logger
	// TERM_FORMAT is the default logging format for terminal
	TERM_FORMAT logging.Formatter

	// backendMap is a collection of logging backend
	backendMap = make(map[string]logging.Backend, 5)
	// m used to avoid race condition involving backendMap
	m sync.Mutex
)

// Init setup default logging for terminal backend
func Init() {
	if log != nil {
		return
	}

	log = logging.MustGetLogger("global")

	log.ExtraCalldepth = 1
	backend := logging.NewLogBackend(os.Stdout, "", 0)

	TERM_FORMAT = logging.MustStringFormatter(
		`%{color}%{level:.4s}: %{time:2006-01-02 15:04:05} %{shortfile} ▶%{color:reset} %{message}`,
	)
	backendFormatter := logging.NewBackendFormatter(backend, TERM_FORMAT)
	backendMap["terminal"] = backendFormatter

	reloadBackend()
}

func AddExistingBackend(key string, backend logging.Backend) error {
	if key == "terminal" {
		return fmt.Errorf("can not override backend with key '%s'", key)
	}
	backendMap[key] = backend

	reloadBackend()

	return nil
}

// AddBackend logging
func AddBackend(key string, level logging.Level, w io.Writer, format string) error {
	if key == "terminal" {
		return fmt.Errorf("can not override backend with key '%s'", key)
	}

	logFormat, err := logging.NewStringFormatter(format)
	if err != nil {
		return err
	}

	backend := logging.NewLogBackend(w, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, logFormat)

	backendLevel := logging.AddModuleLevel(backendFormatter)
	backendLevel.SetLevel(level, "")

	backendMap[key] = backendLevel

	reloadBackend()

	return nil
}

// RemoveBackend logging
func RemoveBackend(key string) error {
	if key == "terminal" {
		return fmt.Errorf("can not remove '%s' logging", key)
	}

	if _, ok := backendMap[key]; !ok {
		return fmt.Errorf("logging backend with key %s does not exists", key)
	}

	delete(backendMap, key)
	reloadBackend()
	return nil
}

// reloadBackend updates the current backend map
func reloadBackend() {
	var backends []logging.Backend
	for _, v := range backendMap {
		backends = append(backends, v)
	}

	m.Lock()
	logging.SetBackend(backends...)
	m.Unlock()
}

func Debugf(format string, args ...interface{}) { log.Debugf(format, args...) }
func Debug(args ...interface{})                 { log.Debug(args...) }

func Printf(format string, args ...interface{}) { log.Infof(format, args...) }
func Print(args ...interface{})                 { log.Info(args...) }

func Infof(format string, args ...interface{}) { log.Noticef(format, args...) }
func Info(args ...interface{})                 { log.Notice(args...) }

func Warnf(format string, args ...interface{}) { log.Warningf(format, args...) }
func Warn(args ...interface{})                 { log.Warning(args...) }

func Errorf(format string, args ...interface{}) { log.Errorf(format, args...) }
func Error(args ...interface{})                 { log.Error(args...) }

func Criticalf(format string, args ...interface{}) { log.Criticalf(format, args...) }
func Critical(args ...interface{})                 { log.Critical(args...) }

func Fatalf(format string, args ...interface{}) { log.Fatalf(format, args...) }
func Fatal(args ...interface{})                 { log.Fatal(args...) }

func Panicf(format string, args ...interface{}) { log.Panicf(format, args...) }
func Panic(args ...interface{})                 { log.Panic(args...) }

type Backtrace struct {
	logFn func(...interface{})
	msg   strings.Builder
}

func EnableBacktrace(logFn func(...interface{})) *Backtrace {
	return &Backtrace{
		logFn: logFn,
	}
}

func (b *Backtrace) Add(m string, args ...interface{}) {
	b.msg.WriteString(fmt.Sprintf(m, args...))
}

func (b Backtrace) Dump() {
	b.logFn(b.msg.String())
}
