package controllers

import (
	"log"
	"sort"
	"strconv"
)

var (
	dictionary = make(map[string]struct{})
	stringElem = make(map[string]struct{})
	numberElem = make(map[string]struct{})
	specialElem = make(map[string][]string)
	logger chan string
	level string
	ruleFile []string
)
func SetRuleFile(input map[string]bool){
	ruleFile = []string{}
	for name := range input{
		if input[name] == true{
			ruleFile = append(ruleFile,name)
		}
	}
}
func SetStringElem(input map[string]struct{}){
	stringElem = input
}
func SetNumberElem(input map[string]struct{}){
	numberElem = input
}
func SetSpecialElem(input map[string][]string){
	specialElem = input
}
func SetLevel(l string){
	level = l
}
func SetLogger(l *chan string){
	logger = *l
}
func ExitErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
func SortMap(a map[string]bool)([]string,[]bool){
	var (
		key []string
		value []bool
	)
	for k := range a{
		key = append(key,k)
	}
	sort.Strings(key)
	for _,k := range key{
		value = append(value,a[k])
	}
	return key,value
}
func loginfo(info string){
	logger <- info
}
func Start(mode string){
	rulelist := ruleList{}
	elemlist := elemList{}
	rulelist.Level = level
	elemlist.Level = level
	if mode == "RuleMode"{
		rulelist.readRuleFiles(ruleFile)
		elemlist.makeElemMap(rulelist.RuleMap)
		elemlist.AddElem("SPECIAL",seStringElem,stringElem)
		elemlist.AddElem("SPECIAL",seNumberElem,numberElem)
		tree := rulelist.makeTree(elemlist.ElemMap)
		generate(tree,elemlist,rulelist.RuleList)
		loginfo("total lines:"+strconv.Itoa(len(dictionary)))
		saveDics()
		loginfo("***********************************************************")
	}
}