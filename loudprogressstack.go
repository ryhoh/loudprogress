package loudprogress

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

type LoudProgressStack struct { // stack of LoudProgress
	lps   *[]*LoudProgress
	mutex sync.Mutex
	// is_runnable bool // default: true
}

// Generate new LoudProgressStack
func NewLoudProgressStack(loudProgresses *[]*LoudProgress) *LoudProgressStack {
	res := new(LoudProgressStack)
	res.lps = loudProgresses
	res.mutex = sync.Mutex{}
	// res.is_runnable = true
	// for _, lp := range *(res.lps) {
	// 	lp.writer = os.Stdout
	for i, _ := range *(res.lps) {
		(*(*res).lps)[i].writer = os.Stdout
		(*(*res).lps)[i].is_running = true // disable individual Start()
	}
	return res
}

// Set writer for rendering multiple progress bar
func (lps *LoudProgressStack) SetWriter(writer io.Writer) {
	for i, _ := range *(lps.lps) {
		(*(*lps).lps)[i].writer = writer
	}
}

// Get is_finished for all LoudProgress
func (lps *LoudProgressStack) IsFinished() bool {
	finished := true
	for _, lp := range *(lps.lps) {
		if !lp.IsFinished() {
			finished = false
		}
	}
	return finished
}

func (lps *LoudProgressStack) IsRunning() bool {
	is_running := false
	for _, lp := range *(lps.lps) {
		if !lp.IsRunning() {
			is_running = true
		}
	}
	return is_running
}

// Function for rendering multiple progress bar
func (lps *LoudProgressStack) Start() error {
	if !lps.IsRunning() {
		fmt.Fprintf((*(*lps).lps)[0].writer, "%s", strings.Repeat("\n", len(*(*lps).lps)-1))
		for i, lp := range *(lps.lps) {
			go lps.render_multi(lp, uint(len(*(lps.lps))-i-1))
		}
		return nil
	}
	return fmt.Errorf("this LoudProgressStack can't runnable")
}

// Function for rendering progress bar
func (lps *LoudProgressStack) render_multi(lp *LoudProgress, lines_from_bottom uint) {
	func_render_main, func_render_post := lp.render_func(lp.size, lp.writer)

	/* render loop */
	for {
		for len(lp.ch) > 0 { // update progress if increased
			lp.current = <-lp.ch
		}
		if lp.current >= lp.size {
			break
		}

		// Lock and write
		lps.mutex.Lock()
		if lines_from_bottom > 0 { // not bottom line
			fmt.Fprintf(lp.writer, "\033[%dA", lines_from_bottom)
		}
		func_render_main(lp.current)
		if lines_from_bottom > 0 { // not bottom line
			fmt.Fprintf(lp.writer, "\033[%dB", lines_from_bottom)
		}
		lps.mutex.Unlock()

		/* wait for next render */
		time.Sleep(lp.wait)
	}

	// Lock and write
	lps.mutex.Lock()
	if lines_from_bottom > 0 { // not bottom line
		fmt.Fprintf(lp.writer, "\033[%dA", lines_from_bottom)
	}
	func_render_post(lp.current)
	if lines_from_bottom > 0 { // not bottom line
		fmt.Fprintf(lp.writer, "\033[%dB", lines_from_bottom)
	}
	lps.mutex.Unlock()

	// make finished
	lps.mutex.Lock()
	lp.is_finished = true
	if lps.IsFinished() {
		fmt.Fprintf(lp.writer, "\n")
	}
	lps.mutex.Unlock()
}
