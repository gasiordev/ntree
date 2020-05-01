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

func (w *TUIWidgetTree) substrName(s string, i int) string {
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

func (w *TUIWidgetTree) addColorAnsiCodes(n string, h string, cmp int, f os.FileInfo) string {
	c := ""
	if f.Mode().IsDir() {
		if cmp > -1 {
			c = "\u001b[37" + h + n + "\u001b[0m"
		} else {
			c = "\u001b[33" + h + n + "\u001b[0m"
		}
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

func (w *TUIWidgetTree) getFileDetails(file os.FileInfo, rootPath string, substrValue int, colours bool) (string, string, string, int, bool) {
	name := file.Name()
	path := filepath.Join(w.rootDir, rootPath, name)
	substr := w.substrName(name, substrValue)
	cmp := w.comparePaths(path, w.workDir)

	if colours {
		hl := w.getHighlightAnsiCode(name)
		substr = w.addColorAnsiCodes(substr, hl, cmp, file)
	}

	filters := w.isMatchFilters(name, file)

	return name, path, substr, cmp, filters
}

func (w *TUIWidgetTree) getWorkDirDepth() int {
	return len(strings.Split(w.workDir, "/"))
}

func (w *TUIWidgetTree) getWorkDirDepthCounts(depthCounts *[]int, maxDepth *int, p *tui.TUIPane, fs []os.FileInfo, rootPath string, depth int) {

	availableWidth := p.GetWidth() - p.GetStyle().V() - depth

	cnt := 0
	for _, file := range fs {
		fileName, filePath, _, fileCmp, fileMatchFilters := w.getFileDetails(file, rootPath, availableWidth, true)
		if fileMatchFilters || fileCmp > -1 {
			cnt++
			if fileCmp > -1 {
				fileInfo, err := ioutil.ReadDir(filePath)
				if err == nil {
					w.getWorkDirDepthCounts(depthCounts, maxDepth, p, fileInfo, filepath.Join(rootPath, fileName), depth+1)
				}
			}
		}
	}
	(*depthCounts)[depth] = cnt
	if depth > *maxDepth {
		*maxDepth = depth
	}
}

func (w *TUIWidgetTree) printDir(p *tui.TUIPane, fs []os.FileInfo, rootPath string, depth int, displayed int, depthCounts []int, maxDepth int) int {
	cntDisplayed := displayed

	availableWidth := p.GetWidth() - p.GetStyle().V() - depth
	availableHeight := p.GetHeight() - p.GetStyle().H() - depth

	depthCountSum := 0
	if depth < maxDepth {
		for j := depth+1; j <= maxDepth; j++ {
			depthCountSum += depthCounts[j]
		}
	}

	workDirOpened := false
	cntBefore := 0
	cntAfter := 0

	for i, file := range fs {
		fileName, filePath, fileDisplayName, fileCmp, fileMatchFilters := w.getFileDetails(file, rootPath, availableWidth, true)

		if !fileMatchFilters {
			continue
		}

		if depth == maxDepth {
			if cntDisplayed < availableHeight {
				p.Write(0, cntDisplayed, strings.Repeat(" ", depth)+fileDisplayName, false)
				cntDisplayed++
				cntAfter++
			}
		} else {
			if cntDisplayed < availableHeight {
				if workDirOpened || fileCmp > -1 {
					if fileCmp > -1 && cntBefore > 0 {
						p.Write(0, cntDisplayed, strings.Repeat(" ", depth)+"("+strconv.Itoa(cntBefore)+")... "+fileDisplayName, false)
					} else {
						if cntDisplayed+1 == availableHeight && i+1 < len(fs) {
							p.Write(0, cntDisplayed, strings.Repeat(" ", depth)+fileDisplayName+" ...("+strconv.Itoa(len(fs)-i-1)+")", false)
						} else {
							p.Write(0, cntDisplayed, strings.Repeat(" ", depth)+fileDisplayName, false)
						}
					}
					cntDisplayed++
					cntAfter++
				} else {
					if cntDisplayed + depthCountSum < availableHeight {
						p.Write(0, cntDisplayed, strings.Repeat(" ", depth)+fileDisplayName, false)
						cntDisplayed++
						cntBefore++
					}
				}
			}
		}

		if fileCmp > -1 {
			fileInfo, err := ioutil.ReadDir(filePath)
			if err == nil {
				subDisplayed := w.printDir(p, fileInfo, filepath.Join(rootPath, fileName), depth+1, cntDisplayed, depthCounts, maxDepth)
				cntDisplayed = subDisplayed
				workDirOpened = true
			}
		}
	}
	return cntDisplayed
}

// Run is main function which just prints out the current time.
func (w *TUIWidgetTree) Run(p *tui.TUIPane) int {
	fileInfo, err := ioutil.ReadDir(w.rootDir)
	if err != nil {
		return 0
	}

	depthCounts := make([]int, 30)
	maxDepth := 0
	w.getWorkDirDepthCounts(&depthCounts, &maxDepth, p, fileInfo, "", 0)

	w.clearBox(p)
	w.printDir(p, fileInfo, "", 0, 0, depthCounts, maxDepth)
	return 1
}

// NewTUIWidgetTree returns instance of TUIWidgetTree struct
func NewTUIWidgetTree() *TUIWidgetTree {
	w := &TUIWidgetTree{}
	return w
}
