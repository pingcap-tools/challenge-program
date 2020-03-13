package types

import "time"

type Pick struct {
	ID         int `gorm:"column:id"`
	Season     int `gorm:"column:season"`
	Owner      string `gorm:"column:owner"`
	Repo       string `gorm:"column:repo"`
	TaskID     int `gorm:"column:task_id"`
	TeamID     int `gorm:"column:teamID"`
	User       string `gorm:"column:user"`
	PullNumber int `gorm:"column:pull_number"`
	Score      int `gorm:"column:score"`
	Status     string `gorm:"column:status"`
	CreatedAt  time.Time `gorm:"column:created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at"`
	ClosedAt   time.Time `gorm:"column:closed_at"`
}

func (p *Pick)GetID() int {
	return p.ID
}

func (p *Pick)GetSeason() int {
	return p.Season
}

func (p *Pick)GetOwner() string {
	return p.Owner
}

func (p *Pick)GetRepo() string {
	return p.Repo
}

func (p *Pick)GetTaskID() int {
	return p.TaskID
}

func (p *Pick)GetTeamID() int {
	return p.TeamID
}

func (p *Pick)GetUser() string {
	return p.User
}

func (p *Pick)GetPullNumber() int {
	return p.PullNumber
}

func (p *Pick)GetScore() int {
	return p.Score
}

func (p *Pick)GetStatus() string {
	return p.Status
}

func (p *Pick)GetCreatedAt() time.Time {
	return p.CreatedAt
}

func (p *Pick)GetUpdatedAt() time.Time {
	return p.UpdatedAt
}

func (p *Pick)GetClosedAt() time.Time {
	return p.ClosedAt
}
