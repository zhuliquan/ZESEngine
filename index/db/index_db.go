package db

import (
	"ZESEngine/index/post_item"
	"ZESEngine/setting"
	"database/sql"
	"fmt"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

type IndexDB struct {
	db        *sql.DB
	operators map[uint32]DBOperateFunc // 数据库基础操作集合
}

func NewIndexDB() *IndexDB {
	index_db := new(IndexDB)
	if err := index_db.Connect(); err != nil {
		panic(err)
	}
	index_db.operators = map[uint32]DBOperateFunc{
		INSERT: SQLInsert,
		DELETE: SQLDelete,
		UPDATE: SQLUpdate,
		SELECT: SQLSelect,
	}
	index_db.BuildTables() // 建立表
	fmt.Println("数据库建立成功")
	return index_db
}

func (this *IndexDB) Connect() error {
	command := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", setting.MYSQL_USERNAME, setting.MYSQL_PASSWORD, setting.MYSQL_HOST, setting.MYSQL_PORT, setting.MYSQL_DATABASE)
	if db, db_err := sql.Open("mysql", command); db_err != nil {
		return db_err
	} else {
		fmt.Println("数据库连接成功")
		this.db = db
		return nil
	}
}

func (this *IndexDB) Close() error {
	if err := this.db.Close(); err != nil {
		return err
	} else {
		fmt.Println("数据库关闭")
		return nil
	}
}

// 创建所有的表
// TODO 以后再将表的id 改为自增变量
func (this *IndexDB) BuildTables() {
	fmt.Println("建表开始")
	this.db.Exec(fmt.Sprintf(
		`create table if not exists %s (
			realDocId int primary key auto_increment,
			filePath varchar(1000) unique,
			avaliable bit(1) default 1)`, setting.INDEX_TABLE_REAL_DOC))
	this.db.Exec(fmt.Sprintf(
		`create table if not exists %s (
			docId int primary key auto_increment,
			content varchar(10000),
			avaliable bit(1) default 1,
			realDocId int,
			foreign key(realDocId) references %s(realDocId))`, setting.INDEX_TABLE_FORWARD_INDEX, setting.INDEX_TABLE_REAL_DOC))
	this.db.Exec(fmt.Sprintf(
		`create table if not exists %s (
			wordId int primary key ,
			word varchar(1000) unique,
			avaliable bit(1) default 1)`, setting.INDEX_TABLE_KEY_WORD))
	this.db.Exec(fmt.Sprintf(
		`create table if not exists %s (
			postItemId int primary key ,
			docId int,
			wordId int, 
			termFrequence int default 0, -- 词频
			wordWeight float default 0.0, -- 单词权重
			avaliable bit(1) default 1,
			foreign key(wordId) references %s(wordId),
			foreign key(docId) references %s(docId))`, setting.INDEX_TABLE_POST_ITEM, setting.INDEX_TABLE_KEY_WORD, setting.INDEX_TABLE_FORWARD_INDEX))
	fmt.Println("建表结束")
}

// 销毁所有的表
func (this *IndexDB) DestroyTables() {
	fmt.Println("清空所有的表")
	this.db.Exec(fmt.Sprintf("drop table %s", setting.INDEX_TABLE_POST_ITEM))
	this.db.Exec(fmt.Sprintf("drop table %s", setting.INDEX_TABLE_KEY_WORD))
	this.db.Exec(fmt.Sprintf("drop table %s", setting.INDEX_TABLE_FORWARD_INDEX))
	this.db.Exec(fmt.Sprintf("drop table %s", setting.INDEX_TABLE_REAL_DOC))
	fmt.Println("已经清空所有的表")
}

// 数据操作相关的函数
type ProcessOneFunc func(tx *sql.Tx) *SQLResult

func ProcessOne(indexDB *IndexDB, processFunc ProcessOneFunc) *SQLResult {
	tx, _ := indexDB.db.Begin()
	result := processFunc(tx)
	if !result.success {
		tx.Rollback()
		return NewSQLResult(nil, false)
	} else {
		tx.Commit()
		return result
	}
}

func (this *IndexDB) DeleteTable(tableName string) *SQLResult {
	return ProcessOne(this, func(tx *sql.Tx) *SQLResult {
		return this.operators[DELETE](tx, NewSQLInfo(tableName, nil, nil, nil, nil, nil, ""))
	})
}

func (this *IndexDB) AddKeyWord(wordId string, word string) *SQLResult {
	return ProcessOne(this, func(tx *sql.Tx) *SQLResult {
		return this.operators[INSERT](tx, NewSQLInfo(
			setting.INDEX_TABLE_KEY_WORD, []string{"wordId", "word"}, [][]string{[]string{wordId, fmt.Sprintf("'%s'", word)}},
			nil, nil, nil, "",
		))
	})
}

/* func (this *IndexDB) AddKeyWord(word string) (int, bool) { */
// ProcessOne(this, func(tx *sql.Tx) *SQLResult {
//     return this.operators[INSERT](tx, NewSQLInfo(
//         setting.INDEX_TABLE_KEY_WORD, []string{"word"}, [][]string{[]string{fmt.Sprintf("'%s'", word)}},
//         nil, nil, nil, "",
//     ))
// })
// r, _ := this.db.Query(fmt.Sprintf("select max(wordId) from %s", setting.INDEX_TABLE_KEY_WORD))
// tmp := 89
// r.Next()
// r.Scan(&tmp)
// return tmp, true
/* } */

func (this *IndexDB) ADDKeyWords(wordIds []int, words []string) *SQLResult {
	for i, n := uint32(0), uint32(len(words)); i < n; i += setting.MAX_INSERT_NUMBER {
		tx, _ := this.db.Begin()
		values := make([][]string, 0)
		for j := i; j < n && j < i+setting.MAX_INSERT_NUMBER; j++ {
			values = append(values, []string{strconv.Itoa(wordIds[j]), fmt.Sprintf("'%s'", words[j])})
		}
		result := this.operators[INSERT](tx, NewSQLInfo(
			setting.INDEX_TABLE_KEY_WORD,
			[]string{"wordId", "word"}, values,
			nil, nil, nil, "",
		))
		if !result.success {
			// 回滚
			tx.Rollback()
			return NewSQLResult(nil, false)
		}
		tx.Commit()
	}
	return NewSQLResult(nil, true)

}

func (this *IndexDB) RemoveAllKeyWords() *SQLResult {
	return this.DeleteTable(setting.INDEX_TABLE_KEY_WORD)
}
func (this *IndexDB) AddRealDoc(filePath string) *SQLResult {
	return ProcessOne(this, func(tx *sql.Tx) *SQLResult {
		return this.operators[INSERT](tx, NewSQLInfo(
			setting.INDEX_TABLE_REAL_DOC,
			[]string{"filePath"}, [][]string{[]string{fmt.Sprintf("'%s'", filePath)}},
			nil, nil, nil, "",
		))
	})

}

/* func (this *IndexDB) AddRealDoc(filePath string) (int, bool) { */
// ProcessOne(this, func(tx *sql.Tx) *SQLResult {
//     return this.operators[INSERT](tx, NewSQLInfo(
//         setting.INDEX_TABLE_REAL_DOC,
//         []string{"filePath"}, [][]string{[]string{fmt.Sprintf("'%s'", filePath)}},
//         nil, nil, nil, "",
//     ))
// })
// r, _ := this.db.Query(fmt.Sprintf("select max(realDocId) from %s", setting.INDEX_TABLE_REAL_DOC))
// tmp := 89
// r.Next()
// r.Scan(&tmp)
// return tmp, true
//
/* } */

func (this *IndexDB) ADDRealDocs(words []string) *SQLResult {
	for i, n := uint32(0), uint32(len(words)); i < n; i += setting.MAX_INSERT_NUMBER {
		tx, _ := this.db.Begin()
		values := make([][]string, 0)
		for j := i; j < n && j < i+setting.MAX_INSERT_NUMBER; j++ {
			values = append(values, []string{fmt.Sprintf("'%s'", words[j])})
		}
		result := this.operators[INSERT](tx, NewSQLInfo(
			setting.INDEX_TABLE_REAL_DOC,
			[]string{"filePath"}, values,
			nil, nil, nil, "",
		))
		if !result.success {
			// 回滚
			tx.Rollback()
			return NewSQLResult(nil, false)
		}
		tx.Commit()
	}
	return NewSQLResult(nil, true)

}

func (this *IndexDB) RemoveAllRealDoc() *SQLResult {
	return this.DeleteTable(setting.INDEX_TABLE_REAL_DOC)
}
func (this *IndexDB) AddLogicDoc(content string, realDocId int) *SQLResult {
	return ProcessOne(this, func(tx *sql.Tx) *SQLResult {
		return this.operators[INSERT](tx, NewSQLInfo(
			setting.INDEX_TABLE_FORWARD_INDEX,
			[]string{"content", "realDocId"}, [][]string{[]string{fmt.Sprintf("'%s'", content), strconv.Itoa(realDocId)}},
			nil, nil, nil, "",
		))
	})
}

/* func (this *IndexDB) AddLogicDoc(content string, realDocId int) (int, bool) { */
// ProcessOne(this, func(tx *sql.Tx) *SQLResult {
//     return this.operators[INSERT](tx, NewSQLInfo(
//         setting.INDEX_TABLE_FORWARD_INDEX,
//         []string{"content", "realDocId"}, [][]string{[]string{fmt.Sprintf("'%s'", content), strconv.Itoa(realDocId)}},
//         nil, nil, nil, "",
//     ))
// })
// r, _ := this.db.Query(fmt.Sprintf("select max(docId) from %s", setting.INDEX_TABLE_FORWARD_INDEX))
// tmp := 89
// r.Next()
// r.Scan(&tmp)
// return tmp, true
/* } */

func (this *IndexDB) AddLogicDocs(contents []string, realDocId int) *SQLResult {
	for i, n := uint32(0), uint32(len(contents)); i < n; i += setting.MAX_INSERT_NUMBER {
		tx, _ := this.db.Begin()
		values := make([][]string, 0)
		for j := i; j < n && j < i+setting.MAX_INSERT_NUMBER; j++ {
			values = append(values, []string{fmt.Sprintf("'%s'", contents[j]), strconv.Itoa(realDocId)})
		}
		result := this.operators[INSERT](tx, NewSQLInfo(
			setting.INDEX_TABLE_FORWARD_INDEX,
			[]string{"content", "realDocId"}, values,
			nil, nil, nil, "",
		))

		if !result.success {
			// 回滚
			tx.Rollback()
			return NewSQLResult(nil, false)
		}
		tx.Commit()
	}
	return NewSQLResult(nil, true)

}

func (this *IndexDB) RemoveAllLogicDoc() *SQLResult {
	return this.DeleteTable(setting.INDEX_TABLE_FORWARD_INDEX)
}
func (this *IndexDB) AddPostItem(itemId int, docId int, wordId int) *SQLResult {
	return ProcessOne(this, func(tx *sql.Tx) *SQLResult {
		return this.operators[INSERT](tx, NewSQLInfo(
			setting.INDEX_TABLE_POST_ITEM,
			[]string{"postItemId", "docId", "wordId"}, [][]string{[]string{strconv.Itoa(itemId), strconv.Itoa(docId), strconv.Itoa(wordId)}},
			nil, nil, nil, "",
		))
	})
}

/* func (this *IndexDB) AddPostItem(docId int, wordId int) (int, bool) { */
// ProcessOne(this, func(tx *sql.Tx) *SQLResult {
//     return this.operators[INSERT](tx, NewSQLInfo(
//         setting.INDEX_TABLE_POST_ITEM,
//         []string{"docId", "wordId"}, [][]string{[]string{strconv.Itoa(docId), strconv.Itoa(wordId)}},
//         nil, nil, nil, "",
//     ))
// })
// r, _ := this.db.Query(fmt.Sprintf("select max(docId) from %s", setting.INDEX_TABLE_POST_ITEM))
// tmp := 89
// r.Next()
// r.Scan(&tmp)
// return tmp, true
//
/* } */

func (this *IndexDB) AddPostItems(itemIds []int, docIds []int, wordIds []int) *SQLResult {
	for i, n := uint32(0), uint32(len(wordIds)); i < n; i += setting.MAX_INSERT_NUMBER {
		tx, _ := this.db.Begin()
		values := make([][]string, 0)
		for j := i; j < n && j < i+setting.MAX_INSERT_NUMBER; j++ {
			values = append(values, []string{strconv.Itoa(itemIds[j]), strconv.Itoa(docIds[j]), strconv.Itoa(wordIds[j])})
		}
		result := this.operators[INSERT](tx, NewSQLInfo(
			setting.INDEX_TABLE_POST_ITEM,
			[]string{"postItemId", "docId", "wordId"}, values,
			nil, nil, nil, "",
		))

		if !result.success {
			// 回滚
			tx.Rollback()
			return NewSQLResult(nil, false)
		}
		tx.Commit()
	}
	return NewSQLResult(nil, true)
}

// TODO 后续的这些都有问题，以后再去修改
func (this *IndexDB) QueryRealDocIdByFilePaths(filePaths []string) []int {
	results := make([]int, 0)
	for i, n := 0, len(filePaths); i < n; i++ {
		rows, err := this.db.Query(fmt.Sprintf("select realDocId from %s where filePath='%s'", setting.INDEX_TABLE_REAL_DOC, filePaths[i]))
		if err != nil {
			fmt.Println(fmt.Sprintf("select realDocId from %s where filePath='%s'", setting.INDEX_TABLE_REAL_DOC, filePaths[i]))
			fmt.Println(err.Error())
			//panic("存在错误")
			return nil
		}

		for rows.Next() {
			temp := 0
			rows.Scan(&temp)
			results = append(results, temp)
		}
	}

	return results
}

func (this *IndexDB) QueryFilePathsByRealDocId(docIds []int) []string {
	paths := make([]string, 0)
	for i, n := 0, len(docIds); i < n; i++ {
		rows, err := this.db.Query(fmt.Sprintf("select filePath from %s where realDocId=%d", setting.INDEX_TABLE_REAL_DOC, docIds[i]))
		if err != nil {
			fmt.Println(err.Error())
			//panic("存在错误")
			return nil
		}
		for rows.Next() {
			temp := ""
			rows.Scan(&temp)
			paths = append(paths, temp)
		}
	}
	return paths
}

func (this *IndexDB) QueryAllForwardIndx() []post_item.ForwardIndexItem {
	items := make([]post_item.ForwardIndexItem, 0)
	if rows, err := this.db.Query(fmt.Sprintf("select docId, content from %s", setting.INDEX_TABLE_FORWARD_INDEX)); err != nil {
		fmt.Println(err.Error())
		return nil
	} else {
		for rows.Next() {
			temp1, temp2 := 0, ""
			rows.Scan(&temp1, &temp2)
			items = append(items, *post_item.NewForwardIndexItem(temp1, temp2))
		}
	}
	return items
}

// 工厂模式进行改造
func (this *IndexDB) QueryPostItemById(PostItemIds []int) []post_item.PostItem {
	items := make([]post_item.PostItem, 0)
	// SQL 还需要优化 可以保存SQL连接的表
	for i, n := 0, len(PostItemIds); i < n; i++ {
		rows, err := this.db.Query(fmt.Sprintf("select content, filePath from %s join %s on %s.docId=%s.docId join %s on %s.realDocId=%s.realDocId where postItemId=%d",
			setting.INDEX_TABLE_POST_ITEM, setting.INDEX_TABLE_FORWARD_INDEX, setting.INDEX_TABLE_POST_ITEM, setting.INDEX_TABLE_FORWARD_INDEX,
			setting.INDEX_TABLE_REAL_DOC, setting.INDEX_TABLE_FORWARD_INDEX, setting.INDEX_TABLE_REAL_DOC,
			PostItemIds[i],
		))
		if err != nil {
			fmt.Println(err.Error())
			return nil
		} else {
			for rows.Next() {
				temp1, temp2 := "", ""
				rows.Scan(&temp1, &temp2)
				items = append(items, *post_item.NewPostItem(temp1, temp2))
				// fmt.Println(temp1, temp2)
			}
		}
	}
	return items
}
