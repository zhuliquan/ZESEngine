package index

import (
	"ZESEngine/index/db"
	"ZESEngine/index/document"
	"ZESEngine/index/post_item"
	"ZESEngine/index/segement"
	"ZESEngine/setting"
	"fmt"
	"strconv"
	"strings"

	"github.com/garyburd/redigo/redis"
	"github.com/gomodule/redigo"
)

// type InvertIndex map[utils.Word]utils.DocId
type InvertIndex struct {
	engineDB     *db.IndexDB
	segmenter    *segement.ChineseSegementer
	keyWords     map[string]int
	globalWordId int
	globalItemId int
	engineRedis  redis.Conn // TODO redis 持久化
}

func NewInvertIndex() *InvertIndex {
	invertIndex := new(InvertIndex)
	invertIndex.engineDB = db.NewIndexDB()
	invertIndex.segmenter = segement.NewChineseSegementer()
	invertIndex.keyWords = make(map[string]int)
	invertIndex.globalWordId = 0
	invertIndex.globalItemId = 0
	fmt.Println("开启Redis....")
	invertIndex.engineRedis, _ = redigo.Dial("tcp", fmt.Sprintf("%s:%s", setting.REDIS_HOST, setting.REDIS_PORT))
	fmt.Println("已经开启了Redis")
	return invertIndex
}

func (this *InvertIndex) Close() {
	fmt.Println("引擎索引开始关闭")
	this.engineDB.DestroyTables() // TODO 以后修改 销毁所有的表
	this.engineDB.Close()
	this.engineRedis.Do("flushall")
	this.engineRedis.Close()
	fmt.Println("引擎索引以经关闭")
}

func (this *InvertIndex) BuildForwardIndex() {
	fmt.Println("开始构建正向索引")
	filePathChan := make(chan string)
	// TODO 这里存在一点问题 感觉中英文需要统一接口
	go document.ReadFileRepository(filePathChan)
	for file := range filePathChan {
		if file == "finish" {
			break
		} else {
			this.engineDB.AddRealDoc(file)
			var realDocId int = this.engineDB.QueryRealDocIdByFilePaths([]string{file})[0]
			realDoc := document.NewRealDocumnet(fmt.Sprintf("%s/%s", setting.REAL_DOCUMENT_REPOSITORY, file))
			content := realDoc.Read()                            // TODO 没有考虑大文件的情况
			sentences, _ := this.segmenter.SentenceCut(&content) // TODO 需要考虑大文件读写
			for i, n := 0, len(sentences); i < n; i++ {
				newSentence := strings.Replace(sentences[i], "\n", "", 100)
				newSentence = strings.Replace(newSentence, "”", "", 100)
				newSentence = strings.Replace(newSentence, "'", "", 100)
				newSentence = strings.Replace(newSentence, "\"", "", 100)
				sentences[i] = newSentence
			}
			this.engineDB.AddLogicDocs(sentences, realDocId)
		}
	}
	fmt.Println("完成构建正向索引")
}

func (this *InvertIndex) BuildInvertIndex() {
	fmt.Println("开始构建倒排索引")
	forwardIndexs := this.engineDB.QueryAllForwardIndx()
	for _, items := range forwardIndexs {
		docId, content := items.DocId, items.Content
		words := this.segmenter.WordCut(&content)
		postItemIds := make([]int, 0)
		wordIds := make([]int, 0)
		docIds := make([]int, 0)
		for i, n := 0, len(words); i < n; i++ {
			if _, ok := this.keyWords[words[i]]; !ok {
				this.keyWords[words[i]] = this.globalWordId
				wordIds = append(wordIds, this.globalWordId)
				this.globalWordId += 1
			} else {
				wordIds = append(wordIds, this.keyWords[words[i]])
			}
			postItemIds = append(postItemIds, this.globalItemId)
			docIds = append(docIds, docId)
			this.engineRedis.Do("LPUSH", words[i], this.globalItemId)
			this.globalItemId += 1
		}
		// fmt.Println(words[0])
		this.engineDB.ADDKeyWords(wordIds, words)
		this.engineDB.AddPostItems(postItemIds, docIds, wordIds)
	}
	fmt.Println("完成倒排索引构建")
}

func (this *InvertIndex) Search(word string) []post_item.PostItem {
	values, _ := redis.Values(this.engineRedis.Do("lrange", word, "0", "10000"))
	postItemIds := make([]int, 0)
	for _, v := range values {
		i, _ := strconv.Atoi(string(v.([]byte)))
		postItemIds = append(postItemIds, i)
	}
	return this.engineDB.QueryPostItemById(postItemIds)
}

// 以后需要实现的功能
/* func (this *InvertIndex) SearchTop(word string) []post_item.PostItem { */
//
/* } */
