package manager

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/google/go-github/github"
	"github.com/jinzhu/gorm"
	"github.com/juju/errors"
	"github.com/pingcap/community/pkg/types"
)

const (
	signUpTeamRegex            = `^-\sTeam Name:\s(.*)`
	signUpTeamLeaderRegex      = `^-\sTeam Leader:\s@(.*)`
	signUpTeamMemberRegex      = `^-\sTeam Member:\s@(.*)`
	signUpTeamLeaderEmailRegex = `^-\sTeam Leader:\s@(.*)\((.*)\)`
	signUpTeamMemberEmailRegex = `^-\sTeam Member:\s@(.*)\((.*)\)`
	projectUrl                 = `https://github.com/pingcap/community/blob/master/challenge-programs/README-CN.md`
	cannotInheriteComment      = `You can not use this team name.`
)

var (
	teamPattern        = regexp.MustCompile(signUpTeamRegex)
	memberPattern      = regexp.MustCompile(signUpTeamMemberRegex)
	leaderPattern      = regexp.MustCompile(signUpTeamLeaderRegex)
	memberEmailPattern = regexp.MustCompile(signUpTeamMemberEmailRegex)
	leaderEmailPattern = regexp.MustCompile(signUpTeamLeaderEmailRegex)
)

func (mgr *Manager) GetUserById(id, season int) (*types.User, error) {
	var user types.User
	if err := mgr.storage.FindOne(&user, "id=? AND season=?", id, season); err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		} else {
			return nil, errors.Trace(err)
		}
	}
	return &user, nil
}

func (mgr *Manager) GetUserByLogin(login string, season int) (*types.User, error) {
	var user types.User
	if err := mgr.storage.FindOne(&user, "user=? AND season=? AND status=?",
		login, season, "opened"); err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		} else {
			return nil, errors.Trace(err)
		}
	}
	return &user, nil
}

func (mgr *Manager) GetUsersByIssueURL(login string, issueURL string, season int) ([]*types.User, error) {
	var users []*types.User
	if err := mgr.storage.Find(&users, "issue_url=? AND season=? AND status=?",
		issueURL, season, "opened"); err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return []*types.User{}, nil
		} else {
			return []*types.User{}, errors.Trace(err)
		}
	}
	return users, nil
}

func (mgr *Manager) CloseUsersByIssueURL(issueURL string) error {
	return errors.Trace(mgr.storage.Update(&types.User{}, "issue_url=?", []interface{}{issueURL},
		map[string]interface{}{"status": "closed"}))
}

func (mgr *Manager) CloseTeamsByIssueURL(issueURL string) error {
	return errors.Trace(mgr.storage.Update(&types.Team{}, "issue_url=?", []interface{}{issueURL},
		map[string]interface{}{"status": "closed"}))
}

func (mgr *Manager) GetUsersByTeam(teamID int) ([]*types.User, error) {
	var users []*types.User
	if err := mgr.storage.Find(&users, "team_id=?", teamID); err != nil {
		return []*types.User{}, errors.Trace(err)
	}
	return users, nil
}

func (mgr *Manager) UpdateUser(user *types.User) error {
	return errors.Trace(mgr.storage.Save(user))
}

func (mgr *Manager) UpdateTeam(team *types.Team) error {
	return errors.Trace(mgr.storage.Save(team))
}

func (mgr *Manager) ProcessSignUp(owner, repo, action string, issue *github.Issue) error {
	if action == "opened" || action == "reopened" {
		return errors.Trace(mgr.CreateSignUp(owner, repo, issue))
	} else if action == "closed" {
		return errors.Trace(mgr.CloseSignUp(owner, repo, issue))
	}
	return nil
}

func (mgr *Manager) CloseSignUp(owner, repo string, issue *github.Issue) error {
	if err := mgr.CloseUsersByIssueURL(issue.GetHTMLURL()); err != nil {
		return errors.Trace(err)
	}
	return errors.Trace(mgr.CloseTeamsByIssueURL(issue.GetHTMLURL()))
}

func (mgr *Manager) CreateSignUp(owner, repo string, issue *github.Issue) error {
	user, team := mgr.parseSignupIssue(issue)
	if user != nil {
		err := mgr.UpdateUser(user)
		if err != nil {
			return errors.Trace(err)
		}
		return errors.Trace(mgr.createPersonalSignUpSuccess(owner, repo, issue, user))
	}
	if team != nil {
		if len(team.Users) > 3 {
			comment := "Sign up failed, up to three members in a team."
			return errors.Trace(mgr.CommentIssue(owner, repo, issue.GetNumber(), comment))
		}
		// if the team can be inherited by them
		var (
			lastSeason     = mgr.Config.Season
			lastSeasonTeam *types.Team
			lastSeasonErr  error
		)
		for lastSeason > 0 && lastSeasonTeam == nil {
			lastSeason--
			lastSeasonTeam, lastSeasonErr = mgr.GetSeasonTeamByName(lastSeason, team.Name)
			if lastSeasonErr != nil {
				return errors.Trace(lastSeasonErr)
			}
		}
		if lastSeasonTeam != nil {
			canInherite := false
			for _, lastUser := range lastSeasonTeam.Users {
				for _, nowUser := range team.Users {
					if nowUser.User == lastUser.User {
						canInherite = true
					}
				}
			}
			if !canInherite {
				return errors.Trace(mgr.CommentIssue(owner, repo, issue.GetNumber(), cannotInheriteComment))
			}
		}

		err := mgr.UpdateTeam(team)
		if err != nil {
			return errors.Trace(err)
		}
		for _, user := range team.GetUsers() {
			user.TeamID = team.GetID()
			err := mgr.UpdateUser(user)
			if err != nil {
				return errors.Trace(err)
			}
		}
		return errors.Trace(mgr.createTeamSignUpSuccess(owner, repo, issue, team))
	}

	return nil
}

func (mgr *Manager) createPersonalSignUpSuccess(owner, repo string, issue *github.Issue, user *types.User) error {
	scoreReport, err := mgr.GetScoreReport(user.GetUser())
	if err != nil {
		return errors.Trace(err)
	}
	comment := fmt.Sprintf("You've signing up successfully.\n\n%s", scoreReport)
	err = mgr.CommentIssue(owner, repo, issue.GetNumber(), comment)
	if err != nil {
		return errors.Trace(err)
	}
	return errors.Trace(mgr.AddLabels(owner, repo, issue, []string{
		fmt.Sprintf("challenge-program-%d", mgr.Config.Season),
		"challenge-program-personal",
	}))
}

func (mgr *Manager) createTeamSignUpSuccess(owner, repo string, issue *github.Issue, team *types.Team) error {
	scoreReport, err := mgr.GetScoreReport(team.GetUsers()[0].GetUser())
	if err != nil {
		return errors.Trace(err)
	}
	comment := fmt.Sprintf("You've signing up successfully.\n\n%s", scoreReport)
	comment = fmt.Sprintf("%s\n\nTeam Member:", comment)
	for _, user := range team.Users {
		comment = fmt.Sprintf("%s\n- @%s", comment, user.GetUser())
	}
	err = mgr.CommentIssue(owner, repo, issue.GetNumber(), comment)
	if err != nil {
		return errors.Trace(err)
	}
	return errors.Trace(mgr.AddLabels(owner, repo, issue, []string{
		fmt.Sprintf("challenge-program-%d", mgr.Config.Season),
		"challenge-program-team",
	}))
}

func (mgr *Manager) parseSignupIssue(issue *github.Issue) (*types.User, *types.Team) {
	team := types.Team{
		Season:   mgr.Config.Season,
		IssueURL: issue.GetHTMLURL(),
		Status:   "opened",
	}
	lines := strings.Split(issue.GetBody(), "\r\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		teamMatch := teamPattern.FindStringSubmatch(line)
		if len(teamMatch) == 2 {
			team.Name = teamMatch[1]
		}
		memberEmailPattern := memberEmailPattern.FindStringSubmatch(line)
		if len(memberEmailPattern) == 3 {
			appendUserNotExist(&team.Users, &types.User{
				Season:   mgr.Config.Season,
				User:     memberEmailPattern[1],
				Email:    memberEmailPattern[2],
				IssueURL: issue.GetHTMLURL(),
				Status:   "opened",
			})
			continue
		}
		memberMatch := memberPattern.FindStringSubmatch(line)
		if len(memberMatch) == 2 {
			appendUserNotExist(&team.Users, &types.User{
				Season:   mgr.Config.Season,
				User:     memberMatch[1],
				IssueURL: issue.GetHTMLURL(),
				Status:   "opened",
			})
			continue
		}
		leaderEmailPattern := leaderEmailPattern.FindStringSubmatch(line)
		if len(leaderEmailPattern) == 3 {
			appendUserNotExist(&team.Users, &types.User{
				Season:   mgr.Config.Season,
				User:     leaderEmailPattern[1],
				Email:    leaderEmailPattern[2],
				IssueURL: issue.GetHTMLURL(),
				Status:   "opened",
				Leader:   true,
			})
			continue
		}
		leaderMatch := leaderPattern.FindStringSubmatch(line)
		if len(leaderMatch) == 2 {
			appendUserNotExist(&team.Users, &types.User{
				Season:   mgr.Config.Season,
				User:     leaderMatch[1],
				IssueURL: issue.GetHTMLURL(),
				Status:   "opened",
				Leader:   true,
			})
			continue
		}
	}
	issueAuthor := issue.GetUser().GetLogin()
	issueAuthorEmail := issue.GetUser().GetEmail()
	appendUserNotExist(&team.Users, &types.User{
		Season:   mgr.Config.Season,
		User:     issueAuthor,
		Email:    issueAuthorEmail,
		IssueURL: issue.GetHTMLURL(),
		Status:   "opened",
	})
	if team.GetName() == "" {
		team.Name = issueAuthor
	}
	// if len(team.Users) == 1 {
	// 	return team.Users[0], nil
	// }
	return nil, &team
}

func appendUserNotExist(userSlice *[]*types.User, user *types.User) {
	user.User = strings.Trim(user.User, " ")
	user.Email = strings.Trim(user.Email, " ")
	user.IssueURL = strings.Trim(user.IssueURL, " ")
	for _, u := range *userSlice {
		if u.GetSeason() == user.GetSeason() &&
			strings.ToLower(u.GetUser()) == strings.ToLower(user.GetUser()) {
			return
		}
	}
	*userSlice = append(*userSlice, user)
}
