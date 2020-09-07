package document

import (
	_ "github.com/go-sql-driver/mysql"
)

// read interface
type DocReader interface {
	Read() string
}

// write interface
type DocWriter interface {
	Write() error
}
