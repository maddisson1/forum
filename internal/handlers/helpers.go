package handlers

import (
	"bytes"
	"fmt"
	"forum/models"
	"net/http"
	"time"
)

func (s *handlers) render(w http.ResponseWriter, status int, page string, data *models.TemplateData) {
	t, ok := s.app.TemplateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		s.app.ServerError(w, err)
		return
	}

	buf := new(bytes.Buffer)
	err := t.ExecuteTemplate(buf, "base", data)
	if err != nil {
		s.app.ServerError(w, err)
		return
	}

	buf.WriteTo(w)
}

func (s *handlers) newTemplateData(r *http.Request) *models.TemplateData {
	return &models.TemplateData{
		CurrentYear: time.Now().Year(),
	}
}

func (s *handlers) authenticatedUser(r *http.Request) int {
	sessionCookie, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			return 0
		}

		if err == sessionCookie.Valid() {
			return 0
		}
	}

	return s.service.GetUserIDBySessionToken(sessionCookie.Value)
}

func (h *handlers) ErrorHandler(w http.ResponseWriter, r *http.Request, errorCode int, msg string) {
	// t, err := template.ParseFiles("./ui/templates/pages/error.html")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	errorki := h.newTemplateData(r)
	errorki.ErrorCode = errorCode
	errorki.ErrorMsg = msg
	w.WriteHeader(errorCode)
	h.render(w, http.StatusOK, "error.html", errorki)
	// t.Execute(w, errorki)
}
