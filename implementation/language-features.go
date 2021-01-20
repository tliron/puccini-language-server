package implementation

import (
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

// protocol.TextDocumentDocumentSymbolFunc signature
func TextDocumentDocumentSymbol(context *glsp.Context, params *protocol.DocumentSymbolParams) (interface{}, error) {
	documentState := validateDocumentState(params.TextDocument.URI, context.Notify)
	return documentState.Symbols, nil
}
