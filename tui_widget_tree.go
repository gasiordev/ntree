package main

import (
	"github.com/gasiordev/go-tui"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type TUIWidgetTree struct {
	workDir    string
	hideFiles  bool
	hideDirs   bool
	showHidden bool
	filter     string
	highlight  string
}

func (w *TUIWidgetTree) GetWorkDir() string {
	return w.workDir
}

func (w *TUIWidgetTree) SetWorkDir(workDir string) {
	w.workDir = workDir
}

func (w *TUIWidgetTree) SetHideFiles(b bool) {
	w.hideFiles = b
}

func (w *TUIWidgetTree) SetHideDirs(b bool) {
	w.hideDirs = b
}

func (w *TUIWidgetTree) SetShowHidden(b bool) {
	w.showHidden = b
}

func (w *TUIWidgetTree) SetFilter(s string) {
	w.filter = s
}

func (w *TUIWidgetTree) SetHighlight(s string) {
	w.highlight = s
}

// InitPane sets pane minimal width and height that's necessary for the pane
// to work.
func (w *TUIWidgetTree) InitPane(p *tui.TUIPane) {
	p.SetMinWidth(5)
	p.SetMinHeight(3)
}

// Run is main function which just prints out the current time.
func (w *TUIWidgetTree) Run(p *tui.TUIPane) int {
	fileInfo, err := ioutil.ReadDir(w.workDir)
	if err != nil {
		return 0
	}
	for i := 0; i < p.GetHeight()-p.GetStyle().H(); i++ {
		p.Write(0, i, strings.Repeat(" ", p.GetWidth()-p.GetStyle().V()), false)
	}
	i := 0
	cntDisplayed := 0
	cntHidden := 1

	for _, file := range fileInfo {
		n := file.Name()
		origN := file.Name()
		if len(n) > p.GetWidth()-p.GetStyle().V() {
			n = n[0 : p.GetWidth()-p.GetStyle().V()]
		}
		highlight := "m"
		if w.highlight != "" {
			m, err := regexp.MatchString(w.highlight, origN)
			if m && err == nil {
				highlight = ";1m"
			}
		}
		if file.IsDir() {
			n = "\u001b[33" + highlight + n + "\u001b[0m"
		} else if file.Mode()&os.ModeSymlink != 0 {
			n = "\u001b[36" + highlight + n + "\u001b[0m"
		} else if file.Mode().IsRegular() {
			n = "\u001b[32" + highlight + n + "\u001b[0m"
		} else {
			n = "\u001b[35" + highlight + n + "\u001b[0m"
		}

		through := true
		if file.IsDir() && w.hideDirs {
			through = false
		}
		if file.Mode().IsRegular() && w.hideFiles {
			through = false
		}
		if file.Name()[0] == '.' && !w.showHidden {
			through = false
		}
		if w.filter != "" {
			m, err := regexp.MatchString(w.filter, origN)
			if !m || err != nil {
				through = false
			}
		}

		if through {
			if cntDisplayed < p.GetHeight()-p.GetStyle().H() {
				p.Write(0, cntDisplayed, n, false)
				cntDisplayed++
			} else {
				cntHidden++
			}
		}
		i++
	}
	if cntHidden > 1 {
		p.Write(0, p.GetHeight()-p.GetStyle().H()-1, strings.Repeat(" ", p.GetWidth()-p.GetStyle().V()), false)
		p.Write(0, p.GetHeight()-p.GetStyle().H()-1, "... and other "+strconv.Itoa(cntHidden), false)
	}
	return 1
}

// NewTUIWidgetTree returns instance of TUIWidgetTree struct
func NewTUIWidgetTree() *TUIWidgetTree {
	w := &TUIWidgetTree{}
	return w
}
