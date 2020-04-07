package main

import (
	"fmt"
	"github.com/gasiordev/go-cli"
	"log"
	"os"
)

func getConfigFilePath(c *cli.CLI) string {
	h, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Cannot get home dir")
	}
	f := h + "/.ntree.json"
	if c.Flag("config") != "" {
		f = c.Flag("config")
	}
	return f
}

func getCLIStartHandler(n *NTree) func(*cli.CLI) int {
	fn := func(c *cli.CLI) int {
		f := getConfigFilePath(c)
		n.Init(f)
		return n.Start(c.Flag("workdir"))
	}

	return fn
}

func getCLISendHandler(n *NTree) func(*cli.CLI) int {
	fn := func(c *cli.CLI) int {
		f := getConfigFilePath(c)
		n.Init(f)
		return n.SendCmd(c.Arg("cmd"), c.Arg("val"))
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
	cmdStart.AddFlag("config", "c", "file", "Config file", cli.TypePathFile|cli.MustExist)
	cmdStart.AddFlag("workdir", "w", "dir", "Directory", cli.TypePathFile|cli.Required)

	cmdMsg := nTreeCLI.AddCmd("send", "Sends command to already running agent", getCLISendHandler(n))
	cmdMsg.AddFlag("config", "c", "file", "Config file", cli.TypePathFile|cli.MustExist)
	cmdMsg.AddArg("cmd", "COMMAND", "Command to be send to agent", cli.TypeString|cli.Required)
	cmdMsg.AddArg("val", "VALUE", "Optional value for the command", cli.TypeString)

	_ = nTreeCLI.AddCmd("version", "Prints version", getCLIVersionHandler(n))

	if len(os.Args) == 2 && (os.Args[1] == "-v" || os.Args[1] == "--version") {
		os.Args = []string{"ntree", "version"}
	}
	return nTreeCLI
}
