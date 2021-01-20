package implementation

import (
	"strings"

	"github.com/op/go-logging"
	protocol "github.com/tliron/glsp/protocol_3_16"
	urlpkg "github.com/tliron/kutil/url"
)

var log = logging.MustGetLogger("implementation")

func findIndex(content string, row int, column int) int {
	index := 0
	for r := 0; r < row; r++ {
		content_ := content[index:]
		if next := strings.Index(content_, "\n"); next != -1 {
			index += next + 1
		} else {
			return 0
		}
	}
	return index + column
}

func positionToIndex(content string, position protocol.Position) int {
	return findIndex(content, int(position.Line), int(position.Character))
}

func rangeToIndex(content string, range_ *protocol.Range) (int, int) {
	return positionToIndex(content, range_.Start), positionToIndex(content, range_.End)
}

func positionEol(content string, position protocol.Position) protocol.Position {
	index := positionToIndex(content, position)
	content_ := content[index:]
	if eol := strings.Index(content_, "\n"); eol != -1 {
		return protocol.Position{
			Line:      protocol.UInteger(position.Line),
			Character: protocol.UInteger(position.Character) + protocol.UInteger(eol),
		}
	} else {
		return position
	}
}

func uriToInternalPath(uri protocol.DocumentUri) string {
	return "language-server:" + string(uri)
}

func setDocument(uri protocol.DocumentUri, content string) {
	urlpkg.UpdateInternalURL(uriToInternalPath(uri), content)
	documentStates.Delete(uri)
}

func getDocument(uri protocol.DocumentUri) (string, bool) {
	if url, err := urlpkg.NewValidInternalURL(uriToInternalPath(uri), nil); err == nil {
		return url.Content, true
	} else {
		return "", false
	}
}

func deleteDocument(uri protocol.DocumentUri) {
	urlpkg.DeregisterInternalURL(uriToInternalPath(uri))
	documentStates.Delete(uri)
}

func renameDocument(oldUri protocol.DocumentUri, newUri protocol.DocumentUri) {
	if content, ok := getDocument(oldUri); ok {
		setDocument(newUri, content)
	}
	deleteDocument(oldUri)
}
