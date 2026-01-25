package repository

import (
	"log"
	"time"

	"go.etcd.io/bbolt"
)

type Bolt struct {
	db              *bbolt.DB
	TransactionRepo *TransactionRepository
	ConfigRepo      *ConfigRepository
}

func check(err error) error {
	if err != nil {
		log.Panic(err)
	}
	return nil
}

func InitBolt(file string) (*Bolt, error) {
	db, err := bbolt.Open(file, 0600, &bbolt.Options{Timeout: 2 * time.Second})

	return &Bolt{
		ConfigRepo:      &ConfigRepository{DB: db},
		TransactionRepo: &TransactionRepository{DB: db},
		db:              db,
	}, err
}

func (b *Bolt) CloseBolt() {
	_ = b.db.Close()
}

func (b *Bolt) InitialBuckets() error {
	err := b.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("config"))
		if b == nil {
			log.Printf("buckets not found, creating...")

			_, err := tx.CreateBucket([]byte("config"))
			_ = check(err)
			log.Printf("config bucket created.")

			_, err = tx.CreateBucket([]byte("transactions"))
			_ = check(err)
			log.Printf("transactions bucket created.")

			b = tx.Bucket([]byte("config"))
			err = b.Put([]byte("client_token"), []byte(""))
			_ = check(err)

			err = b.Put([]byte("secret_key"), []byte(""))
			_ = check(err)

			err = b.Put([]byte("base_url"), []byte("https://vpospayment-preprod.halkode.com.tr"))
			_ = check(err)
			return nil
		}
		return nil
	})
	return err
}
