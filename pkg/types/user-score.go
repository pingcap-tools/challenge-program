package types

import "time"

type RankItem struct {
	Rank       int       `json:"rank"`
	Season     int       `json:"season"`
	Type       string    `json:"type"`
	Name       string    `json:"name"`
	Community  bool      `json:"community"`
	Url        string    `json:"url"`
	Score      int       `json:"score"`
	LastUpdate time.Time `json:"-"`
	DoingTask  string    `json:"doing-task"`
}

type Rank []*RankItem

type UserScore struct {
	Rank   int    `gorm:"column:rank"`
	Score  int    `gorm:"column:score"`
	Type   string `gorm:"column:data_type"`
	TeamID int    `gorm:"column:team_id"`
	Name   string `gorm:"column:complete_user"`
}

func (u *UserScore) GetUser() string {
	return u.Name
}

func (u *UserScore) GetScore() int {
	return u.Score
}

func (u *UserScore) GetTeamID() int {
	return u.TeamID
}

func (u *UserScore) GetName() string {
	return u.Name
}

func (r Rank) Len() int {
	return len(r)
}

func (r Rank) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func (r Rank) Less(i, j int) bool {
	if r[i].Score > r[j].Score {
		return true
	} else if r[i].Score < r[j].Score {
		return false
	}

	return r[i].LastUpdate.Before(r[j].LastUpdate)
}
