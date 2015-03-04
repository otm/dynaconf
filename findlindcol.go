package main

import (
	"bytes"
	"io"
)

func findLineCol(r io.Reader, offset int64) (line int, col int, err error) {
	buf := make([]byte, 1)
	line = 1
	col = 0
	lineSep := []byte{'\n'}

	for offset > 0 {
		_, err = r.Read(buf)
		if err != nil {
			return -1, -1, err
		}

		col++
		if bytes.Equal(buf, lineSep) {
			line++
			col = 0
		}

		offset--
	}

	return line, col, nil
}
