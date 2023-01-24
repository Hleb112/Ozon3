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
	link := s.service.GetUrl(vars["key"])
	io.WriteString(w, "Оригинальная ссылка: ")
	io.WriteString(w, link)
}

func (s *Server) getPageCache(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	link := s.service.GetUrlCache(vars["key"])
	io.WriteString(w, "Оригинальная ссылка: ")
	io.WriteString(w, string(link))
}

func (s *Server) indexPage(w http.ResponseWriter, r *http.Request) {
	templ, _ := template.ParseFiles("template/index.html")
	result := models.Result{}
	if r.Method == "POST" {
		fmt.Println(r.FormValue("s"))
		if s.service.CheckUrl(r.FormValue("s"), &result) != true {
			templ.Execute(w, result)
			return
		}
		if s.service.UniqueUrl(r.FormValue("s")) == true {
			result.Link = r.FormValue("s")
			result.ShortLink = s.service.Shorting()
			s.service.SaveUrl(&result)
		} else {
			fmt.Println("копия")
			link := s.service.GetShortUrl(r.FormValue("s"))
			io.WriteString(w, "короткая ссылка: ")
			io.WriteString(w, link)
			return
		}
	}
	templ.Execute(w, result)
}

func (s *Server) indexPageCache(w http.ResponseWriter, r *http.Request) {
	templ, _ := template.ParseFiles("template/index.html")
	result := models.Result{}
	if r.Method == "POST" {
		if s.service.UniqueUrlCache(r.FormValue("s")) == true {
			result.Link = r.FormValue("s")
			result.ShortLink = s.service.Shorting()
			s.service.SaveUrlinCache(&result)
		} else {
			fmt.Println("копия")
			link := s.service.GetShortUrl(r.FormValue("s"))
			io.WriteString(w, "короткая ссылка: ")
			io.WriteString(w, link)
			return
		}
	}
	templ.Execute(w, result)
}
