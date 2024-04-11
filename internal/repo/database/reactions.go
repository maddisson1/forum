package database

import (
	"database/sql"
	"errors"
	"fmt"
	"forum/models"
)

func (m *Storage) CreateReaction(userid, postid, reaction int) error {
	stmt := `SELECT reaction_status FROM reactions WHERE post_id = ? AND user_id = ?`
	var is int
	err := m.DB.QueryRow(stmt, postid, userid).Scan(&is)
	if err != nil {
		if err == sql.ErrNoRows {
			stmt = `INSERT INTO reactions (user_id, post_id, reaction_status) 
			VALUES (?, ?, ?)`

			_, err = m.DB.Exec(stmt, userid, postid, reaction)
			if err != nil {
				fmt.Println(err)
				return err
			}

		} else {
			return err
		}
	} else {
		if is == reaction {
			stmt = `DELETE FROM reactions WHERE user_id = ? AND post_id = ?`
			_, err := m.DB.Exec(stmt, userid, postid)
			if err != nil {
				return err
			}
		} else {
			stmt = `UPDATE reactions SET reaction_status = ? WHERE user_id = ? AND post_id = ?`
			_, err := m.DB.Exec(stmt, reaction, userid, postid)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (m *Storage) GetLikes(postid int) (int, error) {
	stmt := `SELECT COUNT(reaction_status) FROM reactions WHERE reaction_status = 1 AND post_id = ?`

	row := m.DB.QueryRow(stmt, postid)

	var num int

	err := row.Scan(&num)
	if err != nil {
		fmt.Println(err)
		if errors.Is(err, sql.ErrNoRows) {
			return 0, models.ErrNoRecord
		} else {
			return 0, err
		}
	}
	return num, nil
}

func (m *Storage) GetDislikes(postid int) (int, error) {
	stmt := `SELECT COUNT(reaction_status) FROM reactions WHERE reaction_status = -1 AND post_id = ?`

	row := m.DB.QueryRow(stmt, postid)

	var num int

	err := row.Scan(&num)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, models.ErrNoRecord
		} else {
			return 0, err
		}
	}
	return num, nil
}
