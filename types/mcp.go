package types

type ModelContextProtocolServer interface {
	StdioTransport() Transport
	Start(transport Transport) error
}
