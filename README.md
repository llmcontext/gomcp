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

The `prompts` section is used to define the path to the YAML file containing the prompts to expose to the LLM. See below the definition of the prompts.

The `tools` section is used to define the tools that will be exposed to the LLM. This is an array of tool providers, each provider is an object with a `name` and a `description` field. The `configuration` field is an object that contains the configuration for the tool provider.

In our case, we have a single tool provider called `notion` that has a single tool to retrieve the content of a Notion page.

The configuration for the `notion` tool provider is the Notion token.






## Changelog

### [0.1.0](https://github.com/llmcontext/gomcp/tree/v0.1.0)

- Initial release
- Support Tools for the Model Context Protocol

### [0.2.0](https://github.com/llmcontext/gomcp/tree/v0.2.0)

- Change signature of `mcp.Start(serverName, serverVersion, transport)` to `mcp.Start(transport)`, the server name and version are now read from the configuration file
- Add support for prompts stored in a YAML file. File path is read from the configuration file.
