package tosca

import (
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func getWordAt(content string, position protocol.Position) string {
	return ""
}

// protocol.TextDocumentCompletionFunc signature
// Returns: []CompletionItem | CompletionList | nil
func TextDocumentCompletion(context *glsp.Context, params *protocol.CompletionParams) (interface{}, error) {
	var completionItems []protocol.CompletionItem
	if documentState := getDocumentState(params.TextDocument.URI); (documentState != nil) && (documentState.Problems != nil) {
		for _, problem := range documentState.Problems.Slice() {
			if int(params.Position.Line+1) == problem.Row {
				completionItems = append(completionItems, protocol.CompletionItem{
					Label: "server",
				})
				log.Infof("############ completion!")
			}
		}
	}
	return completionItems, nil
}

// protocol.TextDocumentDocumentSymbolFunc signature
// Returns: []DocumentSymbol | []SymbolInformation | nil
func TextDocumentDocumentSymbol(context *glsp.Context, params *protocol.DocumentSymbolParams) (interface{}, error) {
	documentState := validateDocumentState(params.TextDocument.URI, context.Notify)
	return documentState.Symbols, nil
}
