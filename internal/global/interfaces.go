package global

import "ants/pkg"

type Algorithm interface {
	Do(fields [9]pkg.FieldType) (field uint8, action pkg.Action)
}

type Storage interface {
	Close()
	Get(collection string, key string) ([]byte, error)
	Put(collection string, key string, value []byte) error
	CreateCollectionIfNotExist(collection string)
}
