package transport

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/llmcontext/gomcp/inspector"
	"github.com/llmcontext/gomcp/logger"
	"github.com/llmcontext/gomcp/types"
)

type StdioTransport struct {
	debug             bool
	protocolDebugFile string
	inspector         *inspector.Inspector
	onMessage         func(json.RawMessage)
	onClose           func()
	onError           func(error)
	done              chan struct{}
}

func NewStdioTransport(protocolDebugFile string, inspector *inspector.Inspector) types.Transport {
	return &StdioTransport{
		debug:             protocolDebugFile != "",
		protocolDebugFile: protocolDebugFile,
		inspector:         inspector,
	}
}

func (t *StdioTransport) Start(ctx context.Context) error {
	t.done = make(chan struct{})

	// Start goroutine to read from stdin
	go t.readLoop()
	return nil
}

func (t *StdioTransport) Send(message json.RawMessage) error {
	// Write message followed by newline to stdout
	if t.debug {
		t.logProtocolMessages(string(message), "sending")
	}

	if t.inspector != nil {
		t.inspector.EnqueueMessage(inspector.MessageInfo{
			Timestamp: time.Now().Format(time.RFC3339),
			Direction: inspector.MessageDirectionResponse,
			Content:   string(message),
		})
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
					t.logProtocolMessages(string(line), "receiving")
				}

				if t.inspector != nil {
					t.inspector.EnqueueMessage(inspector.MessageInfo{
						Timestamp: time.Now().Format(time.RFC3339),
						Direction: inspector.MessageDirectionRequest,
						Content:   string(line),
					})
				}

				t.onMessage(json.RawMessage(line))
			}
		}
	}

	if err := scanner.Err(); err != nil && t.onError != nil {
		t.onError(err)
	}
}

func (t *StdioTransport) logProtocolMessages(rawMessage string, direction string) {
	// open log file and append
	file, err := os.OpenFile(t.protocolDebugFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logger.Error("error opening protocol debug file", logger.Arg{"error": err})
	}
	defer file.Close()

	// write to file
	_, err = file.WriteString(fmt.Sprintf(":%s: %s\n", direction, rawMessage))
	if err != nil {
		logger.Error("error writing to protocol debug file", logger.Arg{"error": err})
	}

	// try to parse the message as JSON
	var jsonMessage json.RawMessage
	err = json.Unmarshal([]byte(rawMessage), &jsonMessage)
	if err != nil {
		file.WriteString(fmt.Sprintf("!error parsing message as JSON: %s\n", err))
	} else {
		// pretty print the JSON message
		prettyJSON, err := json.MarshalIndent(jsonMessage, "", "  ")
		if err != nil {
			file.WriteString(fmt.Sprintf("!error formatting json message: %s\n", err))
		} else {
			file.WriteString(fmt.Sprintf("%s\n", string(prettyJSON)))
		}
	}
}
