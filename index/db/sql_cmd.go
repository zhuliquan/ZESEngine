package db

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

// 插入类型
const (
	SELECT uint32 = iota
	DELETE
	UPDATE
	INSERT
)

// 传入sql操作信息
type SQLInfo struct {
	tableName  string
	selectItem []string
	setInfoK   []string
	setInfoV   [][]string
	whereInfo  map[string]string
	joinInfo   map[string]string
	limitInfo  string
}

// sql 结果类型
type SQLResult struct {
	rows    *sql.Rows
	success bool
}

func NewSQLInfo(tableName string, setInfoK []string, setInfoV [][]string, selectItem []string, whereInfo map[string]string, joinInfo map[string]string, limitInfo string) *SQLInfo {
	sqlInfo := new(SQLInfo)
	sqlInfo.selectItem = selectItem
	sqlInfo.setInfoK = setInfoK
	sqlInfo.setInfoV = setInfoV
	sqlInfo.tableName = tableName
	sqlInfo.joinInfo = joinInfo
	sqlInfo.whereInfo = whereInfo
	sqlInfo.limitInfo = limitInfo
	return sqlInfo
}

func NewSQLResult(content *sql.Rows, success bool) *SQLResult {
	sqlResult := new(SQLResult)
	sqlResult.rows = content
	sqlResult.success = success
	return sqlResult
}

// 数据库操作函数接口
type DBOperateFunc func(tx *sql.Tx, sql_info *SQLInfo) *SQLResult

// 对于sql 返回数据的处理函数接口
type ResultProcessFunc func(*sql.Result) *sql.Rows

// 对于结果什么都不做，直接返回nil
func ReturnNothing(result *sql.Result) *sql.Rows {
	return nil
}

// 通用的处理sql处理函数
func NeedReturnOperate(tx *sql.Tx, command string, process ResultProcessFunc) *SQLResult {
	if result, err := tx.Exec(command); err != nil {
		fmt.Println(err.Error())
		fmt.Println(command)
		return NewSQLResult(nil, false)
	} else {
		return NewSQLResult(process(&result), true)
	}
}

// 插入操作
func SQLInsert(tx *sql.Tx, sqlInfo *SQLInfo) *SQLResult {
	command := fmt.Sprintf("insert ignore into %s\n", sqlInfo.tableName)
	if sqlInfo == nil || len(sqlInfo.setInfoK) == 0 || len(sqlInfo.setInfoK) != len(sqlInfo.setInfoV[0]) {
		panic("这个设置列表不可以为空")
	} else if len(sqlInfo.setInfoK) != len(sqlInfo.setInfoV[0]) {
		panic("字段长度不一致")
	} else {
		command += fmt.Sprintf("(%s) values\n", strings.Join(sqlInfo.setInfoK, ", "))
		bs := make([]string, 0)
		for i, n := 0, len(sqlInfo.setInfoV); i < n; i++ {
			bs = append(bs, fmt.Sprintf("(%s)", strings.Join(sqlInfo.setInfoV[i], " ,")))
		}
		command += strings.Join(bs, " ,\n")
	}
	return NeedReturnOperate(tx, command, ReturnNothing)
}

// 删除操作
func SQLDelete(tx *sql.Tx, sqlInfo *SQLInfo) *SQLResult {
	command := fmt.Sprintf("delete from %s\n", sqlInfo.tableName)
	if !(sqlInfo.whereInfo == nil || len(sqlInfo.whereInfo) == 0) {
		command += " where "
		bs := make([]string, 0)
		for k, v := range sqlInfo.whereInfo {
			bs = append(bs, fmt.Sprintf("%s = %s", k, v))
		}
		command += strings.Join(bs, " and ")

	}
	return NeedReturnOperate(tx, command, ReturnNothing)
}

// 更新操作
// TODO 这里有问题需要修改 但是目前不想管它
func SQLUpdate(tx *sql.Tx, sqlInfo *SQLInfo) *SQLResult {
	command := fmt.Sprintf("update %s", sqlInfo.tableName)
	if sqlInfo == nil || len(sqlInfo.setInfoK) == 0 || len(sqlInfo.setInfoK) != len(sqlInfo.setInfoV[0]) {
		panic("这个设置列表不可以为空")
	} else if len(sqlInfo.setInfoK) != len(sqlInfo.setInfoV[0]) {
		panic("字段长度不一致")
	} else {
		command += fmt.Sprintf("(%s) values\n", strings.Join(sqlInfo.setInfoK, " , "))
		bs := make([]string, 0)
		for i, n := 0, len(sqlInfo.setInfoV); i < n; i++ {
			bs = append(bs, fmt.Sprintf("(%s)", strings.Join(sqlInfo.setInfoV[i], " ,")))
		}
		command += strings.Join(bs, " ,\n")
	}
	if !(sqlInfo.whereInfo == nil || len(sqlInfo.whereInfo) == 0) {
		command += " where "
		bs := make([]string, 0)
		for k, v := range sqlInfo.whereInfo {
			bs = append(bs, fmt.Sprintf("%s = %s,", k, v))
		}
		command += strings.Join(bs, " and ")
	}
	return NeedReturnOperate(tx, command, ReturnNothing)
}

// 查询操作
func SQLSelect(tx *sql.Tx, sqlInfo *SQLInfo) *SQLResult {
	command := fmt.Sprintf("select %s from %s", strings.Join(sqlInfo.selectItem, " , "), sqlInfo.tableName)
	if !(sqlInfo.joinInfo == nil || len(sqlInfo.joinInfo) == 0) {
		bs := make([]string, 0)
		for k, v := range sqlInfo.joinInfo {
			bs = append(bs, fmt.Sprintf("join %s on %s", k, v))
		}
		command += strings.Join(bs, " ")
	}
	if !(sqlInfo.whereInfo == nil || len(sqlInfo.whereInfo) == 0) {
		command += " where "
		bs := make([]string, 0)
		for k, v := range sqlInfo.whereInfo {
			bs = append(bs, fmt.Sprintf("%s = %s", k, v))
		}
		command += strings.Join(bs, " and ")
	}
	if len([]rune(sqlInfo.limitInfo)) == 0 {
		command += fmt.Sprintf(" limit %s", sqlInfo.limitInfo)
	}
	// TODO 后续还需要修改
	if Rows, err := tx.Query(command); err != nil {
		fmt.Println(err.Error())
		fmt.Println(command)
		return NewSQLResult(nil, false)
	} else {
		return NewSQLResult(Rows, true)
	}
}
