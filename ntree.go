package main

import (
	"github.com/gasiordev/go-cli"
	"github.com/gasiordev/go-tui"
	"io/ioutil"
	"log"
	"net"
	"os"
)

type NTree struct {
	config   Config
	cli      *cli.CLI
	listener net.Listener
	tui      *tui.TUI
	cwd      string
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

func (n *NTree) GetCwd() string {
	return n.cwd
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
		n.cwd = string(data)
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

func (n *NTree) Start(cwd string) int {
	n.cwd = cwd

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

func (n *NTree) SendCwd(cwd string) int {
	c, err := net.Dial("unix", n.config.GetUnixSocket())
	if err != nil {
		panic(err)
	}
	defer c.Close()

	_, err = c.Write([]byte(cwd))
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
