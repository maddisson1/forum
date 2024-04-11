package repo

import (
	"forum/internal/repo/database"
	"forum/models"
	"net/http"
)

type PostsRepo interface {
	CreatePost(title string, content string, userid int) (int, error)
	GetPost(id int) (*models.Posts, error)
	UserPosts(userid int) ([]*models.Posts, error)
	LatestPosts() ([]*models.Posts, error)
}

type UsersRepo interface {
	CreateUser(username, email, password string) error
	Authenticate(email, password string) (int, error)
	Exitsts(id int) (bool, error)
	GetUser(id int) (string, error)
}

type ReactionsRepo interface {
	CreateReaction(userid, postid, reaction int) error
	GetLikes(postid int) (int, error)
	GetDislikes(postid int) (int, error)
}

type CommentsRepo interface {
	CreateComment(postid, userid int, text string) (int, error)
	GetComment(id int) (*models.Comment, error)
	GetComments(id int) ([]*models.Comment, error)
}

type SessionsRepo interface {
	IsValidToken(token string) (bool, error)
	GetUserIDBySessionToken(sessionToken string) int
	DeleteSession(sessionID string) error
	CreateSession(w http.ResponseWriter, r *http.Request, UserID int) error
}

type RepoI interface {
	PostsRepo
	UsersRepo
	ReactionsRepo
	CommentsRepo
	SessionsRepo
}

func New(storagePath string) (RepoI, error) {
	return database.New(storagePath)
}
