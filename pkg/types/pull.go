package types

import (
	"time"

	"github.com/pingcap/challenge-program/util"
)

type Pull struct {
	ID          int       `gorm:"column:id"`
	Owner       string    `gorm:"column:owner"`
	Repo        string    `gorm:"column:repo"`
	Number      int       `gorm:"column:pull_number"`
	Title       string    `gorm:"column:title"`
	Body        string    `gorm:"column:body"`
	User        string    `gorm:"column:user"`
	Association string    `gorm:"column:association"`
	Relation    string    `gorm:"column:relation"`
	Label       string    `gorm:"column:label"`
	Status      string    `gorm:"column:status"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
	ClosedAt    time.Time `gorm:"column:closed_at"`
	MergedAt    time.Time `gorm:"column:merged_at"`
}

func (p *Pull) GetID() int {
	return p.ID
}

func (p *Pull) GetOwner() string {
	return p.Owner
}

func (p *Pull) GetRepo() string {
	return p.Repo
}

func (p *Pull) GetNumber() int {
	return p.Number
}

func (p *Pull) GetTitle() string {
	return p.Title
}

func (p *Pull) GetBody() string {
	return p.Body
}

func (p *Pull) GetUser() string {
	return p.User
}

func (p *Pull) GetAssociation() string {
	return p.Association
}

func (p *Pull) GetRelation() string {
	return p.Relation
}

func (p *Pull) GetLabel() []string {
	return util.ParseStringSlice(p.Label)
}

func (p *Pull) GetStatus() string {
	return p.Status
}

func (p *Pull) GetCreatedAt() time.Time {
	return p.CreatedAt
}

func (p *Pull) GetUpdatedAt() time.Time {
	return p.UpdatedAt
}

func (p *Pull) GetClosedAt() time.Time {
	return p.ClosedAt
}

func (p *Pull) GetMergedAt() time.Time {
	return p.MergedAt
}
