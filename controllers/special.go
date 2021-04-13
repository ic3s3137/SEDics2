package controllers

import "strings"

func makeDomainElems(input []string)map[string]struct{}{
	var list = make(map[string]struct{})
	for _,e := range input{
		//type1
		list[e] = struct{}{}
		//type2
		field := strings.Split(e,".")
		for _,s := range field[:len(field)-1]{
			list[s] = struct{}{}
		}
		//type3
		list[strings.ReplaceAll(e,".","_")] = struct{}{}
		//list[strings.ReplaceAll(e,".","-")] = struct{}{}
		list[strings.ReplaceAll(e,".","")] = struct{}{}
		list[field[len(field)-2]+"."+field[len(field)-1]] = struct{}{}
		list[field[len(field)-2]+field[len(field)-1]] = struct{}{}
		list[field[len(field)-2]+"_"+field[len(field)-1]] = struct{}{}
		if len(field) > 2{
			list[field[len(field)-3]+"."+field[len(field)-2]] = struct{}{}
			list[field[len(field)-3]+field[len(field)-2]] = struct{}{}
			list[field[len(field)-3]+"_"+field[len(field)-2]] = struct{}{}
		}
	}
	return list
}

func makePinyinElems(elename string,seString []string)map[string]struct{}{
	//var list = make(map[string]struct{})
	var tmplist  = make(map[string]struct{})
	for _,seStr := range seString {
		tmplist[strings.ToUpper(seStr[0:1])+seStr[1:]] = struct{}{}
		tmplist[strings.ToUpper(seStr[0:1])+strings.ToLower(seStr[1:])] = struct{}{}
		tmplist[strings.ToLower(seStr)] = struct{}{}
		var (
			pinyinFirst []string
			pinyin      []string
			flag        string
			num         int
		)
		for i, _ := range seStr {
			n := seStr[i : i+1]
			if n == strings.ToUpper(n) {
				if flag != "" {
					pinyinFirst = append(pinyinFirst, flag)
					pinyin = append(pinyin, seStr[num:i])
					num = i
				}
				flag = n
			}
			if i == len(seStr)-1 && flag != "" {
				pinyinFirst = append(pinyinFirst, flag)
				pinyin = append(pinyin, seStr[num:i+1])
			}
		}
		if len(pinyinFirst) != 0 {
			var t string
			for _, k := range pinyinFirst {
				t = t + k
			}
			tmplist[t] = struct{}{}
			tmplist[strings.ToUpper(t[0:1])+strings.ToLower(t[1:])] = struct{}{}
			tmplist[strings.ToLower(t)] = struct{}{}
			if elename == "Pinyin(name)"{
				tmplist[pinyin[0]] = struct{}{}
				tmplist[strings.ToLower(pinyin[0])] = struct{}{}
				if len(pinyinFirst) == 2{
					tmplist[pinyin[1]] = struct{}{}
					tmplist[strings.ToLower(pinyin[1])] = struct{}{}
				}else{
					tmplist[pinyin[len(pinyin)-2]+pinyin[len(pinyin)-1]] = struct{}{}
					tmplist[strings.ToLower(pinyin[len(pinyin)-2]+pinyin[len(pinyin)-1])] = struct{}{}
				}
			}

		}
	}
	for k := range tmplist {
		tmplist[k] = struct{}{}
	}
	return tmplist
}