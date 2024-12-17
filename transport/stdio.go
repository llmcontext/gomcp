package transport

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/llmcontext/gomcp/channels/hubinspector"
	"github.com/llmcontext/gomcp/types"
)

type StdioTransport struct {
	debug             bool
	protocolDebugFile string
	inspector         *hubinspector.Inspector
	logger            types.Logger
	onStarted         func()
	onMessage         func(json.RawMessage)
	onClose           func()
	onError           func(error)

	// we need to keep track of the pipe reader
	pipeReader *io.PipeReader
}

func NewStdioTransport(protocolDebugFile string, inspector *hubinspector.Inspector, logger types.Logger) types.Transport {
	return &StdioTransport{
		debug:             protocolDebugFile != "",
		protocolDebugFile: protocolDebugFile,
		inspector:         inspector,
		logger:            logger,
		pipeReader:        nil,
	}
}

func (t *StdioTransport) Start(ctx context.Context) error {
	// we create a channel to report errors
	errChan := make(chan error, 1)

	// Start goroutine to read from stdin
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

func (t *StdioTransport) Send(message json.RawMessage) error {
	// Write message followed by newline to stdout
	if t.debug {
		t.logProtocolMessages(string(message), "sending")
	}

	if t.inspector != nil {
		t.inspector.EnqueueMessage(hubinspector.MessageInfo{
			Timestamp: time.Now().Format(time.RFC3339),
			Direction: hubinspector.MessageDirectionResponse,
			Content:   string(message),
		})
	}

	_, err := fmt.Fprintf(os.Stdout, "%s\n", message)
	return err
}

func (t *StdioTransport) OnStarted(callback func()) {
	t.onStarted = callback
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
	// close the pipe reader
	if t.pipeReader != nil {
		t.pipeReader.Close()
		t.pipeReader = nil
	}

	// report the close
	if t.onClose != nil {
		t.onClose()
	}

}

func (t *StdioTransport) readLoop(ctx context.Context, errChan chan error) {
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
		_, err := io.Copy(w, os.Stdin)
		if err != nil && err != io.EOF {
			if t.onError != nil {
				t.onError(fmt.Errorf("error reading from stdin: %w", err))
			}
			// send the error to the errChan
			errChan <- fmt.Errorf("error reading from stdin: %w", err)
		}
	}()

	// call the onStarted callback
	if t.onStarted != nil {
		t.onStarted()
	}

	// we create a scanner to read from the pipe
	scanner := bufio.NewScanner(r)
	readCh := make(chan struct{})
	for {
		go func() {
			if scanner.Scan() {
				if t.onMessage != nil {
					line := scanner.Text()
					if t.debug {
						t.logProtocolMessages(line, "receiving")
					}

					if t.inspector != nil {
						t.inspector.EnqueueMessage(hubinspector.MessageInfo{
							Timestamp: time.Now().Format(time.RFC3339),
							Direction: hubinspector.MessageDirectionRequest,
							Content:   line,
						})
					}
					t.onMessage(json.RawMessage(line))
				}
			}
			close(readCh)
		}()

		select {
		// check if we need to stop
		case <-ctx.Done():
			// close the pipe
			t.Close()
			return
		case <-readCh:
			// we've read a message, let's read another one
			readCh = make(chan struct{})
			continue
		}
	}
}

func (t *StdioTransport) logProtocolMessages(rawMessage string, direction string) {
	// open log file and append
	file, err := os.OpenFile(t.protocolDebugFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		t.logger.Error("error opening protocol debug file", types.LogArg{"error": err})
	}
	defer file.Close()

	// write to file
	_, err = file.WriteString(fmt.Sprintf(":%s: %s\n", direction, rawMessage))
	if err != nil {
		t.logger.Error("error writing to protocol debug file", types.LogArg{"error": err})
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
