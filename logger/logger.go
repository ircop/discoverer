package logger

import (
	"os"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"time"
	"runtime/debug"
)

type logger struct {
	dir			string
	logFile		*os.File
	logPanic	*os.File
	logError	*os.File
	debug		bool
}

var mlog = logger{}

// InitLogger initializes variables, opens required files.
// Returns error or nil
func InitLogger(debug bool, dir string) error {
	mlog.debug = debug

	if dir == "" {
		return fmt.Errorf("log dir is nil")
	}
	mlog.dir = dir

	// open files
	if err := mlog.openFiles(); err != nil {
		return errors.Wrap(err, "Cannot open log files for writing")
	}

	return nil
}

func (l *logger) openFiles() error {
	fi, err := os.Stat(l.dir)
	if err != nil {
		// dir doesn't exist. Try to create it
		errMkdir := os.Mkdir(l.dir, os.FileMode(0770))
		if errMkdir != nil {
			return errors.Wrap(errMkdir, "Cannot create log directory")
		}
	} else {
		// directory exists
		if !fi.Mode().IsDir() {
			return fmt.Errorf("Log path is not a directory")
		}
	}

	// open regular and panic files
	l.logFile, err = os.OpenFile(fmt.Sprintf("%s/nms.log", l.dir),
		os.O_APPEND | os.O_WRONLY | os.O_CREATE, 0600)
	if err != nil {
		return errors.Wrap(err, "Cannot open log file")
	}

	l.logPanic, err = os.OpenFile(fmt.Sprintf("%s/panic.log", l.dir),
		os.O_APPEND | os.O_WRONLY | os.O_CREATE, 0600)
	if err != nil {
		return errors.Wrap(err, "Cannot open panic logfile")
	}

	l.logError, err = os.OpenFile(fmt.Sprintf("%s/error.log", l.dir),
		os.O_APPEND | os.O_WRONLY | os.O_CREATE, 0600)
	if err != nil {
		return errors.Wrap(err,"Cannot open error log")
	}

	return nil
}

func (l *logger) format(prefix string, message string, args ...interface{}) string {
	t := time.Now()
	timeString := fmt.Sprintf("%d.%02d.%02d %02d:%02d:%02d", t.Day(), t.Month(), t.Year(), t.Hour(), t.Minute(), t.Second())

	msg := fmt.Sprintf(message, args...)

	format := fmt.Sprintf("[%s][%s]: %s\n", timeString, prefix, msg)

	return format
}

// write() returns error, but we dont catch it in regular code: because we cannot write it to log, if writing to log fails :)
// We need this error only for testing purposes
func (l *logger) write(message string, w io.Writer) error {
	_, err := w.Write([]byte(message))

	if l.debug {
		fmt.Printf(message)
	}

	if err != nil {
		fmt.Printf("Error! Cannot write log: %s" , err.Error())
		return errors.Wrap(err, "Cannot write log")
	}

	return nil
}

// Err writes log with ERROR tag
func Err(message string, args ...interface{}) {
	mlog.write(mlog.format("ERROR", message, args...), mlog.logFile) // nolint:errcheck
	mlog.write(mlog.format("ERROR", message, args...), mlog.logError) // nolint:errcheck
}

// Panic writes log with PANIC tag into two files: regular log and panic log
func Panic(message string, args ...interface{}) {
	message = fmt.Sprintf("%s%s", message, debug.Stack())
	mlog.write(mlog.format("PANIC", message, args...), mlog.logFile)  // nolint:errcheck
	mlog.write(mlog.format("PANIC", message, args...), mlog.logPanic) // nolint:errcheck
}

// Debug writes log with DEBUG tag, only if debug variable is set to true
func Debug(message string, args ...interface{}) {
	if !mlog.debug {
		return
	}

	mlog.write(mlog.format("DEBUG", message, args...), mlog.logFile) // nolint:errcheck
}

// Log writes regular log with INFO tag
func Log(message string, args ...interface{}) {
	mlog.write(mlog.format("INFO", message, args...), mlog.logFile) // nolint:errcheck
}
