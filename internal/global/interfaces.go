package global

type Storage interface {
	Close()
	Get(collection string, key string) ([]byte, error)
	Put(collection string, key string, value []byte) error
	CreateCollectionIfNotExist(collection string)
}
