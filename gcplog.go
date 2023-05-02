// The package gcplog is a very simple logging library which allows you to use GPC's Structured Logging.
//
// Whilst Google provides an excellent Golang Logging Library (<https://pkg.go.dev/cloud.google.com/go/logging>).
// I've generally found it's a lot more than I need when working on Google Cloud Functions.
//
// I like the simplicity of just being able to use the standard `fmt.Print` functions, but this will result in all of the output
// being recorded in Cloud Logging with a level of "DEFAULT".
//
// I want the flexibility to set my own levels! This library provides a simple way to do that.
//
// It both allows you to log at specific levels, and wraps your log messages in JSON, so that they format correctly in Cloud Logging.
//
// The simplest way to use the library it to create Logger objects at the levels you want to use. Then you can call the various
// `Print` functions on these objects to get logs at your required levels.
//
// For example:
//
//	package main
//
//	import (
//		"github.com/tinyinput/gcplog"
//	)
//
//	func main() {
//		// Create a WARN and an ERROR logger
//		logWarn := New(gcplog.WARNING)
//		logError := New(gcplog.ERROR)
//		// Then simply call `Print` on those object to write your log messages
//		logWarn.Print("This is a Warning Message")
//		logError.Print("This is an Error Message")
//	}
//
// The other (potentially) useful methods are the `PrefixPrint` variants.
//
// These will prefix your log messages (the text), with the severity level that you're logging at.
// I've found this helpful both for searching in Cloud Logging and if I'm routing those logs to another system (like Splunk, Elastic, etc).
// In that next system I often want to parse out the severity level, so this really helps.
//
// This code for example:
//
//	package main
//
//	import (
//		"github.com/tinyinput/gcplog"
//	)
//
//	func main() {
//		// Create a WARN and an ERROR logger
//		logWarn := New(gcplog.WARNING)
//		// Then simply call `PrefixPrint` on that object to write your log messages
//		logWarn.Print("This is a Warning Message")
//	}
//
// Will create a the log message:
//
//	{"severity":"WARNING","message":"WARNING: Hello World"}
//
// So in Cloud Logging, you'll see this in the message field:
//
//	"WARNING: Hello World"
//
// That's all there is too it. Use `Print` and `Printf` in the same way as you would in the `fmt` package.
//
// You can read more about Google's Structured Logging here: <https://cloud.google.com/logging/docs/structured-logging>
package gcplog

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

const (
	DEFAULT   string = "DEFAULT"
	INFO      string = "INFO"
	NOTICE    string = "NOTICE"
	WARN      string = "WARNING"
	WARNING   string = "WARNING"
	ERR       string = "ERROR"
	ERROR     string = "ERROR"
	CRIT      string = "CRITICAL"
	CRITICAL  string = "CRITICAL"
	ALERT     string = "ALERT"
	EMERGENCY string = "EMERGENCY"
	DEBUG     string = "DEBUG"
)

var (
	severityAll = [9]string{DEFAULT, INFO, NOTICE, WARNING, ERROR, CRITICAL, ALERT, EMERGENCY, DEBUG} // A variable to contain all valid severity levels
)

// gcpLogMessage is a simple struct type to represent part of the standard GCP logging structure.
type gcpLogMessage struct {
	Severity string `json:"severity"`
	Message  string `json:"message"`
}

// Logger is the main logging object.
type Logger struct {
	severity string
}

// New returns a pointer to a new Logger.
func New(s ...string) *Logger {
	if len(s) >= 1 {
		if isValidSeverity(s[0]) {
			return &Logger{
				severity: s[0],
			}
		}
	}
	return defaultLogger()
}

// Print uses the same format as fmt.Print to write a log message with the severity of the Logger.
func (l *Logger) Print(v ...any) {
	l.output(fmt.Sprint(v...))
}

// Printf uses the same format as fmt.Printf to write a log message with the severity of the Logger.
func (l *Logger) Printf(format string, v ...any) {
	l.output(fmt.Sprintf(format, v...))
}

// Fatal uses the same format as fmt.Fatal to write a log message with the severity of the Logger and then exit, with exit code 1.
func (l *Logger) Fatal(v ...any) {
	l.output(fmt.Sprint(v...))
	os.Exit(1)
}

// Fatalf uses the same format as fmt.Fatalf to write a log message with the severity of the Logger and then exit, with exit code 1.
func (l *Logger) Fatalf(format string, v ...any) {
	l.output(fmt.Sprintf(format, v...))
	os.Exit(1)
}

// PrefixPrint prefixes the provided message element with severity level of the logger.
// It then uses the same format as fmt.Print to write a log message with the severity of the Logger.
func (l *Logger) PrefixPrint(v ...any) {
	l.output(fmt.Sprint(l.prefix(v...)...))
}

// PrefixPrintf prefixes the provided message element with severity level of the logger.
// It then uses the same format as fmt.Printf to write a log message with the severity of the Logger.
func (l *Logger) PrefixPrintf(format string, v ...any) {
	l.output(fmt.Sprintf("%s%s"+format, l.prefix(v...)...))
}

// PrefixFatal prefixes the provided message element with severity level of the logger.
// It then uses the same format as fmt.Fatal to write a log message with the severity of the Logger and then exit, with exit code 1.
func (l *Logger) PrefixFatal(v ...any) {
	l.output(fmt.Sprint(l.prefix(v...)...))
	os.Exit(1)
}

// PrefixFatalf prefixes the provided message element with severity level of the logger.
// It then uses the same format as fmt.Fatalf to write a log message with the severity of the Logger and then exit, with exit code 1.
func (l *Logger) PrefixFatalf(format string, v ...any) {
	l.output(fmt.Sprintf("%s%s"+format, l.prefix(v...)...))
	os.Exit(1)
}

// Severity returns the current severity of the Logger object, as a string.
func (l *Logger) Severity() string {
	return l.severity
}

// SetSeverity will set the severity of the Logger object to the provided string, if that string is a valid severity level.
// If the provided string is not valid, then the severity level will remain unchanged.
func (l *Logger) SetSeverity(s string) {
	if isValidSeverity(strings.ToUpper(s)) {
		l.severity = strings.ToUpper(s)
	}
}

// output is a method to write to resulting log message to GCP logging.
func (l *Logger) output(s string) {
	jsonBytes, _ := json.Marshal(gcpLogMessage{
		Severity: l.severity,
		Message:  strings.TrimSpace(s),
	})
	fmt.Print(string(jsonBytes))
}

// prefix returns the provided any slice, but with the severity of the logger object as the first element
func (l *Logger) prefix(v ...any) []any {
	p := []any{l.severity, ": "}
	return append(p, v...)
}

// defaultLogger returns a Logger object with all elements set to defaults.
func defaultLogger() *Logger {
	return &Logger{severity: DEFAULT}
}

// isValidSeverity checks to see if the provided string is a valid severity level.
func isValidSeverity(s string) bool {
	s = strings.ToUpper(s)
	for _, sev := range severityAll {
		if s == sev {
			return true
		}
	}
	return false
}
