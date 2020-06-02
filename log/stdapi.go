// Copyright 2017 Microsoft. All rights reserved.
// MIT License

package log

// Standard logger is a pre-defined logger for convenience.
// Set log directory as the current location
var stdLog = NewLogger("azure-container-networking", LevelInfo, TargetStderr, "")

// GetStd - Helper functions for the standard logger.
func GetStd() *Logger {
	return stdLog
}

func SetName(name string) {
	stdLog.SetName(name)
}

func SetLevel(level int) {
	stdLog.SetLevel(level)
}

func SetLogFileLimits(maxFileSize int, maxFileCount int) {
	stdLog.SetLogFileLimits(maxFileSize, maxFileCount)
}

func Close() {
	stdLog.Close()
}

func SetTargetLogDirectory(target int, logDirectory string) error {
	return stdLog.SetTargetLogDirectory(target, logDirectory)
}

func GetLogDirectory() string {
	return stdLog.GetLogDirectory()
}

// SetComponentName sets the component name that appears at the beginning of a message in a log.
// Pass in runtime.Caller(0) as the only argument
func SetComponentName(pc uintptr, fileName string, line int, ok bool) {
	stdLog.SetComponentName(pc, fileName, line, ok)
}

func Request(tag string, request interface{}, err error) {
	stdLog.Request(tag, request, err)
}

func Response(tag string, response interface{}, returnCode int, returnStr string, err error) {
	stdLog.Response(tag, response, returnCode, returnStr, err)
}

// Logf logs to the local log.
func Logf(format string, args ...interface{}) {
	stdLog.Logf(format, args...)
}

// Printf logs to the local log and send the log through the channel.
func Printf(format string, args ...interface{}) {
	stdLog.Printf(format, args...)
}

func Debugf(format string, args ...interface{}) {
	stdLog.Debugf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	stdLog.Errorf(format, args...)
}

// WriteToLog formats a message and outputs it to a log if the specified level
// is as or more important than the logger's current level.
func WriteToLog(level int, format string, args ...interface{}) {
	stdLog.WriteToLog(level, format, args...)
}

// GetComponentString returns the logger's ComponentName surrounded in square brackets,
// or a default string if the ComponentName is empty.
func GetComponentString() string {
	return stdLog.GetComponentString()
}

/*
TODO
- deprecate Printf, Logf, Errorf, Debugf
*/
