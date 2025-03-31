package main

import (
	"crypto/sha256"
	"fmt"
	"github.com/boltdb/bolt"
	"log/slog"
	rand2 "math/rand/v2"
	"os"
)

const (
	dbFile = "pw.db"
	sm     = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~"
)

func GeneratePassword(pw, site string, version string) string {
	hash := sha256.Sum256([]byte(pw + site + version))
	rnd := rand2.NewChaCha8(hash)
	buf := make([]byte, 16)
	n, err := rnd.Read(buf)
	if err != nil {
		panic(err)
	}
	for i, b := range buf[:n] {
		if int(b) >= len(sm) {
			buf[i] = sm[int(b)%len(sm)]
			continue
		}
		buf[i] = sm[b]
	}
	return string(buf)
}

var (
	version int
	login   string
	pass    string
	domain  string
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: pw.gen [gen,logins,domains,version] -h")
		return
	}
	err := InitDB(dbFile)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	defer func(db *bolt.DB) {
		_ = db.Close()
	}(db.db)
	switch os.Args[1] {
	case "gen":
		generate()
	case "logins":
		printLogins()
	case "domains":
		printDomains()
	case "version":
		versionApp()
	default:
		fmt.Println("usage: pw.exe [-v|-version]")

	}
}
