package controller

import (
	"net/http"

	"github.com/ismtabo/sso-poc/auth/pkg/repository"
	"github.com/ismtabo/sso-poc/auth/pkg/service"
)

type PageController interface {
	Home(rw http.ResponseWriter, r *http.Request)
	Login(rw http.ResponseWriter, r *http.Request)
}

type pageController struct {
	sessions service.SessionService
	pages    repository.PageRepository
}

func NewPagesController(sessionSvc service.SessionService, pageRepo repository.PageRepository) PageController {
	return &pageController{sessions: sessionSvc, pages: pageRepo}
}

func (pc *pageController) Home(rw http.ResponseWriter, r *http.Request) {
	if authenticated, err := pc.sessions.IsAuthenticated(rw, r); err != nil || !authenticated {
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		if !authenticated {
			rw.WriteHeader(http.StatusForbidden)
			return
		}
	}
	body, err := pc.pages.Page("index.html")
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	rw.Write(body)
}

func (pc *pageController) Login(rw http.ResponseWriter, r *http.Request) {
	if authenticated, err := pc.sessions.IsAuthenticated(rw, r); err != nil || authenticated {
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		if !authenticated {
			http.Redirect(rw, r, "/", http.StatusFound)
			return
		}
	}
	body, err := pc.pages.Page("login.html")
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	rw.Write(body)
}
