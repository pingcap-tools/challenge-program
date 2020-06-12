package manager

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/go-github/github"
	"github.com/juju/errors"
	"github.com/ngaut/log"
	"github.com/pingcap/challenge-program/pkg/types"
)

const (
	issueNumberRegex         = `ucp:? ?#(\d+)`
	issueNumberWithURLRegex  = `ucp:? ?https:\/\/github\.com\/[a-zA-Z0-9-]+\/[a-zA-Z0-9-]+\/issues\/(\d+)`
	issueNumberWithRepoRegex = `ucp:? ?\[#\d+\].*?.*\/(.*)\/(.*)\/issues\/(.*)\)`
	successOpenedPullComment = "Thanks for your contribution. If your PR get merged, you will be rewarded %d points."
	easyTheshould            = 200
	vectorTaskScore          = 50
	easyTheshouldComment     = "Congratulations, you get %d score from easy level tasks in challenge program season-%d, and if your PRs in reviewed stage all got merged, the score will be %d, try some medium and hard tasks!(you can not reward from easy and vector tasks now)"
	firstTaskComment         = `Congratulation! You have awarded a badge for usability challenge program! Please fill the form to get your reward! http://tidbcommunity.mikecrm.com/QMCv4QL`
)

var (
	tidbVectorIssue            = []int{12101, 12102, 12103, 12104, 12105, 12106, 12176, 12058}
	tikvVectorIssue            = []int{5751}
	testVectorIssue            = []int{63}
	issueNumberPattern         = regexp.MustCompile(issueNumberRegex)
	issueNumberWithURLPattern  = regexp.MustCompile(issueNumberWithURLRegex)
	issueNumberWithRepoPattern = regexp.MustCompile(issueNumberWithRepoRegex)
)

func (mgr *Manager) processPull(repo *types.Repo, pullEvent *github.PullRequestEvent) {
	if pullEvent.GetPullRequest().GetUser().GetLogin() == "sre-bot" {
		return
	}
	switch pullEvent.GetAction() {
	case "opened", "edited", "reopened":
		return
		mgr.createPull(repo, pullEvent.GetPullRequest())
	case "closed":
		mgr.closePull(repo, pullEvent.GetPullRequest())
	}
}

func (mgr *Manager) createPull(repo *types.Repo, pull *github.PullRequest) {
	// log.Info(mgr.mgr.CommentIssue(repo.GetOwner(), repo.GetRepo(), pull.GetNumber(), PCPCloseComment))
	// return

	// filter edit closed pull status
	if pull.GetState() == "closed" {
		return
	}
	// pick already exist
	if pick, err := mgr.mgr.GetRepoPickByNumber(repo.GetOwner(), repo.GetRepo(), pull.GetNumber()); err == nil && pick != nil {
		return
	}
	log.Infof("create pick %s/%s #%d\n", repo.GetOwner(), repo.GetRepo(), pull.GetNumber())
	body := pull.GetBody()
	task, err := mgr.parsePullBody(repo, body)
	if task == nil && err == nil {
		log.Infof("task not found")
		return
	}
	if err != nil {
		log.Errorf("create pull failed %v", errors.ErrorStack(err))
		return
	}
	// TODO: check if signing up here
	user, err := mgr.mgr.GetUserByLogin(pull.GetUser().GetLogin(), mgr.mgr.Config.Season)
	if err != nil {
		log.Errorf("get user error %v", errors.ErrorStack(err))
		return
	} else if user == nil {
		log.Infof("%s/%s #%d file up pull but not sign up", repo.GetOwner(), repo.GetRepo(), pull.GetNumber())
		if err := mgr.mgr.CommentIssue(repo.GetOwner(), repo.GetRepo(), pull.GetNumber(),
			fmt.Sprintf(notSignUpComment, pull.GetUser().GetLogin())); err != nil {
			log.Errorf("get user error %v", errors.ErrorStack(err))
		}
		if err := mgr.mgr.ClosePull(repo.GetOwner(), repo.GetRepo(), pull.GetNumber()); err != nil {
			log.Errorf("close pull failed %v", errors.ErrorStack(err))
		}
		return
	}
	if score, err := mgr.mgr.GetCombinedRepoScore(repo, pull.GetUser().GetLogin()); err != nil {
		log.Errorf("get user score failed %v", errors.ErrorStack(err))
	} else {
		log.Infof("%s/%s #%d current expected easy score %d", repo.GetOwner(), repo.GetRepo(), pull.GetNumber(), score)
		ok, _ := ifVectorIssue(repo, task.GetIssueNumber())
		if (task.GetLevel() == "easy" || ok) && score >= easyTheshould {
			comment := fmt.Sprintf("@%s, you already got %d points from easy level tasks when all pull requests merged. And you will not get score from this PR.", pull.GetUser().GetLogin(), score)
			task.Score = 0
			if err := mgr.mgr.CommentIssue(repo.GetOwner(), repo.GetRepo(), pull.GetNumber(), comment); err != nil {
				log.Errorf("comment score over easy threshould pull %v", errors.ErrorStack(err))
				return
			}
		}
	}
	// not good here but make it work first
	if ok, _ := ifVectorIssue(repo, task.GetIssueNumber()); ok {
		if err := mgr.createVectorPull(repo, task, pull); err != nil {
			log.Errorf("create vector pull %v", errors.ErrorStack(err))
		}
		return
	}
	log.Infof("Task ID %d", task.GetID())
	team, err := mgr.mgr.GetTeamByUser(pull.GetUser().GetLogin(), mgr.mgr.Config.Season)
	if err != nil {
		log.Errorf("get team failed")
		return
	}
	var (
		pick  *types.Pick
		picks []*types.Pick
		er    error
	)
	if team != nil {
		picks, er = mgr.mgr.GetPicksByTeam(team)
	} else {
		picks, er = mgr.mgr.GetPicksByLogin(pull.GetUser().GetLogin())
	}
	if er != nil {
		log.Errorf("get pick error %v", errors.ErrorStack(er))
		return
	}
	valid := false
	for _, p := range picks {
		if p.GetTaskID() == task.GetID() {
			valid = true
			pick = p
		}
	}
	if pick == nil || !valid {
		// comment and close PR
		comment := fmt.Sprintf("Thanks for your pull request. Pick up issue #%d and reopen this PR", task.GetIssueNumber())
		if err := mgr.mgr.CommentIssue(repo.GetOwner(), repo.GetRepo(), pull.GetNumber(), comment); err != nil {
			log.Errorf("comment failed %v", errors.ErrorStack(err))
		}
		if err := mgr.mgr.ClosePull(repo.GetOwner(), repo.GetRepo(), pull.GetNumber()); err != nil {
			log.Errorf("close pull failed %v", errors.ErrorStack(err))
		}
		return
	}
	pick.PullNumber = pull.GetNumber()
	pick.Status = "review"
	pick.ClosedAt = time.Time{}
	if err := mgr.mgr.UpdatePick(pick); err != nil {
		log.Errorf("update pick failed %v", errors.Trace(err))
	}
	if project := mgr.mgr.Config.FindProject(repo.GetOwner(), repo.GetRepo()); project != nil {
		if err := mgr.mgr.CreateProgressCard(project, pull.GetHTMLURL()); err != nil {
			log.Errorf("create PR project card fai %v", errors.ErrorStack(err))
		}
	}
	if err := mgr.mgr.CommentIssue(repo.GetOwner(), repo.GetRepo(), pull.GetNumber(),
		fmt.Sprintf(successOpenedPullComment, pick.GetScore())); err != nil {
		log.Errorf("comment success opened pull failed %v", errors.Trace(err))
	}
}

func (mgr *Manager) closePull(repo *types.Repo, pull *github.PullRequest) {
	if pull.MergedAt == nil {
		return
	}
	log.Infof("task closed %s/%s #%d", repo.GetOwner(), repo.GetRepo(), pull.GetNumber())
	task, err := mgr.parsePullBody(repo, pull.GetBody())
	if err != nil {
		log.Errorf("get task failed %v", errors.Trace(err))
	}
	if task == nil {
		log.Errorf("pull %s/%s#%d task not found", repo.GetOwner(), repo.GetRepo(), pull.GetNumber())
		return
	}
	if ok, _ := ifVectorIssue(repo, task.GetIssueNumber()); ok {
		if err := mgr.closeVectorPull(repo, task, pull); err != nil {
			log.Errorf("close vector pull %v", errors.ErrorStack(err))
		}
		return
	}
	team, err := mgr.mgr.GetTeamByUser(pull.GetUser().GetLogin(), mgr.mgr.Config.Season)
	if err != nil {
		log.Errorf("get team failed")
		return
	}
	var (
		pick *types.Pick
		er   error
	)
	if team != nil {
		pick, er = mgr.mgr.GetTeamPickByPullNumber(team.GetID(), "review", pull.GetNumber())
	} else {
		pick, er = mgr.mgr.GetPickByPullNumber(pull.GetUser().GetLogin(), "review", pull.GetNumber())
	}
	if er != nil {
		log.Errorf("get pick failed %v", errors.Trace(err))
		return
	}
	if pick == nil {
		// Finish job without picking up task
		// should close PR in earlier stage, something wrong
		log.Errorf("finish job without picking up task, pull %s/%s #%d",
			repo.GetOwner(), repo.GetRepo(), pull.GetNumber())
		return
	}
	pick.Status = "success"
	pick.PullNumber = pull.GetNumber()
	now := time.Now()
	pick.UpdatedAt = now
	pick.ClosedAt = now
	if err := mgr.mgr.UpdatePick(pick); err != nil {
		log.Errorf("update pick %d failed %v", pick.GetID(), errors.ErrorStack(err))
	}
	if task.GetPullNumber() != 0 || task.GetStatus() == "success" {
		log.Errorf("task %d already finished", task.GetID())
		return
	}
	if pick.GetTeamID() != 0 {
		task.CompleteTeam = pick.GetTeamID()
	} else {
		task.CompleteUser = pick.GetUser()
	}
	task.Status = "success"
	task.PullNumber = pull.GetNumber()
	if err := mgr.mgr.UpdateTask(task); err != nil {
		log.Errorf("update task %d failed %v", task.GetID(), errors.ErrorStack(err))
	}

	if team != nil {
		picks, err := mgr.mgr.GetPicksByTeam(team)
		if err != nil {
			log.Errorf("get picks faield %v", err)
		} else {
			successCount := 0
			for _, p := range picks {
				if p.Status == "success" {
					successCount++
				}
			}
			// first task only
			if successCount == 1 {
				if err := mgr.mgr.CommentIssue(repo.GetOwner(), repo.GetRepo(), pull.GetNumber(), firstTaskComment); err != nil {
					log.Errorf("comment error %v", errors.ErrorStack(err))
				}
			}
		}
	}

	preComment := ""
	if pick.GetTeamID() != 0 {
		team, err := mgr.mgr.GetTeamByID(pick.GetTeamID())
		if err != nil {
			log.Errorf("get team error %v", errors.ErrorStack(err))
		}
		if team != nil {
			preComment = fmt.Sprintf("Team [%s](%s)", team.GetName(), team.GetIssueURL())
		}
	}
	if preComment == "" {
		preComment = fmt.Sprintf("@%s", pull.GetUser().GetLogin())
	}
	comment := fmt.Sprintf("%s complete task #%d and get %d score",
		preComment, task.GetIssueNumber(), task.GetScore())
	seasonScore, err := mgr.mgr.GetUserSeasonScore(pull.GetUser().GetLogin(), mgr.mgr.Config.Season)
	if err == nil {
		comment = fmt.Sprintf("%s complete task #%d and get %d score, current score %d",
			preComment, task.GetIssueNumber(), task.GetScore(), seasonScore)
	} else {
		log.Errorf("get user score error %v", errors.ErrorStack(err))
	}
	if err := mgr.mgr.CommentIssue(repo.GetOwner(), repo.GetRepo(), pull.GetNumber(), comment); err != nil {
		log.Errorf("comment error %v", errors.ErrorStack(err))
	}
	if err := mgr.mgr.CloseIssue(repo.GetOwner(), repo.GetRepo(), task.GetIssueNumber(), ""); err != nil {
		log.Errorf("close issue error %v", errors.ErrorStack(err))
	}

	// if task.GetLevel() == "easy" {
	// 	if score, err := mgr.mgr.GetCombinedRepoScore(repo, pull.GetUser().GetLogin()); err != nil {
	// 		log.Errorf("get user score failed %v", errors.ErrorStack(err))
	// 	} else if score >= easyTheshould {
	// 		comment := fmt.Sprintf(easyTheshouldComment, seasonScore, mgr.mgr.Config.Season, score)
	// 		if err := mgr.mgr.CommentIssue(repo.GetOwner(), repo.GetRepo(), pull.GetNumber(), comment); err != nil {
	// 			log.Errorf("comment error %v", errors.ErrorStack(err))
	// 		}
	// 	}
	// }

	// if project := mgr.mgr.Config.FindProject(repo.GetOwner(), repo.GetRepo()); project != nil {
	// 	if err := mgr.mgr.CreateDoneCard(project, pull.GetHTMLURL()); err != nil {
	// 		log.Errorf("create PR project card fail %v", errors.ErrorStack(err))
	// 	}
	// }
	// if project := mgr.mgr.Config.FindProject(repo.GetOwner(), repo.GetRepo()); project != nil {
	// 	if err := mgr.mgr.CreateProgressCard(project, pull.GetHTMLURL()); err != nil {
	// 		log.Errorf("create PR project card fai %v", errors.ErrorStack(err))
	// 	}
	// }
	// if err := mgr.mgr.CreateProjectCard("progress", pull.GetHTMLURL()); err != nil {
	// 	log.Errorf("create PR project card fai %v", errors.ErrorStack(err))
	// }
}

func (mgr *Manager) createVectorPull(repo *types.Repo, task *types.Task, pull *github.PullRequest) error {
	pick := &types.Pick{
		Season:     mgr.mgr.Config.Season,
		Owner:      repo.GetOwner(),
		Repo:       repo.GetRepo(),
		TaskID:     0,
		PullNumber: pull.GetNumber(),
		Score:      task.GetScore(),
		Status:     "review",
		ClosedAt:   time.Time{},
	}
	team, err := mgr.mgr.GetTeamByUser(pull.GetUser().GetLogin(), mgr.mgr.Config.Season)
	if err != nil {
		return errors.Trace(err)
	}
	if team != nil {
		pick.TeamID = team.GetID()
	} else {
		pick.User = pull.GetUser().GetLogin()
	}
	if err := mgr.mgr.UpdatePick(pick); err != nil {
		return errors.Trace(err)
	}
	if project := mgr.mgr.Config.FindProject(repo.GetOwner(), repo.GetRepo()); project != nil {
		if err := mgr.mgr.CreateProgressCard(project, pull.GetHTMLURL()); err != nil {
			log.Errorf("create PR project card fail %v", errors.ErrorStack(err))
		}
	}
	return errors.Trace(mgr.mgr.CommentIssue(repo.GetOwner(), repo.GetRepo(), pull.GetNumber(),
		fmt.Sprintf(successOpenedPullComment, pick.GetScore())))
}

func (mgr *Manager) closeVectorPull(repo *types.Repo, task *types.Task, pull *github.PullRequest) error {
	var (
		pick *types.Pick
		err  error
	)
	if team, er := mgr.mgr.GetTeamByUser(pull.GetUser().GetLogin(), mgr.mgr.Config.Season); er != nil {
		log.Info("error")
		return errors.Trace(err)
	} else if team != nil {
		log.Info("team", team)
		pick, err = mgr.mgr.GetTeamPickByPullNumber(team.GetID(), "review", pull.GetNumber())
	} else {
		log.Info("user", pull.GetUser().GetLogin(), pull.GetNumber())
		pick, err = mgr.mgr.GetPickByPullNumber(pull.GetUser().GetLogin(), "review", pull.GetNumber())
	}
	log.Info("got pick", pick, err)
	if err != nil {
		return errors.Trace(err)
	}
	if pick == nil {
		return nil
	}
	pick.Status = "success"
	pick.PullNumber = pull.GetNumber()
	now := time.Now()
	pick.UpdatedAt = now
	pick.ClosedAt = now
	if err := mgr.mgr.UpdatePick(pick); err != nil {
		return errors.Trace(err)
	}
	preComment := ""
	if pick.GetTeamID() != 0 {
		team, err := mgr.mgr.GetTeamByID(pick.GetTeamID())
		if err != nil {
			log.Errorf("get team error %v", errors.ErrorStack(err))
		}
		if team != nil {
			preComment = fmt.Sprintf("Team [%s](%s)", team.GetName(), team.GetIssueURL())
		}
	}
	if preComment == "" {
		preComment = fmt.Sprintf("@%s", pull.GetUser().GetLogin())
	}
	comment := fmt.Sprintf("%s complete task #%d and get %d score.",
		preComment, task.GetIssueNumber(), pick.GetScore())
	seasonScore, err := mgr.mgr.GetUserSeasonScore(pull.GetUser().GetLogin(), mgr.mgr.Config.Season)
	if err == nil {
		comment = fmt.Sprintf("%s complete task #%d and get %d score, current score %d.",
			preComment, task.GetIssueNumber(), pick.GetScore(), seasonScore)
	} else {
		log.Errorf("get user score error %v", errors.ErrorStack(err))
	}
	if err := mgr.mgr.CommentIssue(repo.GetOwner(), repo.GetRepo(), pull.GetNumber(), comment); err != nil {
		log.Errorf("comment error %v", errors.ErrorStack(err))
	}
	if score, err := mgr.mgr.GetUserExpectedEasySeasonScore(pull.GetUser().GetLogin(), mgr.mgr.Config.Season); err != nil {
		log.Errorf("get user score failed %v", errors.ErrorStack(err))
	} else if score >= easyTheshould {
		comment := fmt.Sprintf(easyTheshouldComment, seasonScore, mgr.mgr.Config.Season, score)
		if err := mgr.mgr.CommentIssue(repo.GetOwner(), repo.GetRepo(), pull.GetNumber(), comment); err != nil {
			log.Errorf("comment error %v", errors.ErrorStack(err))
		}
	}
	if project := mgr.mgr.Config.FindProject(repo.GetOwner(), repo.GetRepo()); project != nil {
		if err := mgr.mgr.CreateDoneCard(project, pull.GetHTMLURL()); err != nil {
			log.Errorf("create PR project card fail %v", errors.ErrorStack(err))
		}
	}
	return nil
}

func (mgr *Manager) parsePullBody(repo *types.Repo, body string) (*types.Task, error) {
	lowerBody := strings.ToLower(body)
	m := issueNumberPattern.FindStringSubmatch(lowerBody)
	if len(m) == 0 {
		m = issueNumberWithURLPattern.FindStringSubmatch(lowerBody)
	}
	if len(m) == 0 {
		m = issueNumberWithRepoPattern.FindStringSubmatch(lowerBody)
	}
	if len(m) == 0 {
		return nil, nil
	}

	var (
		o           string
		r           string
		issueNumber int
		err         error
	)

	if len(m) == 2 {
		o = repo.GetOwner()
		r = repo.GetRepo()
		issueNumber, err = strconv.Atoi(m[1])
	}
	if len(m) == 4 {
		o = m[1]
		r = m[2]
		issueNumber, err = strconv.Atoi(m[3])
	}

	if err != nil {
		return nil, errors.Trace(err)
	}
	if issueNumber == 0 {
		return nil, nil
	}

	// hard code, maybe make improvement in the future
	if ok, issueNumber := ifVectorIssue(repo, issueNumber); ok {
		return &types.Task{
			Owner:       o,
			Repo:        r,
			Level:       "easy",
			IssueNumber: issueNumber,
			Score:       vectorTaskScore,
		}, nil
	}

	return mgr.mgr.GetTaskByNumber(o, r, issueNumber, mgr.mgr.Config.Season)
}

func ifVectorIssue(repo *types.Repo, issueNumber int) (bool, int) {
	if repo.GetOwner() == "pingcap" && repo.GetRepo() == "tidb" {
		for _, vectorIssue := range tidbVectorIssue {
			if vectorIssue == issueNumber {
				return true, vectorIssue
			}
		}
	}
	if repo.GetOwner() == "tikv" && repo.GetRepo() == "tikv" {
		for _, vectorIssue := range tikvVectorIssue {
			if vectorIssue == issueNumber {
				return true, vectorIssue
			}
		}
	}
	if repo.GetOwner() == "you06" && repo.GetRepo() == "cherry-pick-playground" {
		for _, vectorIssue := range testVectorIssue {
			if vectorIssue == issueNumber {
				return true, vectorIssue
			}
		}
	}
	return false, 0
}
