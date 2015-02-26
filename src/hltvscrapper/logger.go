package hltvscrapper

import (
	"logger"
)

type LoggerState struct {
  Level int
  LogPath string
}

var log *logger.LoggerState

func NewLogger(logPath string, level int) {
  log = &logger.LoggerState{}
  log.LogPath = logPath
  log.Level = level 
}