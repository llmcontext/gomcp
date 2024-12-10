package transport

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"

	"github.com/llmcontext/gomcp/types"
)

type StdioProxyClientTransport struct {
	programName string
	args        []string
	cmd         *exec.Cmd
	stdin       io.WriteCloser
	stdout      io.ReadCloser
	onMessage   func(json.RawMessage)
	onClose     func()
	onError     func(error)
	done        chan struct{}
}

func NewStdioProxyClientTransport(programName string, args []string) types.Transport {
	return &StdioProxyClientTransport{
		programName: programName,
		args:        args,
	}
}

func (t *StdioProxyClientTransport) Start(ctx context.Context) error {
	fmt.Printf("@@ Starting %s with args %v\n", t.programName, t.args)
	t.cmd = exec.Command(t.programName, t.args...)

	stdin, err := t.cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdin pipe: %v", err)
	}
	t.stdin = stdin

	stdout, err := t.cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %v", err)
	}
	t.stdout = stdout

	t.done = make(chan struct{})
	if err := t.cmd.Start(); err != nil {
		return fmt.Errorf("failed to start command: %v", err)
	}
	go t.readLoop(ctx)
	return nil
}

func (t *StdioProxyClientTransport) readLoop(ctx context.Context) {
	r, w := io.Pipe()
	go func() {
		defer w.Close()
		io.Copy(w, t.stdout)
	}()

	scanner := bufio.NewScanner(r)
	readCh := make(chan struct{})

	for {
		go func() {
			if scanner.Scan() {
				if t.onMessage != nil {
					line := scanner.Text()
					fmt.Printf("@@ [proxy] received message: (%s)\n", line)
					t.onMessage(json.RawMessage(line))
				}
			}
			close(readCh)
		}()

		select {
		// check if we need to stop
		case <-ctx.Done():
			// close the pipe
			r.Close()
			return
		// we'veread a message, let's read another one
		case <-readCh:
			readCh = make(chan struct{})
			continue
		}
	}
}

func (t *StdioProxyClientTransport) OnMessage(callback func(json.RawMessage)) {
	t.onMessage = callback
}

func (t *StdioProxyClientTransport) OnClose(callback func()) {
	t.onClose = callback
}

func (t *StdioProxyClientTransport) OnError(callback func(error)) {
	t.onError = callback
}

func (t *StdioProxyClientTransport) Close() {
	if t.cmd != nil && t.cmd.Process != nil {
		t.cmd.Process.Kill()
	}

	close(t.done)
	if t.onClose != nil {
		t.onClose()
	}
}

func (t *StdioProxyClientTransport) Send(message json.RawMessage) error {
	nlTerminatedMessage := string(message) + "\n"
	_, err := t.stdin.Write([]byte(nlTerminatedMessage))
	return err
}
