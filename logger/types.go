package main

type LogLevel int

const (
	Debug LogLevel = iota
	Info
	Warn
	Error
	Fatal
)

var LogLevelNames = map[LogLevel]string{
	Debug: "DEBUG",
	Info:  "INFO",
	Warn:  "WARN",
	Error: "ERROR",
	Fatal: "FATAL",
}
