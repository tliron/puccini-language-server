package tosca

import (
	"unicode"

	protocol "github.com/tliron/glsp/protocol_3_16"
	"github.com/tliron/puccini/tosca"
)

func createSymbols(context *tosca.Context, content string, documentUri protocol.DocumentUri) []protocol.SymbolInformation {
	var symbols []protocol.SymbolInformation

	context.Namespace.Range(func(entityPtr tosca.EntityPtr) bool {
		context := tosca.GetContext(entityPtr)

		// Filter out names that not in our document
		if documentUri_, ok := urlToDocumentUri(context.URL); ok {
			if documentUri_ != documentUri {
				// Not our URI
				return true
			}
		} else {
			return true
		}

		row, column := context.GetLocation()
		row--
		column--
		if row < 0 {
			row = 0
		}
		if column < 0 {
			column = 0
		}
		position := protocol.Position{
			Line:      protocol.UInteger(row),
			Character: protocol.UInteger(column),
		}
		if position.Line < 0 {
			position.Line = 0
		}
		if position.Character < 0 {
			position.Character = 0
		}

		var kind protocol.SymbolKind
		if clientCapabilities.SupportsSymbolKind(protocol.SymbolKindObject) {
			kind = protocol.SymbolKindObject
		} else {
			kind = protocol.SymbolKindVariable
		}

		// TODO: properly detect types
		if unicode.IsUpper(rune(context.Name[0])) {
			kind = protocol.SymbolKindClass
		}

		/*container := ""
		if parent, ok := tosca.GetParent(entityPtr); ok {
			container = tosca.GetContext(parent).Name
			log.Infof("############## %s", container)
		}*/

		symbols = append(symbols, protocol.SymbolInformation{
			Name: context.Name,
			Kind: kind,
			Location: protocol.Location{
				URI: documentUri,
				Range: protocol.Range{
					Start: position,
					End:   position.EndOfLineIn(content),
				},
			},
		})

		return true
	})

	return symbols
}
