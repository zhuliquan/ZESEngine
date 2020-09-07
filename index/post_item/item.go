package post_item

import "fmt"

type PostItem struct {
	content string
	docPath string
}

func NewPostItem(content string, docPath string) *PostItem {
	pi := new(PostItem)
	pi.content = content
	pi.docPath = docPath
	return pi
}

func (this *PostItem) Show() {
	fmt.Printf("文件 (%s) 的句子:   %s\n", this.docPath, this.content)
}

type ForwardIndexItem struct {
	Content string
	DocId   int
}

func NewForwardIndexItem(docId int, content string) *ForwardIndexItem {
	pi := new(ForwardIndexItem)
	pi.DocId = docId
	pi.Content = content
	return pi
}
