package main

import (
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"log/slog"
)

type Data struct {
	Domain string `json:"domain"`
	Login  string `json:"login"`
	Vesion int    `json:"vesion"`
}

var db *bolt.DB

const bucketName = "pwds"

func InitDB(filename string) (err error) {
	db, err = bolt.Open(filename, 0600, &bolt.Options{})
	if err != nil {
		return err
	}
	_ = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
	return err
}

func save(value Data) error {
	k := []byte(value.Domain)

	v, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return bucket.Put(k, v)
	})
}

func read(domain string) (*Data, error) {
	k := []byte(domain)

	u := new(Data)
	err := db.View(func(tx *bolt.Tx) error {
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

func listLogins() ([]string, error) {
	var logins []string
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		if b == nil {
			return fmt.Errorf("get bucket: FAILED")
		}
		c := b.Cursor()
		var dt Data
		for k, v := c.First(); k != nil; k, v = c.Next() {
			err := json.Unmarshal(v, &dt)
			if err != nil {
				slog.Warn("Unmarshal FAILED", "domain", k)
				continue
			}
			logins = append(logins, dt.Login)
		}
		return nil
	})
	return logins, err
}

func listDomains() ([]string, error) {
	var logins []string
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		if b == nil {
			return fmt.Errorf("get bucket: FAILED")
		}
		c := b.Cursor()
		var dt Data
		for k, v := c.First(); k != nil; k, v = c.Next() {
			err := json.Unmarshal(v, &dt)
			if err != nil {
				slog.Warn("Unmarshal FAILED", "domain", k)
				continue
			}
			logins = append(logins, dt.Login)
		}
		return nil
	})
	return logins, err
}
