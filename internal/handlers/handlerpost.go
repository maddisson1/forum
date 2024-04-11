package handlers

import (
	"fmt"
	"forum/models"
	"net/http"
	"strconv"
)

func (h *handlers) postView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		h.ErrorHandler(w, r, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}
	if r.Method == http.MethodGet {
		data := h.newTemplateData(r)
		data, err = h.service.PostRender(data, r)
		if err != nil {
			h.ErrorHandler(w, r, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
		data.Form = models.CommentForm{}
		data.AuthenticatedUser = h.authenticatedUser(r)

		h.render(w, http.StatusOK, "viewpost.html", data)
	} else if r.Method == http.MethodPost {
		data := h.newTemplateData(r)
		data.AuthenticatedUser = h.authenticatedUser(r)

		data, err := h.service.PostUpdate(data, r)
		if err != nil {
			if err == http.ErrNoCookie {
				http.Redirect(w, r, "/user/login", http.StatusSeeOther)
				return
			}
			h.ErrorHandler(w, r, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
		http.Redirect(w, r, fmt.Sprintf("/post/view?id=%d", id), http.StatusSeeOther)
		return
	}
}

func (h *handlers) createPost(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		if r.URL.Path != "/post/create" {
			h.ErrorHandler(w, r, http.StatusNotFound, http.StatusText(http.StatusNotFound))
			return
		}
		data := h.newTemplateData(r)
		data.Form = models.PostForm{}
		data.AuthenticatedUser = h.authenticatedUser(r)
		h.render(w, http.StatusOK, "createpost.html", data)

	} else if r.Method == http.MethodPost {
		if r.URL.Path != "/post/create" {
			h.ErrorHandler(w, r, http.StatusNotFound, http.StatusText(http.StatusNotFound))
			return
		}

		data := h.newTemplateData(r)
		data.AuthenticatedUser = h.authenticatedUser(r)

		data, id, err := h.service.PostCreate(data, r)
		if err != nil {
			h.ErrorHandler(w, r, http.StatusNotFound, http.StatusText(http.StatusNotFound))
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/post/view?id=%d", id), http.StatusSeeOther)
	} else {
		h.ErrorHandler(w, r, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
}
