package eseascrapper

import (
	"logger"
)

var log *logger.LoggerState

func NewLogger(path string, level int) {
	log = &logger.LoggerState{}
	log.LogPath = path
	log.Level = level
}