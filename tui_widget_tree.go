package main

import (
	"github.com/gasiordev/go-tui"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type TUIWidgetTree struct {
	rootDir    string
	workDir    string
	hideFiles  bool
	hideDirs   bool
	showHidden bool
	filter     string
	highlight  string
}

func (w *TUIWidgetTree) GetRootDir() string {
	return w.rootDir
}

func (w *TUIWidgetTree) GetWorkDir() string {
	return w.workDir
}

func (w *TUIWidgetTree) comparePaths(p1 string, p2 string) int {
	if p1 == p2 {
		return 0
	}
	if strings.HasPrefix(p2, p1) {
		return 1
	}
	return -1
}

func (w *TUIWidgetTree) SetRootDir(rootDir string) {
	d, err := filepath.Abs(rootDir)
	if err != nil {
		w.rootDir = "/directory/that/does/not/exist"
		return
	}
	w.rootDir = d
}

func (w *TUIWidgetTree) SetWorkDir(workDir string) {
	d, err := filepath.Abs(workDir)
	if err != nil {
		w.workDir = "/directory/that/does/not/exist"
		return
	}
	w.workDir = d
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

func (w *TUIWidgetTree) clearBox(p *tui.TUIPane) {
	for i := 0; i < p.GetHeight()-p.GetStyle().H(); i++ {
		p.Write(0, i, strings.Repeat(" ", p.GetWidth()-p.GetStyle().V()), false)
	}
}

func (w *TUIWidgetTree) substrNameToWidth(s string, i int) string {
	if len(s) > i {
		s = s[0:i]
	}
	return s
}

func (w *TUIWidgetTree) getHighlightAnsiCode(n string) string {
	h := "m"
	if w.highlight != "" {
		m, err := regexp.MatchString(w.highlight, n)
		if m && err == nil {
			h = ";1m"
		}
	}
	return h
}

func (w *TUIWidgetTree) addColorAnsiCodes(n string, h string, f os.FileInfo) string {
	c := ""
	if f.Mode().IsDir() {
		c = "\u001b[33" + h + n + "\u001b[0m"
	} else if f.Mode()&os.ModeSymlink != 0 {
		c = "\u001b[36" + h + n + "\u001b[0m"
	} else if f.Mode().IsRegular() {
		c = "\u001b[32" + h + n + "\u001b[0m"
	} else {
		c = "\u001b[35" + h + n + "\u001b[0m"
	}
	return c
}

func (w *TUIWidgetTree) isMatchFilters(n string, f os.FileInfo) bool {
	t := true
	if f.IsDir() && w.hideDirs {
		t = false
	}
	if f.Mode().IsRegular() && w.hideFiles {
		t = false
	}
	if f.Name()[0] == '.' && !w.showHidden {
		t = false
	}
	if w.filter != "" {
		m, err := regexp.MatchString(w.filter, n)
		if !m || err != nil {
			t = false
		}
	}
	return t
}

func (w *TUIWidgetTree) printDir(p *tui.TUIPane, fs []os.FileInfo, depth int) {
	i := 0
	cntDisplayed := 0
	cntHidden := 1

	availableWidth := p.GetWidth() - p.GetStyle().V() - depth

	for _, file := range fs {
		origN := file.Name()
		n := w.substrNameToWidth(origN, availableWidth)
		hl := w.getHighlightAnsiCode(origN)
		n = w.addColorAnsiCodes(n, hl, file)
		if w.isMatchFilters(origN, file) {
			if cntDisplayed < p.GetHeight()-p.GetStyle().H() {
				p.Write(0, cntDisplayed, strings.Repeat(" ", depth)+n, false)
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

}

// Run is main function which just prints out the current time.
func (w *TUIWidgetTree) Run(p *tui.TUIPane) int {
	fileInfo, err := ioutil.ReadDir(w.workDir)
	if err != nil {
		return 0
	}
	w.clearBox(p)
	w.printDir(p, fileInfo, 0)
	return 1
}

// NewTUIWidgetTree returns instance of TUIWidgetTree struct
func NewTUIWidgetTree() *TUIWidgetTree {
	w := &TUIWidgetTree{}
	return w
}
