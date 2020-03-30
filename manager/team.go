package manager

import (
	"github.com/jinzhu/gorm"
	"github.com/juju/errors"
	"github.com/pingcap/challenge-program/pkg/types"
)

func (mgr *Manager) GetTeamByUser(login string, season int) (*types.Team, error) {
	user, err := mgr.GetUserByLogin(login, season)
	if err != nil {
		return nil, err
	}
	if user == nil || user.GetTeamID() == 0 {
		return nil, nil
	}
	return mgr.GetTeamByID(user.GetTeamID())
}

func (mgr *Manager) GetSeasonTeamByName(season int, teamName string) (*types.Team, error) {
	var team types.Team
	if err := mgr.storage.FindOne(&team, "name=? AND season=? AND status=?",
		teamName, season, "opened"); err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		} else {
			return nil, errors.Trace(err)
		}
	}

	return mgr.GetTeamByID(team.ID)
}

func (mgr *Manager) GetTeamByID(id int) (*types.Team, error) {
	if id == 0 {
		return nil, nil
	}
	t := &types.Team{
		Users: []*types.User{},
	}
	if err := mgr.storage.FindOne(&t, "id=?", id); err != nil {
		return t, err
	}
	users, err := mgr.GetUsersByTeam(t.GetID())
	if err != nil {
		return nil, errors.Trace(err)
	}
	t.Users = users
	return t, nil
}

func (mgr *Manager) GetAllTeamsByLogin(login string) ([]int, error) {
	var (
		users   []*types.User
		teamIds []int
	)

	if err := mgr.storage.Find(&users, "team_id!=? and user=?", 0, login); err != nil {
		return teamIds, nil
	}

	for _, user := range users {
		teamIds = append(teamIds, user.GetTeamID())
	}
	return teamIds, nil
}
