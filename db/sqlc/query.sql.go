// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1
// source: query.sql

package db

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/google/uuid"
)

const createOption = `-- name: CreateOption :exec
INSERT INTO options (pollid, optionkey, optionvalue)
VALUES ($1, $2, $3)
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

const createPoll = `-- name: CreatePoll :one
INSERT INTO polls (presentationid, question,pollindex)
VALUES ($1, $2, $3)
RETURNING id, presentationid, question, pollindex, createdat
`

type CreatePollParams struct {
	Presentationid uuid.UUID `db:"presentationid"`
	Question       string    `db:"question"`
	Pollindex      int32     `db:"pollindex"`
}

func (q *Queries) CreatePoll(ctx context.Context, arg CreatePollParams) (Poll, error) {
	row := q.db.QueryRowContext(ctx, createPoll, arg.Presentationid, arg.Question, arg.Pollindex)
	var i Poll
	err := row.Scan(
		&i.ID,
		&i.Presentationid,
		&i.Question,
		&i.Pollindex,
		&i.Createdat,
	)
	return i, err
}

const createPresentation = `-- name: CreatePresentation :one
INSERT INTO presentations (currentpollindex)
VALUES ($1)
RETURNING id
`

func (q *Queries) CreatePresentation(ctx context.Context, currentpollindex sql.NullInt32) (uuid.UUID, error) {
	row := q.db.QueryRowContext(ctx, createPresentation, currentpollindex)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const createPresentationAndPolls = `-- name: CreatePresentationAndPolls :one
WITH presentations_cte AS (
  INSERT INTO presentations (id, currentpollindex)
  VALUES (uuid_generate_v4(), 0)
  RETURNING id
),
data_cte AS (
  SELECT
    to_jsonb($1::jsonb) AS data
),
polls_cte AS (
  INSERT INTO polls (question, presentationid, pollindex)
  SELECT
    (arr.elem ->> 'question')::TEXT AS question,
    pc.id AS presentationid,
    arr.idx AS pollindex
  FROM presentations_cte pc, data_cte dc,
    jsonb_array_elements(dc.data) WITH ORDINALITY AS arr(elem, idx)
  RETURNING id, pollindex
)
SELECT presentations_cte.id
FROM presentations_cte
`

func (q *Queries) CreatePresentationAndPolls(ctx context.Context, dollar_1 json.RawMessage) (uuid.UUID, error) {
	row := q.db.QueryRowContext(ctx, createPresentationAndPolls, dollar_1)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const createVote = `-- name: CreateVote :exec
INSERT INTO votes (id, pollid, optionkey, clientid)
VALUES ($1, $2, $3, $4)
`

type CreateVoteParams struct {
	ID        uuid.UUID `db:"id"`
	Pollid    uuid.UUID `db:"pollid"`
	Optionkey string    `db:"optionkey"`
	Clientid  string    `db:"clientid"`
}

func (q *Queries) CreateVote(ctx context.Context, arg CreateVoteParams) error {
	_, err := q.db.ExecContext(ctx, createVote,
		arg.ID,
		arg.Pollid,
		arg.Optionkey,
		arg.Clientid,
	)
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

const getPoll = `-- name: GetPoll :one
SELECT id, presentationid, question, pollindex, createdat
FROM polls
WHERE id = $1
`

func (q *Queries) GetPoll(ctx context.Context, id uuid.UUID) (Poll, error) {
	row := q.db.QueryRowContext(ctx, getPoll, id)
	var i Poll
	err := row.Scan(
		&i.ID,
		&i.Presentationid,
		&i.Question,
		&i.Pollindex,
		&i.Createdat,
	)
	return i, err
}

const getPollByPID = `-- name: GetPollByPID :one
SELECT id, presentationid, question, pollindex, createdat
FROM polls
WHERE id = $1 and presentationid = $2
`

type GetPollByPIDParams struct {
	ID             uuid.UUID `db:"id"`
	Presentationid uuid.UUID `db:"presentationid"`
}

func (q *Queries) GetPollByPID(ctx context.Context, arg GetPollByPIDParams) (Poll, error) {
	row := q.db.QueryRowContext(ctx, getPollByPID, arg.ID, arg.Presentationid)
	var i Poll
	err := row.Scan(
		&i.ID,
		&i.Presentationid,
		&i.Question,
		&i.Pollindex,
		&i.Createdat,
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

const getPollVotes = `-- name: GetPollVotes :many
SELECT
  votes.optionkey,
  votes.clientid
FROM votes
JOIN polls ON votes.pollid = polls.id
JOIN presentations ON polls.presentationID = presentations.id
WHERE presentations.id = $1
AND polls.id = $2
ORDER BY votes.optionkey
`

type GetPollVotesParams struct {
	ID   uuid.UUID `db:"id"`
	ID_2 uuid.UUID `db:"id_2"`
}

type GetPollVotesRow struct {
	Optionkey string `db:"optionkey"`
	Clientid  string `db:"clientid"`
}

func (q *Queries) GetPollVotes(ctx context.Context, arg GetPollVotesParams) ([]GetPollVotesRow, error) {
	rows, err := q.db.QueryContext(ctx, getPollVotes, arg.ID, arg.ID_2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetPollVotesRow
	for rows.Next() {
		var i GetPollVotesRow
		if err := rows.Scan(&i.Optionkey, &i.Clientid); err != nil {
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

const getPollsCount = `-- name: GetPollsCount :one
SELECT COUNT(*) AS polls_count
FROM polls
WHERE presentationid = $1
`

func (q *Queries) GetPollsCount(ctx context.Context, presentationid uuid.UUID) (int64, error) {
	row := q.db.QueryRowContext(ctx, getPollsCount, presentationid)
	var polls_count int64
	err := row.Scan(&polls_count)
	return polls_count, err
}

const getPresentation = `-- name: GetPresentation :one
SELECT id, currentpollindex
FROM presentations
WHERE id = $1
`

func (q *Queries) GetPresentation(ctx context.Context, id uuid.UUID) (Presentation, error) {
	row := q.db.QueryRowContext(ctx, getPresentation, id)
	var i Presentation
	err := row.Scan(&i.ID, &i.Currentpollindex)
	return i, err
}

const getPresentationPolls = `-- name: GetPresentationPolls :many
SELECT id,question
FROM polls
WHERE presentationid = $1
ORDER BY createdat ASC
`

type GetPresentationPollsRow struct {
	ID       uuid.UUID `db:"id"`
	Question string    `db:"question"`
}

func (q *Queries) GetPresentationPolls(ctx context.Context, presentationid uuid.UUID) ([]GetPresentationPollsRow, error) {
	rows, err := q.db.QueryContext(ctx, getPresentationPolls, presentationid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetPresentationPollsRow
	for rows.Next() {
		var i GetPresentationPollsRow
		if err := rows.Scan(&i.ID, &i.Question); err != nil {
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

const getVote = `-- name: GetVote :many
SELECT v.id AS voteID, v.pollid, o.id AS optionid, o.optionkey, o.optionvalue,v.clientid
FROM votes AS v
JOIN options AS o ON v.optionkey = o.optionkey
WHERE v.pollid = $1
`

type GetVoteRow struct {
	Voteid      uuid.UUID `db:"voteid"`
	Pollid      uuid.UUID `db:"pollid"`
	Optionid    uuid.UUID `db:"optionid"`
	Optionkey   string    `db:"optionkey"`
	Optionvalue string    `db:"optionvalue"`
	Clientid    string    `db:"clientid"`
}

func (q *Queries) GetVote(ctx context.Context, pollid uuid.UUID) ([]GetVoteRow, error) {
	rows, err := q.db.QueryContext(ctx, getVote, pollid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetVoteRow
	for rows.Next() {
		var i GetVoteRow
		if err := rows.Scan(
			&i.Voteid,
			&i.Pollid,
			&i.Optionid,
			&i.Optionkey,
			&i.Optionvalue,
			&i.Clientid,
		); err != nil {
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

const updateCurrPollIndex = `-- name: UpdateCurrPollIndex :one
UPDATE presentations
SET currentpollindex = $1
WHERE id = $2
RETURNING currentpollindex
`

type UpdateCurrPollIndexParams struct {
	Currentpollindex sql.NullInt32 `db:"currentpollindex"`
	ID               uuid.UUID     `db:"id"`
}

func (q *Queries) UpdateCurrPollIndex(ctx context.Context, arg UpdateCurrPollIndexParams) (sql.NullInt32, error) {
	row := q.db.QueryRowContext(ctx, updateCurrPollIndex, arg.Currentpollindex, arg.ID)
	var currentpollindex sql.NullInt32
	err := row.Scan(&currentpollindex)
	return currentpollindex, err
}
