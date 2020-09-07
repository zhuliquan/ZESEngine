package document

import (
	"ZESEngine/setting"
	"io/ioutil"
)

// real document in disk
type RealDocument struct {
	filePath       string
	realDocumentId int32
}

func NewRealDocumnet(filePath string) *RealDocument {
	var obj = new(RealDocument)
	obj.filePath = filePath
	obj.realDocumentId = -1
	return obj
}

// 读取硬盘的文件内容
func (this *RealDocument) Read() string {
	// 利用协程处理
	if data, err := ioutil.ReadFile(string(this.filePath)); err != nil {
		panic(err.Error())
	} else {
		return string(data)
	}
}

// 大规模读取file_repository下的文件
func ReadFileRepository(filePathChan chan string) {
	if files, err := ioutil.ReadDir(setting.REAL_DOCUMENT_REPOSITORY); err != nil {
		panic(err.Error())
	} else {
		for _, file := range files {
			filePathChan <- file.Name()
		}
		filePathChan <- "finish"
	}
}
