package transport

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
)

// Command is a transport that launches an MCP server as a child process and
// speaks newline-delimited JSON-RPC over its standard streams. This is the usual
// way a client connects to a stdio server such as mcpkit.
type Command struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	writeM sync.Mutex

	frames chan frame
	done   chan struct{}
	closeO sync.Once
}

type frame struct {
	data []byte
	err  error
}

// CommandOptions configure a Command transport.
type CommandOptions struct {
	// Dir is the working directory for the child process.
	Dir string
	// Env, when non-nil, replaces the child environment.
	Env []string
	// Stderr receives the child's standard error; defaults to os.Stderr.
	Stderr io.Writer
}

// StartCommand launches name with args and returns a ready transport. The child
// is terminated by Close. A read goroutine buffers inbound frames so Receive can
// honour context cancellation.
func StartCommand(name string, args []string, opts CommandOptions) (*Command, error) {
	cmd := exec.Command(name, args...)
	cmd.Dir = opts.Dir
	cmd.Env = opts.Env
	if opts.Stderr != nil {
		cmd.Stderr = opts.Stderr
	} else {
		cmd.Stderr = os.Stderr
	}

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("start %q: %w", name, err)
	}

	c := &Command{
		cmd:    cmd,
		stdin:  stdin,
		frames: make(chan frame),
		done:   make(chan struct{}),
	}
	go c.readLoop(bufio.NewReaderSize(stdout, 64*1024))
	return c, nil
}

func (c *Command) readLoop(br *bufio.Reader) {
	for {
		data, err := readFrame(br)
		select {
		case c.frames <- frame{data: data, err: err}:
		case <-c.done:
			return
		}
		if err != nil {
			return
		}
	}
}

// Receive implements Transport.
func (c *Command) Receive(ctx context.Context) ([]byte, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-c.done:
		return nil, io.EOF
	case f := <-c.frames:
		return f.data, f.err
	}
}

// Send implements Transport.
func (c *Command) Send(_ context.Context, frame []byte) error {
	if bytes.IndexByte(frame, '\n') >= 0 {
		return errors.New("transport: frame contains an embedded newline")
	}
	c.writeM.Lock()
	defer c.writeM.Unlock()
	if _, err := c.stdin.Write(frame); err != nil {
		return err
	}
	_, err := c.stdin.Write(newline)
	return err
}

// Close terminates the child process and releases resources.
func (c *Command) Close() error {
	var err error
	c.closeO.Do(func() {
		close(c.done)
		_ = c.stdin.Close()
		if c.cmd.Process != nil {
			_ = c.cmd.Process.Kill()
		}
		err = c.cmd.Wait()
	})
	// A killed process reports a non-nil exit error, which is expected on Close.
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return nil
		}
	}
	return err
}
