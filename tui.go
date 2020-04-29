package main

import (
	"github.com/gasiordev/go-tui"
)

func getOnTUIPaneDraw(n *NTree, w *TUIWidgetTree, p *tui.TUIPane) func(*tui.TUIPane) int {
	fn := func(x *tui.TUIPane) int {
		w.SetRootDir(n.GetRootDir())
		w.SetWorkDir(n.GetWorkDir())
		w.SetHideFiles(n.GetHideFiles())
		w.SetHideDirs(n.GetHideDirs())
		w.SetShowHidden(n.GetShowHidden())
		w.SetFilter(n.GetFilter())
		w.SetHighlight(n.GetHighlight())
		return w.Run(x)
	}
	return fn
}

func getOnTUIKeyPress(n *NTree) func(*tui.TUI, []byte) {
	fn := func(x *tui.TUI, b []byte) {
		ch := string(b)
		if ch == "q" || ch == "Q" {
			x.Exit(0)
		}
		if ch == "r" || ch == "R" {
			n.ResetFilter()
			n.ResetHighlight()
		}
		if ch == "d" || ch == "D" {
			n.ToggleHideDirs()
		}
		if ch == "f" || ch == "F" {
			n.ToggleHideFiles()
		}
		if ch == "h" || ch == "H" {
			n.ToggleShowHidden()
		}
	}
	return fn
}

func NewNTreeTUI(n *NTree) *tui.TUI {
	nTreeTUI := tui.NewTUI("Ntree", "Project tree widget", "Mikolaj Gasior")

	nTreeTUI.SetOnKeyPress(getOnTUIKeyPress(n))

	p0 := nTreeTUI.GetPane()
	s1 := tui.NewTUIPaneStyleNone()
	p0.SetStyle(s1)

	w := NewTUIWidgetTree()
	w.InitPane(p0)

	p0.SetOnDraw(getOnTUIPaneDraw(n, w, p0))
	p0.SetOnIterate(getOnTUIPaneDraw(n, w, p0))

	return nTreeTUI
}
