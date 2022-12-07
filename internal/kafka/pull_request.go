package kafka

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type PullRequestOpened struct {
	traceId             uuid.UUID
	id                  int
	openedAt            time.Time
	openedBy            string
	timeSinceLastCommit time.Duration
	totalCommitTime     time.Duration
}

func NewPullRequestOpened(id int, openedAt time.Time, openedBy string, timeSinceLastCommit, totalCommitTime time.Duration) *PullRequestOpened {
	return &PullRequestOpened{
		traceId:             uuid.New(),
		id:                  id,
		openedAt:            openedAt,
		openedBy:            openedBy,
		timeSinceLastCommit: timeSinceLastCommit,
		totalCommitTime:     totalCommitTime,
	}
}

func (pro *PullRequestOpened) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		TraceId             uuid.UUID     `json:trace_id`
		Id                  int           `json:id`
		OpenedAt            time.Time     `json:opened_at`
		OpenedBy            string        `json:opened_by`
		TimeSinceLastCommit time.Duration `json:time_since_last_commit`
		TotalCommitTime     time.Duration `json:total_commit_time`
	}{
		TraceId:             pro.traceId,
		Id:                  pro.id,
		OpenedAt:            pro.openedAt,
		OpenedBy:            pro.openedBy,
		TimeSinceLastCommit: pro.timeSinceLastCommit,
		TotalCommitTime:     pro.totalCommitTime,
	})
}
