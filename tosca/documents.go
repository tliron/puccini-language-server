package tosca

import (
	"strings"

	protocol "github.com/tliron/glsp/protocol_3_16"
	urlpkg "github.com/tliron/kutil/url"
)

const INTERNAL_PATH_PREFIX = "language-server:"

func documentUriToInternalPath(uri protocol.DocumentUri) string {
	return INTERNAL_PATH_PREFIX + string(uri)
}

func urlToDocumentUri(url urlpkg.URL) (protocol.DocumentUri, bool) {
	url_ := url.String()
	if documentUri := strings.TrimPrefix(url_, "internal:"+INTERNAL_PATH_PREFIX); documentUri != url_ {
		return protocol.DocumentUri(documentUri), true
	} else {
		return protocol.DocumentUri(""), false
	}
}

func setDocument(documentUri protocol.DocumentUri, content string) {
	urlpkg.UpdateInternalURL(documentUriToInternalPath(documentUri), content)
	documentStates.Delete(documentUri)
}

func getDocument(documentUri protocol.DocumentUri) (string, bool) {
	if url, err := urlpkg.NewValidInternalURL(documentUriToInternalPath(documentUri), nil); err == nil {
		return url.Content, true
	} else {
		return "", false
	}
}

func deleteDocument(documentUri protocol.DocumentUri) {
	urlpkg.DeregisterInternalURL(documentUriToInternalPath(documentUri))
	documentStates.Delete(documentUri)
}

func renameDocument(oldDocumentUri protocol.DocumentUri, newDocumentUri protocol.DocumentUri) {
	if content, ok := getDocument(oldDocumentUri); ok {
		setDocument(newDocumentUri, content)
	}
	deleteDocument(oldDocumentUri)
}
