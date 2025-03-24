package main

import (
	"fmt"
	"os"
	"sync"
	"time"
)

type LogMessage struct {
	Message   string
	Level     LogLevel
	Timestamp time.Time
}

type LogOutput interface {
	Write(message *LogMessage)
}

type ConsoleLogOutput struct{}

func (c *ConsoleLogOutput) Write(message *LogMessage) {
	fmt.Printf("%s %s %s\n", message.Timestamp.Format(time.RFC3339), message.Message, LogLevelNames[message.Level])
}

type FileLogOutput struct {
	FileName string
	mu       sync.Mutex
}

func NewFileLogOutput(fileName string) *FileLogOutput {
	return &FileLogOutput{
		FileName: fileName,
	}
}

func (f *FileLogOutput) Write(message *LogMessage) {
	f.mu.Lock()
	defer f.mu.Unlock()

	file, err := os.OpenFile(f.FileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	fmt.Fprintf(file, "%s %s %s\n", message.Timestamp.Format(time.RFC3339), message.Message, LogLevelNames[message.Level])
}

type Logger struct {
	mu      sync.Mutex
	level   LogLevel
	outputs []LogOutput
}

var loggerInstance *Logger
var loggerOnce sync.Once

func GetLogger() *Logger {
	loggerOnce.Do(func() {
		loggerInstance = &Logger{
			level: Info,
			outputs: []LogOutput{
				&ConsoleLogOutput{},
			},
		}
	})
	return loggerInstance
}

func (l *Logger) AddOutput(output LogOutput) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.outputs = append(l.outputs, output)
}

func (l *Logger) SetLevel(level LogLevel) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}

func (l *Logger) Log(level LogLevel, message string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if level < l.level {
		return
	}

	msg := &LogMessage{
		Message:   message,
		Level:     level,
		Timestamp: time.Now(),
	}

	for _, output := range l.outputs {
		output.Write(msg)
	}
}

func main() {
	logger := GetLogger()
	logger.SetLevel(Debug)

	logger.AddOutput(NewFileLogOutput("app.log"))

	logger.Log(Debug, "Hello World")
	logger.Log(Error, "Hello Error")
}
