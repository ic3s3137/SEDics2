package gui

import "fyne.io/fyne"

type Tutorial struct {
	Title, Intro string
	View         func(w fyne.Window) fyne.CanvasObject
}

var (
	Tutorials = map[string]Tutorial{
		"RuleMode": {"RuleMode", "", ruleMode},
		//"QuickMode": {"QuickMode","",quickMode },
		"RuleEdit":{"RuleEdit","",ruleEdit},
		"ElemEdit":{"ElemEdit","",ElemEdit},

	}

	TutorialIndex = map[string][]string{
		"":            {"RuleMode", "RuleEdit","ElemEdit"},
	}
)
