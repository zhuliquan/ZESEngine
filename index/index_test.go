package index

import (
	"fmt"
	"testing"
)

func TestIndex(t *testing.T) {
	ii := NewInvertIndex()
	ii.BuildForwardIndex()
	ii.BuildInvertIndex()
	fmt.Println(len(ii.Search("我们")))
	for _, v := range ii.Search("我们") {
		v.Show()
	}
	ii.Close()
}
