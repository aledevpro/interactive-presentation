// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1
// source: option_query.sql

package db

import (
	"context"

	"github.com/google/uuid"
)

const createOption = `-- name: CreateOption :exec
INSERT INTO options
  (pollid, optionkey, optionvalue)
VALUES
  ($1, $2, $3)
`

type CreateOptionParams struct {
	Pollid      uuid.UUID `db:"pollid"`
	Optionkey   string    `db:"optionkey"`
	Optionvalue string    `db:"optionvalue"`
}

func (q *Queries) CreateOption(ctx context.Context, arg CreateOptionParams) error {
	_, err := q.db.ExecContext(ctx, createOption, arg.Pollid, arg.Optionkey, arg.Optionvalue)
	return err
}

const getOptionByKey = `-- name: GetOptionByKey :one
SELECT id, pollid, optionkey, optionvalue
FROM options
WHERE optionkey = $1
`

func (q *Queries) GetOptionByKey(ctx context.Context, optionkey string) (Option, error) {
	row := q.db.QueryRowContext(ctx, getOptionByKey, optionkey)
	var i Option
	err := row.Scan(
		&i.ID,
		&i.Pollid,
		&i.Optionkey,
		&i.Optionvalue,
	)
	return i, err
}

const getPollOptions = `-- name: GetPollOptions :many
SELECT optionkey, optionvalue
FROM options
WHERE pollid = $1
ORDER BY optionkey
`

type GetPollOptionsRow struct {
	Optionkey   string `db:"optionkey"`
	Optionvalue string `db:"optionvalue"`
}

func (q *Queries) GetPollOptions(ctx context.Context, pollid uuid.UUID) ([]GetPollOptionsRow, error) {
	rows, err := q.db.QueryContext(ctx, getPollOptions, pollid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetPollOptionsRow
	for rows.Next() {
		var i GetPollOptionsRow
		if err := rows.Scan(&i.Optionkey, &i.Optionvalue); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
