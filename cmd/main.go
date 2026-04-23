package main

import "github.com/Olprog59/golog"

func init() {
	golog.SetLanguage("fr") // Sets the date and time format to French
	golog.SetTimePrecision(golog.MICRO)
	golog.EnableFileNameLogging()
	err := golog.InitSyslog("tcp", "10.81.210.150:5140", "test golog")
	if err != nil {
		golog.Err("Error initializing syslog:", err)
	}
}

func main() {
	golog.Debug("Debug message")
	golog.Info("Info message")
	golog.InfoWithID("Info message with ID", "session12", "")
	golog.ErrorWithID("Error message with ID", "session12", "sam-4")
	golog.WarnWithID("Warning message with ID", "", "sam-4")
}
