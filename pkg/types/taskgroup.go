package types

import "time"

// Taskgroup struct
type Taskgroup struct {
	ID          int       `gorm:"column:id" json:"-"`
	Season      int       `gorm:"column:season" json:"season"`
	Owner       string    `gorm:"column:owner" json:"owner"`
	Repo        string    `gorm:"column:repo" json:"repo"`
	Title       string    `gorm:"column:title" json:"title"`
	IssueNumber int       `gorm:"column:issue_number" json:"issue-number"`
	Issue       string    `gorm:"-" json:"issue"`
	Bonus       int       `gorm:"column:bonus" json:"bonus"`
	Progress    int       `gorm:"-" json:"progress"`
	Vote        int       `gorm:"column:vote" json:"vote"`
	DoingUsers  []*User   `gorm:"-" json:"doing-users"`
	Tasks       []*Task   `gorm:"-" json:"tasks"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"-"`
}

type SortByVote []*Taskgroup

func (r SortByVote) Len() int {
	return len(r)
}

func (r SortByVote) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func (r SortByVote) Less(i, j int) bool {
	return r[i].Vote > r[j].Vote
}
