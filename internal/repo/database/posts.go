package database

import (
	"database/sql"
	"errors"
	"forum/models"

	_ "github.com/mattn/go-sqlite3"
)

func (m *Storage) CreatePost(title string, content string, userid int) (int, error) {
	stmt := `
    INSERT INTO posts (title, content, created, user_id)
    VALUES(?, ?, DATETIME('now'), ?)
`
	result, err := m.DB.Exec(stmt, title, content, userid)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), err
}

func (m *Storage) GetPost(id int) (*models.Posts, error) {
	stmt := `SELECT post_id, title, content, created, user_id FROM posts
    WHERE post_id = ?`
	row := m.DB.QueryRow(stmt, id)
	var userid int
	s := &models.Posts{}
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &userid)
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

func (m *Storage) UserPosts(userid int) ([]*models.Posts, error) {
	stmt := `SELECT post_id, title, content, created FROM posts 
	WHERE user_id = ? 
	`
	rows, err := m.DB.Query(stmt, userid)
	if err != nil {
		return nil, err
	}
	var posts []*models.Posts
	for rows.Next() {
		s := &models.Posts{}
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created)
		if err != nil {
			return nil, err
		}
		posts = append(posts, s)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (m *Storage) LatestPosts() ([]*models.Posts, error) {
	stmt := `SELECT post_id, title, content, created, user_id FROM posts
    ORDER BY post_id DESC LIMIT 10`
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	var userid int
	defer rows.Close()
	posts := []*models.Posts{}
	for rows.Next() {
		s := &models.Posts{}
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &userid)
		if err != nil {
			return nil, err
		}
		stmt = `SELECT username FROM users WHERE user_id = ?`
		err = m.DB.QueryRow(stmt, userid).Scan(&s.Username)
		if err != nil {
			return nil, err
		}
		s.Likes, err = m.GetLikes(s.ID)
		if err != nil {
			return nil, err
		}
		s.Dislikes, err = m.GetDislikes(s.ID)
		if err != nil {
			return nil, err
		}

		posts = append(posts, s)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}
