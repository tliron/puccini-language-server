package tosca

import (
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

// protocol.TextDocumentCompletionFunc signature
// Returns: []CompletionItem | CompletionList | nil
func TextDocumentCompletion(context *glsp.Context, params *protocol.CompletionParams) (interface{}, error) {
	return nil, nil
}

// protocol.TextDocumentDocumentSymbolFunc signature
// Returns: []DocumentSymbol | []SymbolInformation | nil
func TextDocumentDocumentSymbol(context *glsp.Context, params *protocol.DocumentSymbolParams) (interface{}, error) {
	documentState := validateDocumentState(params.TextDocument.URI, context.Notify)
	return documentState.Symbols, nil
}
