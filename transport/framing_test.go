package transport

import (
	"bufio"
	"errors"
	"io"
	"strings"
	"testing"
)

func TestReadFrameSkipsBlankLines(t *testing.T) {
	br := bufio.NewReader(strings.NewReader("first\n\n\nsecond\n"))
	f1, err := readFrame(br)
	if err != nil || string(f1) != "first" {
		t.Fatalf("frame1=%q err=%v", f1, err)
	}
	f2, err := readFrame(br)
	if err != nil || string(f2) != "second" {
		t.Fatalf("frame2=%q err=%v", f2, err)
	}
}

func TestReadFrameStripsCR(t *testing.T) {
	br := bufio.NewReader(strings.NewReader("data\r\n"))
	f, err := readFrame(br)
	if err != nil || string(f) != "data" {
		t.Fatalf("frame=%q err=%v", f, err)
	}
}

func TestReadFrameEOF(t *testing.T) {
	br := bufio.NewReader(strings.NewReader(""))
	if _, err := readFrame(br); !errors.Is(err, io.EOF) {
		t.Fatalf("expected EOF, got %v", err)
	}
}

func TestReadFrameTrailingWithoutNewline(t *testing.T) {
	br := bufio.NewReader(strings.NewReader("nolinebreak"))
	f, err := readFrame(br)
	if err != nil || string(f) != "nolinebreak" {
		t.Fatalf("frame=%q err=%v", f, err)
	}
}
