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

const path = "../gitsync.log"

var AbsPath string

var levels = map[string]bool{
	"ERROR":   true,
	"INFO":    true,
	"WARNING": true,
}

func SetBaseDir(dir string) {
	AbsPath = fmt.Sprintf("%s/gitsync.log", dir)
}

func StaticWrite(message, level string) {
	if len(AbsPath) < 1 {
		fmt.Println("base directory not set")
		return
	}
	file, err := getLog(AbsPath)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	logged := parseMessage(message, level)
	_, err = file.Write(logged)
	if err != nil {
		fmt.Println(err)
	}
}

func NewLog() (*Log, error) {
	fullPath, err := getFullPath()
	if err != nil {
		return nil, err
	}
	return &Log{
		lastLogged: []byte{},
		file:       nil,
		fullPath:   fullPath,
	}, nil
}

func getFullPath() (string, error) {
	if len(AbsPath) > 0 {
		return AbsPath, nil
	}
	return filepath.Abs(path)
}

// Write writes a message
func (l *Log) Write(message string, level string) error {
	file, err := getLog(l.fullPath)
	if err != nil {
		return err
	}
	l.file = file
	defer l.cleanup()
	return l.writeMessage(message, level)
}

func getLog(fullPath string) (*os.File, error) {
	file, err := os.OpenFile(fullPath, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil && os.IsNotExist(err) {
		_, err = os.Create(fullPath)
		if err != nil {
			return nil, err
		}
		return getLog(fullPath)
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

func parseMessage(message, level string) []byte {
	now := time.Now().Format(time.RFC3339)
	if levels[level] == false {
		level = "ERROR"
	}
	return []byte(fmt.Sprintf("[%s] %s: %s\n", now, level, message))
}

func (l *Log) writeMessage(message, level string) error {
	l.lastLogged = parseMessage(message, level)
	_, err := l.file.Write(l.lastLogged)
	return err
}
