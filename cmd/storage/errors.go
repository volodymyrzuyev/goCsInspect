package storage

import "fmt"

var (
	DbError = fmt.Errorf("DB error")
	NoItem  = fmt.Errorf("Item does not exist")
)
