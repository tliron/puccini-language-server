package implementation

import (
	protocol "github.com/tliron/glsp/protocol_3_16"
)

var Handler protocol.Handler

func init() {
	// General Messages
	Handler.Initialize = Initialize

	// Workspace
	Handler.WorkspaceDidRenameFiles = WorkspaceDidRenameFiles

	// Text Document Synchronization
	Handler.TextDocumentDidOpen = TextDocumentDidOpen
	Handler.TextDocumentDidChange = TextDocumentDidChange
	Handler.TextDocumentDidSave = TextDocumentDidSave
	Handler.TextDocumentDidClose = TextDocumentDidClose

	// Language Features
	Handler.TextDocumentDocumentSymbol = TextDocumentDocumentSymbol
}
