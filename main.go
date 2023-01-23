package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/allegro/bigcache/v3"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var useCache = true

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "admin"
	dbname   = "links"
)

func shorting() string {
	b := make([]byte, 10)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func isValidUrl(token string) bool {
	_, err := url.ParseRequestURI(token)
	if err != nil {
		return false
	}
	u, err := url.Parse(token)
	if err != nil || u.Host == "" {
		return false
	}
	return true
}

type Result struct {
	Link      string
	ShortLink string
	Status    string
}

func indexPage(w http.ResponseWriter, r *http.Request) {
	templ, _ := template.ParseFiles("template/index.html")
	result := Result{}
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	if r.Method == "POST" {
		if !isValidUrl(r.FormValue("s")) {
			fmt.Println("Что-то не так")
			result.Status = "Ссылка имеет неправильный формат!"
			result.Link = ""
		} else {
			result.Link = r.FormValue("s")
			result.ShortLink = shorting()
			if useCache == true {
				cache.Set(result.Link, []byte(result.ShortLink))
				result.Status = "Сокращение было выполнено успешно"
			} else {
				db, err := sql.Open("postgres", psqlInfo)
				if err != nil {
					panic(err)
				}
				defer db.Close()
				db.Exec("insert into links (link, shortlink) values ($1, $2)", result.Link, result.ShortLink)
				result.Status = "Сокращение было выполнено успешно"
				err = db.Ping()
				if err != nil {
					panic(err)
				}
			}
		}
	}
	templ.Execute(w, result)
}

func redirectTo(w http.ResponseWriter, r *http.Request) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	var link string
	vars := mux.Vars(r)
	if useCache == true {
		value, _ := cache.Get(vars["key"])
		link = string(value)
		fmt.Fprintf(w, "<script>location='%s';</script>", link)
		return
	}
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	rows := db.QueryRow("select link from links where shortLink=$1 limit 1", vars["key"])
	rows.Scan(&link)
	fmt.Fprintf(w, "<script>location='%s';</script>", link)
}

var cache, _ = bigcache.New(context.Background(), bigcache.DefaultConfig(10*time.Minute))

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", indexPage)
	router.HandleFunc("/to/{key}", redirectTo)
	log.Fatal(http.ListenAndServe(":8000", router))
}
