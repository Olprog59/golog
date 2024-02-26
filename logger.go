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

func getFileAndLine() string {
	_, file, line, ok := runtime.Caller(4)
	if !ok {
		log.Println("Failed to get the caller information")
	}
	return file[strings.LastIndex(file, "/")+1:] + ":" + strconv.Itoa(line)
}

func color(colorString string) func(message, status string, args ...any) string {
	return func(message, status string, args ...any) string {
		formattedMessage := fmt.Sprintf(message, args...)
		timestamp := time.Now().Format("02-01-2006 15:04:05")
		fileLine := getFileAndLine()
		fullMessage := fmt.Sprintf("%s: %s: %s: %s", timestamp, status, fileLine, formattedMessage)
		msg := fmt.Sprintf(colorString, fullMessage)

		return msg
	}
}

func format(colorFn func(message, status string, args ...any) string, message, status string, params ...any) {
	fmt.Println(colorFn(message, status, params...))
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
