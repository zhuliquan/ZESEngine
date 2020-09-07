package db

import (
	"testing"
)

func TestDB(t *testing.T) {
	EngineDB := NewIndexDB()
	EngineDB.BuildTables()
	//EngineDB.DestroyTables() // TODO 找个时候把这个运行一下
	// i, _ := EngineDB.AddKeyWord("dsada")
	// i, _ = EngineDB.AddKeyWord("xxxxx")
	/* fmt.P */
	/*  words := make([]string, 0) */
	// ids := make([]int, 0)
	// n := 20000
	// for i := 0; i < n; i++ {
	//     word := make([]byte, 100)
	//     for j := 0; j < 100; j++ {
	//         word[j] = 23
	//     }
	//     ids = append(ids, i)
	//     words = append(words, strconv.Itoa(i)+string(word))
	// }
	// // 非常耗时读入操作
	// if result := EngineDB.ADDKeyWords(ids, words); !result.success {
	//     t.Fail()
	// } else {
	//     fmt.Printf("success 导入%d条数据\n", n)
	// }
	// //查看读的能力
	// result := EngineDB.QueryRealDocIdByFilePaths(words)
	//
	// if err := EngineDB.Close(); err != nil {
	//     t.Fail()
	/* } */

}
