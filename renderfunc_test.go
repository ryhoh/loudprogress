package loudprogress

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/fatih/color"
)

func TestRender_Simple(t *testing.T) {
	writer := &bytes.Buffer{}
	render_main_func, render_post_func := Render_Simple(40, writer)

	// current = 0
	writer.Reset()
	render_main_func(0)
	expect := fmt.Sprintf("\r[%s] %d/%d", strings.Repeat(" ", 40), 0, 40)
	actual := writer.String()
	if actual != expect {
		t.Errorf("expected\n%v\nbut given\n%v\n", expect, actual)
	}

	// current = 1
	writer.Reset()
	render_main_func(1)
	expect = fmt.Sprintf("\r[%s] %d/%d", ">"+strings.Repeat(" ", 39), 1, 40)
	actual = writer.String()
	if actual != expect {
		t.Errorf("expected\n%v\nbut given\n%v\n", expect, actual)
	}

	// current = 2
	writer.Reset()
	render_main_func(2)
	expect = fmt.Sprintf("\r[%s] %d/%d", "=>"+strings.Repeat(" ", 38), 2, 40)
	actual = writer.String()
	if actual != expect {
		t.Errorf("expected\n%v\nbut given\n%v\n", expect, actual)
	}

	// current = 3
	writer.Reset()
	render_main_func(3)
	expect = fmt.Sprintf("\r[%s] %d/%d", "==>"+strings.Repeat(" ", 37), 3, 40)
	actual = writer.String()
	if actual != expect {
		t.Errorf("expected\n%v\nbut given\n%v\n", expect, actual)
	}

	// current = 4
	writer.Reset()
	render_main_func(4)
	expect = fmt.Sprintf("\r[%s] %d/%d", ">==>"+strings.Repeat(" ", 36), 4, 40)
	actual = writer.String()
	if actual != expect {
		t.Errorf("expected\n%v\nbut given\n%v\n", expect, actual)
	}

	// finalize
	writer.Reset()
	render_post_func(40)
	expect = fmt.Sprintf("\r[%s] %d/%d", strings.Repeat("=", 40), 40, 40)
	actual = writer.String()
	if actual != expect {
		t.Errorf("expected\n%v\nbut given\n%v\n", expect, actual)
	}
}

func TestRender_Bold(t *testing.T) {
	writer := &bytes.Buffer{}
	render_main_func, render_post_func := Render_Bold(40, writer)

	// current = 0
	writer.Reset()
	render_main_func(0)
	expect := fmt.Sprintf("\r[%s] %d/%d", strings.Repeat(" ", 40), 0, 40)
	actual := writer.String()
	if actual != expect {
		t.Errorf("expected\n%v\nbut given\n%v\n", expect, actual)
	}

	// current = 1
	writer.Reset()
	render_main_func(1)
	expect = fmt.Sprintf("\r[%s] %d/%d", "█"+strings.Repeat(" ", 39), 1, 40)
	actual = writer.String()
	if actual != expect {
		t.Errorf("expected\n%v\nbut given\n%v\n", expect, actual)
	}

	// finalize
	writer.Reset()
	render_post_func(40)
	expect = fmt.Sprintf("\r[%s] %d/%d", strings.Repeat("█", 40), 40, 40)
	actual = writer.String()
	if actual != expect {
		t.Errorf("expected\n%v\nbut given\n%v\n", expect, actual)
	}
}

func TestRender_BoldRainbow(t *testing.T) {
	writer := &bytes.Buffer{}
	render_main_func, render_post_func := Render_BoldRainbow(40, writer)

	// current = 0
	writer.Reset()
	render_main_func(0)
	expect := fmt.Sprintf("\r[%s] %d/%d", strings.Repeat(" ", 40), 0, 40)
	actual := writer.String()
	if actual != expect {
		t.Errorf("expected\n%v\nbut given\n%v\n", expect, actual)
	}

	// current = 1
	writer.Reset()
	render_main_func(1)
	expect = fmt.Sprintf("\r[%s] %d/%d", color.RedString("█")+strings.Repeat(" ", 39), 1, 40)
	actual = writer.String()
	if actual != expect {
		t.Errorf("expected\n%v\nbut given\n%v\n", expect, actual)
	}

	// current = 2
	writer.Reset()
	render_main_func(2)
	expect = fmt.Sprintf("\r[%s] %d/%d", strings.Repeat(color.RedString("█"), 2)+strings.Repeat(" ", 38), 2, 40)
	actual = writer.String()
	if actual != expect {
		t.Errorf("expected\n%v\nbut given\n%v\n", expect, actual)
	}

	// current = 3
	writer.Reset()
	render_main_func(3)
	expect = fmt.Sprintf("\r[%s] %d/%d", strings.Repeat(color.RedString("█"), 3)+strings.Repeat(" ", 37), 3, 40)
	actual = writer.String()
	if actual != expect {
		t.Errorf("expected\n%v\nbut given\n%v\n", expect, actual)
	}

	// current = 4
	writer.Reset()
	render_main_func(4)
	expect = fmt.Sprintf("\r[%s] %d/%d", strings.Repeat(color.RedString("█"), 4)+strings.Repeat(" ", 36), 4, 40)
	actual = writer.String()
	if actual != expect {
		t.Errorf("expected\n%v\nbut given\n%v\n", expect, actual)
	}

	// current = 5
	writer.Reset()
	render_main_func(5)
	expect = fmt.Sprintf("\r[%s] %d/%d", color.YellowString("█")+strings.Repeat(color.RedString("█"), 4)+strings.Repeat(" ", 35), 5, 40)
	actual = writer.String()
	if actual != expect {
		t.Errorf("expected\n%v\nbut given\n%v\n", expect, actual)
	}

	// finalize
	writer.Reset()
	render_post_func(40)
	expect = fmt.Sprintf("\r[%s] %d/%d",
		strings.Repeat(color.YellowString("█"), 2)+strings.Repeat(color.RedString("█"), 4)+strings.Repeat(color.MagentaString("█"), 4)+
			strings.Repeat(color.BlueString("█"), 4)+strings.Repeat(color.CyanString("█"), 4)+strings.Repeat(color.GreenString("█"), 4)+
			strings.Repeat(color.YellowString("█"), 4)+strings.Repeat(color.RedString("█"), 4)+strings.Repeat(color.MagentaString("█"), 4)+
			strings.Repeat(color.BlueString("█"), 4)+strings.Repeat(color.CyanString("█"), 2),
		40, 40,
	)
	actual = writer.String()
	if actual != expect {
		t.Errorf("expected\n%v\nbut given\n%v\n", expect, actual)
	}
}
