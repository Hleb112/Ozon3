package repository

import (
	"Ozon/models"
	"database/sql"
)

type Repository struct {
	db *sql.DB
}

func New(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (rep Repository) SaveUrl(result *models.Result) error {
	_, err := rep.db.Exec("insert into links3 (link, short) values ($1, $2)", result.Link, result.ShortLink)
	if err != nil {
		return err
	}
	result.Status = "Сокращение было выполнено успешно(postgres)"

	return nil
}

func (rep Repository) GetUrl(vars string) (string, error) {
	rows := rep.db.QueryRow("select link from links3 where short = $1", vars)

	var link string

	err := rows.Scan(&link)
	if err != nil {
		return "", err
	}

	return link, nil
}

func (rep Repository) GetUrlDouble(vars string) (string, error) {
	rows := rep.db.QueryRow("select link from links3 where link = $1", vars)

	var link string

	err := rows.Scan(&link)
	if err != nil {
		return "", err
	}

	return link, nil
}

func (rep Repository) GetShortUrl(vars string) (string, error) {
	rows := rep.db.QueryRow("select short from links3 where link = $1", vars)

	var link string

	err := rows.Scan(&link)
	if err != nil {
		return "", err
	}

	return link, nil
}
