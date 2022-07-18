package loudprogress

import (
	"io"
	"os"
	"testing"
	"time"
)

func TestSetWriter(t *testing.T) {
	lp := NewLoudProgress(2, nil)
	lp.SetWriter(os.Stderr)
	if lp.writer != os.Stderr {
		t.Errorf("expected %#v but given %#v", os.Stderr, lp.writer)
	}
}

func TestStart(t *testing.T) {
	stub_func_render_main_count := 0
	stub_func_render_post_count := 0
	stub_supplier_func := func(_ int64, _ io.Writer) (func(int64), func(int64)) {
		return func(_ int64) {
				stub_func_render_main_count++
			}, func(_ int64) {
				stub_func_render_post_count++
			}
	}

	lp := NewLoudProgress(2, stub_supplier_func)
	lp.ch <- 2
	lp.Start() // and break immediately
	if stub_func_render_main_count != 0 && stub_func_render_post_count != 1 {
		t.Errorf("expected (%v, %v) but given (%v, %v)", 0, 1, stub_func_render_main_count, stub_func_render_post_count)
	}
}

func TestRender(t *testing.T) {
	const eps = 10 * time.Millisecond // wait short time for execution in another goroutine
	stub_func_render_main_count := 0
	stub_func_render_post_count := 0
	ch := make(chan int64, 20)

	stub_func_render_main := func(_ int64) {
		stub_func_render_main_count++
	}
	stub_func_post_main := func(_ int64) {
		stub_func_render_post_count++
	}

	// render start
	go render(0, 2, ch, stub_func_render_main, stub_func_post_main)

	// first cycle
	time.Sleep(eps)
	if stub_func_render_main_count != 1 && stub_func_render_post_count != 0 {
		t.Errorf("expected (%v, %v) but given (%v, %v)", 1, 0, stub_func_render_main_count, stub_func_render_post_count)
	}
	ch <- 1
	time.Sleep(WAIT)

	// second cycle
	if stub_func_render_main_count != 2 && stub_func_render_post_count != 0 {
		t.Errorf("expected (%v, %v) but given (%v, %v)", 2, 0, stub_func_render_main_count, stub_func_render_post_count)
	}
	ch <- 2
	time.Sleep(WAIT)

	// third cycle (break)
	if stub_func_render_main_count != 2 && stub_func_render_post_count != 1 {
		t.Errorf("expected (%v, %v) but given (%v, %v)", 2, 1, stub_func_render_main_count, stub_func_render_post_count)
	}
}

func TestIncrement(t *testing.T) {
	lp := NewLoudProgress(2, nil)

	if cur, err := lp.Increment(); cur != 1 || err != nil {
		t.Errorf("expected (%v, %v) but given (%v, %v)", 1, nil, cur, err)
	}

	if cur, err := lp.Increment(); cur != 2 || err != nil {
		t.Errorf("expected (%v, %v) but given (%v, %v)", 1, nil, cur, err)
	}

	if cur, err := lp.Increment(); cur != 2 || err == nil {
		t.Errorf("expected (%v, %v) but given (%v, %v)", 1, nil, cur, err)
	}
}
