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

func NewNTreeTUI(n *NTree) *tui.TUI {
	nTreeTUI := tui.NewTUI("Ntree", "Project tree widget", "Mikolaj Gasior")

	p0 := nTreeTUI.GetPane()
	s1 := tui.NewTUIPaneStyleNone()
	p0.SetStyle(s1)

	w := NewTUIWidgetTree()
	w.InitPane(p0)

	p0.SetOnDraw(getOnTUIPaneDraw(n, w, p0))
	p0.SetOnIterate(getOnTUIPaneDraw(n, w, p0))

	return nTreeTUI
}
