package storage

import "github.com/volodymyrzuyev/goCsInspect/cmd/globalTypes"

type Storage interface {
	GetItem(itemId int64) (globalTypes.Item, error)
	InsertItem(item globalTypes.Item) error
}
