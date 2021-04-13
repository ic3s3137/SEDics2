package controllers

const (
	RuleFilePath = "./rule"
	ElemFilePath = "./element"
	commonElem = "COMMON"
	seStringElem = "$String$"
	seNumberElem = "$Number$"
	grepElemReg = "\\$(.+?)\\$"

)
var (
	levelTag = map[string]bool{"[small]":true,"[medium]":true,"[big]":true}
	levelSpilt = map[string]string{
		"small":"[medium]",
		"medium":"[big]",
		"big":"[laslksjflksefjlfsasdf2312]",
	}
	systemElems = map[string]bool{}
	SpecialElems = map[string]bool{"Domain":true,"FileExtension":true,"Pinyin(name)":true,"Pinyin(other)":true}
)
