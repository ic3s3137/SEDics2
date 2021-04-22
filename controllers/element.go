package controllers

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

type elem struct {
	ElemType string
	ElemName string
	List     map[string]struct{}
}
type elemList struct {
	ElemList []elem
	ElemMap  map[string]map[string]bool
	Level    string
}
type elemIndex struct {
	IndexName []string
	Count     int
}

const (
	SPECIAL_TYPE       = "SPECIAL"
	SYSTEM_TYPE        = "SYSTEM"
	COMMON_TYPE        = "COMMON"
	SPECIAL_STRING_KEY = "$String$"
	SPECIAL_NUMBER_KEY = "$Number$"
	SPECIAL_PINYIN_KEY = "$Pinyin$"
	SPECIAL_TEMP_KEY   = "$Temp$"
	EXIST_ELEMENT      = "12490812347"
)

var (
	specialElemFunc = map[string]func([]string) map[string]struct{}{
		SPECIAL_PINYIN_KEY: makePinyinElems,
		SPECIAL_TEMP_KEY:   makeTempElems,
	}
	SpecialElems = map[string]bool{"Pinyin": true, "Temp": true}
)

func (el *elemList) makeElemMap(ruleMap map[string]map[string]bool) {
	specialMap := make(map[string]bool)
	systemMap := make(map[string]bool)
	commonMap := make(map[string]bool)
	elems := make(map[string]map[string]map[string]struct{})
	elems[commonElem] = make(map[string]map[string]struct{})
	elemMap := make(map[string]map[string]bool)
	grepElemVars := regexp.MustCompile(grepElemReg)
	//提取出每个规则文件所需用到的元素地图
	for ruleName, ruleList := range ruleMap {
		if len(elemMap[ruleName]) == 0 {
			elemMap[ruleName] = make(map[string]bool)
		}
		for r := range ruleList {
			elemVal := grepElemVars.FindAllString(r, -1)
			for _, v := range elemVal {
				elemMap[ruleName][v] = true
			}
		}
		for k := range specialElem {
			elemMap[ruleName]["$"+k+"$"] = true
		}
	}
	//判断元素地图中的每个元素所对应的类型：规则专属元素，特殊元素，通用元素
	for ruleName, elemNameList := range elemMap {
		for elemName := range elemNameList {
			elemType := checkElemType(elemName, ruleName)
			if elemType == SPECIAL_TYPE { //如果为特殊元素
				delete(elemMap[ruleName], elemName)
				specialMap[elemName] = true
				eleString := elemName[1 : len(elemName)-1]
				if elemName == SPECIAL_PINYIN_KEY {
					el.AddElem(SPECIAL_TYPE, SPECIAL_STRING_KEY, specialElemFunc[elemName](specialElem[eleString]))
				} else if elemName == SPECIAL_TEMP_KEY {
					for _, v := range specialElem[eleString] {
						elemTmp := strings.SplitN(v, ":", 2)
						if len(elemTmp) != 2 {
							loginfo("错误的直定义变量:" + v)
							os.Exit(1)
						}
						el.AddElem(SPECIAL_TYPE, "$"+elemTmp[0]+"$", specialElemFunc[elemName]([]string{elemTmp[1]}))
						specialMap["$"+elemTmp[0]+"$"] = true
					}
				} else {
					el.AddElem(SPECIAL_TYPE, elemName, specialElemFunc[elemName](specialElem[eleString]))
				}
				//TODO 添加特殊元素
				fmt.Println("已定义的元素:", elemName)
			} else if elemType == SYSTEM_TYPE {
				delete(elemMap[ruleName], elemName)
				systemMap[elemName] = true
				//TODO 添加系统元素
				fmt.Println("system:", elemName)
			} else if elemType == COMMON_TYPE {
				delete(elemMap[ruleName], elemName)
				commonMap[elemName] = true
				el.readElemFromFile(elemName, elemType)
			} else if elemType == ruleName {
				el.readElemFromFile(elemName, elemType)
			} else if elemType == "" {
				//TODO 播报无效变量到logger界面中
				if elemName != SPECIAL_STRING_KEY && elemName != SPECIAL_NUMBER_KEY {
					loginfo("忽略无效变量:" + elemName)
				}
				delete(elemMap[ruleName], elemName)
			} else if elemType == EXIST_ELEMENT {
				delete(elemMap[ruleName], elemName)
			}
		}
	}
	elemMap[SPECIAL_TYPE] = specialMap
	elemMap[SYSTEM_TYPE] = systemMap
	elemMap[COMMON_TYPE] = commonMap
	el.ElemMap = elemMap
	//fmt.Println(el.ElemList)
}
func (el *elemList) readElemFromFile(ename string, etype string) {
	var (
		file string
	)
	e := elem{ElemName: ename, ElemType: etype}
	ename = ename[1 : len(ename)-1]
	e.List = make(map[string]struct{})
	if etype == COMMON_TYPE {
		file = ElemFilePath + "/" + ename + ".txt"
	} else {
		file = ElemFilePath + "/" + etype[0:len(etype)-5] + "." + ename + ".txt"
	}
	content, err := ioutil.ReadFile(file)
	ExitErr(err)
	strContent := strings.Split(string(content), levelSpilt[el.Level])[0]
	lines := strings.Fields(strContent)
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if levelTag[line] == false && line != "" {
			e.List[line] = struct{}{}
		}
	}
	el.ElemList = append(el.ElemList, e)
}
func (el *elemList) AddElem(elemType string, elemName string, elemlist map[string]struct{}) {
	var (
		existed = false
		index   int
	)
	for i, e := range el.ElemList {
		if e.ElemName == elemName {
			existed = true
			index = i
			break
		}
	}
	if existed == true {
		for a := range elemlist {
			el.ElemList[index].List[a] = struct{}{}
		}
	} else {
		e := elem{ElemName: elemName, ElemType: elemType, List: elemlist}
		el.ElemList = append(el.ElemList, e)
	}
	if el.ElemMap != nil {
		if len(el.ElemMap[elemType]) == 0 {
			el.ElemMap[elemType] = make(map[string]bool)
		}
		el.ElemMap[elemType][elemName] = true
	}

}
func (el *elemList) GetElemByIndex(Index []string) map[string]map[string]struct{} {
	var elemContent = make(map[string]map[string]struct{})
	for _, i := range Index {
		for _, e := range el.ElemList {
			if i == e.ElemName || i == e.ElemType+":::"+e.ElemName {
				i = rename(elemContent, i)
				elemContent[i] = make(map[string]struct{})
				elemContent[i] = e.List
				break
			}
		}
	}
	return elemContent
}
func checkElemType(ElemName string, ruleName string) string {
	//优先级 special > rule > common > system
	tempName := make(map[string]bool)
	for _, n := range specialElem["Temp"] {
		k := strings.SplitN(n, ":", 2)[0]
		tempName[k] = true
	}
	ElemName = ElemName[1 : len(ElemName)-1]
	if SpecialElems[ElemName] == true {
		return SPECIAL_TYPE
	}
	if tempName[ElemName] {
		return EXIST_ELEMENT
	}
	fileName := ruleName[0:len(ruleName)-5] + "." + ElemName + ".txt"
	elemFileList, err := ioutil.ReadDir(ElemFilePath)
	ExitErr(err)
	for _, f := range elemFileList {
		if f.Name() == fileName {
			return ruleName
		}
	}
	fileName = ElemName + ".txt"
	for _, f := range elemFileList {
		if f.Name() == fileName {
			return COMMON_TYPE
		}
	}
	if systemElems[ElemName] {
		return SYSTEM_TYPE
	}
	return ""
}
