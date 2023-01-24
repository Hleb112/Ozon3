package main

import (
	repository2 "Ozon/repository"
	"Ozon/server"
	"Ozon/service"
	"context"
	"database/sql"
	"github.com/allegro/bigcache/v3"
	_ "github.com/lib/pq"
	"log"
	"time"
)

func main() {
	useCache := false
	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres "+
		"password=admin dbname=links sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	cache, _ := bigcache.New(context.Background(), bigcache.DefaultConfig(10*time.Minute))
	repository := repository2.New(db)
	servicE := service.New(repository, cache)

	srv := server.New(servicE)
	err = srv.Start(useCache)
	if err != nil {
		log.Fatalln(err)
	}

}
