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
// Command line parameters: <format> <width> [separator stride] [color scheme]
//
// <format>: u1 or u8
//
// If Separator Stride is set, then a line break is inserted after Stride lines of "pixels" (it has to be even number).
// This is useful to display multiple "heatmaps" with the same geometry, contained in the data stream.
//
// Color scheme:
// - If set to "256" and format is "u8", then bytes are colored using ANSI 256-color scheme.
// - Otherwise, bytes are colored using grayscale scheme.
func main() {
	run()
}

func run() {
	var width = 64
	var ySeparatorStride = int(^uint(0) >> 1) // large number?
	var dataFormat = "u8"
	var ansi256 = false

	w := bufio.NewWriter(os.Stdout)

	if len(os.Args) > 1 {
		dataFormat = os.Args[1]
		if len(os.Args) > 2 {
			u, _ := strconv.ParseInt(os.Args[2], 10, 32)
			width = int(u)
			if len(os.Args) > 3 {
				v, _ := strconv.ParseInt(os.Args[3], 10, 32)
				ySeparatorStride = int(v)
			}
			if len(os.Args) > 4 {
				ansi256 = true
			}
		}
	}

	data, _ := ioutil.ReadFile("/dev/stdin")

	if dataFormat == "u1" {
		if width%8 != 0 {
			os.Exit(1)
		}
		heatMapU1(w, data, width, ySeparatorStride)
	} else {
		heatMapU8(w, data, width, ySeparatorStride, ansi256)
	}

	_ = w.Flush()
}

func heatMapU1(w *bufio.Writer, data []byte, widthBits int, ySeparatorStride int) {
	widthBytes := widthBits / 8
	height := (len(data) + widthBytes - 1) / widthBytes

	for y := 0; y < height; y += 2 {
		if y > 0 && (y%ySeparatorStride == 0) {
			_, _ = fmt.Fprintln(w)
		}

		topByteOffset := widthBytes * y
		botByteOffset := topByteOffset + widthBytes
		for x := 0; x < widthBytes; x++ {
			var topByte byte = 0
			if topByteOffset < len(data) {
				topByte = data[topByteOffset]
			}
			var btmByte byte = 0
			if botByteOffset < len(data) {
				btmByte = data[botByteOffset]
			}

			renderBytes(w, topByte, btmByte)

			topByteOffset++
			botByteOffset++
		}
		_, _ = fmt.Fprintln(w)
	}
}

func renderBytes(w *bufio.Writer, topByte byte, btmByte byte) {
	renderBits(w, topByte, btmByte, 0)
	renderBits(w, topByte, btmByte, 1)
	renderBits(w, topByte, btmByte, 2)
	renderBits(w, topByte, btmByte, 3)
	renderBits(w, topByte, btmByte, 4)
	renderBits(w, topByte, btmByte, 5)
	renderBits(w, topByte, btmByte, 6)
	renderBits(w, topByte, btmByte, 7)
}

func renderBits(w *bufio.Writer, topByte byte, btmByte byte, bit int) {
	bitH := 0
	if topByte&(1<<bit) != 0 {
		bitH = 255
	}
	bitL := 0
	if btmByte&(1<<bit) != 0 {
		bitL = 255
	}
	_, _ = fmt.Fprintf(w, "\x1B[38;2;%d;%d;%dm\x1B[48;2;%d;%d;%dm\u2580\x1B[0m", bitH, bitH, bitH, bitL, bitL, bitL)
}

func heatMapU8(w *bufio.Writer, data []byte, width int, ySeparatorStride int, ansi256 bool) {
	height := (len(data) + width - 1) / width

	for y := 0; y < height; y += 2 {
		if y > 0 && (y%ySeparatorStride == 0) {
			_, _ = fmt.Fprintln(w)
		}

		hiByteOffset := width * y
		loByteOffset := hiByteOffset + width
		for x := 0; x < width; x++ {
			var vh byte = 0
			if hiByteOffset < len(data) {
				vh = data[hiByteOffset]
			}
			var vl byte = 0
			if loByteOffset < len(data) {
				vl = data[loByteOffset]
			}

			if ansi256 {
				_, _ = fmt.Fprintf(w, "\x1B[38;5;%dm\x1B[48;5;%dm\u2580\x1B[0m", vh, vl)
			} else {
				_, _ = fmt.Fprintf(w, "\x1B[38;2;%d;%d;%dm\x1B[48;2;%d;%d;%dm\u2580\x1B[0m", vh, vh, vh, vl, vl, vl)
			}

			hiByteOffset++
			loByteOffset++
		}
		_, _ = fmt.Fprintln(w)
	}
}
