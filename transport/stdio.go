package transport

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

type StdioTransport struct {
	onMessage func(json.RawMessage)
	onClose   func()
	onError   func(error)
	done      chan struct{}
}

func NewStdioTransport() *StdioTransport {
	return &StdioTransport{}
}

func (t *StdioTransport) Start() error {
	t.done = make(chan struct{})

	// Start goroutine to read from stdin
	go t.readLoop()
	return nil
}

func (t *StdioTransport) Send(message json.RawMessage) error {
	// Write message followed by newline to stdout
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
				t.onMessage(json.RawMessage(scanner.Bytes()))
			}
		}
	}

	if err := scanner.Err(); err != nil && t.onError != nil {
		t.onError(err)
	}
}
