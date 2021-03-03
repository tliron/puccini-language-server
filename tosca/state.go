package tosca

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

var documentStates sync.Map // protocol.DocumentUri to DocumentState

func getDocumentState(documentUri protocol.DocumentUri) *DocumentState {
	if documentState, ok := documentStates.Load(documentUri); ok {
		return documentState.(*DocumentState)
	} else {
		return nil
	}
}

func validateDocumentState(documentUri protocol.DocumentUri, notify glsp.NotifyFunc) *DocumentState {
	documentState, created := _getOrCreateDocumentState(documentUri)

	if created {
		go notify(protocol.ServerTextDocumentPublishDiagnostics, &protocol.PublishDiagnosticsParams{
			URI:         documentUri,
			Diagnostics: documentState.Diagnostics,
		})
	}

	return documentState
}

func deleteDocumentState(documentUri protocol.DocumentUri) {
	documentStates.Delete(documentUri)
}

func resetDocumentStates() {
	documentStates.Range(func(protocolUri interface{}, documentState interface{}) bool {
		documentStates.Delete(protocolUri)
		urlpkg.DeregisterInternalURL(documentUriToInternalPath(protocolUri.(protocol.DocumentUri)))
		return true
	})
}

func _getOrCreateDocumentState(documentUri protocol.DocumentUri) (*DocumentState, bool) {
	if documentState, ok := documentStates.Load(documentUri); ok {
		return documentState.(*DocumentState), false
	} else {
		documentState := NewDocumentState(documentUri)
		if existing, loaded := documentStates.LoadOrStore(documentUri, documentState); loaded {
			return existing.(*DocumentState), false
		} else {
			return documentState, true
		}
	}
}

//
// DocumentState
//

type DocumentState struct {
	Content  string
	Context  *parser.Context
	Problems *problems.Problems

	DocumentURI protocol.DocumentUri
	Symbols     []protocol.SymbolInformation
	Diagnostics []protocol.Diagnostic
}

func NewDocumentState(documentUri protocol.DocumentUri) *DocumentState {
	self := DocumentState{DocumentURI: documentUri}

	var err error
	var url urlpkg.URL
	var serviceTemplate *normal.ServiceTemplate
	var clout *cloutpkg.Clout

	urlContext := urlpkg.NewContext()
	defer urlContext.Release()

	path := documentUriToInternalPath(documentUri)
	self.Content, _ = getDocument(documentUri)

	if url, err = urlpkg.NewValidInternalURL(path, urlContext); err != nil {
		log.Errorf("%s", err.Error())
		self.Fill()
		return &self
	}

	if self.Context, serviceTemplate, self.Problems, err = parser.Parse(url, nil, nil); err != nil {
		log.Errorf("%s", err.Error())
		self.Fill()
		return &self
	}

	log.Debugf("%T", serviceTemplate)

	if clout, err = compiler.Compile(serviceTemplate, true); (err != nil) || !self.Problems.Empty() {
		if err != nil {
			log.Errorf("%s", err.Error())
		}
		self.Fill()
		return &self
	}

	compiler.Resolve(clout, self.Problems, urlContext, true, "yaml", true, false, false)
	if !self.Problems.Empty() {
		if err != nil {
			log.Errorf("%s", err.Error())
		}
		self.Fill()
		return &self
	}

	self.Fill()
	return &self
}

func (self *DocumentState) Fill() {
	self.Diagnostics = createDiagnostics(self.Problems, self.Content)
	if self.Context.Root != nil {
		self.Symbols = createSymbols(self.Context.Root.GetContext(), self.Content, self.DocumentURI)
	}
}
