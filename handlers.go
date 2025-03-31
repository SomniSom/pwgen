package main

import (
	"crypto/sha256"
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
)

func versionApp() {
	fmt.Println("pw.gen v0.0.1")
}

func printLogins() {
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

	data, err := db.list()
	if err != nil {
		slog.Error("List database", "error", err.Error())
		os.Exit(1)
	}
	key := sha256.Sum256([]byte(pass))
	var dbdata = new(Data)
	for _, j := range data {
		dt, err := decrypt(j, key[:32])
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

func printDomains() {
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

	data, err := db.list()
	if err != nil {
		slog.Error("List database", "error", err.Error())
		os.Exit(1)
	}
	key := sha256.Sum256([]byte(pass))
	var dbdata = new(Data)
	for _, j := range data {
		dt, err := decrypt(j, key[:32])
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

func generate() {
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
	password := GeneratePassword(pass, domain, ver)
	fmt.Printf("Generated password for %s: %s\n", domain, password)
	err = saveToDB(Data{Domain: domain, Login: login, Vesion: version}, pass)
	if err != nil {
		slog.Error("Save database", "error", err.Error())
		os.Exit(1)
	}
}

func saveToDB(d Data, pass string) error {
	key := sha256.Sum256([]byte(pass))
	k := sha256.Sum256([]byte(d.Domain + d.Login))

	val, err := json.Marshal(d)
	if err != nil {
		return err
	}
	v, err := encrypt(val, key[:32])
	if err != nil {
		return err
	}
	err = db.save(k[:32], v)
	if err != nil {
		return err
	}
	return nil
}
