package document

import (
	"fmt"
	"testing"
)

func TestRealDocument(t *testing.T) {
	filePathChan := make(chan string)

	go ReadFileRepository(filePathChan)
	for dir := range filePathChan {
		fmt.Println(dir)
		if dir == "finish" {
			break
		}
	}
	/* if err := ReadFileRepository(filePathChan); err != nil { */
	// fmt.Println("tet read file sad")
	// t.Fail()
	/* } */
}
