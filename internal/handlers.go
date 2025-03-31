package internal

import (
	"crypto/sha256"
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"pwgen/internal/cipher"
	"pwgen/internal/db"
	"pwgen/internal/gen"
	"strconv"
	"strings"
)

var (
	version int
	login   string
	pass    string
	domain  string
)

func VersionApp() {
	fmt.Println("pw.gen v0.0.1")
}

func PrintLogins() {
	set := flag.NewFlagSet("", flag.ExitOnError)
	set.StringVar(&domain, "d", "", "Domain filter logins")
	set.StringVar(&pass, "p", "", "Password to generate sub password")
	err := set.Parse(os.Args[2:])
	if err != nil {
		slog.Error("Parse params", "error", err.Error())
	}
	if pass == "" {
		slog.Error("You must provide a password")
		os.Exit(1)
	}

	data, err := db.DB.List()
	if err != nil {
		slog.Error("List database", "error", err.Error())
		os.Exit(1)
	}
	key := sha256.Sum256([]byte(pass))
	var dbdata = new(db.Data)
	for _, j := range data {
		dt, err := cipher.Decrypt(j, key[:32])
		if err != nil {
			slog.Error("decrypt data", "error", err.Error())
			continue
		}
		err = json.Unmarshal(dt, dbdata)
		if err != nil {
			slog.Error("decrypt data", "error", err.Error())
			continue
		}
		if domain != "" {
			if strings.Contains(dbdata.Domain, domain) {
				fmt.Println("Domain: ", dbdata.Domain, "Login:", dbdata.Login, "Versions:", dbdata.Vesion)
				continue
			}
		}
		if dbdata.Login != "" {
			fmt.Println("Domain: ", dbdata.Domain, "Login:", dbdata.Login, "Versions:", dbdata.Vesion)
		}
	}
}

func PrintDomains() {
	set := flag.NewFlagSet("", flag.ExitOnError)
	set.StringVar(&login, "l", "", "Login filter domains")
	set.StringVar(&pass, "p", "", "Password to generate sub password")
	err := set.Parse(os.Args[2:])
	if err != nil {
		slog.Error("Parse params", "error", err.Error())
	}
	if pass == "" {
		slog.Error("You must provide a password")
		os.Exit(1)
	}

	data, err := db.DB.List()
	if err != nil {
		slog.Error("List database", "error", err.Error())
		os.Exit(1)
	}
	key := sha256.Sum256([]byte(pass))
	var dbdata = new(db.Data)
	for _, j := range data {
		dt, err := cipher.Decrypt(j, key[:32])
		if err != nil {
			slog.Error("decrypt data", "error", err.Error())
			continue
		}
		err = json.Unmarshal(dt, dbdata)
		if err != nil {
			slog.Error("decrypt data", "error", err.Error())
			continue
		}
		if login != "" {
			if strings.Contains(dbdata.Login, login) {
				fmt.Println("Domain: ", dbdata.Domain, "Login:", dbdata.Login, "Versions:", dbdata.Vesion)
				continue
			}
		}
		fmt.Println("Domain: ", dbdata.Domain, "Login:", dbdata.Login, "Versions:", dbdata.Vesion)
	}
}

func Generate() {
	set := flag.NewFlagSet("", flag.ExitOnError)
	set.StringVar(&domain, "d", "", "Domain to generate sub password")
	set.IntVar(&version, "v", 0, "Version to generate sub password")
	set.StringVar(&login, "l", "", "Login to generate sub password")
	set.StringVar(&pass, "p", "", "Password to generate sub password")
	err := set.Parse(os.Args[2:])
	if err != nil {
		slog.Error("Parse params", "error", err.Error())
	}
	if domain == "" {
		flag.Usage()
		os.Exit(1)
	}
	if pass == "" {
		flag.Usage()
		os.Exit(1)
	}
	var ver string
	if version > 0 {
		ver = strconv.Itoa(version)
	}
	password := gen.GeneratePassword(pass, domain, ver)
	fmt.Printf("Generated password for %s: %s\n", domain, password)
	err = saveToDB(db.Data{Domain: domain, Login: login, Vesion: version}, pass)
	if err != nil {
		slog.Error("Save database", "error", err.Error())
		os.Exit(1)
	}
}

func saveToDB(d db.Data, pass string) error {
	key := sha256.Sum256([]byte(pass))
	k := sha256.Sum256([]byte(d.Domain + d.Login))

	val, err := json.Marshal(d)
	if err != nil {
		return err
	}
	v, err := cipher.Encrypt(val, key[:32])
	if err != nil {
		return err
	}
	err = db.DB.Save(k[:32], v)
	if err != nil {
		return err
	}
	return nil
}
