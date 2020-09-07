package setting

// 各种配置的字段
const (
	REAL_DOCUMENT_REPOSITORY string = "../data/file_repository/"
	DATABASE_DIRECTORY       string = "./data/data_base/"
	REDIS_HOST               string = "127.0.0.1"
	REDIS_PORT               string = "6379"
	MYSQL_USERNAME           string = "root"
	MYSQL_PASSWORD           string = "123456"
	MYSQL_HOST               string = "127.0.0.1"
	MYSQL_PORT               string = "3306"
	MYSQL_DATABASE           string = "zesengine_db"
	MYSQL_CHARSET            string = "utf8"
	// 构建index的表
	INDEX_TABLE_KEY_WORD      string = "key_word"
	INDEX_TABLE_FORWARD_INDEX string = "forward_index"
	INDEX_TABLE_POST_ITEM     string = "post_item"
	INDEX_TABLE_REAL_DOC      string = "real_doc"

	MAX_INSERT_NUMBER uint32 = 3000 // 每一条事务最大输入的词条数目
	BUFFER_SIZE       uint32 = 4096 // 读写文件的缓冲区大小
)
