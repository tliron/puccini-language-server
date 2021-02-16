package tosca

import (
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
	"github.com/tliron/kutil/format"
)

type Configuration struct {
	Debug ConfigurationDebug `json:"debug"`
}

type ConfigurationDebug struct {
	Verbosity int `json:"verbosity"`
}

func Configure(context *glsp.Context, scope *protocol.DocumentUri) {
	var results []Configuration
	var tosca = "tosca"
	context.Call(protocol.ServerWorkspaceConfiguration, &protocol.ConfigurationParams{
		Items: []protocol.ConfigurationItem{
			{
				ScopeURI: scope,
				Section:  &tosca,
			},
		},
	}, &results)
	s, _ := format.EncodeJSON(results, "")
	log.Errorf("****** %s", s)
}
