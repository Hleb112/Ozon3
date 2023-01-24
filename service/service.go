package service

import (
	"Ozon/models"
	"Ozon/repository"
	"github.com/allegro/bigcache/v3"
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

func (s Service) CheckUrl(value string, result *models.Result) bool {
	if !s.repo.IsValidUrl(value) {
		result.Status = "Ссылка имеет неправильный формат!"
		result.Link = ""
		return false
	}
	return true
}

func (s Service) Shorting() string {
	sh := s.repo.Shorting()
	return sh
}

func (s Service) SaveUrl(result *models.Result) {
	s.repo.SaveUrl(result)
}

func (s Service) SaveUrlinCache(result *models.Result) {
	s.cache.Set(result.Link, []byte(result.ShortLink))
	s.cache.Set(result.ShortLink, []byte(result.Link))
	result.Status = "Сокращение было выполнено успешно"
}

func (s Service) GetUrl(vars string) string {
	link := s.repo.GetUrl(vars)
	return link
}

func (s Service) GetUrlCache(vars string) []byte {
	link, _ := s.cache.Get(vars)
	return link
}
