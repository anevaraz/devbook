package repositories

import (
	"database/sql"
	"devbook/src/models"
)

type posts struct {
	db *sql.DB
}

func Posts(db *sql.DB) *posts {
	return &posts{db}
}

func (repository posts) Create(post models.Post) (uint64, error) {
	statement, err := repository.db.Prepare(
		"INSERT INTO posts (title, content, author_id) VALUES (?, ?, ?)",
	)
	if err != nil {
		return 0, err
	}
	defer statement.Close()

	result, err := statement.Exec(post.Title, post.Content, post.AuthorID)
	if err != nil {
		return 0, err
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint64(lastInsertID), nil
}

func (repository posts) Find(tokenUserID uint64) ([]models.Post, error) {

	rows, err := repository.db.Query(`
		SELECT
			DISTINCT p.*,
			u.nick
		FROM 
			posts p
		INNER JOIN users u ON
			u.id = p.author_id
		INNER JOIN followers f ON
			p.author_id = f.user_id
		WHERE
			u.id = ? OR f.follower_id = ?
		ORDER BY 1 DESC
		`, tokenUserID, tokenUserID,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var posts []models.Post

	for rows.Next() {
		var post models.Post

		if err = rows.Scan(
			&post.ID,
			&post.Title,
			&post.Content,
			&post.Likes,
			&post.AuthorID,
			&post.AuthorNick,
			&post.CreatedAt,
		); err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func (repository posts) FindOneById(postID uint64) (models.Post, error) {
	rows, err := repository.db.Query(`
		SELECT
			p.*
		FROM 
			posts p
		INNER JOIN users u ON
			u.id = p.author_id
		WHERE p.id = ?
		`, postID,
	)

	if err != nil {
		return models.Post{}, err
	}
	defer rows.Close()

	var post models.Post

	if rows.Next() {
		if err = rows.Scan(
			&post.ID,
			&post.Title,
			&post.Content,
			&post.Likes,
			&post.AuthorID,
			&post.AuthorNick,
			&post.CreatedAt,
		); err != nil {
			return models.Post{}, err
		}
	}

	return post, nil
}

func (repository posts) FindByUser(userID uint64) ([]models.Post, error) {

	rows, err := repository.db.Query(`
		SELECT
			p.*,
			u.nick
		FROM 
			posts p
		JOIN users u ON
			u.id = p.author_id
		WHERE
			p.author_id = ?
		`, userID,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var posts []models.Post

	for rows.Next() {
		var post models.Post

		if err = rows.Scan(
			&post.ID,
			&post.Title,
			&post.Content,
			&post.Likes,
			&post.AuthorID,
			&post.AuthorNick,
			&post.CreatedAt,
		); err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func (repository posts) Update(postID uint64, post models.Post) error {
	statement, err := repository.db.Prepare(
		"UPDATE posts SET title = ?, content = ? WHERE id = ?",
	)
	if err != nil {
		return err
	}
	defer statement.Close()

	if _, err = statement.Exec(post.Title, post.Content, postID); err != nil {
		return err
	}

	return nil
}

func (repository posts) Delete(postID uint64) error {
	statement, err := repository.db.Prepare("DELETE FROM posts WHERE id = ?")
	if err != nil {
		return err
	}
	defer statement.Close()

	if _, err = statement.Exec(postID); err != nil {
		return err
	}

	return nil
}

func (repository posts) Like(postID uint64) error {
	statement, err := repository.db.Prepare(
		"UPDATE posts SET likes = likes + 1 WHERE id = ?",
	)
	if err != nil {
		return err
	}
	defer statement.Close()

	if _, err = statement.Exec(postID); err != nil {
		return err
	}

	return nil
}

func (repository posts) Unlike(postID uint64) error {
	statement, err := repository.db.Prepare(`
		UPDATE 
			posts 
		SET 
			likes = CASE
						WHEN likes > 0 likes -1
						ELSE 0
					END
		WHERE id = ?
	`)
	if err != nil {
		return err
	}
	defer statement.Close()

	if _, err = statement.Exec(postID); err != nil {
		return err
	}

	return nil
}
