package controller

import (
	"net/http"

	"github.com/ismtabo/sso-poc/auth/pkg/service"
)

// AuthController implements controller for authorization
type AuthController interface {
	Login(rw http.ResponseWriter, r *http.Request)
	Logout(rw http.ResponseWriter, r *http.Request)
}

type authController struct {
	auth     service.AuthService
	sessions service.SessionService
}

func NewAuthController(authSvc service.AuthService, sessionSvc service.SessionService) AuthController {
	return &authController{auth: authSvc, sessions: sessionSvc}
}

func (ac *authController) Login(rw http.ResponseWriter, r *http.Request) {
	if authenticated, err := ac.sessions.IsAuthenticated(rw, r); err != nil || authenticated {
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		http.Redirect(rw, r, "/", http.StatusFound)
		return
	}
	r.ParseForm()
	uid := r.FormValue("username")
	if uid == "" {
		rw.Write([]byte("empty username form field"))
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	password := r.FormValue("password")
	if password == "" {
		rw.Write([]byte("empty password form field"))
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := ac.auth.Authenticate(uid, password); err != nil {
		rw.WriteHeader(http.StatusForbidden)
		return
	}
	if err := ac.sessions.Login(rw, r); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.Redirect(rw, r, "/", http.StatusFound)
}

func (ac *authController) Logout(rw http.ResponseWriter, r *http.Request) {
	if authenticated, err := ac.sessions.IsAuthenticated(rw, r); err != nil || !authenticated {
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		http.Redirect(rw, r, "/login", http.StatusFound)
		return
	}
	if err := ac.sessions.Logout(rw, r); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.Redirect(rw, r, "/login", http.StatusFound)
}
