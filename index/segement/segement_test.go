package segement

import (
	"fmt"
	"strings"
	"testing"
)

func Equal(as, bs []string) bool {
	if len(as) != len(bs) {
		return false
	} else {
		for i := 0; i < len(as); i++ {
			if as[i] != bs[i] {
				return false
			}
		}
		return true
	}
}

func splitFunc(r rune) bool {
	return strings.ContainsRune(",-.%", r)
}

func TestSegmenter(t *testing.T) {
	s := "我们是一体的"
	cs := NewChineseSegementer()
	fmt.Println(cs.WordCut(&s))
	s = "我们是什么，我们就近是真么。我峨嵋你www。dsad"
	fmt.Println(cs.SentenceCut(&s))
}
