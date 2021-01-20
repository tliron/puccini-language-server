package implementation

import (
	"strings"
	"unicode"

	protocol "github.com/tliron/glsp/protocol_3_16"
	"github.com/tliron/puccini/tosca"
)

func createSymbols(context *tosca.Context, content string, uri protocol.DocumentUri) []protocol.SymbolInformation {
	var symbols []protocol.SymbolInformation

	context.Namespace.Range(func(forType tosca.EntityPtr, entityPtr tosca.EntityPtr) bool {
		context := tosca.GetContext(entityPtr)

		// Filter out names that not in our document
		url := context.URL.String()
		url_ := strings.TrimPrefix(url, "internal:language-server:")
		if url_ == url {
			// Doesn't have the prefix
			return true
		}
		if url_ != string(uri) {
			// Not our URI
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

		kind := protocol.SymbolKindObject

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
				URI: uri,
				Range: protocol.Range{
					Start: position,
					End:   positionEol(content, position),
				},
			},
		})

		return true
	})

	return symbols
}
