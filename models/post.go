package models

import (
	"forum/pkg/validator"
	"time"
)

type Posts struct {
	ID       int
	Title    string
	Content  string
	Created  time.Time
	Username string
	Comment  []*Comment
	Likes    int
	Dislikes int
}

type Reaction struct {
	ReactionID     int
	UserID         int
	PostID         int
	ReactionStatus int
}

type Comment struct {
	CommentId int
	UserID    int
	Username  string
	PostID    int
	Text      string
	Created   time.Time
}

type CommentForm struct {
	Text   string
	Userid int
	validator.Validator
}

type PostForm struct {
	Title               string
	Content             string
	validator.Validator `form:"-"`
}
