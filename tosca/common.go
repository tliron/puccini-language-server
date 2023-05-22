package tosca

import (
	"github.com/tliron/commonlog"
	"github.com/tliron/puccini/tosca/parser"
)

var log = commonlog.GetLogger("puccini-language-server.tosca")

var parserContext = parser.NewContext()
