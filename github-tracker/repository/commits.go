package repository

import (
	"context"
	"database/sql"
	"github-tracker/github-tracker/repository/entity"
)

type Commit interface {
	Insert(ctx context.Context, commit *entity.Commit) (err error)
	GetCommitsByAuthorEmail(ctx context.Context, email string) (commits []entity.Commit, err error) // para que el enginetin manager tenga info de quien hace los commits
}

type commit struct {
	Conn *sql.DB
}

func NewCommit(conn *sql.DB) commit {
	return commit{
		Conn: conn,
	}
}

func (m commit) Insert(ctx context.Context, commit *entity.Commit) (err error) {
	query := `
		INSERT INTO commits (repo_name, commit_id, commit_message, author_username, author_email, payload, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	//aca se proteje de ; en casos de concatenar, se usan los place holder para que no se concatene string por string
	// es decir que los ataque los tomara como string simples mas no como codigos de sql

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(
		ctx,
		commit.RepoName,
		commit.CommitID,
		commit.CommitMessage,
		commit.AuthorUsername,
		commit.AuthorEmail,
		commit.Payload,
		commit.CreatedAt,
		commit.UpdatedAt,
	).Err()

	return err
}

func (m commit) GetCommitsByAuthorEmail(ctx context.Context, email string) ([]entity.Commit, error) {
	query := `
        SELECT *
        FROM commits
        WHERE author_email = $1
    `

	rows, err := m.Conn.QueryContext(ctx, query, email)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var commits []entity.Commit
	for rows.Next() {
		var commit entity.Commit
		err := rows.Scan(
			&commit.ID,
			&commit.RepoName,
			&commit.CommitID,
			&commit.CommitMessage,
			&commit.AuthorUsername,
			&commit.AuthorEmail,
			&commit.Payload,
			&commit.CreatedAt,
			&commit.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		commits = append(commits, commit)
	}

	return commits, nil
}
