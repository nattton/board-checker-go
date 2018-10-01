package main

import (
	"database/sql"
	"flag"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	// log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	// log.SetLevel(log.WarnLevel)
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP Network Address")
	dsn := flag.String("dsn", os.Getenv("BC_DSN"), "Database DSN")
	htmlDir := flag.String("html-dir", os.Getenv("GOPATH")+"/src/gitlab.com/code-mobi/board-checker/ui/html", "Path to static assets")
	secret := flag.String("secret", "y8eETRxBcYUyvv9x6c6Pk7JsWf7bpC37", "Secret key")
	staticDir := flag.String("static-dir", os.Getenv("GOPATH")+"/src/gitlab.com/code-mobi/board-checker/ui/static", "Path to static assets")
	storeDir := flag.String("store-dir", os.Getenv("BC_STORE"), "Path to store files")

	flag.Parse()

	sessionManager := scs.NewCookieManager(*secret)
	sessionManager.Lifetime(12 * time.Hour)
	sessionManager.Persist(true)

	app := &App{
		DSN:       *dsn,
		Sessions:  sessionManager,
		HTMLDir:   *htmlDir,
		StaticDir: *staticDir,
		StoreDir:  *storeDir,
		SecretKey: *secret,
	}

	log.Println("Starting server on " + *addr)
	err := http.ListenAndServe(*addr, app.Routes())
	log.Fatal(err)
}

func connect(dsn string) *sql.DB {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	return db
}
