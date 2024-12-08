# gomcp

## Description
An unofficial Golang implementation of the Model Context Protocol defined by Anthropic.

The officialannouncement of the Model Context Protocol is [here](https://www.anthropic.com/news/model-context-protocol).

The Model Context Protocol (MCP) provides a standardized, secure mechanism for AI models to interact with external tools and data sources. 

By defining a precise interface and communication framework, MCP allows AI assistants, such as the Claude desktop application, to safely 
extend their capabilities.

## Reference documentation

* the full documentation of the Model Context Protocol is [here](https://modelcontextprotocol.io/introduction)
* the official TypeScript SDK is available[here](https://github.com/modelcontextprotocol/typescript-sdk)
* the official Python SDK is available[here](https://github.com/modelcontextprotocol/python-sdk)

## Installation

```
go get github.com/llmcontext/gomcp
```

Direct dependencies:

* [github.com/invopop/jsonschema](https://github.com/invopop/jsonschema): to generate JSON Schemas from Go types through reflection
* [github.com/xeipuuv/gojsonschema](https://github.com/xeipuuv/gojsonschema): for JSON schema validation
* [gopkg.in/yaml.v3](https://github.com/go-yaml/yaml): for YAML parsing of the prompts definitionfile
* [github.com/stretchr/testify](https://github.com/stretchr/testify): for testing
* [go.uber.org/zap](https://github.com/uber-go/zap): for logging    

## Usage

Let's consider a simple example where we want to use the `mcp` package to create a server that can retrieve the content of a Notion page so that you can use it in a Claude chat.

The way to do this is to define a set of tools that can be exposed to the LLM through the Model Context Protocol and to implement them in Go.

The first step is to define the tools that will be exposed to the LLM.

In Mcp, you define a set of `Tool providers` and each provider is a set of `Tools`.

In our case, we have a single provider called `notion` that has a single tool to retrieve the content of a Notion page.

### configuration file

Once compiled, the `mcp` command needs a configuration file to start the server.

An example of configuration file would be:

```json
{
    "serverInfo": {
        "name": "gomcp",
        "version": "0.1.0"
    },
    "logging": {
        "file": "/var/log/gomcp/mcpnotion.log",
        "level": "debug",
        "withStderr": false
    },
    "prompts": {
        "file": "/etc/gomcp/prompts.yaml"
    },
    "tools": [
        {
            "name": "notion",
            "description": "Get a notion document",
            "configuration": {
                "notionToken": "ntn_<redacted>"
            }
        }
    ]
}
```

The `serverInfo` section is used to identify the server and its version, it is mandatory as they are used in the MCP protocolto identify the server.

The `logging` section is used to configure the logging system. The `file` field is the path to the log file, the `level` field is the logging level (debug, info, warn, error) and the `withStderr` field is used to redirect the logging to the standard error stream.

The `prompts` section is used to define the path to the YAML file containing the prompts to expose to the LLM. See below for a description of the YAML syntax to define the prompts.

The `tools` section is used to define the tools that will be exposed to the LLM. This is an array of tool providers, each provider is an object with a `name` and a `description` field. The `configuration` field is an object that contains the configuration for the tool provider.

In our case, we have a single tool provider called `notion` that has a single tool to retrieve the content of a Notion page.

The configuration for the `notion` tool provider is the Notion token.

This configuration must be backed by a Golang struct that will be used to parse the configuration file:

```go
type NotionGetDocumentConfiguration struct {
	NotionToken string `json:"notionToken" jsonschema_description:"the notion token for the Notion client."`
}
```
The tags here (`json` and `jsonschema_description`) are used to generate the JSON Schema for the configuration data.
If the configuration is invalid, the `mcp` command will fail to start.

You then create a function that will use those configuration data to generate a `Tool Context`:

```go
func NotionToolInit(ctx context.Context, config *NotionGetDocumentConfiguration) (*NotionGetDocumentContext, error) {
	client := notionapi.NewClient(notionapi.Token(config.NotionToken))

	// we need to initialize the Notion client
	return &NotionGetDocumentContext{NotionClient: client}, nil
}
```

And the definition of the `Tool Context`:

```go
type NotionGetDocumentContext struct {
	// The Notion client.
	NotionClient *notionapi.Client
}
```
This time, you don't need to add tags to your type as this struct is internal to the tool provider: the instance of this struct is created by the `ToolInit` function and is passed to the functions implementing the tool(s).

A tool function is defined like this:

```go
func NotionGetPage(
        ctx context.Context, 
        toolCtx *NotionGetDocumentContext, 
        input *NotionGetDocumentInput, 
        output types.ToolCallResult) error {
	logger := gomcp.GetLogger(ctx)
	logger.Info("NotionGetPage", types.LogArg{
		"pageId": input.PageId,
	})

	content, err := getPageContent(ctx, toolCtx.NotionClient, input.PageId)
	if err != nil {
		return err
	}
	output.AddTextContent(strings.Join(content, "\n"))

	return nil
}
```
The first parameter is the context, it is mandatory as it is used to retrieve the logger.

The second parameter is the tool context, it is the instance of the struct created by the `ToolInit` function.

The third parameter is the input, it is an object that contains the input data for the tool call. Those input parameter are provided by the LLM when the tool is called.

The last parameter is the output, it is an interface that allows the function to construct the data it wants to return to the LLM.

Here again, the input type must be tagged appropriately:

```go
type NotionGetDocumentInput struct {
	PageId string `json:"pageId" jsonschema_description:"the ID of the Notion page to retrieve."`
}
```

Those tags will be used to generate the JSON Schema for the input data:
* they will be returned to the LLM during the discovery phase of the MCP protocol
* they will be used to validate the input data when the tool is called

Once those typea and functions are defined, you can bind them in the MCP server by calling the `RegisterTool` function:


```go
func RegisterTools(toolRegistry types.ToolRegistry) error {
	toolProvider, err := toolRegistry.DeclareToolProvider("notion", NotionToolInit)
	if err != nil {
		return err
	}
	err = toolProvider.AddTool("notion_get_page", "Get the markdown content of a notion page", NotionGetPage)
	if err != nil {
		return err
	}
	return nil
}
```

This `RegisterTools` function is called from the `main()` function of your MCP server:

```go
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/llmcontext/gomcp"
	"github.com/llmcontext/mcpnotion/tools"
)

func main() {
	configFile := flag.String("configFile", "", "config file path (required)")
	flag.Parse()

	if *configFile == "" {
		fmt.Println("Config file is required")
		flag.PrintDefaults()
		os.Exit(1)
	}

	mcp, err := gomcp.NewModelContextProtocolServer(*configFile)
	if err != nil {
		fmt.Println("Error creating MCP server:", err)
		os.Exit(1)
	}
	toolRegistry := mcp.GetToolRegistry()

	err = tools.RegisterTools(toolRegistry)
	if err != nil {
		fmt.Println("Error registering tools:", err)
		os.Exit(1)
	}

	transport := mcp.StdioTransport()

	mcp.Start(transport)
}
```

* `gomcp.NewModelContextProtocolServer(*configFile)` creates a new MCP server with the configuration file
* `mcp.GetToolRegistry()` returns the tool registry, it is used to register the tools. The `tools.RegisterTools` is the function we defined earlier and that binds the tools to the MCP server.
* `mcp.StdioTransport()` creates a new transport based on standard input/output streams. That's the transport used to integrate with the Claude desktop application.
* `mcp.Start(transport)` starts the MCP server with the given transport


## integration with Claude desktop application

## Changelog

### [0.1.0](https://github.com/llmcontext/gomcp/tree/v0.1.0)

- Initial release
- Support Tools for the Model Context Protocol

### [0.2.0](https://github.com/llmcontext/gomcp/tree/v0.2.0)

- Change signature of `mcp.Start(serverName, serverVersion, transport)` to `mcp.Start(transport)`, the server name and version are now read from the configuration file
- Add support for prompts stored in a YAML file. File path is read from the configuration file.
