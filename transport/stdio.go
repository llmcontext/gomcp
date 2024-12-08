package transport

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/llmcontext/gomcp/logger"
)

type StdioTransport struct {
	debug     bool
	onMessage func(json.RawMessage)
	onClose   func()
	onError   func(error)
	done      chan struct{}
}

func NewStdioTransportWithDebug() *StdioTransport {
	return &StdioTransport{
		debug: true,
	}
}

func NewStdioTransport() *StdioTransport {
	return &StdioTransport{
		debug: false,
	}
}

func (t *StdioTransport) Start() error {
	t.done = make(chan struct{})

	// Start goroutine to read from stdin
	go t.readLoop()
	return nil
}

func (t *StdioTransport) Send(message json.RawMessage) error {
	// Write message followed by newline to stdout
	if t.debug {
		logger.Info(string(message), logger.Arg{"direction": "sending"})
	}
	_, err := fmt.Fprintf(os.Stdout, "%s\n", message)
	return err
}

func (t *StdioTransport) OnMessage(callback func(json.RawMessage)) {
	t.onMessage = callback
}

func (t *StdioTransport) OnClose(callback func()) {
	t.onClose = callback
}

func (t *StdioTransport) OnError(callback func(error)) {
	t.onError = callback
}

func (t *StdioTransport) Close() {
	close(t.done)
	if t.onClose != nil {
		t.onClose()
	}
}

func (t *StdioTransport) readLoop() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		select {
		case <-t.done:
			return
		default:
			if t.onMessage != nil {
				line := scanner.Bytes()
				if t.debug {
					logger.Info(string(line), logger.Arg{"direction": "receiving"})
				}

				t.onMessage(json.RawMessage(line))
			}
		}
	}

	if err := scanner.Err(); err != nil && t.onError != nil {
		t.onError(err)
	}
}
