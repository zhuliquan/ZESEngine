package document

// 分片之后存储在数据库的数据
type LogicDocument struct {
	content        string
	RealDocumentId int
}

func NewLogicDocument(content string, RealDocumentId int) *LogicDocument {
	var ld *LogicDocument = new(LogicDocument)
	ld.content = content
	ld.RealDocumentId = RealDocumentId
	return ld
}
