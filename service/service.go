package service

import (
	"Ozon/models"
	"Ozon/repository"
	"github.com/allegro/bigcache/v3"
	"math/rand"
	"net/url"
)

type Service struct {
	repo  *repository.Repository
	cache *bigcache.BigCache
}

func New(repo *repository.Repository, cache *bigcache.BigCache) *Service {
	return &Service{
		repo:  repo,
		cache: cache,
	}
}

func (s Service) CheckUrl(url string, result *models.Result) bool {
	if IsValidUrl(url) {
		result.Status = "Ссылка имеет неправильный формат!"
		result.Link = ""
		return false
	}
	return true
}

func (s Service) Shorting() string {
	return shorting()
}

func (s Service) SaveUrl(result *models.Result) {
	err := s.repo.SaveUrl(result)
	if err != nil {
		return
	}
}

func (s Service) SaveUrlinCache(result *models.Result) error {
	err := s.cache.Set(result.Link, []byte(result.ShortLink))
	if err != nil {
		return err
	}

	err = s.cache.Set(result.ShortLink, []byte(result.Link))
	if err != nil {
		return err
	}

	result.Status = "Сокращение было выполнено успешно"

	return nil
}

func (s Service) GetUrl(url string) (string, error) {
	link, err := s.repo.GetUrl(url)
	if err != nil {
		return "", err
	}

	return link, nil
}

func (s Service) GetShortUrl(url string) (string, error) {
	link, err := s.repo.GetShortUrl(url)
	if err != nil {
		return "", err
	}

	return link, nil
}

func (s Service) GetUrlCache(url string) []byte {
	link, _ := s.cache.Get(url)
	return link
}

func (s Service) UniqueUrl(url string) (bool, error) {
	link, err := s.repo.GetUrlDouble(url)
	if err != nil {
		return false, err
	}

	if link == url {
		return false, nil
	}
	return true, nil
}

func (s Service) UniqueUrlCache(url string) bool {
	shortUrl, _ := s.cache.Get(url)
	strShortUrl := string(shortUrl)
	longUrl, _ := s.cache.Get(strShortUrl)
	strLongUrl := string(longUrl)
	if strLongUrl == url {
		return false
	}
	return true
}

func IsValidUrl(token string) bool {
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

func shorting() string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890_"
	b := make([]byte, 10)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
