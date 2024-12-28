package types

type ModelContextProtocol interface {
	StdioTransport() Transport
	Start(transport Transport) error
}
