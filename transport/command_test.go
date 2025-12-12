package transport

import (
	"context"
	"runtime"
	"strings"
	"testing"
)

func TestStartCommandFrames(t *testing.T) {
	// A process that prints two lines proves the transport spawns the child and
	// reads newline-delimited frames from its stdout.
	var name string
	var args []string
	if runtime.GOOS == "windows" {
		name, args = "cmd", []string{"/c", "echo one& echo two"}
	} else {
		name, args = "sh", []string{"-c", "printf 'one\\ntwo\\n'"}
	}
	tr, err := StartCommand(name, args, CommandOptions{})
	if err != nil {
		t.Fatal(err)
	}
	defer tr.Close()

	ctx := context.Background()
	first, err := tr.Receive(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(first), "one") {
		t.Fatalf("first frame=%q", first)
	}
	second, err := tr.Receive(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(second), "two") {
		t.Fatalf("second frame=%q", second)
	}
}

func TestStartCommandNotFound(t *testing.T) {
	if _, err := StartCommand("this-binary-does-not-exist-xyz", nil, CommandOptions{}); err == nil {
		t.Fatal("expected error starting a missing binary")
	}
}

func TestSendRejectsEmbeddedNewline(t *testing.T) {
	var name string
	var args []string
	if runtime.GOOS == "windows" {
		name, args = "cmd", []string{"/c", "more"}
	} else {
		name, args = "cat", nil
	}
	tr, err := StartCommand(name, args, CommandOptions{})
	if err != nil {
		t.Fatal(err)
	}
	defer tr.Close()
	if err := tr.Send(context.Background(), []byte("a\nb")); err == nil {
		t.Fatal("expected error for embedded newline")
	}
}
