package loudprogress

import (
	"fmt"
	"io"
	"os"
	"time"
)

type RenderFunc func(size int64, writer io.Writer) (func(int64), func(int64)) // alias for rendering function

type LoudProgress struct { // progress manager
	size        int64         // total size of progress
	current     int64         // current progress
	render_func RenderFunc    // function for rendering
	ch          chan int64    // channel for receive current number
	writer      io.Writer     // writer for rendering (default: os.Stdout)
	wait        time.Duration // duration between renderings
	is_running  bool          // default: false
	is_finished bool          // default: false
}

// Generate new LoudProgress
func NewLoudProgress(size int64, render_func RenderFunc) *LoudProgress {
	res := new(LoudProgress)
	res.current = 0
	res.size = size
	res.render_func = render_func
	res.ch = make(chan int64, 20)
	res.writer = os.Stdout
	res.wait = 250 * time.Millisecond
	res.is_running = false
	res.is_finished = false
	return res
}

// Set writer
func (lp *LoudProgress) SetWriter(writer io.Writer) {
	lp.writer = writer
}

// Set wait
func (lp *LoudProgress) SetWait(wait time.Duration) error {
	if wait < time.Millisecond {
		return fmt.Errorf("too small wait (needs at least 1 millisecond but given %v)", wait)
	}
	lp.wait = wait
	return nil
}

// Set size
func (lp *LoudProgress) ExpandSize(size int64) error {
	// allow expand
	if size >= lp.size {
		lp.size = size
		return nil
	}
	// deny shrink
	return fmt.Errorf("shrinking size is prohibited (current size is %d and given %d)", lp.size, size)
}

// Get wait
func (lp *LoudProgress) GetWait() time.Duration {
	return lp.wait
}

// Get is_running
func (lp *LoudProgress) IsRunning() bool {
	return lp.is_running
}

// Get is_finished
func (lp *LoudProgress) IsFinished() bool {
	return lp.is_finished
}

// Start rendering progress
func (lp *LoudProgress) Start() error {
	if !lp.is_running {
		go lp.render()
		lp.is_running = true
		return nil
	}
	return fmt.Errorf("this LoudProgress can't runnable")
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
func (lp *LoudProgress) render() {
	func_render_main, func_render_post := lp.render_func(lp.size, lp.writer)

	/* render loop */
	for {
		for len(lp.ch) > 0 { // update progress if increased
			lp.current = <-lp.ch
		}
		if lp.current >= lp.size {
			break
		}

		func_render_main(lp.current)

		/* wait for next render */
		time.Sleep(lp.wait)
	}

	func_render_post(lp.current)

	// make finished
	lp.is_finished = true
	lp.is_running = false
	fmt.Fprintf(lp.writer, "\n")
}
