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
	if len(os.Args) > 1 {
		u, _ := strconv.ParseInt(os.Args[1], 10, 32)
		width = int(u)
	} else {
		width = 64
	}

	data, _ := ioutil.ReadFile("/dev/stdin")
	height := (len(data) + width - 1) / width;

	w := bufio.NewWriter(os.Stdout)
	for y := 0; y < height; y += 2 {
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
