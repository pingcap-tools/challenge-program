package manager

import (
	"github.com/juju/errors"
	"github.com/pingcap/community/pkg/types"
)

// GetScoreBySeason get user's season score
func (mgr *Manager) GetScoreBySeason(season int, login string, repo *types.Repo) (int, error) {
	team, err := mgr.GetTeamByUser(login, season)
	if err != nil {
		return 0, errors.Trace(err)
	}
	userScore := types.UserScore{}
	if team == nil {
		err := mgr.storage.Scan(&types.Pick{}, &userScore,
			"sum(score) as score", "season=? AND status=? AND owner=? AND repo=? AND user=?",
			season, "success", repo.GetOwner(), repo.GetRepo(), login)
		if err != nil {
			return 0, errors.Trace(err)
		}
	} else {
		err := mgr.storage.Scan(&types.Pick{}, &userScore,
			"sum(score) as score", "season=? AND status=? AND owner=? AND repo=? AND teamID=?",
			season, "success", repo.GetOwner(), repo.GetRepo(), team.GetID())
		if err != nil {
			return 0, errors.Trace(err)
		}
	}

	return userScore.GetScore(), nil
}
