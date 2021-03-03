package tosca

import (
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
	"github.com/tliron/kutil/version"
)

var clientCapabilities *protocol.ClientCapabilities

// protocol.InitializeFunc signature
// Returns: InitializeResult | InitializeError
func Initialize(context *glsp.Context, params *protocol.InitializeParams) (interface{}, error) {
	clientCapabilities = &params.Capabilities

	if params.Trace != nil {
		protocol.SetTraceValue(*params.Trace)
	}

	serverCapabilities := Handler.CreateServerCapabilities()
	serverCapabilities.TextDocumentSync = protocol.TextDocumentSyncKindIncremental
	serverCapabilities.CompletionProvider = &protocol.CompletionOptions{}

	return &protocol.InitializeResult{
		Capabilities: serverCapabilities,
		ServerInfo: &protocol.InitializeResultServerInfo{
			Name:    "puccini-language-server",
			Version: &version.GitVersion,
		},
	}, nil
}

// protocol.InitializedFunc signature
func Initialized(context *glsp.Context, params *protocol.InitializedParams) error {
	return nil
}

// protocol.ShutdownFunc signature
func Shutdown(context *glsp.Context) error {
	protocol.SetTraceValue(protocol.TraceValueOff)
	resetDocumentStates()
	return nil
}

// protocol.LogTraceFunc signature
func LogTrace(context *glsp.Context, params *protocol.LogTraceParams) error {
	return nil
}

// protocol.SetTraceFunc signature
func SetTrace(context *glsp.Context, params *protocol.SetTraceParams) error {
	protocol.SetTraceValue(params.Value)
	return nil
}
