package mux

type JsonRpcResponseProxyRegisterResult struct {
	SessionId  string `json:"sessionId"`
	ProxyId    string `json:"proxyId"`
	Persistent bool   `json:"persistent"`
	Denied     bool   `json:"denied"`
}
