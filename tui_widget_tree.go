package main

import (
	"github.com/gasiordev/go-tui"
	"io/ioutil"
    "strings"
)

type TUIWidgetTree struct {
	cwd string
}

func (w *TUIWidgetTree) GetCwd() string {
	return w.cwd
}

func (w *TUIWidgetTree) SetCwd(cwd string) {
	w.cwd = cwd
}

// InitPane sets pane minimal width and height that's necessary for the pane
// to work.
func (w *TUIWidgetTree) InitPane(p *tui.TUIPane) {
	p.SetMinWidth(5)
	p.SetMinHeight(3)
}

// Run is main function which just prints out the current time.
func (w *TUIWidgetTree) Run(p *tui.TUIPane) int {
	fileInfo, err := ioutil.ReadDir(w.cwd)
	if err != nil {
		return 0
	}
    for i:=0; i<p.GetHeight()-p.GetStyle().H(); i++ {
        p.Write(0, i, strings.Repeat(" ", p.GetWidth()-p.GetStyle().V()), false)
    }
	i := 0
	for _, file := range fileInfo {
		if i < p.GetHeight()-p.GetStyle().H() {
            n := file.Name()
            if len(n) > p.GetWidth()-p.GetStyle().V() {
                n = n[0:p.GetWidth()-p.GetStyle().V()]
            }
			p.Write(0, i, n, false)
		}
		i++
	}
	return 1
}

// NewTUIWidgetTree returns instance of TUIWidgetTree struct
func NewTUIWidgetTree() *TUIWidgetTree {
	w := &TUIWidgetTree{}
	return w
}
