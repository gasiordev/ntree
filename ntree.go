package main

import (
	"github.com/gasiordev/go-cli"
	"github.com/gasiordev/go-tui"
	"io/ioutil"
	"log"
	"net"
	"os"
	"regexp"
	"strconv"
)

type NTree struct {
	config     Config
	cli        *cli.CLI
	listener   net.Listener
	tui        *tui.TUI
	rootDir    string
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

func (n *NTree) GetRootDir() string {
	return n.rootDir
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

func (n *NTree) ToggleHideDirs() {
	if n.hideDirs {
		n.hideDirs = false
	} else {
		n.hideDirs = true
	}
}

func (n *NTree) ToggleHideFiles() {
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

func (n *NTree) ToggleShowHidden() {
	if n.showHidden {
		n.showHidden = false
	} else {
		n.showHidden = true
	}
}

func (n *NTree) ResetFilter() {
	n.filter = ""
}

func (n *NTree) ResetHighlight() {
	n.highlight = ""
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

		m, err := regexp.Match(`^ROOTDIR .+`, data)
		if m && !n.freeze {
			n.rootDir = string(buf[8:nr])
			continue
		}
		m, err = regexp.Match(`^WORKDIR .+`, data)
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
			n.ToggleHideDirs()
			continue
		}
		if string(data) == "FILES" && !n.freeze {
			n.ToggleHideFiles()
			continue
		}
		if string(data) == "HIDDEN" && !n.freeze {
			n.ToggleShowHidden()
			continue
		}
		if string(data) == "RESET-FILTER" && !n.freeze {
			n.ResetFilter()
			continue
		}
		if string(data) == "RESET-HIGHLIGHT" && !n.freeze {
			n.ResetHighlight()
			continue
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

func (n *NTree) Start(rootDir string, workDir string) int {
	n.rootDir = rootDir
	n.workDir = workDir

	i, err := os.Stat(n.config.GetUnixSocket())
	if err == nil && (i.Mode()&os.ModeSocket != 0) {
		err = os.Remove(n.config.GetUnixSocket())
		if err != nil {
			log.Fatal("socket file rm error: ", err)
		}
	}

	l, err := net.Listen("unix", n.config.GetUnixSocket())
	if err != nil {
		log.Fatal("listen error: ", err)
	}
	n.listener = l

	go n.goAccept()

	t := NewNTreeTUI(n)

	ls, err := strconv.Atoi(n.config.GetLoopSleep())
	if err != nil {
		log.Fatal("loop_sleep has invalid value")
	}
	if ls < 0 {
		ls = 1000
	}

	t.SetLoopSleep(ls)
	n.tui = t
	return t.Run(os.Stdout, os.Stderr)
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
