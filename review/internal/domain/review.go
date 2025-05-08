package domain

import (
	"fmt"
	"time"
)

const (
	MaxScore = 10
	MinScore = 0
)

var ErrInvalidScore = fmt.Errorf("score out of range: %d <= Score <= %d", MinScore, MaxScore)

func NewScore(s int) (Score, error) {
	if s < MinScore || s > MaxScore {
		return Score(0), ErrInvalidScore
	}

	return Score(s), nil
}

type Score int

type Review struct {
	ID        int
	UserID    int
	BookID    int
	CreatedAt time.Time
	Title     string
	Text      string
	Score     Score
}
