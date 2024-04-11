package database

import (
	"database/sql"
	"errors"
	"forum/pkg/cookie"
	"net/http"
	"time"

	"github.com/gofrs/uuid"
)

const (
	insertSessionSQL  = `INSERT INTO sessions (user_id, session_token, expires_at) VALUES ($1, $2, $3);`
	getSessionSQL     = `SELECT session_token, expires_at FROM sessions WHERE user_id = ?;`
	deleteSessionSQL  = `DELETE FROM sessions WHERE session_token = ?;`
	deleteSessionByID = `DELETE FROM sessions WHERE user_id = ?`
	updateSessionSQL  = `UPDATE sessions SET session_token = ?, expires_at = ? WHERE user_id = ?`
)

func (sm *Storage) CreateSession(w http.ResponseWriter, r *http.Request, UserID int) error {
	var err error

	sessionToken, err := generateToken()
	if err != nil {
		return err
	}
	expiryAt := time.Now().Add(10 * time.Second)

	// Проверка наличия записи с UserID
	var existingSessionToken string
	err = sm.DB.QueryRow(getSessionSQL, UserID).Scan(&existingSessionToken)
	if err == sql.ErrNoRows {
		// Записи с UserID не существует, выполняем вставку новой записи
		_, err := sm.DB.Exec(insertSessionSQL, UserID, sessionToken, expiryAt)
		if err != nil {
			return err
		}
	} else {
		// Запись с UserID существует, выполняем обновление
		sessionToken, err = generateToken()
		if err != nil {
			return err
		}

		expiryAt = time.Now().Add(30 * time.Minute)
		_, err = sm.DB.Exec(updateSessionSQL, sessionToken, expiryAt, UserID)
		if err != nil {
			return err
		}
	}
	cookie.SetSessionCookie(w, sessionToken, expiryAt)

	return nil
}

func (sm *Storage) GetUserIDBySessionToken(sessionToken string) int {
	// Выполните запрос к базе данных, чтобы найти UserID по sessionToken
	var userID int
	err := sm.DB.QueryRow("SELECT user_id FROM sessions WHERE session_token = ?", sessionToken).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			// Сессия с указанным токеном не найдена
			return 0
		}
		// Обработка других ошибок базы данных
		return 0
	}

	// Возврат найденного UserID
	return userID
}

func (sm *Storage) IsValidToken(token string) (bool, error) {
	stmt := `SELECT expires_at FROM sessions WHERE session_token = ?`
	var expTimeStr string

	err := sm.DB.QueryRow(stmt, token).Scan(&expTimeStr)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	expTime, err := time.Parse("2006-01-02 15:04:05.999999999-07:00", expTimeStr)
	if err != nil {
		return false, err
	}

	if !expTime.After(time.Now()) {
		return false, nil
	}
	return true, nil
}

func (sm *Storage) DeleteSession(sessionID string) error {
	_, err := sm.DB.Exec(deleteSessionSQL, sessionID)
	return err
}

func (sm *Storage) DeleteSessionByID(userid string) error {
	_, err := sm.DB.Exec(deleteSessionByID, userid)
	return err
}

func generateToken() (string, error) {
	token, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	return token.String(), nil
}
