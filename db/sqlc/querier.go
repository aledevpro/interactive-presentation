// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1

package db

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/google/uuid"
)

type Querier interface {
	CreateOption(ctx context.Context, arg CreateOptionParams) error
	CreatePoll(ctx context.Context, arg CreatePollParams) (Poll, error)
	CreatePresentation(ctx context.Context, currentpollindex sql.NullInt32) (uuid.UUID, error)
	CreatePresentationAndPolls(ctx context.Context, dollar_1 json.RawMessage) (uuid.UUID, error)
	CreateVote(ctx context.Context, arg CreateVoteParams) error
	GetOptionByKey(ctx context.Context, optionkey string) (Option, error)
	GetPoll(ctx context.Context, id uuid.UUID) (Poll, error)
	GetPollByPID(ctx context.Context, arg GetPollByPIDParams) (Poll, error)
	GetPollOptions(ctx context.Context, pollid uuid.UUID) ([]GetPollOptionsRow, error)
	GetPollVotes(ctx context.Context, arg GetPollVotesParams) ([]GetPollVotesRow, error)
	GetPollsCount(ctx context.Context, presentationid uuid.UUID) (int64, error)
	GetPresentation(ctx context.Context, id uuid.UUID) (Presentation, error)
	GetPresentationCurrentPoll(ctx context.Context, presentationid uuid.UUID) (GetPresentationCurrentPollRow, error)
	GetPresentationCurrentPoll2(ctx context.Context, presentationid uuid.UUID) (GetPresentationCurrentPoll2Row, error)
	GetPresentationPolls(ctx context.Context, presentationid uuid.UUID) ([]GetPresentationPollsRow, error)
	GetVote(ctx context.Context, pollid uuid.UUID) ([]GetVoteRow, error)
	MoveBackwardToPreviousPoll(ctx context.Context, presentationid uuid.UUID) (MoveBackwardToPreviousPollRow, error)
	MoveForwardToNextPoll(ctx context.Context, presentationid uuid.UUID) (MoveForwardToNextPollRow, error)
	UpdateCurrPollIndexBackward(ctx context.Context, presentationid uuid.UUID) (Presentation, error)
	UpdateCurrPollIndexForward(ctx context.Context, presentationid uuid.UUID) (Presentation, error)
}

var _ Querier = (*Queries)(nil)
