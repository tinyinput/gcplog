# GCP Logging Library

The package **gcplog** is a very simple logging library which allows you to use GCP's Structured Logging. Whilst Google provides an excellent Golang Logging Library (<https://pkg.go.dev/cloud.google.com/go/logging>), I've generally found it's a lot more than I need when working on Google Cloud Functions.

I like the simplicity of just being able to use the standard `fmt.Print` functions, but this will result in all of the output being recorded in Cloud Logging with a level of "DEFAULT".

I want the flexibility to set _my own levels!_ This library provides a simple way to do that.

It allows you to log at any of the supported levels, and wraps your log messages in JSON, so that they format correctly in Cloud Logging.

The simplest way to use the library it to create Logger objects at the levels you want to use. Then you can call the various `Print` functions on these objects to get logs at your required levels.

For example:
```go
package main
import (
	"github.com/tinyinput/gcplog"
)
func main() {
	// Create a WARN and an ERROR logger
	logWarn := New(gcplog.WARNING)
	logError := New(gcplog.ERROR)
	// Then simply call `Print` on those object to write your log messages
	logWarn.Print("This is a Warning Message")
	logError.Print("This is an Error Message")
}
```

The other (potentially) useful methods are the `PrefixPrint` variants.

These will prefix your log messages (the text), with the severity level that you're logging at.

I've found this helpful both for searching in Cloud Logging and if I'm routing those logs to another system (like Splunk, Elastic, etc). As in that next system I often want to parse out the severity level, so this really helps.

This code, for example:
```go
package main
import (
	"github.com/tinyinput/gcplog"
)
func main() {
	// Create a WARN and an ERROR logger
	logWarn := New(gcplog.WARNING)
	// Then simply call `PrefixPrint` on that object to write your log messages
	logWarn.Print("This is a Warning Message")
}
```

Will create a the log message:
```text
{"severity":"WARNING","message":"WARNING: Hello World"}
```

So in Cloud Logging, you'll see this in the message field:
```text
"WARNING: Hello World"
```

That's all there is too it. Use `Print` and `Printf` in the same way as you would in the `fmt` package.

You can read more about Google's Structured Logging here: <https://cloud.google.com/logging/docs/structured-logging>
	