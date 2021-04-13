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
	elemtext *widget.Entry
	elemtitle *widget.Label
	elemfilepath string
	elist *widget.List
	elemadd *widget.Button
	elemdel *widget.Button
	elemShow *container.Scroll
	elemtextShow *container.Scroll
	elemleftShow *fyne.Container
	elemallShow *fyne.Container
	saveTextE string
	lastSelectE widget.ListItemID
)
func ElemEdit(_ fyne.Window) fyne.CanvasObject {
	saveTextE = ""
	elemtext = widget.NewMultiLineEntry()
	elemtitle = widget.NewLabel("")
	elemlist := showElemList()
	elemtextShow = container.NewHScroll(container.NewVScroll(elemtext))
	elemtextShow.SetMinSize(fyne.NewSize(410,430))
	save := widget.NewButton("save",func (){
		cnf := dialog.NewConfirm("Confirmation", "Are you sure change this file?", ElemConfirmCallback, Windows)
		cnf.SetDismissText("Cacel")
		cnf.SetConfirmText("Save")
		cnf.Show()
	})
	elemadd = widget.NewButton("new",func(){
		filename := widget.NewEntry()
		filename.SetPlaceHolder("File Name")
		content := widget.NewForm(widget.NewFormItem("", filename))
		dialog.ShowCustomConfirm("Create elemfile","Create", "Cancel", content, func(b bool) {
			if b {
				name := controllers.ElemFilePath+"/"+filename.Text
				if Exists(name){
					cnf := dialog.NewConfirm("Confirmation", "File existed!", func(b bool){
						if b == true{
							err := ioutil.WriteFile(name,[]byte(""),666)
							if err != nil{
								dialog.ShowInformation("Attention", "Error!"+err.Error(), Windows)
							}else{
								ElemRewriteLeftShow()
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
						ElemRewriteLeftShow()
					}
				}
			}
		}, Windows)
	})
	elemdel = widget.NewButton("delete",func(){
		if elemfilepath == ""{
			dialog.ShowInformation("Attention", "No elem file selected!", Windows)
		}else{
			cnf := dialog.NewConfirm("Confirmation", "Are you sure delete this file?", ElemDeleteCallback, Windows)
			cnf.SetDismissText("Cacel")
			cnf.SetConfirmText("Delete")
			cnf.Show()
		}

	})

	elemShow = container.NewVScroll(elemlist)
	elemShow.SetMinSize(fyne.NewSize(70,450))
	elemleftShow = container.NewVBox(elemShow,elemadd,elemdel)
	rigthShow := container.NewVBox(elemtitle,widget.NewSeparator(),elemtextShow,save)
	elemallShow = container.NewHBox(elemleftShow,widget.NewSeparator(),rigthShow)
	return elemallShow
}
func showElemList()fyne.CanvasObject{
	var data []string
	data,_ = controllers.SortMap(controllers.GetElemFiles())
	elist = widget.NewList(
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
	elist.OnSelected = func(id widget.ListItemID) {
		if id == lastSelectE{
			return
		}
		if saveTextE != elemtext.Text && saveTextE != ""{
			cnf := dialog.NewConfirm("Confirmation", "File changed,sure to quit?", func(b bool){
				if b == true{
					elemtitle.SetText(data[id])
					elemfilepath = controllers.ElemFilePath + "/"+ data[id]
					content,err := ioutil.ReadFile(elemfilepath)
					if err != nil{
						dialog.ShowInformation("Attention", "Error!"+err.Error(), Windows)
					}else{
						elemtext.Text = string(content)
						saveTextE = elemtext.Text
						elemtext.Refresh()
						elemtextShow.Refresh()
					}
					lastSelectE = id
				}else{
					elist.Select(lastSelectE)
				}
			}, Windows)
			cnf.SetDismissText("Cacel")
			cnf.SetConfirmText("Quit")
			cnf.Show()
		}else{
			elemtitle.SetText(data[id])
			elemfilepath = controllers.ElemFilePath + "/"+ data[id]
			content,err := ioutil.ReadFile(elemfilepath)
			if err != nil{
				dialog.ShowInformation("Attention", "Error!"+err.Error(), Windows)
			}else{
				elemtext.Text = string(content)
				saveTextE = elemtext.Text
				elemtext.Refresh()
				elemtextShow.Refresh()
			}
			lastSelectE = id
		}


	}
	//list.OnUnselected = func(id widget.ListItemID) {
	//	label.SetText("Select An Item From The List")
	//}
	return elist
	//return list
}
func ElemConfirmCallback(response bool){
	if response == true{
		newContent := elemtext.Text
		err := ioutil.WriteFile(elemfilepath,[]byte(newContent),666)
		if err != nil{
			dialog.ShowInformation("Attention", "Error!"+err.Error(), Windows)
		}
	}
}

func ElemDeleteCallback(response bool){
	if response == true{
		err := os.Remove(elemfilepath)
		if err != nil{
			dialog.ShowInformation("Attention", "Error!"+err.Error(), Windows)
		}else{
			ElemRewriteLeftShow()
		}
	}
}
func ElemRewriteLeftShow(){
	elemleftShow.Remove(elemadd)
	elemleftShow.Remove(elemdel)
	elemleftShow.Remove(elemShow)
	elemlist := showElemList()
	elemShow = container.NewVScroll(elemlist)
	elemShow.SetMinSize(fyne.NewSize(150,450))
	elemleftShow.Add(elemShow)
	elemleftShow.Add(elemadd)
	elemleftShow.Add(elemdel)
	elemleftShow.Refresh()
}
