package gui

import (
	"SEDics2/controllers"
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/container"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/widget"
	"strings"
)
var (
	stringElem = make(map[string]struct{})
	stringElemMap = make(map[string]bool)
	numberElem = make(map[string]struct{})
	numberElemMap = make(map[string]bool)
	specialElem = make(map[string][]string)
	specialElemMap = make(map[string]bool)
	strElemWidget = make(map[string]fyne.CanvasObject)
	numElemWidget = make(map[string]fyne.CanvasObject)
	speElemWidget = make(map[string]fyne.CanvasObject)
	ruleFileSetting = make(map[string]bool)
	specialSelectBar []string
	level string   //字典的大小,分为small,medium,big
	specialType string
	strremovelist []string
	numremovelist []string
	totalremovelist []string
	speremovelist []string
	msgch = make(chan string)
	loginfo string
	a = app.New()
)
var Windows = a.NewWindow("SEDics2 author by ic3s")
func ruleMode(_ fyne.Window) fyne.CanvasObject {
	stringElem = make(map[string]struct{})
	numberElem = make(map[string]struct{})
	specialElem = make(map[string][]string)
	stringElemMap = make(map[string]bool)
	numberElemMap = make(map[string]bool)
	specialElemMap = make(map[string]bool)
	strElemWidget = make(map[string]fyne.CanvasObject)
	numElemWidget = make(map[string]fyne.CanvasObject)
	speElemWidget = make(map[string]fyne.CanvasObject)
	strremovelist = []string{}
	numremovelist = []string{}
	totalremovelist = []string{}
	speremovelist = []string{}
	specialType = ""
	stringShow := container.NewVBox()
	numberShow := container.NewVBox()
	specialShow := container.NewVBox()
	totalShow := container.NewVBox()
	showTabs := widget.NewTabContainer(
		widget.NewTabItem("total",container.NewVScroll(totalShow)),
		widget.NewTabItem("string",container.NewVScroll(stringShow)),
		widget.NewTabItem("number",container.NewVScroll(numberShow)),
		widget.NewTabItem("special",container.NewVScroll(specialShow)),
	)
	ruleShow := makeRuleList()
	seString := widget.NewEntry()
	seString.SetPlaceHolder("String Type  ")
	seNum := widget.NewEntry()
	seNum.SetPlaceHolder("Number Type  ")
	input := container.NewVBox(container.NewHScroll(seString),container.NewHScroll(seNum))
	specialSelectBar,_ = controllers.SortMap(controllers.SpecialElems)
	typeSelect := widget.NewSelect(specialSelectBar, func(s string) { specialType = s })
	seSpecial := widget.NewEntry()
	seSpecial.SetPlaceHolder("Special Type")
	add,del,start,reset := makeButton(totalShow,stringShow,numberShow,specialShow,seString,seNum,seSpecial)
	button2 := container.NewAdaptiveGrid(3,add,del,reset)
	specialBar := container.NewVBox(typeSelect,container.NewHScroll(seSpecial))
	logger := widget.NewLabel(loginfo)
	loggerBox := container.NewHScroll(container.NewVScroll(logger))
	loggerBox.SetMinSize(fyne.NewSize(400,300))
	go messageLogger(logger)
	//showlog := container.NewVScroll(container.NewVBox(logger))
	//showlog.SetMinSize(fyne.NewSize(80,400))
	inputBox := container.NewVBox(input,specialBar,button2,start)
	rightTop := container.NewHBox(inputBox,widget.NewSeparator(),showTabs)
	right := container.NewVBox(rightTop,widget.NewSeparator(),loggerBox)
	return container.NewHBox(ruleShow,widget.NewSeparator(),right)
}
func makeRuleList()fyne.CanvasObject{
	var ruleFileWidget []fyne.CanvasObject
	levelSelect := widget.NewSelect([]string{"small","medium","big"},func(s string){level = s})
	levelSelect.SetSelected("small")
	ruleFileWidget = append(ruleFileWidget,levelSelect)
	ruleFileSetting = controllers.GetRuleFiles()
	data,_ := controllers.SortMap(ruleFileSetting)
	for _,f := range data{
		w := widget.NewCheck(f, func(on bool) {
		})
		w.OnChanged = func(on bool){
			ruleFileSetting[w.Text] = on
		}
		ruleFileWidget = append(ruleFileWidget, w)
	}
	//ruleShow := container.NewVScroll(widget.NewVBox(ruleFileWidget...))
	ruleShow := container.NewVBox(widget.NewVBox(ruleFileWidget...))
	//ruleShow := container.NewVScroll(container.NewVBox(ruleFileWidget...))
	return ruleShow
}
func makeButton(totalShow *fyne.Container,stringShow *fyne.Container,numberShow *fyne.Container,specialShow *fyne.Container,seString *widget.Entry,seNum *widget.Entry,seSpecial *widget.Entry)(fyne.CanvasObject,fyne.CanvasObject,fyne.CanvasObject,fyne.CanvasObject){
	add := widget.NewButton("add",func(){
		var text string
		if seString.Text != ""{
			text = seString.Text
			seString.SetText("")
			if stringElemMap[text] == false {
				stringElem[text] = struct{}{}
				stringElemMap[text] = true
				w := widget.NewCheck(text, func(on bool) {
					if on == true {
						strremovelist = append(strremovelist, text)
						totalremovelist = append(totalremovelist, text)
					}
				})
				strElemWidget[text] = w
				stringShow.Add(w)
				totalShow.Add(w)
			}
		}else if seNum.Text != ""{
			text = seNum.Text
			seNum.SetText("")
			if numberElemMap[text] == false{
				numberElem[text] = struct{}{}
				numberElemMap[text] = true
				w := widget.NewCheck(text, func(on bool) {
					if on == true {
						numremovelist = append(numremovelist, text)
						totalremovelist = append(totalremovelist, text)
					}
				})
				numElemWidget[text] = w
				numberShow.Add(w)
				totalShow.Add(w)
			}
		}else if seSpecial.Text != ""{
			if specialType == ""{
				dialog.ShowInformation("Attention", "No special type selected", Windows)
			}else{
				text = seSpecial.Text
				seSpecial.SetText("")
				if specialElemMap[specialType+":::"+text] == false {
					specialElemMap[specialType+":::"+text] = true
					tag := "[" + specialType + "]" + text
					specialElem[specialType] = append(specialElem[specialType], text)
					w := widget.NewCheck(tag, func(on bool) {
						if on == true {
							speremovelist = append(speremovelist, specialType+":::"+text)
							totalremovelist = append(totalremovelist, specialType+":::"+text)
						}
					})
					speElemWidget[specialType+":::"+text] = w
					specialShow.Add(w)
					totalShow.Add(w)
				}
			}
		}
	})
	del := widget.NewButton("del", func() {
		for _, e := range strremovelist {
			stringShow.Remove(strElemWidget[e])
			delete(stringElem,e)
			strremovelist = []string{}
			stringElemMap[e] = false
		}
		for _, e := range numremovelist {
			numberShow.Remove(numElemWidget[e])
			delete(numberElem,e)
			numremovelist = []string{}
			numberElemMap[e] = false
		}
		for _,e := range speremovelist{
			specialShow.Remove(speElemWidget[e])
			specialElemMap[e] = false
			setype := strings.Split(e,":::")[0]
			ele := strings.Split(e,":::")[1]
			for i,s := range specialElem[setype]{
				if s == ele{
					specialElem[setype] = append(specialElem[setype][:i],specialElem[setype][i+1:]...)
				}
			}
			speremovelist = []string{}
		}
		for _,e := range totalremovelist{
			totalShow.Remove(strElemWidget[e])
			totalShow.Remove(numElemWidget[e])
			totalShow.Remove(speElemWidget[e])
		}
	})
	reset := widget.NewButton("reset",func(){
		for i := range strElemWidget{
			stringShow.Remove(strElemWidget[i])
			totalShow.Remove(strElemWidget[i])
			delete(strElemWidget,i)
		}
		for i := range numElemWidget{
			numberShow.Remove(numElemWidget[i])
			totalShow.Remove(numElemWidget[i])
			delete(numElemWidget,i)
		}
		for i := range speElemWidget{
			specialShow.Remove(speElemWidget[i])
			totalShow.Remove(speElemWidget[i])
			delete(speElemWidget,i)
		}
		stringElem = make(map[string]struct{})
		numberElem = make(map[string]struct{})
		specialElem = make(map[string][]string)
		stringElemMap = make(map[string]bool)
		numberElemMap = make(map[string]bool)
		specialElemMap = make(map[string]bool)
		strElemWidget = make(map[string]fyne.CanvasObject)
		numElemWidget = make(map[string]fyne.CanvasObject)
		speElemWidget = make(map[string]fyne.CanvasObject)
		strremovelist = []string{}
		numremovelist = []string{}
		totalremovelist = []string{}
		speremovelist = []string{}
		specialType = ""

	})
	start := widget.NewButton("start", func() {
		selected := false
		for _,v := range ruleFileSetting{
			if v == true{
				selected = true
				break
			}
		}
		if selected == false{
				dialog.ShowInformation("Attention", "No rule file selected!", Windows)
		}else{
			controllers.SetLevel(level)
			controllers.SetRuleFile(ruleFileSetting)
			controllers.SetStringElem(stringElem)
			controllers.SetNumberElem(numberElem)
			controllers.SetSpecialElem(specialElem)
			controllers.SetLogger(&msgch)
			controllers.Start("RuleMode")
		}
	})
	return add,del,start,reset
}
func messageLogger(logger *widget.Label){
	for msg := range msgch{
		loginfo = loginfo + msg + "\n"
		logger.SetText(loginfo)
		logger.Refresh()
	}
}
