package main

import (
	"github.com/tliron/kutil/util"

	_ "net/http/pprof"

	_ "github.com/tliron/kutil/logging/simple"
)

func main() {
	err := command.Execute()
	util.FailOnError(err)
	util.Exit(0)
}
