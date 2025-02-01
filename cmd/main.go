package main

import (
	"context"
	"fmt"
	"os"

	"github.com/llmcontext/gomcp"
	"github.com/llmcontext/gomcp/types"
)

type DummyToolConfiguration struct {
	Name string
}

type DummyContext struct {
	Name string
}

func DummyToolInit(ctx context.Context, config *DummyToolConfiguration) (*DummyContext, error) {
	return &DummyContext{
		Name: config.Name,
	}, nil
}

type DummyPingInput struct {
	Message string `json:"message" jsonschema_description:"the message to ping."`
}

func DummyPing(ctx context.Context, toolCtx *DummyContext, input *DummyPingInput, output types.ToolCallResult) error {
	output.AddTextContent(fmt.Sprintf("pong %s from %s", input.Message, toolCtx.Name))
	return nil
}

func main() {
	// create the mcpServerDefinition
	mcpServerDefinition := gomcp.NewMcpServerDefinition("dummy", "0.0.1")
	mcpServerDefinition.SetDebugLevel("debug", "debug.log")

	mcpToolsDefinition := mcpServerDefinition.WithTools(&DummyToolConfiguration{
		Name: "dummy",
	}, DummyToolInit)

	mcpToolsDefinition.AddTool("ping", "A ping function", DummyPing)

	mcp, err := gomcp.NewModelContextProtocolServer(mcpServerDefinition)
	if err != nil {
		fmt.Println("Error creating MCP server:", err)
		os.Exit(1)
	}
	// start the server
	transport := mcp.StdioTransport()
	err = mcp.Start(transport)
	if err != nil {
		fmt.Println("Error starting MCP server:", err)
		os.Exit(1)
	}
}
