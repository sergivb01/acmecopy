package utils

import (
	"log"
	"time"
)

func Start() time.Time {
	return time.Now()
}

func Track(name string, startTime time.Time) {
	log.Printf("Executed %s in %s!\n", name, time.Since(startTime))
}
