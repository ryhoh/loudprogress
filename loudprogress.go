package loudprogress

import (
	"fmt"
	"io"
	"os"
	"time"
)

type RenderFunc func(size int64, writer io.Writer) (func(int64), func(int64)) // alias for rendering function

type LoudProgress struct { // progress manager
	size        int64      // total size of progress
	current     int64      // current progress
	render_func RenderFunc // function for rendering
	ch          chan int64 // channel for receive current number
	writer      io.Writer  // writer for rendering (default: os.Stdout)
}

const WAIT = 250 * time.Millisecond

// Generate new LoudProgress
func NewLoudProgress(size int64, render_func RenderFunc) *LoudProgress {
	res := new(LoudProgress)
	res.current = 0
	res.size = size
	res.render_func = render_func
	res.ch = make(chan int64, 20)
	res.writer = os.Stdout
	return res
}

// Set writer
func (lp *LoudProgress) SetWriter(writer io.Writer) {
	lp.writer = writer
}

// Start rendering progress
func (lp *LoudProgress) Start() {
	func_render_main, func_render_post := lp.render_func(lp.size, lp.writer)
	go render(lp.current, lp.size, lp.ch, func_render_main, func_render_post)
}

// Method for increment current number
// () -> (lp.current, err)
func (lp *LoudProgress) Increment() (int64, error) {
	// can't increase current if already finished
	if lp.current >= lp.size {
		return lp.current, fmt.Errorf("this progress already finished [%d/%d]", lp.current, lp.size)
	}

	// increase if available
	lp.current += 1
	lp.ch <- lp.current
	return lp.current, nil
}

// Function for rendering progress bar
func render(current, size int64, ch chan int64, func_render_main, func_render_post func(int64)) {
	/* render loop */
	for {
		if len(ch) > 0 { // update progress if increased
			current = <-ch
		}
		if current >= size {
			break
		}

		func_render_main(current)
		/* wait for next render */
		time.Sleep(WAIT)
	}
	func_render_post(current)
}
