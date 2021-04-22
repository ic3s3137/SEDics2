package controllers

import (
	"fmt"
	"log"
	"sort"
	"strconv"
)

var (
	dictionary  = make(map[string]struct{})
	stringElem  = make(map[string]struct{})
	numberElem  = make(map[string]struct{})
	specialElem = make(map[string][]string)
	LevelString = map[int]string{1: "small", 2: "medium", 3: "big"}
	SaveFile    string
	level       string
	ruleFile    []string
)

func SetRuleFile(input []string) {
	for _, v := range unique(input) {
		ruleFile = append(ruleFile, v+".rule")
	}
}
func SetStringElem(input []string) {
	stringElem = string2map(input)
}
func SetNumberElem(input []string) {
	numberElem = string2map(input)
}
func SetSpecialElem(input map[string][]string) {
	specialElem = input
}
func SetLevel(l int) {
	level = LevelString[l]
	loginfo("字典大小:" + level)
}

func ExitErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
func SortMap(a map[string]bool) ([]string, []bool) {
	var (
		key   []string
		value []bool
	)
	for k := range a {
		key = append(key, k)
	}
	sort.Strings(key)
	for _, k := range key {
		value = append(value, a[k])
	}
	return key, value
}
func loginfo(info string) {
	fmt.Println(info)
}
func Start() {
	rulelist := ruleList{}
	elemlist := elemList{}
	rulelist.Level = level
	elemlist.Level = level
	rulelist.readRuleFiles(ruleFile)
	elemlist.makeElemMap(rulelist.RuleMap)
	elemlist.AddElem("SPECIAL", seStringElem, stringElem)
	elemlist.AddElem("SPECIAL", seNumberElem, numberElem)
	tree := rulelist.makeTree(elemlist.ElemMap)
	generate(tree, elemlist, rulelist.RuleList)
	loginfo("字典条数:" + strconv.Itoa(len(dictionary)))
	saveDics()
}
func unique(input []string) []string {
	var tmp = make(map[string]struct{})
	var news []string
	for _, i := range input {
		if i != "" {
			tmp[i] = struct{}{}
		}
	}
	for k := range tmp {
		news = append(news, k)
	}
	return news
}
func string2map(input []string) map[string]struct{} {
	var news = make(map[string]struct{})
	for _, i := range input {
		news[i] = struct{}{}
	}
	return news
}
