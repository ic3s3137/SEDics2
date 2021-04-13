package controllers

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
)
type elem struct{
	ElemType string
	ElemName string
	List map[string]struct{}
}
type elemList struct{
	ElemList []elem
	ElemMap map[string]map[string]bool
	Level string

}
type elemIndex struct{
	IndexName []string
	Count int
}
func (el *elemList)makeElemMap(ruleMap map[string]map[string]bool){
	specialMap := make(map[string]bool)
	systemMap := make(map[string]bool)
	commonMap := make(map[string]bool)
	elems := make(map[string]map[string]map[string]struct{})
	elems[commonElem] = make(map[string]map[string]struct{})
	elemMap := make(map[string]map[string]bool)
	grepElemVars := regexp.MustCompile(grepElemReg)
	for ruleName,ruleList := range ruleMap{
		if len(elemMap[ruleName]) == 0{
			elemMap[ruleName] = make(map[string]bool)
		}
		for r := range ruleList{
			elemVal := grepElemVars.FindAllString(r,-1)
			for _,v := range elemVal{
				elemMap[ruleName][v] = true
			}
		}
		for k := range specialElem{
			elemMap[ruleName]["$"+k+"$"] = true
		}
	}
	for ruleName,elemNameList := range elemMap{
		for elemName := range elemNameList {
			elemType := checkElemType(elemName, ruleName)
			if elemType == "SPECIAL" { //如果为特殊元素
				delete(elemMap[ruleName], elemName)
				specialMap[elemName] = true
				if strings.Index(elemName,"$Pinyin") == 0{
					el.AddElem("SPECIAL","$String$",makeSpecialELem(elemName))
				}else{
					el.AddElem("SPECIAL",elemName,makeSpecialELem(elemName))
				}
				//TODO 添加特殊元素
				fmt.Println("special元素:", elemName)
			} else if elemType == "SYSTEM" {
				delete(elemMap[ruleName], elemName)
				systemMap[elemName] = true
				//TODO 添加系统元素
				fmt.Println("system:", elemName)
			}else if elemType == "COMMON"{
				delete(elemMap[ruleName], elemName)
				commonMap[elemName] = true
				el.readElemFromFile(elemName, elemType)
			}else if elemType == ruleName{
				el.readElemFromFile(elemName, elemType)
			}else if elemType == ""{
				//TODO 播报无效变量到logger界面中
				if elemName != "$String$" && elemName != "$Number$"{
					loginfo("Invaild element: "+elemName)
				}
				delete(elemMap[ruleName],elemName)
			}
		}
	}
	elemMap["SPECIAL"] = specialMap
	elemMap["SYSTEM"] = systemMap
	elemMap["COMMON"] = commonMap
	el.ElemMap = elemMap
	//fmt.Println(el.ElemMap)
}
func (el *elemList)readElemFromFile(ename string,etype string){
	var (
		file string
	)
	e := elem{ElemName: ename,ElemType: etype}
	ename = ename[1:len(ename)-1]
	e.List = make(map[string]struct{})
	if etype == "COMMON"{
		file = ElemFilePath + "/" + ename + ".txt"
	}else{
		file = ElemFilePath + "/" + etype[0:len(etype)-5] + "." + ename + ".txt"
	}
	content,err := ioutil.ReadFile(file)
	ExitErr(err)
	strContent := strings.Split(string(content),levelSpilt[el.Level])[0]
	lines := strings.Fields(strContent)
	for _,line := range lines{
		line = strings.TrimSpace(line)
		if levelTag[line] == false && line != ""{
			e.List[line] = struct{}{}
		}
	}
	el.ElemList = append(el.ElemList,e)
}
func (el *elemList)AddElem(elemType string,elemName string,elemlist map[string]struct{}){
	var (
		existed = false
		index int
	)
	for i,e := range el.ElemList{
		if e.ElemName == elemName{
			existed = true
			index = i
			break
		}
	}
	if existed == true{
		for a := range elemlist{
			el.ElemList[index].List[a] = struct{}{}
		}
	}else{
		e := elem{ElemName: elemName,ElemType: elemType,List: elemlist}
		el.ElemList = append(el.ElemList,e)
	}
	if el.ElemMap != nil{
		if len(el.ElemMap[elemType]) == 0{
			el.ElemMap[elemType] = make(map[string]bool)
		}
		el.ElemMap[elemType][elemName] = true
	}

}
func (el *elemList)GetElemByIndex(Index []string)map[string]map[string]struct{}{
	var elemContent = make(map[string]map[string]struct{})
	for _,i := range Index{
		for _,e := range el.ElemList{
			if i == e.ElemName || i == e.ElemType+":::"+e.ElemName{
				i = rename(elemContent,i)
				elemContent[i] = make(map[string]struct{})
				elemContent[i] = e.List
				break
			}
		}
	}
	return elemContent
}
func checkElemType(ElemName string,ruleName string)string{
	//优先级 special > rule > common > system
	ElemName = ElemName[1 : len(ElemName)-1]
	if SpecialElems[ElemName] == true{
		return "SPECIAL"
	}
	fileName := ruleName[0:len(ruleName)-5] + "." + ElemName + ".txt"
	elemFileList,err := ioutil.ReadDir(ElemFilePath)
	ExitErr(err)
	for _,f := range elemFileList{
		if f.Name() == fileName{
			return ruleName
		}
	}
	fileName = ElemName + ".txt"
	for _,f := range elemFileList{
		if f.Name() == fileName{
			return "COMMON"
		}
	}
	if systemElems[ElemName]{
		return "SYSTEM"
	}
	return ""
}
func makeSpecialELem(elename string)map[string]struct{}{
	elename = elename[1:len(elename)-1]
	var list map[string]struct{}
	if elename == "Domain" {
		list = makeDomainElems(specialElem[elename])
	}else if elename[:6] == "Pinyin"{
		list = makePinyinElems(elename,specialElem[elename])
	}else{
		list = make(map[string]struct{})
		for _,e := range specialElem[elename]{
			list[e] = struct{}{}
		}
	}
	return list
}

func GetElemFiles()map[string]bool{
	filelist := make(map[string]bool)
	dirInfo,err := ioutil.ReadDir(ElemFilePath)
	ExitErr(err)
	for i := range dirInfo{
		filelist[dirInfo[i].Name()] = false
	}
	return filelist
}
