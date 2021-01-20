package implementation

import (
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

// protocol.WorkspaceDidRenameFilesFunc signature
func WorkspaceDidRenameFiles(context *glsp.Context, params *protocol.RenameFilesParams) error {
	for _, file := range params.Files {
		renameDocument(file.OldURI, file.NewURI)
	}
	return nil
}
