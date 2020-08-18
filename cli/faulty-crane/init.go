package main

import (
	log "github.com/sirupsen/logrus"
)

func initLogging() {
	log.SetFormatter(&log.TextFormatter{
		DisableLevelTruncation: true,
	})
}
