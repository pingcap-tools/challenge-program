package types

import "time"

// Task struct
type Task struct {
	ID              int       `gorm:"column:id" json:"-"`
	TaskgroupID     int       `gorm:"column:taskgroup_id" json:"-"`
	Season          int       `gorm:"column:season" json:"season"`
	CompleteUser    string    `gorm:"column:complete_user" json:"-"`
	DoingUsers      []User    `gorm:"-" json:"inprogress-user"`
	APICompleteUser *User     `gorm:"-" json:"complete-user"`
	CompleteTeam    int       `gorm:"column:complete_team" json:"-"`
	DoingTeams      []Team    `gorm:"-" json:"inprogress-team"`
	APICompleteTeam *Team     `gorm:"-" json:"complete-team"`
	Owner           string    `gorm:"column:owner" json:"owner"`
	Repo            string    `gorm:"column:repo" json:"repo"`
	Title           string    `gorm:"column:title" json:"title"`
	IssueNumber     int       `gorm:"column:issue_number" json:"-"`
	Issue           string    `gorm:"-" json:"issue"`
	PullNumber      int       `gorm:"column:pull_number" json:"-"`
	Level           string    `gorm:"column:level" json:"level"`
	MinScore        int       `gorm:"column:min_score" json:"-"`
	Score           int       `gorm:"column:score" json:"score"`
	Status          string    `gorm:"column:status" json:"status"`
	CreatedAt       time.Time `gorm:"column:created_at" json:"-"`
	Expired         string    `gorm:"column:expired" json:"-"`
}

func (t *Task) GetID() int {
	return t.ID
}

func (t *Task) GetSeason() int {
	return t.Season
}

func (t *Task) GetCompleteUser() string {
	return t.CompleteUser
}

func (t *Task) GetCompleteTeam() int {
	return t.CompleteTeam
}

func (t *Task) GetOwner() string {
	return t.Owner
}

func (t *Task) GetRepo() string {
	return t.Repo
}

func (t *Task) GetTitle() string {
	return t.Title
}

func (t *Task) GetIssueNumber() int {
	return t.IssueNumber
}

func (t *Task) GetPullNumber() int {
	return t.PullNumber
}

func (t *Task) GetLevel() string {
	return t.Level
}

func (t *Task) GetScore() int {
	return t.Score
}

func (t *Task) GetMinScore() int {
	return t.MinScore
}

func (t *Task) GetStatus() string {
	return t.Status
}

func (t *Task) GetCreatedAt() time.Time {
	return t.CreatedAt
}

func (t *Task) GetExpired() string {
	return t.Expired
}
