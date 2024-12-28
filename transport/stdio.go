package transport

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/llmcontext/gomcp/channels/hubinspector"
	"github.com/llmcontext/gomcp/types"
)

type StdioTransport struct {
	isClosed  bool
	inspector *hubinspector.Inspector
	logger    types.Logger
	onStarted func()
	onMessage func(json.RawMessage)
	onClose   func()
	onError   func(error)
}

func NewStdioTransport(inspector *hubinspector.Inspector, logger types.Logger) types.Transport {
	return &StdioTransport{
		inspector: inspector,
		logger:    logger,
		isClosed:  false,
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
	// check if we are already closed
	if t.isClosed {
		return
	}
	t.isClosed = true

	// close the stdin
	os.Stdin.Close()

	// report the close
	if t.onClose != nil {
		t.onClose()
	}
}

func (t *StdioTransport) readLoop(ctx context.Context, errChan chan error) {
	// call the onStarted callback
	if t.onStarted != nil {
		t.onStarted()
	}

	// Start a goroutine to read from stdin
	go func() {
		// we create a scanner to read from the pipe
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			select {
			case <-ctx.Done():
				t.logger.Info("@@ stdio transport - context done in scanner", types.LogArg{})
				t.Close()
				return
			default:
				line := scanner.Text()
				if t.onMessage != nil {
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
		}
		if err := scanner.Err(); err != nil {
			t.logger.Error("error reading from stdin", types.LogArg{"error": err})
			errChan <- fmt.Errorf("error reading from stdin: %w", err)
		}
		// we reach that when the MCP client (eg Claude) closes the connection
		t.logger.Info("stdio transport - readLoop() done", types.LogArg{})
		errChan <- fmt.Errorf("MCP client closed the connection")
	}()
}

// func (t *StdioTransport) logProtocolMessages(rawMessage string, direction string) {
// 	// open log file and append
// 	file, err := os.OpenFile(t.protocolDebugFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
// 	if err != nil {
// 		t.logger.Error("error opening protocol debug file", types.LogArg{"error": err})
// 	}
// 	defer file.Close()

// 	// write to file
// 	_, err = file.WriteString(fmt.Sprintf(":%s: %s\n", direction, rawMessage))
// 	if err != nil {
// 		t.logger.Error("error writing to protocol debug file", types.LogArg{"error": err})
// 	}

// 	// try to parse the message as JSON
// 	var jsonMessage json.RawMessage
// 	err = json.Unmarshal([]byte(rawMessage), &jsonMessage)
// 	if err != nil {
// 		file.WriteString(fmt.Sprintf("!error parsing message as JSON: %s\n", err))
// 	} else {
// 		// pretty print the JSON message
// 		prettyJSON, err := json.MarshalIndent(jsonMessage, "", "  ")
// 		if err != nil {
// 			file.WriteString(fmt.Sprintf("!error formatting json message: %s\n", err))
// 		} else {
// 			file.WriteString(fmt.Sprintf("%s\n", string(prettyJSON)))
// 		}
// 	}
// }
