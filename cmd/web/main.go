package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"github.com/jersonsatoru/lets-go/internal/models"
	_ "github.com/lib/pq"
)

func main() {
	port := flag.String("port", ":8080", "Port used by web application")
	dsn := flag.String("dsn", "", "DSN for database connection")
	driverName := flag.String("driverName", "postgres", "Driver name (DB)")

	flag.Parse()

	db, err := getConnecionPool(*driverName, *dsn)
	if err != nil {
		log.Panic(err)
	}

	defer db.Close()
	snippetModel := models.NewSnippetModel(db)

	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Llongfile)
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate)

	templateCache, err := newTemplateCache()
	if err != nil {
		log.Fatal(err)
	}

	sessionManager := scs.New()
	sessionManager.Store = postgresstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour

	formDecoder := form.NewDecoder()
	app := NewApplication(
		infoLog,
		errorLog,
		snippetModel,
		templateCache,
		formDecoder,
		sessionManager,
	)

	srv := http.Server{
		Addr:     *port,
		Handler:  app.routes(),
		ErrorLog: errorLog,
	}

	msg := fmt.Sprintf("Server started at port: %s", *port)
	infoLog.Println(msg)
	err = srv.ListenAndServe()

	if err != nil {
		log.Panic(err)
	}
}
