package service

import (
	"forum/models"
	"forum/pkg/validator"
	"net/http"
	"strconv"
)

func (s service) LatestPosts() ([]*models.Posts, error) {
	posts, err := s.repo.LatestPosts()
	if err != nil {
		return nil, err
	}
	return posts, err
}

func (s service) ReactionDone(used_id, post_id, reaction int) error {
	return s.repo.CreateReaction(used_id, post_id, reaction)
}

func (s *service) PostRender(data *models.TemplateData, r *http.Request) (*models.TemplateData, error) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		return nil, err
	}

	if r.URL.Path != "/post/view" {
		return nil, err
	}

	post, err := s.repo.GetPost(id)
	if err != nil {
		if err == models.ErrNoRecord {
			return nil, err
		} else {
			return nil, err
		}
	}
	comments, err := s.repo.GetComments(id)
	if err != nil {
		return nil, err
	}

	likes, err := s.repo.GetLikes(id)
	if err != nil {
		return nil, err
	}
	dislikes, err := s.repo.GetDislikes(id)
	if err != nil {
		return nil, err
	}
	data.Post = *post
	data.Post.Likes = likes
	data.Post.Dislikes = dislikes
	data.Comments = comments
	return data, nil
}

func (s *service) PostUpdate(data *models.TemplateData, r *http.Request) (*models.TemplateData, error) {
	data, err := s.PostRender(data, r)
	if err != nil {
		return nil, err
	}
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		return nil, err
	}
	err = r.ParseForm()
	if err != nil {
		return nil, err
	}
	post, err := s.repo.GetPost(id)
	if err != nil {
		if err == models.ErrNoRecord {
			return nil, err
		} else {
			return nil, err
		}
	}
	data.Comments, err = s.repo.GetComments(id)
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
	if r.PostForm.Get("reaction") == "1" {
		err := s.repo.CreateReaction(UserID, id, 1)
		if err != nil {
			return nil, err
		}
	} else if r.PostForm.Get("reaction") == "-1" {
		err := s.repo.CreateReaction(UserID, id, -1)
		if err != nil {
			return nil, err
		}
	}

	if r.Form.Has("comment") {
		Form := &models.CommentForm{
			Text:   r.FormValue("comment"),
			Userid: UserID,
		}

		Form.CheckField(validator.NotBlank(Form.Text), "comment", "This field cannot be blank")
		Form.CheckField(validator.MaxChars(Form.Text, 50), "comment", "This field cannot be more than 100 characters long")

		if !Form.Valid() {
			data.Form = Form
			data.Post = *post
			return data, err
		}

		data.Form = Form
		data.Post = *post

		s.repo.CreateComment(id, Form.Userid, Form.Text)
	}

	return data, nil
}

func (s *service) PostCreate(data *models.TemplateData, r *http.Request) (*models.TemplateData, int, error) {
	err := r.ParseForm()
	if err != nil {
		return nil, 0, err
	}

	Form := &models.PostForm{
		Title:   r.FormValue("title"),
		Content: r.FormValue("content"),
	}

	Form.CheckField(validator.NotBlank(Form.Title), "title", "This field cannot be blank")
	Form.CheckField(validator.MaxChars(Form.Title, 100), "title", "This field cannot be more than 100 characters long")
	Form.CheckField(validator.NotBlank(Form.Content), "content", "This field cannot be blank")

	if !Form.Valid() {
		data.Form = Form
		return data, 0, err
	}
	id, err := s.repo.CreatePost(Form.Title, Form.Content, data.AuthenticatedUser)
	if err != nil {
		return nil, 0, err
	}
	data.Form = Form
	return data, id, nil
}
