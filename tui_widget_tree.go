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

// comparePaths checks if two paths are the same or one is part of another
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

// clearBox clears the pane (fills it with spaces)
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

// getHighlightAnsiCode checks for highlight value, compares it with given
// value and if string should be highlighted then returns an apropriate
// ANSI code
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

// addColorAnsiCodes adds look to file name depending on its type
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

// isMatchFilters checks if file matches filters and shouldn't be hidden
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

// getFileDetails takes file and returns its various details. Function is used
// when iterating a directory
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

// getWorkDirDepthCounts walks through work dir and counts its files (and dirs)
// that will be displayed (so checks if they match any filters etc.)
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

// printDir iterates over directory items, filters them and prints out things
// onto screen. It's called recursively. Main logic is here.
func (w *TUIWidgetTree) printDir(p *tui.TUIPane, fs []os.FileInfo, rootPath string, depth int, displayed int, depthCounts []int, maxDepth int) int {
	cntDisplayed := displayed

	availableWidth := p.GetWidth() - p.GetStyle().V() - depth
	availableHeight := p.GetHeight() - p.GetStyle().H() - depth

	// sum up items in all subdirectories
	depthCountSum := 0
	if depth < maxDepth {
		for j := depth + 1; j <= maxDepth; j++ {
			depthCountSum += depthCounts[j]
		}
	}

	workDirOpened := false // have we already iterated through directory from work dir
	cntBefore := 0         // number of items before work dir
	cntAfter := 0          // number of items after work dir

	j := 0
	for i, file := range fs {
		fileName, filePath, fileDisplayName, fileCmp, fileMatchFilters := w.getFileDetails(file, rootPath, availableWidth, true)

		// ignore files that do not match filter or should be hidden (hide files, show hidden, hide dirs etc.)
		if !fileMatchFilters && fileCmp < 0 {
			continue
		}

		// if we are in the last directory from work dir then print all the items
		if depth == maxDepth {
			if cntDisplayed < availableHeight {
				p.Write(0, cntDisplayed, strings.Repeat(" ", depth)+fileDisplayName, false)
				cntDisplayed++
				cntAfter++
			}
		} else {
			if cntDisplayed < availableHeight {
				// if directory that is part of work dir was already printed or it is that directory then we can print
				if workDirOpened || fileCmp > -1 {
					// if it's dir on the work dir path and there were files in this directory that were not displayed, let's prepend the filename with number of them
					if fileCmp > -1 && j-cntBefore > 0 {
						p.Write(0, cntDisplayed, strings.Repeat(" ", depth)+"("+strconv.Itoa(j-cntBefore)+")... "+fileDisplayName, false)
					} else {
						// if there is no place to print more directories but still there are some, let's append number of them at the end of filename
						if cntDisplayed+1 == availableHeight && i+1 < len(fs) {
							p.Write(0, cntDisplayed, strings.Repeat(" ", depth)+fileDisplayName+" ...("+strconv.Itoa(len(fs)-i-1)+")", false)
						} else {
							p.Write(0, cntDisplayed, strings.Repeat(" ", depth)+fileDisplayName, false)
						}
					}
					cntDisplayed++
					cntAfter++
				} else {
					// if we can display the file, let's do that
					if cntDisplayed+depthCountSum < availableHeight {
						p.Write(0, cntDisplayed, strings.Repeat(" ", depth)+fileDisplayName, false)
						cntDisplayed++
						cntBefore++
					}
				}
			}
		}

		// if the directory is on work dir path, open it and print out its contents
		if fileCmp > -1 {
			fileInfo, err := ioutil.ReadDir(filePath)
			if err == nil {
				subDisplayed := w.printDir(p, fileInfo, filepath.Join(rootPath, fileName), depth+1, cntDisplayed, depthCounts, maxDepth)
				cntDisplayed = subDisplayed
				workDirOpened = true
			}
		}

		j++
	}
	return cntDisplayed
}

// Run is main function that opens root directory and calls the recursive print
// function to print the directories on the stdout
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
