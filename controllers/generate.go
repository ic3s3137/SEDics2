package controllers

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type cache struct{
	Index []string
	List []map[string]string
}
var Counter = 1
func generate(tree []elemIndex,elemlist elemList,rulelist []rule)map[string]struct{}{
	var elemCache []cache
	for _,t := range tree {
		for true {
			//fmt.Println(t.IndexName)
			var (
				ele cache
				contained = false
				sRule  rule
				isEqual = false
			)
			for i,r := range rulelist{
				if equal(r.Index,t.IndexName){
					contained = true
					isEqual = true
					sRule = rulelist[i]
					rulelist = append(rulelist[:i],rulelist[i+1:]...)
					break
				}else if contained == true || contain(r.Index,t.IndexName){
					contained = true
				}
			}
			if contained == true {
				//*********算法1*********
				//ele = cache{Index:t.IndexName,List:makeCache(elemlist.GetElemByIndex(t.IndexName),make([]map[string]string,0))}
				//*********算法2**********
				//if len(elemCache) > 0 && contain(t.IndexName, elemCache[len(elemCache)-1].Index) {
				//	ele = cache{Index: t.IndexName, List: makeCache(elemlist.GetElemByIndex(t.IndexName), elemCache[len(elemCache)-1].List)}
				//} else{
				//	ele = cache{Index: t.IndexName, List: makeCache(elemlist.GetElemByIndex(t.IndexName), make([]map[string]string, 0))}
				//}
				//*******算法3*******
				for true {
					if len(elemCache) > 0 && contain(t.IndexName, elemCache[len(elemCache)-1].Index) && !equal(t.IndexName, elemCache[len(elemCache)-1].Index) {
						ele = cache{Index: t.IndexName, List: makeCache(elemlist.GetElemByIndex(t.IndexName), elemCache[len(elemCache)-1].List)}
						elemCache = append(elemCache, ele)
						break
					} else if len(elemCache) == 0 {
						ele = cache{Index: t.IndexName, List: makeCache(elemlist.GetElemByIndex(t.IndexName), make([]map[string]string, 0))}
						elemCache = append(elemCache, ele)
						break
					}else if equal(t.IndexName, elemCache[len(elemCache)-1].Index){
						break
					}else {
						elemCache = elemCache[:len(elemCache)-1]
					}
				}
				//fmt.Println(ele.Index, len(ele.List))
			}
			if isEqual == true{
				//生成字典
				generateDics(sRule,elemCache[len(elemCache)-1])
				//fmt.Println(sRule.Index,"====",elemCache[len(elemCache)-1].Index)
			}else{
				break
			}
		}
	}
	for _,r := range rulelist{
		l := elemlist.GetElemByIndex(r.Index)
		index := r.Index[0]
		if strings.Index(index,r.RuleName+":::") == 0{
			index = strings.Split(index,r.RuleName+":::")[1]
		}
		for e := range l[r.Index[0]]{
			dictionary[strings.ReplaceAll(r.Rule,index,e)] = struct{}{}
		}
	}
	fmt.Println("运算次数",Counter)

	return nil
}
func makeCache(el map[string]map[string]struct{},previous []map[string]string)[]map[string]string{
	var (
	left []map[string]string
	)

	left = previous
	if len(left) == 0{
		for t := range el{
			for t1 := range el[t]{
				left = append(left,map[string]string{t:t1})
			}
			delete(el,t)
			break
		}
	}else{
		for t := range left[0]{
			delete(el,t)
		}
	}
	for true{
		if len(el) == 0{
			break
		}
		var right []map[string]string
		for t := range el{
			for t1 := range el[t]{
				right = append(right,map[string]string{t:t1})
			}
			delete(el,t)
			break
		}
		//fmt.Println("rigth>>",len(right),right)
		//fmt.Println("el>>>",el)
		var newContainer []map[string]string
		var sub map[string]string
		for _,r := range right{
			for _,l := range left{
				Counter ++
				sub = make(map[string]string)
				for n := range r{
					//n = rename(sub,n)
					sub[n] = r[n]
				}
				for n2 := range l{
					//n2 = rename(sub,n2)
					sub[n2] = l[n2]
				}
				newContainer = append(newContainer,sub)
			}
		}
		left = newContainer
	}
	//fmt.Println(len(left))//,left)
	return left

}
func rename(sub map[string]map[string]struct{},n string)string{
	var name = regexp.MustCompile("\\$"+n[1:len(n)-1]+"\\$")
	var allk string
	for k := range sub{
		allk = allk + k
	}
	c := len(name.FindAllString(allk,-1))
	if c == 0{
		return n
	}else{
		return n+strconv.Itoa(c)
	}
}
func generateDics(r rule,c cache){
	var (
		nameAlias = make(map[string]string)
		formatRule string
	)
	//fmt.Println(r,"===========",c)
	formatRule = r.Rule
	counter := make(map[string]int)
	for _,i := range r.Index{
		if strings.Index(i,r.RuleName+":::") == 0{
			i = strings.Split(i,r.RuleName+":::")[1]
		}
		grepName := regexp.MustCompile(grepElemReg)
		namelist := grepName.FindAllString(formatRule,-1)
		prefix := ""
		suffix := formatRule
		for _,n := range namelist{
			if n != i{
				prefix = prefix + strings.Split(suffix,n)[0] + n
				suffix = strings.Split(suffix,n)[1]
			}else{
				counter[i] ++
				if counter[i] == 1{
					basestr := "$"+base64.StdEncoding.EncodeToString([]byte(i[1:len(i)-1]))+"$"
					formatRule = prefix + strings.Replace(suffix,i,basestr,1)
					nameAlias[i] = basestr
				}else{
					basestr := "$"+base64.StdEncoding.EncodeToString([]byte(i[1:len(i)-1]+strconv.Itoa(counter[i]-1)))+"$"
					formatRule = prefix + strings.Replace(suffix,i,basestr,1)
					nameAlias[i+strconv.Itoa(counter[i]-1)] = basestr
				}
				break
			}
		}
	}
	for _,l := range c.List{
		one := formatRule
		for n := range l{
			var n1 string
			if strings.Index(n,r.RuleName+":::") == 0{
				n1 = strings.Split(n,r.RuleName+":::")[1]
			}else{
				n1 = n
			}
			one = strings.Replace(one,nameAlias[n1],l[n],1)
		}
		//fmt.Println(formatRule,one,nameAlias)
		dictionary[one] = struct{}{}
	}
}
func saveDics(){
	content := ""
	for d := range dictionary{
		content = content + d + "\n"
	}
	t := time.Now()
	pathname := fmt.Sprintf("%d%d%d%d%d%d.txt", t.Year(),t.Month(),t.Day(),t.Hour(),t.Minute(),t.Second())
	file:= filepath.Dir(os.Args[0])
	prefix ,err:= filepath.Abs(file)
	ExitErr(err)
	loginfo("[i]Save to "+prefix+"\\"+pathname)
	err = ioutil.WriteFile(pathname,[]byte(content),666)
	ExitErr(err)
}