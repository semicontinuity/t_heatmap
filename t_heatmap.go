package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

// Displays Heat Map for the binary data, fed to STDIN, in the terminal.
// Every terminal character contains 2 vertical "pixels"
//
// Command line parameters: <format> <width> <separator stride>
//
// <format>: u1 or u8
//
// If Separator Stride is set, then a line break is inserted after Stride lines of "pixels" (it has to be even number).
// This is useful to display multiple "heatmaps" with the same geometry, contained in the data stream.
func main() {
	run()
}

func run() {
	var width = 64
	var y_separator_stride = int(^uint(0) >> 1) // large number?
	var dataFormat = "u8"

	w := bufio.NewWriter(os.Stdout)

	if len(os.Args) > 1 {
		dataFormat = os.Args[1]
		if len(os.Args) > 2 {
			u, _ := strconv.ParseInt(os.Args[2], 10, 32)
			width = int(u)
			if len(os.Args) > 3 {
				v, _ := strconv.ParseInt(os.Args[3], 10, 32)
				y_separator_stride = int(v)
			}
		}
	}

	data, _ := ioutil.ReadFile("/dev/stdin")

	if dataFormat == "u1" {
		if width%8 != 0 {
			os.Exit(1)
		}
		heatMapU1(w, data, width, y_separator_stride)
	} else {
		heatMapU8(w, data, width, y_separator_stride)
	}

	_ = w.Flush()
}

func heatMapU1(w *bufio.Writer, data []byte, widthBits int, y_separator_stride int) {
	widthBytes := widthBits / 8
	height := (len(data) + widthBytes - 1) / widthBytes

	for y := 0; y < height; y += 2 {
		if y > 0 && (y%y_separator_stride == 0) {
			_, _ = fmt.Fprintln(w)
		}

		top_byte_offset := widthBytes * y
		bot_byte_offset := top_byte_offset + widthBytes
		for x := 0; x < widthBytes; x++ {
			var top_byte byte = 0
			if top_byte_offset < len(data) {
				top_byte = data[top_byte_offset]
			}
			var bot_byte byte = 0
			if bot_byte_offset < len(data) {
				bot_byte = data[bot_byte_offset]
			}

			renderBytes(w, top_byte, bot_byte)

			top_byte_offset++
			bot_byte_offset++
		}
		_, _ = fmt.Fprintln(w)
	}
}

func renderBytes(w *bufio.Writer, top_byte byte, bot_byte byte) {
	renderBits(w, top_byte, bot_byte, 0)
	renderBits(w, top_byte, bot_byte, 1)
	renderBits(w, top_byte, bot_byte, 2)
	renderBits(w, top_byte, bot_byte, 3)
	renderBits(w, top_byte, bot_byte, 4)
	renderBits(w, top_byte, bot_byte, 5)
	renderBits(w, top_byte, bot_byte, 6)
	renderBits(w, top_byte, bot_byte, 7)
}

func renderBits(w *bufio.Writer, top_byte byte, bot_byte byte, bit int) {
	bitH := 0
	if top_byte&(1<<bit) != 0 {
		bitH = 255
	}
	bitL := 0
	if bot_byte&(1<<bit) != 0 {
		bitL = 255
	}
	_, _ = fmt.Fprintf(w, "\x1B[38;2;%d;%d;%dm\x1B[48;2;%d;%d;%dm\u2580\x1B[0m", bitH, bitH, bitH, bitL, bitL, bitL)
}

func heatMapU8(w *bufio.Writer, data []byte, width int, y_separator_stride int) {
	height := (len(data) + width - 1) / width

	for y := 0; y < height; y += 2 {
		if y > 0 && (y%y_separator_stride == 0) {
			_, _ = fmt.Fprintln(w)
		}

		hi_byte_offset := width * y
		lo_byte_offset := hi_byte_offset + width
		for x := 0; x < width; x++ {
			var vh byte = 0
			if hi_byte_offset < len(data) {
				vh = data[hi_byte_offset]
			}
			var vl byte = 0
			if lo_byte_offset < len(data) {
				vl = data[lo_byte_offset]
			}

			_, _ = fmt.Fprintf(w, "\x1B[38;2;%d;%d;%dm\x1B[48;2;%d;%d;%dm\u2580\x1B[0m", vh, vh, vh, vl, vl, vl)

			hi_byte_offset++
			lo_byte_offset++
		}
		_, _ = fmt.Fprintln(w)
	}
}
