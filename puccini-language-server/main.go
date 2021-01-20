package main

import (
	_ "net/http/pprof"

	"github.com/tebeka/atexit"
	"github.com/tliron/kutil/util"
)

func main() {
	err := command.Execute()
	util.FailOnError(err)
	atexit.Exit(0)
}
