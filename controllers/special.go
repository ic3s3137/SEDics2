package controllers

import (
	"github.com/go-ego/gse"
	"github.com/go-ego/gse/hmm/pos"
	"github.com/mozillazg/go-pinyin"
	"strings"
)

var (
	seg    gse.Segmenter
	posSeg pos.Segmenter
)

type MPinyin struct {
	First      []string
	PinyinList []string
	Pinyin     string
	IsName     bool
	Flag       string
	py         pinyin.Args
	Num        int
}

func (p *MPinyin) Init(hans string, ChineseType []gse.SegPos) {
	p.py = pinyin.NewArgs()
	for _, w := range pinyin.Pinyin(hans, p.py) {
		p.PinyinList = append(p.PinyinList, w[0])
		//p.PinyinList = append(p.PyinyinList)
	}
	for _, w := range p.PinyinList {
		p.First = append(p.First, string(w[0]))
		p.Pinyin = p.Pinyin + w
	}
	for _, i := range ChineseType {
		if i.Pos == "nr" && i.Text == hans {
			p.IsName = true
		}
	}
}

func makePinyinElems(input []string) map[string]struct{} {
	var tmplist = make(map[string]struct{})
	seg.SkipLog = true
	seg.LoadDict()
	for _, i := range input {
		ChineseList, ChineseType := cutPinyin(i)
		if len(ChineseList) > 1 {
			ChineseList = append(ChineseList, i)
		}
		for _, Hans := range ChineseList {
			p := MPinyin{}
			p.Init(Hans, ChineseType)
			tmplist[strings.ToUpper(p.Pinyin[0:1])+p.Pinyin[1:]] = struct{}{}
			tmplist[strings.ToUpper(p.Pinyin[0:1])+strings.ToLower(p.Pinyin[1:])] = struct{}{}
			tmplist[strings.ToLower(p.Pinyin)] = struct{}{}
			var t string
			for _, k := range p.First {
				t = t + k
			}
			tmplist[t] = struct{}{}
			tmplist[strings.ToUpper(t[0:1])+strings.ToLower(t[1:])] = struct{}{}
			tmplist[strings.ToLower(t)] = struct{}{}
			if p.IsName {
				tmplist[p.PinyinList[0]] = struct{}{}
				tmplist[strings.ToLower(p.PinyinList[0])] = struct{}{}
				if len(p.First) == 2 {
					tmplist[p.PinyinList[1]] = struct{}{}
					tmplist[strings.ToLower(p.PinyinList[1])] = struct{}{}
				} else {
					tmplist[p.PinyinList[len(p.PinyinList)-2]+p.PinyinList[len(p.PinyinList)-1]] = struct{}{}
					tmplist[strings.ToLower(p.PinyinList[len(p.PinyinList)-2]+p.PinyinList[len(p.PinyinList)-1])] = struct{}{}
				}
			}

		}
	}
	return tmplist
}
func makeTempElems(input []string) map[string]struct{} {
	var news = make(map[string]struct{})
	for _, v := range strings.Split(input[0], "|") {
		news[v] = struct{}{}
	}
	return news
}
func cutPinyin(text string) ([]string, []gse.SegPos) {
	hmm := seg.Cut(text, true)
	posSeg.WithGse(seg)
	po := posSeg.Cut(text, true)
	return hmm, po
}
