package manager

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/google/go-github/github"
	"github.com/jinzhu/gorm"
	"github.com/juju/errors"
	"github.com/ngaut/log"
	"github.com/pingcap/challenge-program/pkg/types"
)

const (
	difficultRegex               = "^-\\sdifficult:\\s([a-z]+)"
	scoreRegex                   = ".*score:\\s([1-9][0-9]+)"
	scoreLineRegex               = ".*?score:? ?([1-9][0-9]+).*"
	numberRegex                  = ".*?([1-9][0-9]+).*?"
	minScoreRegex                = "^-\\smin\\sscore\\srequired:\\s([0-9]+)"
	expiredRegex                 = "^-\\sexpired:\\s([0-9a-z]+)"
	requestIssueRegex            = `^req-pcp:\s.*`
	difficultPrefix              = `## difficult`
	difficultPrefixPropertyRegex = `[*|-]\s([a-z]+)`
	scorePrefix                  = `## score`
	scorePrefixPropertyRegex     = `[*|-]\s([0-9]+)`
	tidbReqComment               = "@zz-jason @winkyao @jackysp PTAL"
	tikvReqComment               = "@AndreMouche @nolouch @https://github.com/zhangjinpeng1987 PTAL"
)

var (
	teamID2Team            = map[int]*types.Team{}
	username2User          = map[string]*types.User{}
	difficultPattern       = regexp.MustCompile(difficultRegex)
	scorePattern           = regexp.MustCompile(scoreRegex)
	scoreLinePattern       = regexp.MustCompile(scoreLineRegex)
	numberPattern          = regexp.MustCompile(numberRegex)
	minScorePattern        = regexp.MustCompile(minScoreRegex)
	expiredPattern         = regexp.MustCompile(expiredRegex)
	difficultPrefixPattern = regexp.MustCompile(difficultPrefixPropertyRegex)
	scorePrefixPattern     = regexp.MustCompile(scorePrefixPropertyRegex)
)

func (mgr *Manager) GetAllRank(season int) ([]*types.RankItem, error) {
	picks, err := mgr.GetSuccessPicks(season)
	if err != nil {
		return []*types.RankItem{}, errors.Trace(err)
	}

	var rankItems types.Rank

	for _, pick := range picks {
		hasRank := false
		for _, rank := range rankItems {
			if pick.TeamID == 0 {
				// user
				if rank.Name == pick.User && rank.Season == pick.Season {
					hasRank = true
					rank.Score = rank.Score + pick.Score
					rank.LastUpdate = maxTime(rank.LastUpdate, pick.ClosedAt)
				}
			} else {
				// team
				team := mgr.getTeamByID(pick.TeamID)
				if rank.Name == team.Name {
					hasRank = true
					rank.Score = rank.Score + pick.Score
					rank.LastUpdate = maxTime(rank.LastUpdate, pick.ClosedAt)
				}
			}
		}

		if !hasRank {
			if pick.TeamID == 0 {
				user := mgr.getUserByUsername(pick.Season, pick.User)
				isMember, err := mgr.isMember(user.User)
				if err != nil {
					return []*types.RankItem{}, errors.Trace(err)
				}
				rankItems = append(rankItems, &types.RankItem{
					Type:       "user",
					Season:     pick.Season,
					Name:       pick.User,
					Community:  !isMember,
					Url:        user.IssueURL,
					Score:      pick.Score,
					LastUpdate: pick.ClosedAt,
				})
			} else {
				team := mgr.getTeamByID(pick.TeamID)
				community := true
				if team != nil {
					for _, user := range team.Users {
						isMember, err := mgr.isMember(user.GetUser())
						if err != nil {
							return []*types.RankItem{}, errors.Trace(err)
						}
						if isMember {
							community = false
						}
					}
				}
				rankItems = append(rankItems, &types.RankItem{
					Type:       "team",
					Season:     pick.Season,
					Name:       team.Name,
					Community:  community,
					Url:        team.IssueURL,
					Score:      pick.Score,
					LastUpdate: pick.ClosedAt,
				})
			}
		}
	}

	sort.Sort(rankItems)

	for index, r := range rankItems {
		r.Rank = index + 1
	}

	return rankItems, nil
}

func (mgr *Manager) GetSuccessPicks(season int) ([]*types.Pick, error) {
	picks, err := mgr.GetPicks(season)
	if err != nil {
		return picks, errors.Trace(err)
	}
	successPicks := []*types.Pick{}
	for _, pick := range picks {
		if pick.Status == "success" {
			successPicks = append(successPicks, pick)
		}
	}
	return successPicks, nil
}

func (mgr *Manager) getTeamByID(id int) *types.Team {
	team, ok := teamID2Team[id]
	if ok {
		return team
	}
	team, err := mgr.GetTeamByID(id)
	if err != nil {
		log.Error(err)
		return nil
	}
	teamID2Team[id] = team
	return team
}

func (mgr *Manager) getUserByUsername(season int, username string) *types.User {
	user, ok := username2User[username]
	if ok {
		return user
	}
	user, err := mgr.GetUserByLogin(username, season)
	if err != nil {
		log.Error(err)
		return nil
	}
	username2User[username] = user
	return user
}

// func (mgr *Manager)GetAllRank(season int) ([]*types.RankItem, error) {
// 	teamRank, err := mgr.GetTeamRank(season)
// 	if err != nil {
// 		return []*types.RankItem{}, errors.Trace(err)
// 	}
// 	userRank, err := mgr.GetUserRank(season)
// 	if err != nil {
// 		return []*types.RankItem{}, errors.Trace(err)
// 	}
// 	var rank types.Rank
// 	for _, t := range teamRank {
// 		team, err := mgr.GetTeamByID(t.GetTeamID())
// 		if err != nil {
// 			return []*types.RankItem{}, errors.Trace(err)
// 		}
// 		teamName := "Unknown Team"
// 		teamURL := ""
// 		if team != nil {
// 			teamName = team.GetName()
// 			teamURL = team.GetIssueURL()
// 		}
// 		community := true
// 		if team != nil {
// 			for _, user := range team.Users {
// 				isMember, err := mgr.isMember(user.GetUser())
// 				if err != nil {
// 					return []*types.RankItem{}, errors.Trace(err)
// 				}
// 				if isMember {
// 					community = false
// 				}
// 			}
// 		}
// 		rank = append(rank, &types.RankItem{
// 			Rank: 0,
// 			Type: "team",
// 			Name: teamName,
// 			Community: community,
// 			Url: teamURL,
// 			Score: t.GetScore(),
// 		})
// 	}
// 	for _, u := range userRank {
// 		isMember, err := mgr.isMember(u.GetName())
// 		if err != nil {
// 			return []*types.RankItem{}, errors.Trace(err)
// 		}
// 		rank = append(rank, &types.RankItem{
// 			Rank: 0,
// 			Type: "user",
// 			Name: u.GetName(),
// 			Community: !isMember,
// 			Url: fmt.Sprintf("https://github.com/%s", u.GetName()),
// 			Score: u.GetScore(),
// 		})
// 	}
// 	for index, r := range rank {
// 		r.Rank = index + 1
// 	}
// 	sort.Sort(rank)
// 	return rank, nil
// }

// func (mgr *Manager)GetTeamRank(season int) ([]*types.UserScore, error) {
// 	var rank []*types.UserScore
// 	err := mgr.storage.Group(&types.Pick{}, &rank,
// 		"sum(score) as score, teamID as team_id", "teamID",
// 		"season=? AND status=? AND teamID <> 0 AND teamID IN (SELECT id FROM teams WHERE season=? AND status=?)", season, "success", season, "opened")
// 	if err != nil {
// 		if gorm.IsRecordNotFoundError(err) {
// 			return []*types.UserScore{}, nil
// 		} else {
// 			return []*types.UserScore{}, errors.Trace(err)
// 		}
// 	}
// 	return rank, nil
// }

// func (mgr *Manager)GetUserRank(season int) ([]*types.UserScore, error) {
// 	var rank []*types.UserScore
// 	err := mgr.storage.Group(&types.Pick{}, &rank,
// 		"sum(score) as score, user as complete_user", "user",
// 		"season=? AND status=? AND user <> \"\" AND user IN (SELECT user FROM users WHERE season=? AND status=?)", season, "success", season, "opened")
// 	if err != nil {
// 		if gorm.IsRecordNotFoundError(err) {
// 			return []*types.UserScore{}, nil
// 		} else {
// 			return []*types.UserScore{}, errors.Trace(err)
// 		}
// 	}
// 	return rank, nil
// }

// func (mgr *Manager)GetUserRank(season int) ([]*types.UserScore, error) {
// 	var rank []*types.UserScore
// 	err := mgr.storage.Group(&types.Task{}, &rank,
// 		"sum(score) as score, complete_user", "complete_user",
// 		"season=? AND status=? AND complete_user <> ''", season, "success")
// 	if err != nil {
// 		if gorm.IsRecordNotFoundError(err) {
// 			return []*types.UserScore{}, nil
// 		} else {
// 			return []*types.UserScore{}, errors.Trace(err)
// 		}
// 	}
// 	var resRank []*types.UserScore
// 	// filter the users who already have team
// 	for _, r := range rank {
// 		_, err := mgr.GetTeamByUser(r.GetName(), mgr.Config.Season)
// 		if  err != nil {
// 			if gorm.IsRecordNotFoundError(err) {
// 				resRank = append(resRank, r)
// 			} else {
// 				return []*types.UserScore{}, errors.Trace(err)
// 			}
// 		}
// 	}
// 	return resRank, nil
// }

func (mgr *Manager) GetUserAllScore(login string) (int, error) {
	userScore := types.UserScore{}
	err := mgr.storage.Scan(&types.Pick{}, &userScore,
		"sum(score) as score",
		"WHERE status=? AND complete_user=?", "success", login)
	if err != nil {
		return 0, errors.Trace(err)
	}
	return userScore.GetScore(), nil
}

func (mgr *Manager) GetUserSeasonScore(login string, season int) (int, error) {
	team, err := mgr.GetTeamByUser(login, season)
	if err != nil {
		return 0, errors.Trace(err)
	}
	userScore := types.UserScore{}
	if team == nil {
		err := mgr.storage.Scan(&types.Pick{}, &userScore,
			"sum(score) as score",
			"status=? AND user=? AND season=?", "success", login, season)
		if err != nil {
			return 0, errors.Trace(err)
		}
	} else {
		err := mgr.storage.Scan(&types.Pick{}, &userScore,
			"sum(score) as score",
			"status=? AND teamID=? AND season=?", "success", team.GetID(), season)
		if err != nil {
			return 0, errors.Trace(err)
		}
	}
	return userScore.GetScore(), nil
}

func (mgr *Manager) GetUserExpectedSeasonScore(login string, season int) (int, error) {
	team, err := mgr.GetTeamByUser(login, season)
	if err != nil {
		return 0, errors.Trace(err)
	}
	userScore := types.UserScore{}
	if team == nil {
		err := mgr.storage.Scan(&types.Pick{}, &userScore,
			"sum(score) as score",
			"status<>? AND user=? AND season=?", "failed", login, season)
		if err != nil {
			return 0, errors.Trace(err)
		}
	} else {
		err := mgr.storage.Scan(&types.Pick{}, &userScore,
			"sum(score) as score",
			"status<>? AND teamID=? AND season=?", "failed", team.GetID(), season)
		if err != nil {
			return 0, errors.Trace(err)
		}
	}
	return userScore.GetScore(), nil
}

func (mgr *Manager) GetUserExpectedEasySeasonScore(login string, season int) (int, error) {
	team, err := mgr.GetTeamByUser(login, season)
	if err != nil {
		return 0, errors.Trace(err)
	}
	userScore := types.UserScore{}
	if team == nil {
		err := mgr.storage.Scan(&types.Pick{}, &userScore,
			"sum(score) as score",
			"status IN (?, ?) AND user=? AND season=? AND (task_id=0 OR task_id IN (SELECT id FROM tasks WHERE level=?))",
			"review", "success", login, season, "easy")
		if err != nil {
			return 0, errors.Trace(err)
		}
	} else {
		err := mgr.storage.Scan(&types.Pick{}, &userScore,
			"sum(score) as score",
			"status IN (?, ?) AND teamID=? AND season=? AND (task_id=0 OR task_id IN (SELECT id FROM tasks WHERE level=?))",
			"review", "success", team.GetID(), season, "easy")
		if err != nil {
			return 0, errors.Trace(err)
		}
	}
	return userScore.GetScore(), nil
}

func (mgr *Manager) GetSeasonScoreByRepo(login, owner, repo string) (int, error) {
	team, err := mgr.GetTeamByUser(login, mgr.Config.Season)
	if !gorm.IsRecordNotFoundError(err) {
		return 0, errors.Trace(err)
	}
	if team != nil {
		return mgr.GetTeamSeasonScoreByRepo(team, owner, repo)
	}
	return mgr.GetUserSeasonScoreByRepo(login, owner, repo)
}

func (mgr *Manager) GetUserSeasonScoreByRepo(login, owner, repo string) (int, error) {
	userScore := types.UserScore{}
	err := mgr.storage.Scan(&types.Pick{}, &userScore,
		"sum(score) as score",
		"status=? AND user=? AND season=? AND owner=? AND repo=?",
		"success", login, mgr.Config.Season, owner, repo)
	if err != nil {
		return 0, errors.Trace(err)
	}
	return userScore.GetScore(), nil
}

func (mgr *Manager) GetTeamSeasonScoreByRepo(team *types.Team, owner, repo string) (int, error) {
	userScore := types.UserScore{}
	err := mgr.storage.Scan(&types.Task{}, &userScore,
		"sum(score) as score",
		"status=? AND complete_team=? AND season=? AND owner=? AND repo=?",
		"success", team.GetID(), mgr.Config.Season, owner, repo)
	if err != nil {
		return 0, errors.Trace(err)
	}
	return userScore.GetScore(), nil
}

func (mgr *Manager) GetTaskById(id int) (*types.Task, error) {
	var task types.Task
	if err := mgr.storage.FindOne(&task, "id=?", id); err != nil {
		return nil, errors.NotFoundf("task number %d", id)
	}
	return &task, nil
}

func (mgr *Manager) GetTaskByNumber(owner, repo string, number, season int) (*types.Task, error) {
	var task types.Task
	if err := mgr.storage.FindOne(&task, "owner=? AND repo=? AND issue_number=? AND season=?", owner, repo, number, season); err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, errors.NotFoundf("task number %d", number)
	}
	return &task, nil
}

func (mgr *Manager) GetTasksByRepo(owner, repo string) ([]*types.Task, error) {
	var tasks []*types.Task
	if err := mgr.storage.Find(tasks, "owner=? AND repo=? AND season=?",
		owner, repo, mgr.Config.Season); err != nil {
		return []*types.Task{}, errors.Trace(err)
	}
	return tasks, nil
}

func (mgr *Manager) GetTasksByTaskgroup(groupID int) ([]*types.Task, error) {
	var tasks []*types.Task
	if err := mgr.storage.Find(&tasks, "taskgroup_id = ?", groupID); err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return []*types.Task{}, nil
		} else {
			return []*types.Task{}, errors.Trace(err)
		}
	}
	return tasks, nil
}

func (mgr *Manager) CreateTask(owner, repo string, issue *github.Issue) error {
	// r := regexp.MustCompile(requestIssueRegex)
	// if r.MatchString(strings.ToLower(issue.GetTitle())) {
	// 	switch repo {
	// 	case "tidb":
	// 		return errors.Trace(mgr.CommentIssue(owner, repo, issue.GetNumber(), tidbReqComment))
	// 	case "tikv":
	// 		return errors.Trace(mgr.CommentIssue(owner, repo, issue.GetNumber(), tikvReqComment))
	// 	}
	// }
	// isMember, err := mgr.isMember(issue.GetUser().GetLogin())
	// if err != nil {
	// 	return errors.Trace(err)
	// }
	// if !isMember {
	// 	return nil
	// }
	task, err := mgr.GetTaskByNumber(owner, repo, issue.GetNumber(), mgr.Config.Season)
	if err != nil {
		return errors.Trace(err)
	}
	if task == nil {
		task = &types.Task{
			Season:      mgr.Config.Season,
			Owner:       owner,
			Repo:        repo,
			IssueNumber: issue.GetNumber(),
			Status:      "opened",
		}
	}

	task.Title = issue.GetTitle()
	task.CreatedAt = issue.GetCreatedAt()
	task.Expired = ""

	oldScore := task.Score
	if !parseIssue(task, issue) {
		if owner != "tidb-challenge-program" || repo != "bug-hunting-issue" {
			log.Infof("Issue parse failed %d", issue.GetNumber())
			return nil
		}

		hasChallengeLabel := false
		for _, l := range issue.Labels {
			if l.GetName() == "challenge-program-2" {
				hasChallengeLabel = true
			}
		}

		if !hasChallengeLabel {
			return nil
		}
	}
	log.Infof("create task %s/%s#%d", owner, repo, issue.GetNumber())
	if err := mgr.UpdateTask(task); err != nil {
		return errors.Trace(err)
	}
	if owner == "tidb-challenge-program" && repo == "bug-hunting-issue" {
		// update score for picks
		if task.ID != 0 && task.Score != 0 && task.Score != oldScore {
			picks, err := mgr.GetPicksByTasks([]int{task.ID})
			if err != nil {
				return errors.Trace(err)
			}
			for _, pick := range picks {
				if pick.Status == "failed" {
					continue
				}
				if pick.Score == 0 {
					pick.Score = task.Score
					if err := mgr.UpdatePick(pick); err != nil {
						return errors.Trace(err)
					}
				}
			}
			if err := mgr.CommentIssue(owner, repo, issue.GetNumber(),
				fmt.Sprintf("This issue will be awarded %d points.", task.Score)); err != nil {
				return errors.Trace(err)
			}
		}
	}
	// if err := mgr.CreateProjectCard(task.GetLevel(), issue.GetHTMLURL()); err != nil {
	// 	log.Errorf("create project card failed, %v", errors.ErrorStack(err))
	// }
	return errors.Trace(mgr.AddLabel(owner, repo, issue, fmt.Sprintf("challenge-program-%d", mgr.Config.Season)))
}

func (mgr *Manager) UpdateTask(task *types.Task) error {
	return errors.Trace(mgr.storage.Save(task))
}

func parseIssue(task *types.Task, issue *github.Issue) bool {
	var (
	// ifDifficultMatch = false
	// ifMinScoreMatch  = false
	// ifScoreMatch     = false
	//ifExpiredMatch = false
	)

	var (
		lines    = strings.Split(strings.ReplaceAll(strings.ToLower(issue.GetBody()), "\r", ""), "\n")
		lastMark = ""
	)
	for _, line := range lines {
		line = strings.Trim(line, " ")
		if line == "" {
			continue
		}
		scoreLineMatch := scoreLinePattern.FindStringSubmatch(line)
		if len(scoreLineMatch) == 2 {
			score, err := parseScore(scoreLineMatch[1])
			if err != nil {
				log.Errorf("score parse error %+v", err)
				return false
			}
			task.Score = score
		}

		if lastMark == "score" {
			numberMatch := numberPattern.FindStringSubmatch(line)
			if len(numberMatch) == 2 {
				score, err := parseScore(numberMatch[1])
				if err != nil {
					log.Errorf("score parse error %+v", err)
					return false
				}
				task.Score = score
			}
		}

		if strings.HasSuffix(line, "score") && len(line) < 14 {
			lastMark = "score"
		} else {
			lastMark = ""
		}
	}

	if task.Score <= 300 {
		task.Level = "easy"
	} else if task.Score <= 10000 {
		task.Level = "medium"
	} else {
		task.Level = "hard"
	}

	return task.Score != 0
}

func parseScore(scoreStr string) (int, error) {
	return strconv.Atoi(scoreStr)
}

func maxTime(t1, t2 time.Time) time.Time {
	if t1.After(t2) {
		return t1
	}
	return t2
}
