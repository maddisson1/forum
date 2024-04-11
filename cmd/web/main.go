package main

import (
	"forum/app"
	"forum/internal/config"
	"forum/internal/handlers"
	"forum/internal/repo"
	"forum/internal/service"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	conf := config.Loader()

	templateCache, err := app.NewTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	app := app.New(infoLog, errorLog, templateCache)

	repo, err := repo.New(conf.StoragePath)
	if err != nil {
		errorLog.Fatal(err)
	}

	serv := service.New(repo)

	hand := handlers.New(serv, app)

	srv := &http.Server{
		Addr:         conf.Address,
		ErrorLog:     errorLog,
		Handler:      hand.Routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	infoLog.Printf("Starting server on http://localhost%s/", conf.Address)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}
