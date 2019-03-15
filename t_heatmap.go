package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

func main() {
	var width int
	var y_step int = int(^uint(0) >> 1)
	if len(os.Args) > 1 {
		u, _ := strconv.ParseInt(os.Args[1], 10, 32)
		width = int(u)
		if len(os.Args) > 2 {
			v, _ := strconv.ParseInt(os.Args[2], 10, 32)
			y_step = int(v)
		}
	} else {
		width = 64
	}

	data, _ := ioutil.ReadFile("/dev/stdin")
	height := (len(data) + width - 1) / width;

	w := bufio.NewWriter(os.Stdout)
	for y := 0; y < height; y += 2 {
		if y > 0 && (y % y_step == 0) {
			_, _ = fmt.Fprintln(w)
		}

		ih := width * y;
		il := ih + width;
		for x := 0; x < width; x++ {
			var vh byte = 0
			if ih < len(data) {
				vh = data[ih]
			}
			var vl byte = 0
			if il < len(data) {
				vl = data[il]
			}

			_, _ = fmt.Fprintf(w, "\x1B[38;2;%d;%d;%dm\x1B[48;2;%d;%d;%dm\u2580\x1B[0m", vh, vh, vh, vl, vl, vl)

			ih++
			il++
		}
		_, _ = fmt.Fprintln(w)
	}
	_ = w.Flush()
}
