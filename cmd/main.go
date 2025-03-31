package main

import (
	"fmt"
	"log/slog"
	"os"
	"pwgen/internal"
	"pwgen/internal/db"
)

const (
	dbFile = "pw.db"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: pw.gen [gen,logins,domains,version] -h")
		return
	}
	err := db.InitDB(dbFile)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	defer func(d db.Database) {
		d.Close()
	}(db.DB)

	switch os.Args[1] {
	case "gen":
		internal.Generate()
	case "logins":
		internal.PrintLogins()
	case "domains":
		internal.PrintDomains()
	case "version":
		internal.VersionApp()
	default:
		fmt.Println("usage: pw.gen [gen,logins,domains,version] -h")

	}
}
