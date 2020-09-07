/*
分词库
*/
package segement

import (
	"strings"

	"github.com/go-ego/gse"
)

// 可以做到分词
type WordCuter interface {
	WordCut(text *string) []string
}

// 可以做到分句字
type SentenceCutter interface {
	SentenceCut(text *string) ([]string, string)
}

// 具有英文分词
type EnglishSegmenter struct {
	wordSplitSymbol     string
	sentenceSplitSymbol []string
}

func (es *EnglishSegmenter) init() {
	es.wordSplitSymbol = ".,:;"
	es.sentenceSplitSymbol = []string{". ", "; ", "! ", "? "}
}

/* func (es *EnglishSegmenter) checkSeq(s string) bool { */
//     for i := 0; i < len(es.sentenceSplitSymbol); i++ {
//         if (s == )
//     }
// }
//
// func (es *EnglishSegmenter) wordCut(text *string) []string{
//     bs := []byte(*text)
//     words := make([]string, 1)
//     i, j, n := 0, 0, len(bs)
//     for j < n {
//
//     }
//
// }
/*  */
/* func (es *EnglishSegmenter) sentenceCut(text *string) ([]string, string) { */
//
// }

type ChineseSegementer struct {
	seg                 gse.Segmenter
	sentenceSplitSymbol string
}

func NewChineseSegementer() *ChineseSegementer {
	cs := new(ChineseSegementer)
	cs.seg = gse.New("zh")
	cs.sentenceSplitSymbol = "。？！：；"
	return cs
}

func (cs *ChineseSegementer) WordCut(text *string) []string {
	return cs.seg.CutAll(*text)
}

func (cs *ChineseSegementer) SentenceCut(text *string) ([]string, string) {
	bs := []rune(*text)
	sentences := make([]string, 0)
	i, j, n := 0, 0, len(bs)
	for j < n {
		if strings.ContainsRune(cs.sentenceSplitSymbol, bs[j]) {
			if i < j {
				sentences = append(sentences, string(bs[i:j]))
			}
			i = j + 1
		}
		j += 1
	}
	if i < j {
		return sentences, string(bs[i:])
	} else {
		return sentences, ""
	}
}
