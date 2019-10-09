package global

import (
	"github.com/boltdb/bolt"
)

type Bolt struct {
	db *bolt.DB
}

func NewBolt(DBName string) *Bolt {
	db, err := bolt.Open(DBName+".db", 0600, nil)
	if err != nil {
		panic(err)
	}

	return &Bolt{db: db}
}

func (b *Bolt) Close() {
	err := b.db.Close()
	if err != nil {
		panic(err)
	}
}

func (b *Bolt) Get(collection string, key string) ([]byte, error) {
	var data []byte
	err := b.db.View(func(tx *bolt.Tx) error {
		data = tx.Bucket([]byte(collection)).Get([]byte(key))

		return nil
	})

	return data, err
}

func (b *Bolt) Put(collection string, key string, value []byte) error {
	err := b.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(collection))

		err := bucket.Put([]byte(key), value)

		return err
	})

	return err
}

func (b *Bolt) CreateCollectionIfNotExist(name string) {
	err := b.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(name))

		return err
	})

	if err != nil {
		panic(err)
	}
}
