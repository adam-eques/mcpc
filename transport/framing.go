package transport

import (
	"bufio"
	"bytes"
	"errors"
)

// MaxFrameBytes bounds a single inbound frame to protect the client from a
// server that never emits a newline.
const MaxFrameBytes = 16 << 20 // 16 MiB

// ErrFrameTooLarge is returned when an inbound frame exceeds MaxFrameBytes.
var ErrFrameTooLarge = errors.New("transport: frame exceeds maximum size")

var newline = []byte{'\n'}

// readFrame reads a single newline-terminated frame, skipping blank lines and
// stripping the trailing CR/LF. It returns the terminal error (for example
// io.EOF) once no more data is available.
func readFrame(br *bufio.Reader) ([]byte, error) {
	var buf []byte
	for {
		chunk, err := br.ReadSlice('\n')
		buf = append(buf, chunk...)
		if len(buf) > MaxFrameBytes {
			return nil, ErrFrameTooLarge
		}
		if errors.Is(err, bufio.ErrBufferFull) {
			continue
		}
		if err != nil {
			line := bytes.TrimRight(buf, "\r\n")
			if len(line) == 0 {
				return nil, err
			}
			return line, nil
		}
		line := bytes.TrimRight(buf, "\r\n")
		if len(line) == 0 {
			buf = buf[:0]
			continue // skip blank keep-alive lines
		}
		return line, nil
	}
}
