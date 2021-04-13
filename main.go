package main

import (
	"SEDics2/gui"
	"fyne.io/fyne"
	"fyne.io/fyne/container"
	"fyne.io/fyne/widget"
)

func main(){

	w := gui.Windows
	w.SetMaster()
	content := container.NewMax()
	setTutorial := func(t gui.Tutorial) {
		content.Objects = []fyne.CanvasObject{t.View(w)}
		content.Refresh()
	}

	tutorial := container.NewBorder(  nil, nil, nil, content)
	split := container.NewHSplit(makeNav(setTutorial), tutorial)
	//split := container.NewHBox(makeNav(setTutorial), tutorial)
	split.Offset = 0.1
	w.SetContent(split)
	w.Resize(fyne.NewSize(700, 550))
	w.ShowAndRun()
}
func makeNav(setTutorial func(tutorial gui.Tutorial)) fyne.CanvasObject {

	tree := &widget.Tree{
		ChildUIDs: func(uid string) []string {
			return gui.TutorialIndex[uid]
		},
		IsBranch: func(uid string) bool {
			children, ok := gui.TutorialIndex[uid]

			return ok && len(children) > 0
		},
		CreateNode: func(branch bool) fyne.CanvasObject {
			return widget.NewLabel("Collection Widgets")
		},
		UpdateNode: func(uid string, branch bool, obj fyne.CanvasObject) {
			t, ok := gui.Tutorials[uid]
			if !ok {
				fyne.LogError("Missing tutorial panel: "+uid, nil)
				return
			}
			obj.(*widget.Label).SetText(t.Title)
		},
		OnSelected: func(uid string) {
			if t, ok := gui.Tutorials[uid]; ok {
				setTutorial(t)
			}
		},
	}
	tree.Select("RuleMode")

	return container.NewBorder(nil, nil, nil, tree)
}


