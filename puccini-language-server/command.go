package main

import (
	"fmt"

	"github.com/spf13/cobra"
	serverpkg "github.com/tliron/glsp/server"
	"github.com/tliron/kutil/logging"
	"github.com/tliron/kutil/util"
	versionpkg "github.com/tliron/kutil/version"
	"github.com/tliron/puccini-language-server/tosca"
)

var logTo string
var verbose int
var version bool

var protocol string
var address string

func init() {
	command.Flags().StringVarP(&logTo, "log", "l", "", "log to file (defaults to stderr)")
	command.Flags().CountVarP(&verbose, "verbose", "v", "add a log verbosity level (can be used twice)")
	command.Flags().BoolVar(&version, "version", false, "print version")

	command.Flags().StringVarP(&protocol, "protocol", "p", "stdio", "protocol (\"stdio\", \"tcp\", \"websocket\", or \"nodejs\"")
	command.Flags().StringVarP(&address, "address", "a", ":4389", "address (for \"tcp\" and \"websocket\"")
}

var command = &cobra.Command{
	Use:   toolName,
	Short: "Start the Puccini TOSCA language server",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if logTo == "" {
			logging.Configure(verbose, nil)
		} else {
			logging.Configure(verbose, &logTo)
		}

		if verbose > 0 {
			// Reduce Puccini logging even in verbose mode
			logging.SetMaxLevel("puccini.*", logging.Warning)
		}

		if version {
			versionpkg.Print()
			util.Exit(0)
			return
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		err := Run(protocol, address)
		util.FailOnError(err)
	},
}

func Run(protocol string, address string) error {
	log.Infof("version %s", versionpkg.GitVersion)

	server := serverpkg.NewServer(&tosca.Handler, toolName, verbose > 0)

	switch protocol {
	case "stdio":
		return server.RunStdio()

	case "tcp":
		return server.RunTCP(address)

	case "websocket":
		return server.RunWebSocket(address)

	case "nodejs":
		return server.RunNodeJs()

	default:
		return fmt.Errorf("unsupported protocol: %s", protocol)
	}
}
