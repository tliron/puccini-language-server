module github.com/tliron/puccini-language-server

go 1.16

// replace github.com/tliron/kutil => /Depot/Projects/RedHat/kutil

replace github.com/tliron/puccini => /Depot/Projects/RedHat/puccini

replace github.com/tliron/glsp => /Depot/Projects/RedHat/glsp

require (
	github.com/spf13/cobra v1.1.3
	github.com/tebeka/atexit v0.3.0
	github.com/tliron/glsp v0.0.0-00010101000000-000000000000
	github.com/tliron/kutil v0.1.21
	github.com/tliron/puccini v0.0.0-00010101000000-000000000000
)
