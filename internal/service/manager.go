package service

import (
	"forum/internal/repo"
	"forum/models"
	"net/http"
)

type service struct {
	repo repo.RepoI
}

type UserService interface {
	UserPosts(data *models.TemplateData, r *http.Request) (*models.TemplateData, error)
	UserSignUp(data *models.TemplateData, r *http.Request) (*models.TemplateData, error)
	UserLogin(data *models.TemplateData, r *http.Request) (*models.TemplateData, int, error)
}

type Render interface {
	MainLoader(data *models.TemplateData, r *http.Request) (*models.TemplateData, error)
	HomeUpdates(data *models.TemplateData, r *http.Request) (*models.TemplateData, error)
}

type PostsService interface {
	LatestPosts() ([]*models.Posts, error)
	ReactionDone(used_id, post_id, reaction int) error
	PostRender(data *models.TemplateData, r *http.Request) (*models.TemplateData, error)
	PostUpdate(data *models.TemplateData, r *http.Request) (*models.TemplateData, error)
	PostCreate(data *models.TemplateData, r *http.Request) (*models.TemplateData, int, error)
}

type SessionManager interface {
	IsValidToken(token string) (bool, error)
	GetUserIDBySessionToken(sessionToken string) int
	CreateSession(w http.ResponseWriter, r *http.Request, UserID int) error
	DeleteSession(sessionID string) error
}

type ServiceI interface {
	SessionManager
	PostsService
	Render
	UserService
}

func New(r repo.RepoI) ServiceI {
	return &service{
		r,
	}
}
