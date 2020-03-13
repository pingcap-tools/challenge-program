package types

import "time"

type CommentAttach struct {
	CommentType string
	Number      int
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Association string
}

type Comment struct {
	ID          int `gorm:"column:id"`
	Owner       string `gorm:"column:owner"`
	Repo        string `gorm:"column:repo"`
	CommentID   int64 `gorm:"column:comment_id"`
	CommentType string `gorm:"column:comment_type"`
	Number      int `gorm:"column:pull_number"`
	Body        string `gorm:"column:body"`
	User        string `gorm:"column:user"`
	Url         string `gorm:"column:url"`
	Association string `gorm:"column:association"`
	Relation    string `gorm:"column:relation"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
}

func (c *CommentAttach)GetCommentType() string {
	return c.CommentType
}

func (c *CommentAttach)GetNumber() int {
	return c.Number
}

func (c *CommentAttach)GetCreatedAt() time.Time {
	return c.CreatedAt
}

func (c *CommentAttach)GetUpdatedAt() time.Time {
	return c.UpdatedAt
}

func (c *CommentAttach)GetAuthorAssociation() string {
	return c.Association
}

func (c *Comment)GetID() int {
	return c.ID
}

func (c *Comment)GetOwner() string {
	return c.Owner
}

func (c *Comment)GetRepo() string {
	return c.Repo
}

func (c *Comment)GetCommentID() int64 {
	return c.CommentID
}

func (c *Comment)GetCreatedAt() time.Time {
	return c.CreatedAt
}

func (c *Comment)GetUpdatedAt() time.Time {
	return c.UpdatedAt
}
