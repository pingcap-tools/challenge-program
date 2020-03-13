package manager

import (
	"fmt"
	"strings"

	"github.com/juju/errors"
	"github.com/ngaut/log"
	"github.com/pingcap/community/pkg/types"
)

const (
	tidbLink       = "[TiDB](https://github.com/pingcap/tidb/projects/26)"
	tikvLink       = "[TiKV](https://github.com/tikv/tikv/projects/20)"
	pdLink         = "[PD](https://github.com/pingcap/pd/projects/2)"
	chaosMeshLink  = "[Chaos Mesh](https://github.com/pingcap/chaos-mesh/projects/14)"
	dmLink         = "[DM](https://github.com/pingcap/dm/projects/1)"
	brLink         = "[BR](https://github.com/pingcap/br/projects/1)"
	clientRustLink = "[client-rust](https://github.com/tikv/client-rust/projects/3)"
	dashboardLink  = "[TiDB Dashboard](https://github.com/pingcap-incubator/tidb-dashboard/projects/17)"
	cherrybotLink  = "[cherry-bot](https://github.com/pingcap-incubator/cherry-bot/projects/1)"
)

func (mgr *Manager) GetCombinedRepoScore(repo *types.Repo, login string) (int, error) {
	if repo.GetOwner() == "pingcap" && repo.GetRepo() == "tidb" {
		return mgr.GetRepoScore(repo, login)
	}
	if repo.GetOwner() == "pingcap" && repo.GetRepo() == "dm" {
		return mgr.GetRepoScore(repo, login)
	}
	if repo.GetOwner() == "pingcap" && repo.GetRepo() == "pd" {
		return mgr.GetRepoScore(repo, login)
	}
	if repo.GetOwner() == "pingcap" && repo.GetRepo() == "br" {
		return mgr.GetRepoScore(repo, login)
	}
	if repo.GetOwner() == "pingcap" && repo.GetRepo() == "chaos-mesh" {
		return mgr.GetRepoScore(repo, login)
	}
	if repo.GetOwner() == "pingcap-incubator" && repo.GetRepo() == "tidb-dashboard" {
		return mgr.GetRepoScore(repo, login)
	}
	if repo.GetOwner() == "pingcap-incubator" && repo.GetRepo() == "cherry-bot" {
		return mgr.GetRepoScore(repo, login)
	}
	if repo.GetOwner() == "tikv" {
		return mgr.GetOrgScore("tikv", login)
	}
	return mgr.GetOrgScore("pingcap", login)
}

func (mgr *Manager) GetScoreReport(login string) (string, error) {
	var (
		b                     strings.Builder
		errs                  []error
		canPickMediumRepos    []string
		canNotPickMediumRepos []string
	)

	tidbScore, err := mgr.GetCombinedRepoScore(&types.Repo{Owner: "pingcap", Repo: "tidb"}, login)
	errs = append(errs, err)
	pdScore, err := mgr.GetCombinedRepoScore(&types.Repo{Owner: "pingcap", Repo: "pd"}, login)
	errs = append(errs, err)
	dmScore, err := mgr.GetCombinedRepoScore(&types.Repo{Owner: "pingcap", Repo: "dm"}, login)
	errs = append(errs, err)
	// chaosMeshScore, err := mgr.GetCombinedRepoScore(&types.Repo{Owner: "pingcap", Repo: "chaos-mesh"}, login)
	// errs = append(errs, err)
	tikvScore, err := mgr.GetCombinedRepoScore(&types.Repo{Owner: "tikv", Repo: "tikv"}, login)
	errs = append(errs, err)
	brScore, err := mgr.GetCombinedRepoScore(&types.Repo{Owner: "pingcap", Repo: "br"}, login)
	errs = append(errs, err)
	dashboardScore, err := mgr.GetCombinedRepoScore(&types.Repo{Owner: "pingcap-incubator", Repo: "tidb-dashboard"}, login)
	errs = append(errs, err)
	botScore, err := mgr.GetCombinedRepoScore(&types.Repo{Owner: "pingcap-incubator", Repo: "cherry-bot"}, login)
	errs = append(errs, err)

	for _, err := range errs {
		if err != nil {
			log.Error(err)
		}
	}

	if tidbScore >= 200 {
		canPickMediumRepos = append(canPickMediumRepos, tidbLink)
	} else {
		canNotPickMediumRepos = append(canNotPickMediumRepos, tidbLink)
	}
	if tikvScore >= 200 {
		canPickMediumRepos = append(canPickMediumRepos, tikvLink, clientRustLink)
	} else {
		canNotPickMediumRepos = append(canNotPickMediumRepos, tikvLink, clientRustLink)
	}
	// chaos mesh can skip easy tasks by default
	canPickMediumRepos = append(canPickMediumRepos, chaosMeshLink)
	// if chaosMeshScore >= 200 {
	// 	canPickMediumRepos = append(canPickMediumRepos, chaosMeshLink)
	// } else {
	// 	canNotPickMediumRepos = append(canNotPickMediumRepos, chaosMeshLink)
	// }
	if pdScore >= 200 {
		canPickMediumRepos = append(canPickMediumRepos, pdLink)
	} else {
		canNotPickMediumRepos = append(canNotPickMediumRepos, pdLink)
	}
	if dmScore >= 200 {
		canPickMediumRepos = append(canPickMediumRepos, dmLink)
	} else {
		canNotPickMediumRepos = append(canNotPickMediumRepos, dmLink)
	}
	if brScore >= 200 {
		canPickMediumRepos = append(canPickMediumRepos, brLink)
	} else {
		canNotPickMediumRepos = append(canNotPickMediumRepos, brLink)
	}
	if dashboardScore >= 200 {
		canPickMediumRepos = append(canPickMediumRepos, dashboardLink)
	} else {
		canNotPickMediumRepos = append(canNotPickMediumRepos, dashboardLink)
	}
	if botScore >= 200 {
		canPickMediumRepos = append(canPickMediumRepos, cherrybotLink)
	} else {
		canNotPickMediumRepos = append(canNotPickMediumRepos, cherrybotLink)
	}

	if len(canPickMediumRepos) == 0 {
		canPickMediumRepos = append(canPickMediumRepos, "null")
	}
	if len(canNotPickMediumRepos) == 0 {
		canNotPickMediumRepos = append(canNotPickMediumRepos, "null")
	}

	fmt.Fprintf(&b, "Your current score is: %d\n\n", tidbScore+pdScore+dmScore+tikvScore)
	fmt.Fprintf(&b, "Now you can pick up \"medium\" or \"hard\" task in these repos: %s\n\n", strings.Join(canPickMediumRepos, "/"))
	fmt.Fprintf(&b, "Please pick up \"easy\" task in these repos until you get 200 score: %s", strings.Join(canNotPickMediumRepos, "/"))

	return b.String(), nil
}

func (mgr *Manager) GetHistoryTeams(login string) ([]int, error) {
	var (
		currentTeam *types.Team
		teams       []int
	)

	currentTeam, err := mgr.GetTeamByUser(login, mgr.Config.Season)
	if err != nil {
		return []int{}, errors.Trace(err)
	}
	if currentTeam != nil {
		teams = append(teams, currentTeam.ID)
		for season := 1; season < mgr.Config.Season; season++ {
			team, err := mgr.GetSeasonTeamByName(season, currentTeam.Name)
			if err != nil {
				return []int{}, errors.Trace(err)
			}
			if team == nil {
				continue
			}
			teams = append(teams, team.ID)
		}
	}

	return teams, nil
}

func (mgr *Manager) GetRepoScore(repo *types.Repo, login string) (int, error) {
	teams, err := mgr.GetHistoryTeams(login)
	if err != nil {
		return 0, errors.Trace(err)
	}

	// userScore := types.UserScore{}
	// if err := mgr.storage.Scan(&types.Pick{}, &userScore,
	// 	"sum(score) as score", "status=? AND owner=? AND repo=? AND user=?",
	// 	"success", repo.GetOwner(), repo.GetRepo(), login); err != nil {
	// 	return 0, errors.Trace(err)
	// }

	teamScore := types.UserScore{}
	if len(teams) > 0 {
		if err := mgr.storage.Scan(&types.Pick{}, &teamScore,
			"sum(score) as score", "status=? AND owner=? AND repo=? AND teamID IN (?)",
			"success", repo.GetOwner(), repo.GetRepo(), teams); err != nil {
			return 0, errors.Trace(err)
		}
	}

	return teamScore.GetScore(), nil
}

func (mgr *Manager) GetOrgScore(owner, login string) (int, error) {
	teams, err := mgr.GetHistoryTeams(login)
	if err != nil {
		return 0, errors.Trace(err)
	}

	userScore := types.UserScore{}
	if err := mgr.storage.Scan(&types.Pick{}, &userScore,
		"sum(score) as score", "status=? AND owner=? AND user=?",
		"success", owner, login); err != nil {
		return 0, errors.Trace(err)
	}

	teamScore := types.UserScore{}
	if err := mgr.storage.Scan(&types.Pick{}, &teamScore,
		"sum(score) as score", "status=? AND owner=? AND teamID IN (?)",
		"success", owner, teams); err != nil {
		return 0, errors.Trace(err)
	}

	return userScore.GetScore() + teamScore.GetScore(), nil
}
