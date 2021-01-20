package implementation

import (
	"sync"

	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
	"github.com/tliron/kutil/problems"
	urlpkg "github.com/tliron/kutil/url"
	cloutpkg "github.com/tliron/puccini/clout"
	"github.com/tliron/puccini/tosca/compiler"
	"github.com/tliron/puccini/tosca/normal"
	"github.com/tliron/puccini/tosca/parser"
)

type DocumentState struct {
	Symbols     []protocol.SymbolInformation
	Diagnostics []protocol.Diagnostic
}

var documentStates sync.Map // protocol.DocumentUri to DocumentState

func validateDocumentState(uri protocol.DocumentUri, notify glsp.NotifyFunc) *DocumentState {
	documentState, created := _getOrCreateDocumentState(uri)

	if created {
		go notify(protocol.ServerTextDocumentPublishDiagnostics, &protocol.PublishDiagnosticsParams{
			URI:         uri,
			Diagnostics: documentState.Diagnostics,
		})
	}

	return documentState
}

func deleteDocumentState(uri protocol.DocumentUri) {
	documentStates.Delete(uri)
}

func _getOrCreateDocumentState(uri protocol.DocumentUri) (*DocumentState, bool) {
	if documentState, ok := documentStates.Load(uri); ok {
		return documentState.(*DocumentState), false
	} else {
		documentState := _createDocumentState(uri)
		if existing, loaded := documentStates.LoadOrStore(uri, documentState); loaded {
			return existing.(*DocumentState), false
		} else {
			return documentState, true
		}
	}
}

func _createDocumentState(uri protocol.DocumentUri) *DocumentState {
	var documentState DocumentState

	var err error
	var url urlpkg.URL
	var content string
	var context *parser.Context
	var serviceTemplate *normal.ServiceTemplate
	var clout *cloutpkg.Clout
	var problems *problems.Problems

	urlContext := urlpkg.NewContext()
	defer urlContext.Release()

	path := uriToInternalPath(uri)
	content, _ = getDocument(uri)

	if url, err = urlpkg.NewValidInternalURL(path, urlContext); err != nil {
		log.Errorf("%s", err.Error())
		documentState.Diagnostics = createDiagnostics(problems, content)
		return &documentState
	}

	if context, serviceTemplate, problems, err = parser.Parse2(url, nil, nil); err != nil {
		log.Errorf("%s", err.Error())
		documentState.Diagnostics = createDiagnostics(problems, content)
		documentState.Symbols = createSymbols(context.Root.GetContext(), content, uri)
		return &documentState
	}

	log.Debugf("%T", serviceTemplate)

	if clout, err = compiler.Compile(serviceTemplate, true); (err != nil) || !problems.Empty() {
		if err != nil {
			log.Errorf("%s", err.Error())
		}
		documentState.Diagnostics = createDiagnostics(problems, content)
		documentState.Symbols = createSymbols(context.Root.GetContext(), content, uri)
		return &documentState
	}

	log.Debugf("%T", clout)

	/*compiler.Resolve(clout, problems, urlContext, true, "yaml", true, false, false)
	if !problems.Empty() {
		return nil, nil
	}*/

	documentState.Diagnostics = createDiagnostics(problems, content)
	documentState.Symbols = createSymbols(context.Root.GetContext(), content, uri)

	return &documentState
}
