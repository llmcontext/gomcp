# gomcp
Unofficial Golang SDK for Anthropic Model Context Protocol

This is still very much a work in progress.

You'll find a presentation of the project in this [blog post](http://pcarion.com/blog/go_model_context_protocol/)


# Reference documentation

You'll find the documentation of the Model Context Protocol [here](https://modelcontextprotocol.io/introduction)


# Changelog

## [0.1.0](https://github.com/llmcontext/gomcp/tree/v0.1.0)

- Initial release
- Support Tools for the Model Context Protocol

## [0.2.0](https://github.com/llmcontext/gomcp/tree/v0.2.0)

- Change signature of `mcp.Start(serverName, serverVersion, transport)` to `mcp.Start(transport)`, the server name and version are now read from the configuration file
- Add support for prompts stored in a YAML file. File path is read from the configuration file.
