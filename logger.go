package golog

import (
	"fmt"
	"log"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var (
	err     = color("\033[1;31m%s\033[0m")
	warn    = color("\033[1;33m%s\033[0m")
	debug   = color("\033[1;34m%s\033[0m")
	info    = color("\033[1;35m%s\033[0m")
	notice  = color("\033[1;36m%s\033[0m")
	success = color("\033[1;32m%s\033[0m")
)

// separator is a type defined for representing a separator in the log messages.
// It is used to separate different parts of a log message.
var separator string = " | "

// timePrecision is a type defined for representing the precision of seconds in the time format.
type timePrecision string

// Constants for different levels of time precision.
const (
	// MILLI represents millisecond precision in the time format.
	MILLI timePrecision = ".000"
	// MICRO represents microsecond precision in the time format.
	MICRO timePrecision = ".000000"
	// NANO represents nanosecond precision in the time format.
	NANO timePrecision = ".000000000"
)

// timeFormat is a global variable that holds the format of the time string.
var timeFormat string

// language is a global variable that holds the current language setting for the logger.
// It is used to determine the format of the timestamp in the log messages.
// The default language is set to English ("en").
var language string = "en"

// SetLanguage is a function that sets the time format based on the provided language.
// It takes a string representing the language and sets the time format accordingly.
func SetLanguage(lang string) {
	switch lang {
	case "en":
		timeFormat = "01-02-2006 3:04:05 PM"
	case "fr":
		timeFormat = "02-01-2006 15:04:05"
	case "de":
		timeFormat = "02.01.2006 15:04:05"
	case "zh":
		timeFormat = "2006年01月02日 15:04:05"
	case "ja":
		timeFormat = "2006年01月02日 15:04:05"
	case "ru":
		timeFormat = "02.01.2006 15:04:05"
	case "es":
		timeFormat = "02-01-2006 15:04:05"
	case "it":
		timeFormat = "02-01-2006 15:04:05"
	case "ar":
		timeFormat = "02/01/2006 15:04:05"
	default:
		timeFormat = "2006-01-02 15:04:05"
	}
	language = lang
}

// SetTimePrecision is a function that sets the time format to include seconds precision.
// It takes a timePrecision value representing the precision of the seconds (MILLI, MICRO, NANO).
// If the current language is English ("en"), it sets the time format to "01-02-2006 3:04:05" followed by the precision and " PM".
// For other languages, it appends the precision to the current time format.
func SetTimePrecision(precision timePrecision) {
	if language == "en" {
		timeFormat = "01-02-2006 3:04:05" + string(precision) + " PM"
	} else {
		timeFormat = timeFormat + string(precision)
	}
}

// SetCustomTimeFormat is a function that sets the time format to a custom format.
// It takes a string representing the custom format and sets the time format to it.
func SetCustomTimeFormat(customFormat string) {
	timeFormat = customFormat
}

// SetSeparator is a function that sets the separator for the log messages.
// It takes a string representing the separator and sets the global variable 'separator' to it.
// The separator is used to separate different parts of a log message.
// The default separator is " | ".
// Example: "2021-07-01 12:00:00 | Debug | main.go:15 | Server is running on localhost 8080"
func SetSeparator(sep string) {
	separator = sep
}

// includeFileName is a global variable that determines whether the file name should be included in the log.
var includeFileName = false

// EnableFileNameLogging is a function that enables the inclusion of the file name in the log.
func EnableFileNameLogging() {
	includeFileName = true
}

// getFileAndLine is a function that retrieves the file name and line number of the caller.
// It returns a string in the format "file:line".
func getFileAndLine() string {
	_, file, line, ok := runtime.Caller(4)
	if !ok {
		log.Println("Failed to get the caller information")
	}
	return file[strings.LastIndex(file, "/")+1:] + ":" + strconv.Itoa(line)
}

// color is a function that returns a function that formats a message with color.
// It takes a string representing the color and returns a function that takes a message, a level, and optional arguments.
// The returned function formats the message with the color, the current time, the level, and the file name and line number if enabled.
func color(colorString string) func(message, level string, args ...any) string {
	return func(message, level string, args ...any) string {
		formattedMessage := fmt.Sprintf(message, args...)
		timestamp := time.Now().Format(timeFormat)
		fileLine := getFileAndLine()
		var fullMessage string
		if !includeFileName {
			fullMessage = fmt.Sprintf("%s%s %-7s%s %s", timestamp, separator, level, separator, formattedMessage)
		} else {
			fullMessage = fmt.Sprintf("%s%s %-7s%s %-15s%s %s", timestamp, separator, level, separator, fileLine, separator, formattedMessage)
		}
		msg := fmt.Sprintf(colorString, fullMessage)

		return msg
	}
}

// format is a function that prints a formatted message.
// It takes a function that formats the message, a message, a level, and optional arguments.
// It prints the message formatted by the provided function.
func format(colorFn func(message, level string, args ...any) string, message, level string, params ...any) {
	fmt.Println(colorFn(message, level, params...))
}

// Err is a function that prints an error message.
// It takes a message and optional arguments and prints the message formatted as an error.
func Err(message string, params ...any) {
	format(err, message, "Err", params...)
}

// Warn is a function that prints a warning message.
// It takes a message and optional arguments and prints the message formatted as a warning.
func Warn(message string, params ...any) {
	format(warn, message, "Warn", params...)
}

// Debug is a function that prints a debug message.
// It takes a message and optional arguments and prints the message formatted as a debug message.
func Debug(message string, params ...any) {
	format(debug, message, "Debug", params...)
}

// Info is a function that prints an info message.
// It takes a message and optional arguments and prints the message formatted as an info message.
func Info(message string, params ...any) {
	format(info, message, "Info", params...)
}

// Notice is a function that prints a notice message.
// It takes a message and optional arguments and prints the message formatted as a notice.
func Notice(message string, params ...any) {
	format(notice, message, "Notice", params...)
}

// Success is a function that prints a success message.
// It takes a message and optional arguments and prints the message formatted as a success.
func Success(message string, params ...any) {
	format(success, message, "Success", params...)
}
