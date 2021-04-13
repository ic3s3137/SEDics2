package controllers

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
)

type rule struct{
	Index []string
	RuleName string
	Rule string
}
type ruleList struct{
	RuleList []rule
	RuleMap map[string]map[string]bool
	Level string
}

func (r *rule)equal(v []string)bool{
	for _,i := range r.Index{
		flag := false
		for _,e := range v{
			if i == e{
				flag = true
				break
			}
		}
		if flag == false{
			return false
		}
	}
	return true
}
func (r *rule)contain(v []string)bool{
	for _,a := range v{
		flag := false
		for _,b := range r.Index{
			if a == b || r.RuleName+":::"+b == a{
				flag = true
				break
			}
		}
		if flag == false{
			return false
		}
	}
	return true
}

func (rl *ruleList)readRuleFiles(files []string){
	if len(rl.RuleMap) == 0{
		rl.RuleMap = make(map[string]map[string]bool)
	}
	for _,f := range files{
		file := RuleFilePath+"/"+f
		content,err := ioutil.ReadFile(file)
		ExitErr(err)
		strContent := strings.Split(string(content),levelSpilt[rl.Level])[0]
		lines := strings.Fields(strContent)
		if len(rl.RuleMap[f]) == 0{
			rl.RuleMap[f] = make(map[string]bool)
		}
		for _,line := range lines{
			line = strings.TrimSpace(line)
			if levelTag[line] == false && line != ""{
				rl.RuleMap[f][line] = true
			}
		}
	}
}

func GetRuleFiles()map[string]bool{
	filelist := make(map[string]bool)
	dirInfo,err := ioutil.ReadDir(RuleFilePath)
	ExitErr(err)
	for i := range dirInfo{
		filelist[dirInfo[i].Name()] = false
	}
	return filelist
}

func (rl *ruleList)makeTree(elemMap map[string]map[string]bool)[]elemIndex{
	var treeIndex []elemIndex
	grepElemVars := regexp.MustCompile(grepElemReg)
	for ruleName,ruleList := range rl.RuleMap{
		for r := range ruleList{
			elemVal := grepElemVars.FindAllString(r,-1)
			flag := true
			for i,eName := range elemVal{
				if elemMap["COMMON"][eName] == false && elemMap["SPECIAL"][eName] == false && elemMap[ruleName][eName] == false && elemMap["SYSTEM"][eName] == false{
					flag = false
					break
				}else if elemMap[ruleName][eName] == true{
					elemVal[i] = ruleName+":::"+eName
				}
			}
			if flag == true{
				rl.RuleList = append(rl.RuleList,rule{RuleName:ruleName,Rule: r,Index: elemVal})
			}
		}
	}
	for _,ruleList := range rl.RuleList{
		var allElemName [][]string
		for _,i := range ruleList.Index{
			allElemName = append(allElemName,[]string{i})
		}
		index := combina(ruleList.Index,allElemName,len(allElemName))
		for _,i := range index{
			flag := false
			for n,t := range treeIndex{
				if equal(t.IndexName,i){
					flag = true
					treeIndex[n].Count ++
					break
				}
			}
			if flag == false{
				treeIndex = append(treeIndex,elemIndex{IndexName: i,Count: 1})
			}
		}
	}
	//建立生成树
	tree := buildTree(treeIndex)
	fmt.Println("索引树 >>>",len(tree),tree)
	return tree
}
/*
a b c d
ab ac ad bc bd cd
abc abd acd
 */
func combina(rigth []string,left [][]string,length int)[][]string{
	if length == 1{
		return [][]string{}
	}
	var newlist [][]string
	for _,l := range left{
		 n := exclude(rigth,l)
		 for _,n1 := range n{
			 e := append(l,n1)
			 //去除重复的元素
			 if groupContain(newlist,e) == false{
				 newlist = append(newlist,e)
			 }
		 }
	}
	return append(newlist,combina(rigth,newlist,length-1)...)
}
func exclude(big []string,small []string)[]string{
	var tmp1 = make([]string,len(big))
	var tmp2 = make([]string,len(small))
	copy(tmp1,big)
	copy(tmp2,small)
	for _,s := range small{
		for ib,b := range tmp1{
			if s == b{
				tmp1 = append(tmp1[0:ib],tmp1[ib+1:]...)
				break
			}
		}
	}
	return tmp1
}
func groupContain(big [][]string,e []string)bool{
	if len(big) == 0{
		return false
	}
	for _,b := range big{
		if equal(b,e){
			return true
		}
	}
	return false
}
func contain(e []string,v []string)bool{
	//e包含v
	var tmp1 = make([]string,len(e))
	var tmp2 =make([]string,len(v))
	copy(tmp1,e)
	copy(tmp2,v)
	for true{
		if len(tmp2) == 0{
			return true
		}
		for i2,a2 := range tmp2{
			flag := false
			for i1,a1 := range tmp1{
				if a1 == a2 {
					flag = true
					tmp1 = append(tmp1[0:i1],tmp1[i1+1:]...)
					tmp2 = append(tmp2[0:i2],tmp2[i2+1:]...)
					break
				}
			}
			if flag == false{
				return false
			}
			break
		}

	}
	return false
}
func equal(e1 []string,e2 []string)bool{
	if len(e1) != len(e2){
		return false
	}
	var tmp1 = make([]string,len(e1))
	var tmp2 = make([]string,len(e2))
	copy(tmp1,e1)
	copy(tmp2,e2)
	for true{
		if len(tmp1) == 0 && len(tmp2) == 0{
			return true
		}
		for i1,a1 := range tmp1{
			flag := false
			for i2,a2 := range tmp2{
				if a1 == a2 {
					flag = true
					tmp1 = append(tmp1[0:i1],tmp1[i1+1:]...)
					tmp2 = append(tmp2[0:i2],tmp2[i2+1:]...)
					break
				}
			}
			if flag == false{
				return false
			}
			break
		}

	}
	return false
}

func buildTree(newIndex []elemIndex)[]elemIndex{
	tree := []elemIndex{}
	var (
		max = elemIndex{Count: 0,IndexName: []string{}}
		length = -1
	)
	for true{
		if len(newIndex) == 0{
			break
		}
		if len(max.IndexName) == 0{
			 for _,n := range newIndex{
				  if n.Count > max.Count{
					  max = elemIndex{IndexName: n.IndexName,Count: n.Count}
				  }
			 }
		}else{
			 for _,n := range newIndex{
				 if n.Count > max.Count && contain(n.IndexName,max.IndexName){
					 max = elemIndex{IndexName: n.IndexName,Count: n.Count}
				 }
			 }
		}
		if max.Count == 0{
			if length == 0{
				max = elemIndex{Count: 0,IndexName: []string{}}
			}else{
				max = elemIndex{Count: 0,IndexName: tree[length-1].IndexName}
				length --
			}
			continue
		}
		//fmt.Println("Max:",max)
		tree = append(tree,max)
		length ++
		//fmt.Println("=====",tree)
		for i,n := range newIndex{
			 if equal(n.IndexName,max.IndexName){
				  newIndex = append(newIndex[0:i],newIndex[i+1:]...)
				  break
			 }
		}
		max.Count = 0
	}
	return tree
}
