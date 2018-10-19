package lab7

import (
	"log"
	"time"
)

// TimeElapsed measures the time it takes to execute a function.
// Use it as like this with defer:
//     defer TimeElapsed(time.Now(), "FunctionToTime")
//
// For more details see: https://coderwall.com/p/cp5fya
func TimeElapsed(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %v", name, elapsed)
}
