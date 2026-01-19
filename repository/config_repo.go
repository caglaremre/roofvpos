package repository

import "go.etcd.io/bbolt"

type ConfigRepository struct {
	DB *bbolt.DB
}

func (cf *ConfigRepository) GetClientAndSecretKey() (string, string) {
	var clientToken, secretKey string
	err := cf.DB.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("config"))
		clientToken = string(b.Get([]byte("client_token")))
		secretKey = string(b.Get([]byte("secret_key")))
		return nil
	})
	_ = check(err)
	return clientToken, secretKey
}

func (cf *ConfigRepository) UpdateClientAndSecretKey(clientToken, secretKey string) error {
	err := cf.DB.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("config"))
		err := b.Put([]byte("client_token"), []byte(clientToken))
		_ = check(err)

		err = b.Put([]byte("secret_key"), []byte(secretKey))
		_ = check(err)
		return nil
	})
	return err
}

func (cf *ConfigRepository) GetBaseURL() string {
	baseURL := ""
	err := cf.DB.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("config"))
		baseURL = string(b.Get([]byte("base_url")))
		return nil
	})
	_ = check(err)
	return baseURL
}
