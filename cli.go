package main

import (
	"fmt"
	"github.com/gasiordev/go-cli"
	"os"
)

func getCLIStartHandler(n *NTree) func(*cli.CLI) int {
	fn := func(c *cli.CLI) int {
		n.Init(c.Flag("config"))
		return n.Start()
	}

	return fn
}

func getCLISetCwdHandler(n *NTree) func(*cli.CLI) int {
	fn := func(c *cli.CLI) int {
		n.Init(c.Flag("config"))
		return n.SendCwd(c.Flag("cwd"))
	}
	return fn
}

func getCLIVersionHandler(n *NTree) func(*cli.CLI) int {
	fn := func(c *cli.CLI) int {
		fmt.Fprintf(os.Stdout, VERSION+"\n")
		return 0
	}
	return fn
}

func NewNTreeCLI(n *NTree) *cli.CLI {
	nTreeCLI := cli.NewCLI("Ntree", "Project tree widget", "Mikolaj Gasior")

	cmdStart := nTreeCLI.AddCmd("start", "Starts agent", getCLIStartHandler(n))
	cmdStart.AddFlag("config", "Config file", cli.CLIFlagTypePathFile|cli.CLIFlagMustExist|cli.CLIFlagRequired)

	cmdSetCwd := nTreeCLI.AddCmd("set-cwd", "Set current working directory", getCLISetCwdHandler(n))
	cmdSetCwd.AddFlag("config", "Config file", cli.CLIFlagTypePathFile|cli.CLIFlagMustExist|cli.CLIFlagRequired)
	cmdSetCwd.AddFlag("cwd", "Directory", cli.CLIFlagTypePathFile|cli.CLIFlagRequired)

	_ = nTreeCLI.AddCmd("version", "Prints version", getCLIVersionHandler(n))

	if len(os.Args) == 2 && (os.Args[1] == "-v" || os.Args[1] == "--version") {
		os.Args = []string{"ntree", "version"}
	}
	return nTreeCLI
}
