package storage

type Storage interface {
	Close()
	Get(collection string, key string) ([]byte, error)
	GetKeys(collection string) ([]string, error)
	Put(collection string, key string, value []byte) error
	CreateCollectionIfNotExist(collection string)
}
