package main

import (
	"SEDics2/controllers"
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	ifList := flag.Bool("L", false, "show rule list")
	Rule := flag.String("r", "", "rule file")
	SEStrings := flag.String("s", "", "string type element")
	SENumber := flag.String("n", "", "number type element")
	PinYin := flag.String("p", "", "Pinyin")
	SavePath := flag.String("o", "", "saved path")
	TempVars := flag.String("t", "", "temp variables")
	Level := flag.Int("l", 1, "dictionary length level:1,2,3")
	flag.Parse()
	if *ifList {
		rulelist, _ := controllers.SortMap(controllers.GetRuleFiles())
		for _, f := range rulelist {
			fmt.Println(f)
		}
		os.Exit(1)
	}
	if *Rule == "" {
		fmt.Println("必须指定规则名:-r <rule>")
		os.Exit(1)
	}
	controllers.SetRuleFile(strings.Split(*Rule, ","))
	if *SavePath == "" {
		fmt.Println("必须指定字典保持路径:-o <filepath>")
		os.Exit(1)
	}
	controllers.SaveFile = *SavePath
	//设置规则名
	var specialList = make(map[string][]string)
	if *SEStrings != "" {
		controllers.SetStringElem(strings.Split(*SEStrings, ","))
	}
	if *SENumber != "" {
		controllers.SetNumberElem(strings.Split(*SENumber, ","))
	}
	if *PinYin != "" {
		specialList["Pinyin"] = strings.Split(*PinYin, ",")
	}
	if *TempVars != "" {
		specialList["Temp"] = strings.Split(*TempVars, ",")
	}
	controllers.SetSpecialElem(specialList)
	controllers.SetLevel(*Level)
	controllers.Start()

}
