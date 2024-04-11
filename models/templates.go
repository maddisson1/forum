package models

type TemplateData struct {
	// isAuthenticated   bool
	AuthenticatedUser int
	CurrentYear       int
	Post              Posts
	Posts             []*Posts
	Form              any
	User              User
	Username          string
	Comments          []*Comment
	Likes             int
	Dislikes          int
	ErrorCode         int
	ErrorMsg          string
}
