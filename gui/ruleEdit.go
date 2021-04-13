package gui

import (
	"SEDics2/controllers"
	"fyne.io/fyne"
	"fyne.io/fyne/container"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"io/ioutil"
	"os"
)
var (
	ruletext *widget.Entry
	ruletitle *widget.Label
	rulefilepath string
	rlist *widget.List
	ruleadd *widget.Button
	ruledel *widget.Button
	ruleShow *container.Scroll
	ruletextShow *container.Scroll
	ruleleftShow *fyne.Container
	ruleallShow *fyne.Container
	saveText string
	lastSelect widget.ListItemID
)
func ruleEdit(_ fyne.Window) fyne.CanvasObject {
	saveText = ""
	ruletext = widget.NewMultiLineEntry()
	ruletitle = widget.NewLabel("")
	rulelist := showRuleList()
	ruletextShow = container.NewHScroll(container.NewVScroll(ruletext))
	ruletextShow.SetMinSize(fyne.NewSize(410,430))
	save := widget.NewButton("save",func (){
		cnf := dialog.NewConfirm("Confirmation", "Are you sure change this file?", ruleConfirmCallback, Windows)
		cnf.SetDismissText("Cacel")
		cnf.SetConfirmText("Save")
		cnf.Show()
	})
	ruleadd = widget.NewButton("new",func(){
		filename := widget.NewEntry()
		filename.SetPlaceHolder("File Name")
		content := widget.NewForm(widget.NewFormItem("", filename))
		dialog.ShowCustomConfirm("Create rule","Create", "Cancel", content, func(b bool) {
			if b {
				name := controllers.RuleFilePath+"/"+filename.Text
				if Exists(name){
					cnf := dialog.NewConfirm("Confirmation", "File existed!", func(b bool){
						if b == true{
							err := ioutil.WriteFile(name,[]byte(""),666)
							if err != nil{
								dialog.ShowInformation("Attention", "Error!"+err.Error(), Windows)
							}else{
								ruleRewriteLeftShow()
							}
						}
					}, Windows)
					cnf.SetDismissText("Cacel")
					cnf.SetConfirmText("Overwrite")
					cnf.Show()
				}else{
					err := ioutil.WriteFile(name,[]byte(""),666)
					if err != nil{
						dialog.ShowInformation("Attention", "Error!"+err.Error(), Windows)
					}else{
						ruleRewriteLeftShow()
					}
				}
			}
		}, Windows)
	})
	ruledel = widget.NewButton("delete",func(){
		if rulefilepath == ""{
			dialog.ShowInformation("Attention", "No rule file selected!", Windows)
		}else{
			cnf := dialog.NewConfirm("Confirmation", "Are you sure delete this file?", ruleDeleteCallback, Windows)
			cnf.SetDismissText("Cacel")
			cnf.SetConfirmText("Delete")
			cnf.Show()
		}

	})
	ruleShow = container.NewVScroll(rulelist)
	ruleShow.SetMinSize(fyne.NewSize(90,450))
	ruleleftShow = container.NewVBox(ruleShow,ruleadd,ruledel)
	rigthShow := container.NewVBox(ruletitle,widget.NewSeparator(),ruletextShow,save)
	ruleallShow = container.NewHBox(ruleleftShow,widget.NewSeparator(),rigthShow)
	return ruleallShow
}
func showRuleList()fyne.CanvasObject{
	var data []string
	data,_ = controllers.SortMap(controllers.GetRuleFiles())

	 rlist = widget.NewList(
		func() int {
			return len(data)
		},
		func() fyne.CanvasObject {
			return fyne.NewContainerWithLayout(layout.NewHBoxLayout(), widget.NewIcon(theme.DocumentIcon()), widget.NewLabel("12345678901234567891"))
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			item.(*fyne.Container).Objects[1].(*widget.Label).SetText(data[id])
		},
	)
	rlist.OnSelected = func(id widget.ListItemID) {
		if id == lastSelect{
			return
		}
		if saveText != ruletext.Text{
			cnf := dialog.NewConfirm("Confirmation", "File changed,sure to quit?", func(b bool){
				if b == true{
					ruletitle.SetText(data[id])
					rulefilepath = controllers.RuleFilePath + "/"+ data[id]
					content,err := ioutil.ReadFile(rulefilepath)
					if err != nil{
						dialog.ShowInformation("Attention", "Error!"+err.Error(), Windows)
					}else{
						ruletext.Text = string(content)
						saveText = ruletext.Text
						ruletext.Refresh()
					}
					lastSelect = id
				}else{
					rlist.Select(lastSelect)
				}
			}, Windows)
			cnf.SetDismissText("Cacel")
			cnf.SetConfirmText("Quit")
			cnf.Show()
		}else{
			ruletitle.SetText(data[id])
			rulefilepath = controllers.RuleFilePath + "/"+ data[id]
			content,err := ioutil.ReadFile(rulefilepath)
			if err != nil{
				dialog.ShowInformation("Attention", "Error!"+err.Error(), Windows)
			}else{
				ruletext.Text = string(content)
				saveText = ruletext.Text
				ruletext.Refresh()
			}
			lastSelect = id
		}

	}
	return rlist
	//return list
}
func ruleConfirmCallback(response bool){
	if response == true{
		newContent := ruletext.Text
		err := ioutil.WriteFile(rulefilepath,[]byte(newContent),666)
		if err != nil{
			dialog.ShowInformation("Attention", "Error!"+err.Error(), Windows)
		}
	}
}

func ruleDeleteCallback(response bool){
	if response == true{
		err := os.Remove(rulefilepath)
		if err != nil{
			dialog.ShowInformation("Attention", "Error!"+err.Error(), Windows)
		}else{
			ruleRewriteLeftShow()
		}
	}
}
func ruleRewriteLeftShow(){
	ruleleftShow.Remove(ruleadd)
	ruleleftShow.Remove(ruledel)
	ruleleftShow.Remove(ruleShow)
	rulelist := showRuleList()
	ruleShow = container.NewVScroll(rulelist)
	ruleShow.SetMinSize(fyne.NewSize(150,450))
	ruleleftShow.Add(ruleShow)
	ruleleftShow.Add(ruleadd)
	ruleleftShow.Add(ruledel)
	ruleleftShow.Refresh()
}
func Exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}