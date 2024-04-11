package service

import "net/http"

func (s service) GetUserIDBySessionToken(sessionToken string) int {
	id := s.repo.GetUserIDBySessionToken(sessionToken)
	return id
}

func (s service) DeleteSession(sessionID string) error {
	return s.repo.DeleteSession(sessionID)
}

func (s service) CreateSession(w http.ResponseWriter, r *http.Request, UserID int) error {
	return s.repo.CreateSession(w, r, UserID)
}

func (s service) IsValidToken(token string) (bool, error) {
	return s.repo.IsValidToken(token)
}
