package handlers

import (
	"forum/models"
	"forum/pkg/cookie"
	"net/http"
)

func (h *handlers) signup(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		data := h.newTemplateData(r)
		data.Form = models.UserSignupForm{}
		data.AuthenticatedUser = h.authenticatedUser(r)

		h.render(w, http.StatusOK, "signup.html", data)
	} else if r.Method == http.MethodPost {
		data := h.newTemplateData(r)
		data.AuthenticatedUser = h.authenticatedUser(r)
		data, err := h.service.UserSignUp(data, r)
		if err != nil {
			h.ErrorHandler(w, r, http.StatusNotFound, http.StatusText(http.StatusNotFound))
			return
		}

		http.Redirect(w, r, "/user/login", http.StatusSeeOther) // redirect to login page
	} else {
		h.ErrorHandler(w, r, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
}

func (h *handlers) login(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		data := h.newTemplateData(r)
		data.Form = models.UserLoginForm{}
		data.AuthenticatedUser = h.authenticatedUser(r)

		h.render(w, http.StatusOK, "login.html", data)
	} else if r.Method == http.MethodPost {
		data := h.newTemplateData(r)
		data.AuthenticatedUser = h.authenticatedUser(r)
		data, id, err := h.service.UserLogin(data, r)
		if err != nil {
			h.ErrorHandler(w, r, http.StatusNotFound, http.StatusText(http.StatusNotFound))
			return
		}

		err = h.service.CreateSession(w, r, id)
		if err != nil {
			h.ErrorHandler(w, r, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther) // redirect to home page
	} else {
		h.ErrorHandler(w, r, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
}

func (h *handlers) logout(w http.ResponseWriter, r *http.Request) {
	sessionCookie, err := r.Cookie("session_token")

	if err == nil {
		err := h.service.DeleteSession(sessionCookie.Value)
		if err != nil {
			h.app.ServerError(w, err)
			return
		}
	}

	cookie.ExpireSessionCookie(w)
	// http.SetCookie(w, &http.Cookie{
	// 	Name:    "session_token",
	// 	Value:   "",
	// 	Expires: time.Unix(0, 0),
	// 	Path:    "/",
	// 	MaxAge:  -1,
	// })
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
