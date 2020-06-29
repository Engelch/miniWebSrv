package heavenshelp

// small helper class for logging + writing to stderr

import (
	"fmt"
	"log/syslog"
	"os"
	"time"
)

var logInitialised = false
var logger *syslog.Writer
var err error
var stringLog = false

// verifyLogInitialised is a helper function to verify that log is initialised. If not, it is initialised in this function.
// verifyLogInitialised is used by the other Log... functions as a form of precondition checking.
func verifyLogInitialised() {
	if !logInitialised {
		LogInit("uninitialised")
	}
}

// LogInit tries to initialise the logging service.
func LogInit(msg string) {
	logger, err = syslog.New(syslog.LOG_INFO, msg)
	if err != nil {
		panic("could not create connection to logging")
	}
	logInitialised = true
}

// LogStringInit does not use syslog (for dockerised environments. Instead, it writes all messages to stderr)
func LogStringInit() {
	stringLog = true
	logInitialised = true
}

// LogErr creates a message preprended with ERROR to syslog and stderr, but tries to continue execution.
func LogErr(msg string) {
	verifyLogInitialised()
	msgFull := "ERROR:" + time.Now().UTC().Format(time.RFC3339) + ":" + msg
	_, _ = fmt.Fprintln(os.Stderr, msgFull)
	if !stringLog {
		err = logger.Err(msgFull)
		if err != nil {
			panic("PANIC: cannot produce error logging message")
		}
	}
}

// LogWarn creates a syslog and STDERR message labeled with WARNING.
func LogWarn(msg string) {
	verifyLogInitialised()
	msgFull := "WARNING:" + time.Now().UTC().Format(time.RFC3339) + ":" + msg
	_, _ = fmt.Fprintln(os.Stderr, msgFull)
	if !stringLog {
		err = logger.Warning(msgFull)
		if err != nil {
			LogErr("Cannot write warning log for:" + msgFull)
		}
	}
}

// LogInfo creates an info error message to syslog and STDERR.
func LogInfo(msg string) {
	verifyLogInitialised()
	msgFull := "Info:" + time.Now().UTC().Format(time.RFC3339) + ":" + msg
	_, _ = fmt.Fprintln(os.Stderr, msgFull)
	if !stringLog {
		err = logger.Info(msgFull)
		if err != nil {
			LogErr("Cannot write info log for:" + msgFull)
		}
	}
}

// LogPanic creates a message to syslog and STDERR and then stops execution of this thread.
func LogPanic(msg string) {
	verifyLogInitialised()
	msgFull := "PANIC:" + time.Now().UTC().Format(time.RFC3339) + ":" + msg
	_, _ = fmt.Fprintln(os.Stderr, msgFull)
	if !stringLog {
		logger.Alert(msgFull)
	}
	panic(msgFull)
}
