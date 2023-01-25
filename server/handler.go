package server

import (
	"Ozon/models"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"io"
	"net/http"
)

func (s *Server) getPage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	link, err := s.service.GetUrl(vars["key"])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	io.WriteString(w, "Оригинальная ссылка из базы: ")
	io.WriteString(w, link)
}

func (s *Server) getPageCache(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	link := s.service.GetUrlCache(vars["key"])
	io.WriteString(w, "Оригинальная ссылка из кэша: ")
	io.WriteString(w, string(link))
}

func (s *Server) indexPage(w http.ResponseWriter, r *http.Request) {
	templ, _ := template.ParseFiles("template/index.html")
	result := models.Result{}
	if r.Method == http.MethodPost {
		fmt.Println(r.FormValue("s"))
		if s.service.CheckUrl(r.FormValue("s"), &result) != true {

			err := templ.Execute(w, result)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			return
		}
		unique, err := s.service.UniqueUrl(r.FormValue("s"))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if unique == true {
			result.Link = r.FormValue("s")
			result.ShortLink = s.service.Shorting()
			s.service.SaveUrl(&result)
		} else {
			link, err := s.service.GetShortUrl(r.FormValue("s"))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			io.WriteString(w, "короткая ссылка из базы: ")
			io.WriteString(w, link)
			return
		}
	}

	err := templ.Execute(w, result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *Server) indexPageCache(w http.ResponseWriter, r *http.Request) {
	templ, _ := template.ParseFiles("template/index.html")
	result := models.Result{}
	if r.Method == http.MethodPost {
		if s.service.CheckUrl(r.FormValue("s"), &result) != true {

			err := templ.Execute(w, result)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			return
		}
		if s.service.UniqueUrlCache(r.FormValue("s")) == true {
			result.Link = r.FormValue("s")
			result.ShortLink = s.service.Shorting()

			err := s.service.SaveUrlinCache(&result)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

		} else {
			link := s.service.GetUrlCache(r.FormValue("s"))
			strlink := string(link)
			io.WriteString(w, "короткая ссылка из кэша: ")
			io.WriteString(w, strlink)
			return
		}
	}

	err := templ.Execute(w, result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
