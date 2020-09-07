package main

import (
	"ZESEngine/index"
	"fmt"
)

func main() {
	ii := index.NewInvertIndex()
	//ii.BuildForwardIndex()
	//ii.BuildInvertIndex()
	fmt.Println("###############################开始你的操作#####################################################")
	for {
		var word string
		fmt.Println("请输入你要查询的词: (结束按Q)")
		fmt.Scan(&word)
		fmt.Scanf("%s", &word)
		if word == "Q" {
			break
		} else {
			results := ii.Search(word)
			if len(results) == 0 {
				fmt.Println("没有找到要的文档")
			} else {
				for _, v := range results {
					v.Show()
				}
			}
		}
	}
	//ii.Close()

}
