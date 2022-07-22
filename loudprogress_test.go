package loudprogress

import (
	"bytes"
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

func TestSetWait(t *testing.T) {
	wait := 300 * time.Millisecond
	lp := NewLoudProgress(2, nil)
	err := lp.SetWait(wait)
	if err != nil {
		t.Errorf("expected %#v but given %#v", err, nil)
	}
	if lp.wait != wait {
		t.Errorf("expected %#v but given %#v", lp.wait, wait)
	}

	// try setting too short wait
	err = lp.SetWait(time.Nanosecond)
	if err == nil {
		t.Errorf("expected %#v but given %#v", err, nil)
	}
	if lp.wait != wait {
		t.Errorf("expected %#v but given %#v", lp.wait, wait)
	}
}

func TestGetWait(t *testing.T) {
	expected := 300 * time.Millisecond
	lp := NewLoudProgress(2, nil)
	lp.wait = expected
	actual := lp.GetWait()
	if expected != actual {
		t.Errorf("expected %#v but given %#v", expected, actual)
	}
}

func TestIsRunning(t *testing.T) {
	expected := true
	lp := NewLoudProgress(2, nil)
	lp.is_running = expected
	actual := lp.IsRunning()
	if expected != actual {
		t.Errorf("expected %#v but given %#v", expected, actual)
	}
}

func TestIsFinished(t *testing.T) {
	expected := true
	lp := NewLoudProgress(2, nil)
	lp.is_finished = expected
	actual := lp.IsFinished()
	if expected != actual {
		t.Errorf("expected %#v but given %#v", expected, actual)
	}
}

func TestExpandSize(t *testing.T) {
	var size uint64 = 3
	lp := NewLoudProgress(2, nil)

	// expand
	err := lp.ExpandSize(size)
	if err != nil {
		t.Errorf("expected %#v but given %#v", err, nil)
	}
	if lp.size != size {
		t.Errorf("expected %#v but given %#v", lp.size, size)
	}

	// try shrink
	err = lp.ExpandSize(2)
	if err == nil {
		t.Errorf("expected %#v but given %#v", err, nil)
	}
	if lp.size != size {
		t.Errorf("expected %#v but given %#v", lp.size, size)
	}
}

func TestStart(t *testing.T) {
	stub_func_render_main_count := 0
	stub_func_render_post_count := 0
	stub_supplier_func := func(_ uint64, _ io.Writer) (func(uint64), func(uint64)) {
		return func(_ uint64) {
				stub_func_render_main_count++
			}, func(_ uint64) {
				stub_func_render_post_count++
			}
	}

	lp := NewLoudProgress(2, stub_supplier_func)
	lp.writer = &bytes.Buffer{} // dummy writer

	// fail starting
	lp.is_running = true
	err := lp.Start()
	if err == nil {
		t.Errorf("expected %v but given %v", nil, err)
	}

	// success starting
	lp.is_running = false
	lp.ch <- 2
	err = lp.Start() // and break immediately
	if err != nil {
		t.Errorf("expected %v but given %v", nil, err)
	}
	if stub_func_render_main_count != 0 && stub_func_render_post_count != 1 {
		t.Errorf("expected (%v, %v) but given (%v, %v)", 0, 1, stub_func_render_main_count, stub_func_render_post_count)
	}
}

func TestRender(t *testing.T) {
	const (
		eps  = 10 * time.Millisecond  // wait short time for execution in another goroutine
		wait = 250 * time.Millisecond // duration between renderings
	)

	stub_func_render_main_count := 0
	stub_func_render_post_count := 0
	stub_supplier_func := func(_ uint64, _ io.Writer) (func(uint64), func(uint64)) {
		return func(_ uint64) {
				stub_func_render_main_count++
			},
			func(_ uint64) {
				stub_func_render_post_count++
			}
	}
	lp := NewLoudProgress(2, stub_supplier_func)
	lp.writer = &bytes.Buffer{} // dummy writer
	ch := lp.ch

	// render start
	go lp.render()

	// first cycle
	time.Sleep(eps)
	if stub_func_render_main_count != 1 && stub_func_render_post_count != 0 {
		t.Errorf("expected (%v, %v) but given (%v, %v)", 1, 0, stub_func_render_main_count, stub_func_render_post_count)
	}
	ch <- 1
	time.Sleep(wait)

	// second cycle
	if stub_func_render_main_count != 2 && stub_func_render_post_count != 0 {
		t.Errorf("expected (%v, %v) but given (%v, %v)", 2, 0, stub_func_render_main_count, stub_func_render_post_count)
	}
	ch <- 2
	time.Sleep(wait)

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
