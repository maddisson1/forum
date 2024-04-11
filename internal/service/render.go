package service

import (
	"forum/models"
	"forum/pkg"
	"net/http"
)

func (s service) MainLoader(data *models.TemplateData, r *http.Request) (*models.TemplateData, error) {
	posts, err := s.LatestPosts()
	if err != nil {
		return nil, err
	}

	data.Posts = posts
	return data, nil
}

func (s service) HomeUpdates(data *models.TemplateData, r *http.Request) (*models.TemplateData, error) {
	posts, err := s.LatestPosts()
	if err != nil {
		return nil, err
	}
	cookie, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			return nil, err
		} else {
			return nil, err
		}
	}
	UserID := s.GetUserIDBySessionToken(cookie.Value)

	postID := pkg.SplitID(r.URL.Path)
	if r.Form.Has("reaction") {
		if r.PostForm.Get("reaction") == "1" {
			err := s.ReactionDone(UserID, postID, 1)
			if err != nil {
				return nil, err
			}
		} else if r.PostForm.Get("reaction") == "-1" {
			err := s.ReactionDone(UserID, postID, -1)
			if err != nil {
				return nil, err
			}
		}
	}
	data.Posts = posts
	return data, nil
}
