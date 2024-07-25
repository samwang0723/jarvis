package helper

import (
	"log"
	"time"
)

func TrackElapsed(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("> %s took %s", name, elapsed)
}
