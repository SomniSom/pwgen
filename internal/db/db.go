package db

import (
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
)

type Data struct {
	Domain string `json:"domain"`
	Login  string `json:"login"`
	Vesion int    `json:"vesion"`
}

type Database struct {
	db *bolt.DB
}

var DB = Database{}

const bucketName = "pwds"

func InitDB(filename string) (err error) {
	DB.db, err = bolt.Open(filename, 0600, &bolt.Options{})
	if err != nil {
		return err
	}

	_ = DB.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
	return err
}

func (d Database) Close() {
	_ = d.db.Close()
}

func (d Database) Save(k, v []byte) error {
	return d.db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return bucket.Put(k, v)
	})
}

func (d Database) Read(domain string) (*Data, error) {
	k := []byte(domain)

	u := new(Data)
	err := d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		if b == nil {
			return fmt.Errorf("get bucket: FAILED")
		}
		if len(b.Get(k)) == 0 {
			u.Domain = domain
			return nil
		}
		return json.Unmarshal(b.Get(k), u)
	})
	return u, err
}

func (d Database) List() ([][]byte, error) {
	var logins [][]byte
	err := d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		if b == nil {
			return fmt.Errorf("get bucket: FAILED")
		}
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			logins = append(logins, v)
		}
		return nil
	})
	return logins, err
}
