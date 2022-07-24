package main

import (
	"fmt"
	"time"

	loudprogress "github.com/ryhoh/loudprogress"
)

func main() {
	size := 40
	waits := []time.Duration{50 * time.Millisecond, 100 * time.Millisecond, 200 * time.Millisecond, 350 * time.Millisecond}
	// waits := []time.Duration{400 * time.Millisecond, 200 * time.Millisecond, 100 * time.Millisecond, 50 * time.Millisecond}
	lp_s := make([]*loudprogress.LoudProgress, 4)
	for i, _ := range lp_s {
		lp_s[i] = loudprogress.NewLoudProgress( // make progress
			uint64(size),                    // size of progress
			loudprogress.Render_BoldRainbow, // render function
		)
		lp_s[i].SetWait(waits[i])
	}
	lps := loudprogress.NewLoudProgressStack(&lp_s)
	err := lps.Start() // start
	fmt.Printf("%#v", err)
	time.Sleep(time.Second)

	for {
		for _, lp := range lp_s {
			// if rand.Intn(2) == 1 {
			lp.Increment() // step progress
			// }
		}

		time.Sleep(time.Second)
		if lps.IsFinished() {
			break
		}
	}
}
