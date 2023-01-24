package repository

import (
	"Ozon/models"
	"database/sql"
	"math/rand"
	"net/url"
)

type Repository struct {
	db *sql.DB
}

func New(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (rep Repository) IsValidUrl(token string) bool {
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

func (rep Repository) Shorting() string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890_"
	b := make([]byte, 10)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func (rep Repository) SaveUrl(result *models.Result) {
	rep.db.Exec("insert into links3 (link, short) values ($1, $2)", result.Link, result.ShortLink)
	result.Status = "Сокращение было выполнено успешно(postgres)"

}

func (rep Repository) GetUrl(vars string) string {
	var link string
	rows := rep.db.QueryRow("select link from links3 where short = $1 limit 1", vars)
	rows.Scan(&link)
	return link
}

func (rep Repository) GetUrlDouble(vars string) string {
	var link string
	rows := rep.db.QueryRow("select link from links3 where link = $1 limit 1", vars)
	rows.Scan(&link)
	return link
}

func (rep Repository) GetShortUrl(vars string) string {
	var link string
	rows := rep.db.QueryRow("select short from links3 where link = $1 limit 1", vars)
	rows.Scan(&link)
	return link
}
