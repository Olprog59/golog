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

var timeFormat string

func SetLanguage(language string) {
	switch language {
	case "en":
		timeFormat = "01-02-2006 3:04 PM"
	case "fr":
		timeFormat = "02-01-2006 15:04"
	case "de":
		timeFormat = "02.01.2006 15:04"
	case "zh":
		timeFormat = "2006年01月02日 15:04"
	case "ja":
		timeFormat = "2006年01月02日 15:04"
	case "ru":
		timeFormat = "02.01.2006 15:04"
	case "es":
		timeFormat = "02-01-2006 15:04"
	case "it":
		timeFormat = "02-01-2006 15:04"
	case "ar":
		timeFormat = "02/01/2006 15:04"
	default:
		timeFormat = "2006-01-02 15:04"
	}
}

func SetCustomTimeFormat(customFormat string) {
	timeFormat = customFormat
}

var includeFileName = false

func EnableFileNameLogging() {
	includeFileName = true
}

func getFileAndLine() string {
	_, file, line, ok := runtime.Caller(4)
	if !ok {
		log.Println("Failed to get the caller information")
	}
	return file[strings.LastIndex(file, "/")+1:] + ":" + strconv.Itoa(line)
}

func color(colorString string) func(message, level string, args ...any) string {
	return func(message, level string, args ...any) string {
		formattedMessage := fmt.Sprintf(message, args...)
		timestamp := time.Now().Format(timeFormat)
		fileLine := getFileAndLine()
		var fullMessage string
		if !includeFileName {
			fullMessage = fmt.Sprintf("%s: %-7s: %s", timestamp, level, formattedMessage)
		} else {
			fullMessage = fmt.Sprintf("%s: %-7s: %-15s: %s", timestamp, level, fileLine, formattedMessage)
		}
		msg := fmt.Sprintf(colorString, fullMessage)

		return msg
	}
}

func format(colorFn func(message, level string, args ...any) string, message, level string, params ...any) {
	fmt.Println(colorFn(message, level, params...))
}

func Err(message string, params ...any) {
	format(err, message, "Err", params...)
}
func Warn(message string, params ...any) {
	format(warn, message, "Warn", params...)
}
func Debug(message string, params ...any) {
	format(debug, message, "Debug", params...)
}
func Info(message string, params ...any) {
	format(info, message, "Info", params...)
}
func Notice(message string, params ...any) {
	format(notice, message, "Notice", params...)
}
func Success(message string, params ...any) {
	format(success, message, "Success", params...)
}
