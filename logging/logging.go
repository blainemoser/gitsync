package logging

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type Log struct {
	lastLogged []byte
	file       *os.File
	fullPath   string
}

const path = "../log.log"

var levels = map[string]bool{
	"ERROR":   true,
	"INFO":    true,
	"WARNING": true,
}

func NewLog() (*Log, error) {
	fullPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	return &Log{
		lastLogged: []byte{},
		file:       nil,
		fullPath:   fullPath,
	}, nil
}

// Write writes a message
func (l *Log) Write(message string, level string) error {
	file, err := l.getLog()
	if err != nil {
		return err
	}
	l.file = file
	defer l.cleanup()
	return l.writeMessage(message, level)
}

func (l *Log) getLog() (*os.File, error) {
	file, err := os.OpenFile(l.fullPath, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil && os.IsNotExist(err) {
		_, err = os.Create(l.fullPath)
		if err != nil {
			return nil, err
		}
		return l.getLog()
	}
	return file, err
}

func (l *Log) cleanup() {
	if l.file == nil {
		return
	}
	l.file.Close()
	l.file = nil
}

func (l *Log) parseMessage(message, level string) {
	now := time.Now().Format(time.RFC3339)
	if levels[level] == false {
		level = "ERROR"
	}
	l.lastLogged = []byte(fmt.Sprintf("[%s] %s: %s\n", now, level, message))
}

func (l *Log) writeMessage(message, level string) error {
	l.parseMessage(message, level)
	_, err := l.file.Write(l.lastLogged)
	return err
}
