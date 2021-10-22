package service

import (
	"net/http"

	"github.com/gorilla/sessions"
)

const SESSION_PROPERTY_AUTH = "authenticated"

type SessionService interface {
	IsAuthenticated(rw http.ResponseWriter, r *http.Request) (bool, error)
	Login(rw http.ResponseWriter, r *http.Request) error
	Logout(rw http.ResponseWriter, r *http.Request) error
}

type cookieSessionService struct {
	store  *sessions.CookieStore
	cookie string
}

func NewCookieSessionService(store *sessions.CookieStore, cookieName string) SessionService {
	return &cookieSessionService{store: store, cookie: cookieName}
}

func (cs *cookieSessionService) IsAuthenticated(rw http.ResponseWriter, r *http.Request) (bool, error) {
	session, err := cs.store.Get(r, cs.cookie)
	if err != nil {
		return false, err
	}
	authenticated, ok := session.Values[SESSION_PROPERTY_AUTH]
	if !ok {
		return false, nil
	}
	return authenticated.(bool), nil
}

func (cs *cookieSessionService) Login(rw http.ResponseWriter, r *http.Request) error {
	session, err := cs.store.Get(r, cs.cookie)
	if err != nil {
		return err
	}
	session.Values[SESSION_PROPERTY_AUTH] = true
	return session.Save(r, rw)
}

func (cs *cookieSessionService) Logout(rw http.ResponseWriter, r *http.Request) error {
	session, err := cs.store.Get(r, cs.cookie)
	if err != nil {
		return err
	}
	session.Values[SESSION_PROPERTY_AUTH] = false
	return session.Save(r, rw)
}
