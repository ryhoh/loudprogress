package main

import (
	"time"

	loudprogress "github.com/ryhoh/loudprogress"
)

func main() {
	size := 40
	lp := loudprogress.NewLoudProgress( // make progress
		uint64(size),               // size of progress
		loudprogress.Render_Simple, // render function
	)
	lp.Start() // start
	time.Sleep(time.Second)

	for i := 0; i < size; i++ {
		lp.Increment() // step progress
		time.Sleep(time.Second)
	}
}
