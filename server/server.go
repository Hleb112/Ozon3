package server

import (
	"Ozon/service"
	"github.com/allegro/bigcache/v3"
	"github.com/gorilla/mux"
	"net/http"
)

type Server struct {
	service *service.Service
	cache   *bigcache.BigCache
}

func New(service *service.Service) *Server {
	return &Server{
		service: service,
	}
}

func (s *Server) Start(useCache bool) error {

	router := mux.NewRouter()
	if useCache == false {
		router.HandleFunc("/", s.indexPage)
		router.HandleFunc("/{key}", s.getPage)
	} else {
		router.HandleFunc("/", s.indexPageCache)
		router.HandleFunc("/{key}", s.getPageCache)
	}

	//router.HandleFunc("/to/{key}", redirectTo)
	return http.ListenAndServe(":8000", router)
}
