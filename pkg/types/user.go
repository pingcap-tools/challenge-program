package types

type User struct {
	ID       int    `gorm:"column:id"`
	Season   int    `gorm:"column:season" json:"-"`
	User     string `gorm:"column:user" json:"user"`
	Email    string `gorm:"column:email" json:"-"`
	TeamID   int    `gorm:"column:team_id" json:"-"`
	IssueURL string `gorm:"column:issue_url" json:"-"`
	Status   string `gorm:"column:status" json:"-"`
	Avatar   string `gorm:"-" json:"avatar"`
	GitHub   string `gorm:"-" json:"github"`
	Leader   bool   `gorm:"column:leader" json:"leader"`
}

func (u *User) GetID() int {
	return u.ID
}

func (u *User) GetSeason() int {
	return u.Season
}

func (u *User) GetUser() string {
	return u.User
}

func (u *User) GetEmail() string {
	return u.Email
}

func (u *User) GetTeamID() int {
	return u.TeamID
}

func (u *User) GetIssueUrl() string {
	return u.IssueURL
}

func (u *User) GetStatus() string {
	return u.Status
}

type Team struct {
	ID       int    `gorm:"column:id"`
	Season   int    `gorm:"column:season"`
	Name     string `gorm:"column:name"`
	IssueURL string `gorm:"column:issue_url"`
	Users    []*User
	Status   string `gorm:"column:status"`
	Avatar   string `gorm:"-" json:"avatar"`
	GitHub   string `gorm:"-" json:"github"`
}

func (t *Team) GetID() int {
	return t.ID
}

func (t *Team) GetSeason() int {
	return t.Season
}

func (t *Team) GetName() string {
	return t.Name
}

func (t *Team) GetIssueURL() string {
	return t.IssueURL
}

func (t *Team) GetUsers() []*User {
	return t.Users
}
