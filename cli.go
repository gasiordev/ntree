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

func getCLISetWorkDirHandler(n *NTree) func(*cli.CLI) int {
	fn := func(c *cli.CLI) int {
		f := getConfigFilePath(c)
		n.Init(f)
		return n.SendWorkDir(c.Flag("workdir"))
	}
	return fn
}

func getCLIFilterHandler(n *NTree) func(*cli.CLI) int {
	fn := func(c *cli.CLI) int {
		f := getConfigFilePath(c)
		n.Init(f)
		return n.SendFilter(c.Flag("filter"))
	}
	return fn
}

func getCLIHighlightHandler(n *NTree) func(*cli.CLI) int {
	fn := func(c *cli.CLI) int {
		f := getConfigFilePath(c)
		n.Init(f)
		return n.SendHighlight(c.Flag("highlight"))
	}
	return fn
}

func getCLIToggleDirsHandler(n *NTree) func(*cli.CLI) int {
	fn := func(c *cli.CLI) int {
		f := getConfigFilePath(c)
		n.Init(f)
		return n.SendToggleDirs()
	}
	return fn
}

func getCLIToggleFilesHandler(n *NTree) func(*cli.CLI) int {
	fn := func(c *cli.CLI) int {
		f := getConfigFilePath(c)
		n.Init(f)
		return n.SendToggleFiles()
	}
	return fn
}

func getCLIToggleHiddenHandler(n *NTree) func(*cli.CLI) int {
	fn := func(c *cli.CLI) int {
		f := getConfigFilePath(c)
		n.Init(f)
		return n.SendToggleHidden()
	}
	return fn
}

func getCLIToggleFreezeHandler(n *NTree) func(*cli.CLI) int {
	fn := func(c *cli.CLI) int {
		f := getConfigFilePath(c)
		n.Init(f)
		return n.SendToggleFreeze()
	}
	return fn
}

func getCLIResetHighlightHandler(n *NTree) func(*cli.CLI) int {
	fn := func(c *cli.CLI) int {
		f := getConfigFilePath(c)
		n.Init(f)
		return n.SendResetHighlight()
	}
	return fn
}

func getCLIResetFilterHandler(n *NTree) func(*cli.CLI) int {
	fn := func(c *cli.CLI) int {
		f := getConfigFilePath(c)
		n.Init(f)
		return n.SendResetFilter()
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

	cmdSetWorkDir := nTreeCLI.AddCmd("set-workdir", "Set current working directory", getCLISetWorkDirHandler(n))
	cmdSetWorkDir.AddFlag("config", "c", "file", "Config file", cli.TypePathFile|cli.MustExist)
	cmdSetWorkDir.AddFlag("workdir", "w", "dir", "Directory", cli.TypePathFile|cli.Required)

	cmdToggleDirs := nTreeCLI.AddCmd("toggle-dirs", "Toggle directory visibility", getCLIToggleDirsHandler(n))
	cmdToggleDirs.AddFlag("config", "c", "file", "Config file", cli.TypePathFile|cli.MustExist)

	cmdToggleFiles := nTreeCLI.AddCmd("toggle-files", "Toggle file visibility", getCLIToggleFilesHandler(n))
	cmdToggleFiles.AddFlag("config", "c", "file", "Config file", cli.TypePathFile|cli.MustExist)

	cmdToggleFreeze := nTreeCLI.AddCmd("toggle-freeze", "Toggle freeze", getCLIToggleFreezeHandler(n))
	cmdToggleFreeze.AddFlag("config", "c", "file", "Config file", cli.TypePathFile|cli.MustExist)

	cmdToggleHidden := nTreeCLI.AddCmd("toggle-hidden", "Toggle hides files visibility", getCLIToggleHiddenHandler(n))
	cmdToggleHidden.AddFlag("config", "c", "file", "Config file", cli.TypePathFile|cli.MustExist)

	cmdFilter := nTreeCLI.AddCmd("filter", "Filter tree", getCLIFilterHandler(n))
	cmdFilter.AddFlag("config", "c", "file", "Config file", cli.TypePathFile|cli.MustExist)
	cmdFilter.AddFlag("filter", "f", "filter", "Filter", cli.TypeString|cli.Required)

	cmdHighlight := nTreeCLI.AddCmd("highlight", "Highlight string", getCLIHighlightHandler(n))
	cmdHighlight.AddFlag("config", "c", "file", "Config file", cli.TypePathFile|cli.MustExist)
	cmdHighlight.AddFlag("highlight", "h", "highlight", "Highlight", cli.TypeString|cli.Required)

	cmdResetFilter := nTreeCLI.AddCmd("reset-filter", "Reset filter", getCLIResetFilterHandler(n))
	cmdResetFilter.AddFlag("config", "c", "file", "Config file", cli.TypePathFile|cli.MustExist)

	cmdResetHighlight := nTreeCLI.AddCmd("reset-highlight", "Reset highlight", getCLIResetHighlightHandler(n))
	cmdResetHighlight.AddFlag("config", "c", "file", "Config file", cli.TypePathFile|cli.MustExist)

	_ = nTreeCLI.AddCmd("version", "Prints version", getCLIVersionHandler(n))

	if len(os.Args) == 2 && (os.Args[1] == "-v" || os.Args[1] == "--version") {
		os.Args = []string{"ntree", "version"}
	}
	return nTreeCLI
}
