package handlers

import (
	"net/http"
)

func (h *handlers) home(w http.ResponseWriter, r *http.Request) {
	data := h.newTemplateData(r)
	data.AuthenticatedUser = h.authenticatedUser(r)
	data, err := h.service.MainLoader(data, r)
	if err != nil {
		h.ErrorHandler(w, r, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	h.render(w, http.StatusOK, "index.html", data)
}

func (h *handlers) homePost(w http.ResponseWriter, r *http.Request) {
	data := h.newTemplateData(r)
	data, err := h.service.HomeUpdates(data, r)
	if err != nil {
		if err == http.ErrNoCookie {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		} else {
			h.ErrorHandler(w, r, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
	}
	data.AuthenticatedUser = h.authenticatedUser(r)

	h.render(w, http.StatusOK, "index.html", data)
}

func (h *handlers) userPosts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.ErrorHandler(w, r, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
	if r.URL.Path != "/user/posts" {
		h.ErrorHandler(w, r, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}

	data := h.newTemplateData(r)
	data.AuthenticatedUser = h.authenticatedUser(r)
	data, err := h.service.UserPosts(data, r)
	if err != nil {
		h.ErrorHandler(w, r, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	h.render(w, http.StatusOK, "userposts.html", data)
}
