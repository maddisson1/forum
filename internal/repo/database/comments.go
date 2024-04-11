package database

import (
	"database/sql"
	"errors"
	"forum/models"
)

func (m *Storage) CreateComment(postid, userid int, text string) (int, error) {
	stmt := `
		INSERT INTO comments (user_id, post_id, comment, created_at)
		VALUES(?, ?, ?, DATETIME('now'))
	`

	result, err := m.DB.Exec(stmt, userid, postid, text)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), err
}

func (m *Storage) GetComment(id int) (*models.Comment, error) {
	stmt := `
		SELECT comment_id, user_id, post_id, comment, created_at FROM comments
		WHERE comment_id = ?
	`
	row := m.DB.QueryRow(stmt, id)
	var userid int
	s := &models.Comment{}
	err := row.Scan(&s.CommentId, &s.UserID, &s.PostID, &s.Text, &s.Created)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}
	stmt = `SELECT username FROM users WHERE user_id = ?`
	err = m.DB.QueryRow(stmt, userid).Scan(&s.Username)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (m *Storage) GetComments(id int) ([]*models.Comment, error) {
	stmt := `SELECT comment_id, user_id, post_id, comment, created_at 
	FROM comments
	WHERE post_id = ?
	`
	rows, err := m.DB.Query(stmt, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	comments := []*models.Comment{}
	for rows.Next() {
		s := &models.Comment{}
		err := rows.Scan(&s.CommentId, &s.UserID, &s.PostID, &s.Text, &s.Created)
		if err != nil {
			return nil, err
		}
		stmt = `SELECT username FROM users WHERE user_id = ?`
		err = m.DB.QueryRow(stmt, *&s.UserID).Scan(&s.Username)
		if err != nil {
			return nil, err
		}

		comments = append(comments, s)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}
