package core

import (
	"log"
	"os"
	"sync"
	"time"
)

const (
	DATE_FORMAT = "20060102"
)

type Logger struct {
	handle *log.Logger
	date   string
	logDir string
	fp     *os.File
	mutex  sync.RWMutex
}

var logger *Logger

func GetLogInstance() *Logger {
	if logger == nil {
		logger = &Logger{}
		logger.SetLogDir(appPath + "/runtime")
	}

	return logger
}

func (logger *Logger) SetLogDir(dirName string) {
	logger.logDir = dirName
}

func (logger *Logger) check() {
	nowDate := time.Now().Format(DATE_FORMAT)

	if logger.fp == nil || logger.date != nowDate {
		logger.date = nowDate

		if logger.fp != nil {
			logger.fp.Close()
		}

		var err error
		fileName := logger.logDir + "/log_" + logger.date + ".log"
		logger.fp, err = os.OpenFile(fileName, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
		if err != nil {
			log.Panic(err)
		}

		logger.handle = log.New(logger.fp, "", log.LstdFlags)
	}
}

func (logger *Logger) Close() {
	logger.fp.Close()
}

func (logger *Logger) LogInfo(v ...interface{}) {
	logger.mutex.Lock()
	defer logger.mutex.Unlock()

	logger.check()

	logger.handle.SetPrefix("[info]")
	logger.handle.Println(v...)
}

func (logger *Logger) LogWarm(v ...interface{}) {
	logger.mutex.Lock()
	defer logger.mutex.Unlock()

	logger.check()

	logger.handle.SetPrefix("[warm]")
	logger.handle.Println(v...)
}

func (logger *Logger) LogFatal(v ...interface{}) {
	logger.mutex.Lock()
	defer logger.mutex.Unlock()

	logger.check()

	logger.handle.SetPrefix("[fatal]")
	logger.handle.Fatalln(v...)
}

func (logger *Logger) LogPanic(v ...interface{}) {
	logger.mutex.Lock()
	defer logger.mutex.Unlock()

	logger.check()

	logger.handle.SetPrefix("[panic]")
	logger.handle.Panicln(v...)
}
