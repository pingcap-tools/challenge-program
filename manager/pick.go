package manager

import (
	"github.com/jinzhu/gorm"
	"github.com/juju/errors"
	"github.com/ngaut/log"
	"github.com/pingcap/challenge-program/pkg/types"
)

func (mgr *Manager) CreatePick(login string, task *types.Task) error {
	pick := types.Pick{
		Season: mgr.Config.Season,
		Owner:  task.GetOwner(),
		Repo:   task.GetRepo(),
		TaskID: task.GetID(),
		Score:  task.GetScore(),
		Status: "doing",
	}
	team, err := mgr.GetTeamByUser(login, mgr.Config.Season)
	if err != nil {
		return errors.Trace(err)
	}
	// team task or personal task
	if team == nil {
		pick.User = login
	} else {
		pick.TeamID = team.GetID()
	}
	log.Info(pick)
	return errors.Trace(mgr.UpdatePick(&pick))
}

func (mgr *Manager) UpdatePick(pick *types.Pick) error {
	return errors.Trace(mgr.storage.Save(pick))
}

func (mgr *Manager) HasDoingPick(login string) (*types.Pick, error) {
	var pick *types.Pick
	var err error
	team, er := mgr.GetTeamByUser(login, mgr.Config.Season)
	if er != nil {
		return nil, errors.Trace(er)
	}
	if team == nil {
		pick, err = mgr.GetPickByLogin(login, "doing")
	} else {
		pick, err = mgr.GetPickByTeam(team, "doing")
	}
	if err != nil {
		return nil, errors.Trace(err)
	}
	return pick, nil
}

func (mgr *Manager) DoingPickOverLimit(login string) (bool, []*types.Pick, error) {
	var (
		team      *types.Team
		picks     []*types.Pick
		overLimit = false
		err       error
	)
	team, err = mgr.GetTeamByUser(login, mgr.Config.Season)
	if err != nil {
		return false, picks, errors.Trace(err)
	}
	if team == nil {
		picks, err = mgr.GetDoingPicksByLogin(login)
		if len(picks) >= 1 {
			overLimit = true
		}
	} else {
		picks, err = mgr.GetDoingPicksByTeam(team)
		if len(picks) >= len(team.Users) {
			overLimit = true
		}
	}

	if err != nil {
		return false, picks, errors.Trace(err)
	}

	return overLimit, picks, nil
}

func (mgr *Manager) GetPickByLogin(login, status string) (*types.Pick, error) {
	pick := types.Pick{}
	if err := mgr.storage.FindOne(&pick, "season=? AND status=? AND user=?",
		mgr.Config.Season, status, login); err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		} else {
			return nil, errors.Trace(err)
		}
	}
	return &pick, nil
}

func (mgr *Manager) GetDoingPicksByLogin(login string) ([]*types.Pick, error) {
	var picks []*types.Pick
	if err := mgr.storage.Find(&picks, "season=? AND user=? AND status=?",
		mgr.Config.Season, login, "doing"); err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return []*types.Pick{}, nil
		} else {
			return []*types.Pick{}, errors.Trace(err)
		}
	}
	return picks, nil
}

func (mgr *Manager) GetPicks(season int) ([]*types.Pick, error) {
	var (
		picks []*types.Pick
		err   error
	)

	if season == 0 {
		err = mgr.storage.Find(&picks, "")
	} else {
		err = mgr.storage.Find(&picks, "season=?", season)
	}

	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return []*types.Pick{}, nil
		}
		return []*types.Pick{}, errors.Trace(err)
	}
	return picks, nil
}

func (mgr *Manager) GetPicksByLogin(login string) ([]*types.Pick, error) {
	var picks []*types.Pick
	if err := mgr.storage.Find(&picks, "season=? AND user=?",
		mgr.Config.Season, login); err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return []*types.Pick{}, nil
		} else {
			return []*types.Pick{}, errors.Trace(err)
		}
	}
	return picks, nil
}

func (mgr *Manager) GetPickByLoginTask(login, status string, taskID int) (*types.Pick, error) {
	pick := types.Pick{}
	if err := mgr.storage.FindOne(&pick, "season=? AND status=? AND user=? AND task_id=?",
		mgr.Config.Season, status, login, taskID); err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		} else {
			return nil, errors.Trace(err)
		}
	}
	return &pick, nil
}

func (mgr *Manager) GetPickByPullNumber(login, status string, pullNumber int) (*types.Pick, error) {
	pick := types.Pick{}
	if err := mgr.storage.FindOne(&pick, "season=? AND status=? AND user=? AND pull_number=?",
		mgr.Config.Season, status, login, pullNumber); err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		} else {
			return nil, errors.Trace(err)
		}
	}
	return &pick, nil
}

func (mgr *Manager) GetTeamPickByPullNumber(teamID int, status string, pullNumber int) (*types.Pick, error) {
	pick := types.Pick{}
	if err := mgr.storage.FindOne(&pick, "season=? AND status=? AND teamID=? AND pull_number=?",
		mgr.Config.Season, status, teamID, pullNumber); err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		} else {
			return nil, errors.Trace(err)
		}
	}
	return &pick, nil
}

func (mgr *Manager) GetPickByTeam(team *types.Team, status string) (*types.Pick, error) {
	pick := types.Pick{}
	if err := mgr.storage.FindOne(&pick, "season=? AND status=? AND teamID=?",
		mgr.Config.Season, status, team.GetID()); err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		} else {
			return nil, errors.Trace(err)
		}
	}
	return &pick, nil
}

func (mgr *Manager) GetDoingPickByTask(task *types.Task) (*types.Pick, error) {
	pick := types.Pick{}
	if err := mgr.storage.FindOne(&pick, "season=? AND status=? AND task_id=?",
		mgr.Config.Season, "doing", task.GetID()); err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		} else {
			return nil, errors.Trace(err)
		}
	}
	return &pick, nil
}

func (mgr *Manager) GetPicksByTeam(team *types.Team) ([]*types.Pick, error) {
	var picks []*types.Pick
	if err := mgr.storage.Find(&picks, "season=? AND teamID=?",
		mgr.Config.Season, team.GetID()); err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return []*types.Pick{}, nil
		} else {
			return []*types.Pick{}, errors.Trace(err)
		}
	}
	return picks, nil
}

func (mgr *Manager) GetDoingPicks() ([]*types.Pick, error) {
	var picks []*types.Pick
	if err := mgr.storage.Find(&picks, "season=? AND status=?",
		mgr.Config.Season, "doing"); err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return []*types.Pick{}, nil
		} else {
			return []*types.Pick{}, errors.Trace(err)
		}
	}
	return picks, nil
}

func (mgr *Manager) GetDoingPicksByTeam(team *types.Team) ([]*types.Pick, error) {
	var picks []*types.Pick
	if err := mgr.storage.Find(&picks, "season=? AND teamID=? AND status=?",
		mgr.Config.Season, team.GetID(), "doing"); err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return []*types.Pick{}, nil
		} else {
			return []*types.Pick{}, errors.Trace(err)
		}
	}
	return picks, nil
}

func (mgr *Manager) GetPickByPullWithRange(pullNumber int, taskRange []int) (*types.Pick, error) {
	pick := types.Pick{}
	if err := mgr.storage.FindOne(&pick, "season=? AND status=? AND pull_number=? AND task_id IN (?)",
		mgr.Config.Season, "doing", pullNumber, taskRange); err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		} else {
			return nil, errors.Trace(err)
		}
	}
	return &pick, nil
}

func (mgr *Manager) GetRepoPickByNumber(owner, repo string, number int) (*types.Pick, error) {
	pick := types.Pick{}
	if err := mgr.storage.FindOne(&pick, "owner=? AND repo=? AND pull_number=?",
		owner, repo, number); err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		} else {
			return nil, errors.Trace(err)
		}
	}
	return &pick, nil
}

func (mgr *Manager) GetPicksByTasks(tasks []int) ([]*types.Pick, error) {
	var picks []*types.Pick
	if err := mgr.storage.Find(&picks, "task_id IN (?)", tasks); err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return []*types.Pick{}, nil
		} else {
			return []*types.Pick{}, errors.Trace(err)
		}
	}
	return picks, nil
}
