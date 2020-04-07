package main

import (
	"github.com/gasiordev/go-cli"
	"github.com/gasiordev/go-tui"
	"io/ioutil"
	"log"
	"net"
	"os"
	"regexp"
)

type NTree struct {
	config     Config
	cli        *cli.CLI
	listener   net.Listener
	tui        *tui.TUI
	workDir    string
	filter     string
	highlight  string
	hideDirs   bool
	hideFiles  bool
	showHidden bool
	freeze     bool
}

func NewNTree() *NTree {
	n := &NTree{}
	return n
}

func (n *NTree) GetConfig() *Config {
	return &(n.config)
}

func (n *NTree) GetCLI() *cli.CLI {
	return n.cli
}

func (n *NTree) GetListener() net.Listener {
	return n.listener
}

func (n *NTree) GetWorkDir() string {
	return n.workDir
}

func (n *NTree) GetFilter() string {
	return n.filter
}

func (n *NTree) GetHighlight() string {
	return n.highlight
}

func (n *NTree) GetHideDirs() bool {
	return n.hideDirs
}

func (n *NTree) GetHideFiles() bool {
	return n.hideFiles
}

func (n *NTree) GetShowHidden() bool {
	return n.showHidden
}

func (n *NTree) GetFreeze() bool {
	return n.freeze
}

func (n *NTree) toggleHideDirs() {
	if n.hideDirs {
		n.hideDirs = false
	} else {
		n.hideDirs = true
	}
}

func (n *NTree) toggleHideFiles() {
	if n.hideFiles {
		n.hideFiles = false
	} else {
		n.hideFiles = true
	}
}

func (n *NTree) toggleFreeze() {
	if n.freeze {
		n.freeze = false
	} else {
		n.freeze = true
	}
}

func (n *NTree) toggleShowHidden() {
	if n.showHidden {
		n.showHidden = false
	} else {
		n.showHidden = true
	}
}

func (n *NTree) Init(cfgFile string) {
	cfgJSON, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		log.Fatal("Error reading config file")
	}

	var cfg Config
	cfg.SetFromJSON(cfgJSON)
	n.config = cfg
}

func (n *NTree) goReadData(c net.Conn) {
	for {
		buf := make([]byte, 512)
		nr, err := c.Read(buf)
		if err != nil {
			return
		}

		data := buf[0:nr]

		m, err := regexp.Match(`^WORKDIR .+`, data)
		if m && !n.freeze {
			n.workDir = string(buf[8:nr])
			continue
		}
		m, err = regexp.Match(`^FILTER .+`, data)
		if m && !n.freeze {
			n.filter = string(buf[7:nr])
			continue
		}
		m, err = regexp.Match(`^HIGHLIGHT .+`, data)
		if m && !n.freeze {
			n.highlight = string(buf[10:nr])
			continue
		}
		if string(data) == "DIRS" && !n.freeze {
			n.toggleHideDirs()
			continue
		}
		if string(data) == "FILES" && !n.freeze {
			n.toggleHideFiles()
			continue
		}
		if string(data) == "HIDDEN" && !n.freeze {
			n.toggleShowHidden()
			continue
		}
		if string(data) == "RESET-FILTER" && !n.freeze {
			n.filter = ""
			continue
		}
		if string(data) == "RESET-HIGHLIGHT" && !n.freeze {
			n.highlight = ""
		}
		if string(data) == "FREEZE" {
			n.toggleFreeze()
		}
	}
}

func (n *NTree) goAccept() {
	for {
		fd, err := n.listener.Accept()
		if err != nil {
			log.Fatal("Accept error: ", err)
		}

		go n.goReadData(fd)
	}
}

func (n *NTree) Start(workDir string) int {
	n.workDir = workDir

	l, err := net.Listen("unix", n.config.GetUnixSocket())
	if err != nil {
		log.Fatal("listen error: ", err)
	}
	n.listener = l

	go n.goAccept()

	t := NewNTreeTUI(n)
	n.tui = t
	return t.Run(os.Stdout, os.Stderr)
}

func (n *NTree) SendWorkDir(workDir string) int {
	return n.SendCmd("WORKDIR", workDir)
}

func (n *NTree) SendFilter(filter string) int {
	return n.SendCmd("FILTER", filter)
}

func (n *NTree) SendHighlight(highlight string) int {
	return n.SendCmd("HIGHLIGHT", highlight)
}

func (n *NTree) SendToggleDirs() int {
	return n.SendCmd("DIRS", "")
}

func (n *NTree) SendToggleFiles() int {
	return n.SendCmd("FILES", "")
}

func (n *NTree) SendToggleFreeze() int {
	return n.SendCmd("FREEZE", "")
}

func (n *NTree) SendToggleHidden() int {
	return n.SendCmd("HIDDEN", "")
}

func (n *NTree) SendResetFilter() int {
	return n.SendCmd("RESET-FILTER", "")
}

func (n *NTree) SendResetHighlight() int {
	return n.SendCmd("RESET-HIGHLIGHT", "")
}

func (n *NTree) SendCmd(cmd string, val string) int {
	c, err := net.Dial("unix", n.config.GetUnixSocket())
	if err != nil {
		panic(err)
	}
	defer c.Close()

	m := []byte(cmd)
	if val != "" {
		m = []byte(cmd + " " + val)
	}

	_, err = c.Write(m)
	if err != nil {
		log.Fatal("Write error: ", err)
	}
	return 1
}

func (n *NTree) Run() {
	nCLI := NewNTreeCLI(n)
	n.cli = nCLI
	os.Exit(nCLI.Run(os.Stdout, os.Stderr))
}
