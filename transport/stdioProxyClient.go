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

type ProxiedMcpServerDescription struct {
	ProxyId                 string
	ProxyName               string
	CurrentWorkingDirectory string
	ProgramName             string
	ProgramArgs             []string
}

type StdioProxyClientTransport struct {
	options   *ProxiedMcpServerDescription
	cmd       *exec.Cmd
	stdin     io.WriteCloser
	stdout    io.ReadCloser
	onMessage func(json.RawMessage)
	onClose   func()
	onError   func(error)
	onStarted func()
	// we need to keep track of the pipe reader
	pipeReader *io.PipeReader
}

func NewStdioProxyClientTransport(options *ProxiedMcpServerDescription) types.Transport {
	return &StdioProxyClientTransport{
		options:    options,
		pipeReader: nil,
	}
}

func (t *StdioProxyClientTransport) Start(ctx context.Context) error {
	t.cmd = exec.Command(t.options.ProgramName, t.options.ProgramArgs...)

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

	if err := t.cmd.Start(); err != nil {
		return fmt.Errorf("failed to start command: %v", err)
	}

	errChan := make(chan error, 1)
	go t.readLoop(ctx, errChan)

	select {
	case err := <-errChan:
		t.Close()
		return err
	case <-ctx.Done():
		t.Close()
		return ctx.Err()
	}
}

func (t *StdioProxyClientTransport) readLoop(ctx context.Context, errChan chan error) {
	// we create a pipe to read from stdin
	// and push the data to the pipe
	// we need to do this because we need to read from stdin in a non-blocking way
	// if we don't do this, we will block the readLoop
	// if we close the pipe, we will stop the readLoop
	r, w := io.Pipe()
	// we keep track of the pipe reader (to close it later)
	t.pipeReader = r

	go func() {
		defer w.Close()
		_, err := io.Copy(w, t.stdout)
		if err != nil {
			if t.onError != nil {
				t.onError(fmt.Errorf("error reading from process stdout: %w", err))
			}
			// send the error to the errChan
			errChan <- fmt.Errorf("error reading from process stdout: %w", err)
		}
	}()

	scanner := bufio.NewScanner(r)
	// Increase buffer size to handle larger lines
	const maxScanTokenSize = 1024 * 1024 * 10 // 10MB buffer
	buf := make([]byte, maxScanTokenSize)
	scanner.Buffer(buf, maxScanTokenSize)

	readCh := make(chan struct{})

	if t.onStarted != nil {
		t.onStarted()
	}

	for {
		go func() {
			if scanner.Scan() {
				if t.onMessage != nil {
					line := scanner.Text()
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

func (t *StdioProxyClientTransport) OnStarted(callback func()) {
	t.onStarted = callback
}

func (t *StdioProxyClientTransport) Close() {
	if t.cmd != nil && t.cmd.Process != nil {
		t.cmd.Process.Kill()
	}
	if t.pipeReader != nil {
		t.pipeReader.Close()
		t.pipeReader = nil
	}

	if t.onClose != nil {
		t.onClose()
	}
}

func (t *StdioProxyClientTransport) Send(message json.RawMessage) error {
	nlTerminatedMessage := string(message) + "\n"
	_, err := t.stdin.Write([]byte(nlTerminatedMessage))
	return err
}
