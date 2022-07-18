package loudprogress

import (
	"fmt"
	"io"

	color "github.com/fatih/color"
)

/* Progress bar functions */
// (size) -> (func_render_main, func_render_post)

// Simple progress bar
func Render_Simple(size int64, writer io.Writer) (func(int64), func(int64)) {
	const (
		width_all = 40 // progress width on TUI
		wave_span = 3  // wave char '>' appears every wave_span progress charactors
	)
	var (
		wave_pos int = 0 // [0, wave_span] wave char '>' appears each `(i+wave_pos) mod wave_span == 0` position
	)

	return func(current int64) {
			/* culc length on TUI */
			width_done := int(current * int64(width_all) / size)

			/* progress animation */
			fmt.Fprint(writer, "\r[")
			for i := 1; i < width_done; i++ { // render done bar
				if (i+wave_pos)%wave_span == 0 {
					fmt.Fprint(writer, ">")
				} else {
					fmt.Fprint(writer, "=")
				}
			}
			if width_done != 0 {
				fmt.Fprint(writer, ">") // right end must be ">"
			}
			for i := width_done + 1; i < width_all+1; i++ { // render untouched bar
				fmt.Fprint(writer, " ")
			}
			fmt.Fprintf(writer, "] %d/%d", current, size) // print current/size as number

			/* wait for next render */
			wave_pos--
			if wave_pos < 0 {
				wave_pos += wave_span
			}
		},
		func(current int64) {
			/* render finished progress */
			fmt.Fprint(writer, "\r[")
			for i := 1; i < width_all+1; i++ { // render done bar
				fmt.Fprint(writer, "=")
			}
			fmt.Fprintf(writer, "] %d/%d\n", current, size)
		}
}

// Bold progress bar
func Render_Bold(size int64, writer io.Writer) (func(int64), func(int64)) {
	const (
		width_all = 40 // progress width on TUI
	)

	return func(current int64) {
			/* culc length on TUI */
			width_done := int(current * int64(width_all) / size)

			/* progress animation */
			fmt.Fprint(writer, "\r[")
			for i := 1; i < width_done+1; i++ { // render done bar
				fmt.Fprint(writer, "█")
			}
			for i := width_done + 1; i < width_all+1; i++ { // render untouched bar
				fmt.Fprint(writer, " ")
			}
			fmt.Fprintf(writer, "] %d/%d", current, size) // print current/size as number
		},
		func(current int64) {
			/* render finished progress */
			fmt.Fprint(writer, "\r[")
			for i := 1; i < width_all+1; i++ { // render done bar
				fmt.Fprint(writer, "█")
			}
			fmt.Fprintf(writer, "] %d/%d\n", current, size)
		}
}

// Bold rainbow progress bar
func Render_BoldRainbow(size int64, writer io.Writer) (func(int64), func(int64)) {
	const (
		width_all            = 40                           // progress width on TUI
		color_num            = 6                            // number of colors
		color_band_width     = 4                            // width of band for each color
		color_spector_length = color_band_width * color_num // length of color spector (means 1 rainbow cycle)
	)
	var (
		color_func_table = [color_num]func(format string, a ...interface{}) string{ // color table for rainbow
			color.MagentaString,
			color.BlueString,
			color.CyanString,
			color.GreenString,
			color.YellowString,
			color.RedString,
		}
		head_pos int = 0 // [0, color_spector_length - 1] head position of rainbow
	)

	return func(current int64) {
			/* culc length on TUI */
			width_done := int(current * int64(width_all) / size)

			/* progress animation */
			fmt.Fprint(writer, "\r[")
			for i := 1; i < width_done+1; i++ { // render done bar
				color_index := ((head_pos + i - 1) % color_spector_length) / color_band_width
				color_func := color_func_table[color_index]
				fmt.Fprint(writer, color_func("█"))
			}
			for i := width_done + 1; i < width_all+1; i++ { // render untouched bar
				fmt.Fprint(writer, " ")
			}
			fmt.Fprintf(writer, "] %d/%d", current, size) // print current/size as number

			/* wait for next render */
			head_pos--
			if head_pos < 0 {
				head_pos += color_spector_length
			}
		},
		func(current int64) {
			/* render finished progress */
			fmt.Fprint(writer, "\r[")
			for i := 1; i < width_all+1; i++ { // render done bar
				color_index := ((head_pos + i - 1) % color_spector_length) / color_band_width
				color_func := color_func_table[color_index]
				fmt.Fprint(writer, color_func("█"))
			}
			fmt.Fprintf(writer, "] %d/%d\n", current, size)
		}
}
