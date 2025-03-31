package main

import (
	"crypto/sha256"
	"fmt"
	rand2 "math/rand/v2"
	"os"
)

const (
	dbFile = "pw.db"
	sm     = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~"
)

var dbKey = "masterkey"

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
		fmt.Println("usage: pw.gen [-v|-version]")
		return
	}
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

	//key := sha256.Sum256([]byte("1234567890"))
	//enc, err := encrypt([]byte(`{"test": "test"}`), key[:32])
	//if err != nil {
	//	slog.Error("Encrypt", "err", err)
	//	return
	//}
	//dec, err := decrypt(enc, key[:32])
	//fmt.Println(string(dec))
	//return

}
