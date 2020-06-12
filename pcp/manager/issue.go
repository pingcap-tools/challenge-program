package manager

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/go-github/github"
	"github.com/ngaut/log"
	"github.com/pingcap/challenge-program/pkg/types"
	"github.com/pingcap/errors"
)

const (
	pickNotRequiredComment = `You can file a PR directly without picking up because this issue is a collection of multiple tasks. Each PR related to this issue will rewarded 50 score when merged.`
	notSignUpComment       = "@%s Please sign up before pick a challenge.\nYou can signing up by file a issue here https://github.com/tidb-challenge-program/register"
	alreadyPickUp          = "@%s you have already picked up this issue."
	pickedByOther          = "@%s this task has been picked by other, you can pick up this after %s"
)

var (
	awardPointPattern = regexp.MustCompile(`^\/award-point ([0-9]+).*`)
)

func (mgr *Manager) processIssue(repo *types.Repo, issuesEvent *github.IssuesEvent) {
	issue := issuesEvent.GetIssue()
	if repo.GetOwner() == "tidb-perf-challenge" && repo.GetRepo() == "pcp" &&
		issuesEvent.GetAction() != "closed" {
		log.Info(mgr.mgr.CommentIssue(repo.GetOwner(), repo.GetRepo(), issuesEvent.GetIssue().GetNumber(), RegisterMoveComment))
		return
	}
	// handle sign up
	if repo.GetOwner() == "tidb-challenge-program" && repo.GetRepo() == "register" {
		// log.Info(mgr.mgr.CommentIssue(repo.GetOwner(), repo.GetRepo(), issuesEvent.GetIssue().GetNumber(), PCPCloseComment))
		// return
		if err := mgr.mgr.ProcessSignUp(repo.GetOwner(), repo.GetRepo(), issuesEvent.GetAction(), issue); err != nil {
			log.Infof("create task failed %v", errors.ErrorStack(err))
		}
		return
	}
	// handle create task
	if issuesEvent.GetAction() == "opened" ||
		issuesEvent.GetAction() == "reopened" ||
		issuesEvent.GetAction() == "edited" ||
		issuesEvent.GetAction() == "labeled" {
		isMember, err := mgr.mgr.IsMember(issuesEvent.GetSender().GetLogin())
		if err != nil {
			log.Info(err)
		} else if isMember {
			if err := mgr.mgr.CreateTask(repo.GetOwner(), repo.GetRepo(), issue); err != nil {
				log.Infof("create task failed %v", errors.ErrorStack(err))
			}
		} else {
			log.Info(issue.GetUser().GetLogin(), "is not member, don't create task")
		}
	}
}

func (mgr *Manager) processIssueComment(repo *types.Repo, issueCommentEvent *github.IssueCommentEvent) {
	issueComment := issueCommentEvent.GetComment()
	command := strings.Trim(issueComment.GetBody(), " ")
	if strings.HasPrefix(command, "/pick-up-challenge") {
		log.Info(mgr.mgr.CommentIssue(repo.GetOwner(), repo.GetRepo(), issueCommentEvent.GetIssue().GetNumber(), PCPCloseComment))
		return
		if err := mgr.tryPickUp(repo, issueCommentEvent.GetIssue(), issueComment.GetUser().GetLogin()); err != nil {
			log.Errorf("pick up challenge failed, %v", errors.ErrorStack(err))
		}
	}
	if strings.HasPrefix(command, "/give-up-challenge") {
		log.Info(mgr.mgr.CommentIssue(repo.GetOwner(), repo.GetRepo(), issueCommentEvent.GetIssue().GetNumber(), PCPCloseComment))
		return
		if err := mgr.tryGiveUp(repo, issueCommentEvent.GetIssue(), issueComment.GetUser().GetLogin()); err != nil {
			log.Errorf("pick up challenge failed, %v", errors.ErrorStack(err))
		}
	}
	if strings.HasPrefix(command, "/award-point") {
		log.Info(mgr.mgr.CommentIssue(repo.GetOwner(), repo.GetRepo(), issueCommentEvent.GetIssue().GetNumber(), PCPCloseComment))
		return
		if err := mgr.awardPoint(repo, issueCommentEvent.GetIssue(), issueComment, issueComment.GetUser().GetLogin()); err != nil {
			log.Errorf("award point failed, %v", errors.ErrorStack(err))
		}
	}
}

func (mgr *Manager) tryPickUp(repo *types.Repo, issue *github.Issue, login string) error {
	log.Infof("%s pick up %s/%s #%d", login, repo.GetOwner(), repo.GetRepo(), issue.GetNumber())
	user, err := mgr.mgr.GetUserByLogin(login, mgr.mgr.Config.Season)
	if err != nil {
		return errors.Trace(err)
	} else if user == nil {
		return errors.Trace(mgr.mgr.CommentIssue(repo.GetOwner(), repo.GetRepo(), issue.GetNumber(),
			fmt.Sprintf(notSignUpComment, login)))
	}
	if ok, _ := ifVectorIssue(repo, issue.GetNumber()); ok {
		return errors.Trace(mgr.mgr.CommentIssue(repo.GetOwner(), repo.GetRepo(), issue.GetNumber(), pickNotRequiredComment))
	}
	task, err := mgr.mgr.GetTaskByNumber(repo.GetOwner(), repo.GetRepo(), issue.GetNumber(), mgr.mgr.Config.Season)
	if err != nil {
		return errors.Trace(err)
	}
	if task == nil && (repo.Repo != "bug-hunting-issue" && repo.Owner != "tidb-challenge-program") {
		return errors.Trace(mgr.mayNotChallengeIssue(repo, issue, login))
	}

	// if doingPick, err := mgr.mgr.HasDoingPick(login); err != nil {
	// 	return errors.Trace(err)
	// } else if doingPick != nil {
	// 	var comment string
	// 	if doingPick.GetTaskID() == task.GetID() {
	// 		comment = fmt.Sprintf("@%s you have already picked up this issue.", login)
	// 	} else {
	// 		doingTask, er := mgr.mgr.GetTaskById(doingPick.GetTaskID())
	// 		if er != nil {
	// 			return errors.Trace(er)
	// 		}
	// 		comment = fmt.Sprintf("@%s already had picked %s/%s#%d, finish it before pick up a new one.",
	// 			login, doingTask.GetOwner(), doingTask.GetRepo(), doingTask.GetIssueNumber())
	// 	}
	// 	return errors.Trace(mgr.mgr.CommentIssue(repo.GetOwner(), repo.GetRepo(),
	// 		issue.GetNumber(), comment))
	// }

	// If pick count over limit
	overlimit, doingPicks, err := mgr.mgr.DoingPickOverLimit(login)
	pickUpThis := false
	for _, pick := range doingPicks {
		if pick.GetTaskID() == task.GetID() {
			pickUpThis = true
		}
	}

	if overlimit || pickUpThis {
		var comment string
		if pickUpThis {
			comment = fmt.Sprintf(alreadyPickUp, login)
		} else {
			var commentPicks []string
			for _, pick := range doingPicks {
				doingTask, err := mgr.mgr.GetTaskById(pick.GetTaskID())
				if err != nil {
					return errors.Trace(err)
				}
				commentPicks = append(commentPicks, fmt.Sprintf("%s/%s#%d", doingTask.GetOwner(), doingTask.GetRepo(), doingTask.GetIssueNumber()))
			}
			comment = fmt.Sprintf("@%s already had picked %s, finish one before pick up a new one.", login, strings.Join(commentPicks, ", "))
		}
		return errors.Trace(mgr.mgr.CommentIssue(repo.GetOwner(), repo.GetRepo(),
			issue.GetNumber(), comment))
	}

	// If the task is picked up by other
	doingPick, err := mgr.mgr.GetDoingPickByTask(task)
	if err != nil {
		return errors.Trace(err)
	}
	if doingPick != nil {
		comment := fmt.Sprintf(pickedByOther, login, doingPick.GetCreatedAt().AddDate(0, 0, 7).Format("2006-01-02"))
		return errors.Trace(mgr.mgr.CommentIssue(repo.GetOwner(), repo.GetRepo(),
			issue.GetNumber(), comment))
	}

	score, err := mgr.mgr.GetCombinedRepoScore(repo, login)
	if err != nil {
		return errors.Trace(err)
	}
	// isMember, err := mgr.mgr.IsMember(login)
	// if err != nil {
	// 	return errors.Trace(err)
	// }

	var memberCount float64 = 0
	team, err := mgr.mgr.GetTeamByUser(login, mgr.mgr.Config.Season)
	if err != nil {
		return errors.Trace(err)
	}
	for _, teamUser := range team.Users {
		ifMember, err := mgr.mgr.IsMember(teamUser.User)
		if err != nil {
			return errors.Trace(err)
		}
		if ifMember {
			memberCount++
		}
	}

	var (
		comment        string
		groupUserCount int
	)
	groupUserCount = len(team.Users)
	if groupUserCount == 0 {
		groupUserCount = 1
	}
	switch task.GetLevel() {
	case "easy":
		{
			if score >= easyTheshould {
				comment = fmt.Sprintf("\"easy\" issue is not available since you have got %d score in this repo. "+
					"Please pickup \"medium\" or \"hard\" directly. Or you can pick up \"easy\" Issue in other repos. "+
					"Thank you.", easyTheshould)
			}
		}
	case "medium", "hard":
		{
			if score < easyTheshould && repo.Repo != "chaos-mesh" {
				comment = fmt.Sprintf("@%s don't have enough score, pick up failed\n"+
					"\nProgress `%d/%d`\nYou may pick up some easy issues first", login, score, easyTheshould)
			}
			if task.GetLevel() == "hard" && memberCount/float64(groupUserCount) > 0.6 {
				comment = `Team with PingCAPers more than 2/3 can only pick "hard" level issues.`
			}
		}
	}

	// bug challenge issues do not have score limitation
	if repo.Owner == "tidb-challenge-program" && repo.Repo == "bug-hunting-issue" {
		comment = ""
	}

	if comment != "" {
		return errors.Trace(mgr.mgr.CommentIssue(repo.GetOwner(), repo.GetRepo(), issue.GetNumber(),
			comment))
	}

	task.Status = "doing"
	if err := mgr.mgr.UpdateTask(task); err != nil {
		return errors.Trace(err)
	}
	if err := mgr.mgr.CreatePick(login, task); err != nil {
		return errors.Trace(err)
	}

	// bug challenge issues may need to give point
	if repo.Owner == "tidb-challenge-program" && repo.Repo == "bug-hunting-issue" {
		if task.Score == 0 {
			return errors.Trace(mgr.mgr.CommentIssue(repo.GetOwner(), repo.GetRepo(), issue.GetNumber(),
				fmt.Sprintf("@%s pick up issue success, PTAL @tidb-challenge-program/point-team", login)))
		}
	}
	return errors.Trace(mgr.mgr.CommentIssue(repo.GetOwner(), repo.GetRepo(), issue.GetNumber(),
		fmt.Sprintf("@%s pick up issue success", login)))
}

func (mgr *Manager) tryGiveUp(repo *types.Repo, issue *github.Issue, login string) error {
	log.Infof("%s give up %s/%s #%d", login, repo.GetOwner(), repo.GetRepo(), issue.GetNumber())
	user, err := mgr.mgr.GetUserByLogin(login, mgr.mgr.Config.Season)
	if err != nil {
		return errors.Trace(err)
	} else if user == nil {
		return errors.Trace(mgr.mgr.CommentIssue(repo.GetOwner(), repo.GetRepo(), issue.GetNumber(),
			fmt.Sprintf(notSignUpComment, login)))
	}
	task, err := mgr.mgr.GetTaskByNumber(repo.GetOwner(), repo.GetRepo(), issue.GetNumber(), mgr.mgr.Config.Season)
	if err != nil {
		return errors.Trace(err)
	}
	if task == nil {
		return errors.Trace(mgr.mayNotChallengeIssue(repo, issue, login))
	}

	pick, err := mgr.mgr.HasDoingPick(login)
	if err != nil {
		return errors.Trace(err)
	}
	if pick == nil || pick.GetTaskID() != task.GetID() {
		return errors.Trace(mgr.mgr.CommentIssue(repo.GetOwner(), repo.GetRepo(), issue.GetNumber(),
			fmt.Sprintf("@%s has not pick up this issue", login)))
	}
	pick.Status = "failed"
	if err := mgr.mgr.UpdatePick(pick); err != nil {
		return errors.Trace(err)
	}

	return errors.Trace(mgr.mgr.CommentIssue(repo.GetOwner(), repo.GetRepo(), issue.GetNumber(),
		fmt.Sprintf("@%s give up issue success", login)))
}

func (mgr *Manager) awardPoint(repo *types.Repo, issue *github.Issue, comment *github.IssueComment, login string) error {
	pick, err := mgr.mgr.GetRepoPickByNumber(repo.GetOwner(), repo.GetRepo(), issue.GetNumber())
	if err != nil || pick == nil {
		return errors.Trace(err)
	}
	if pick.Status != "review" {
		return nil
	}
	matches := awardPointPattern.FindStringSubmatch(comment.GetBody())
	if len(matches) != 2 {
		return nil
	}
	points, err := strconv.Atoi(matches[1])
	if err != nil {
		return errors.Trace(err)
	}
	if points > 500 && points > pick.GetScore() {
		return errors.Trace(mgr.mgr.CommentIssue(repo.GetOwner(), repo.GetRepo(), issue.GetNumber(), "Not allowed to award more than 500 points."))
	}
	if isMember, err := mgr.mgr.IsMember(login); err != nil {
		return errors.Trace(err)
	} else if !isMember {
		return nil
	}

	pick.Score = points
	if err := errors.Trace(mgr.mgr.UpdatePick(pick)); err != nil {
		return errors.Trace(mgr.mgr.CommentIssue(repo.GetOwner(), repo.GetRepo(), issue.GetNumber(), "Update score failed, contact admin for help."))
	}
	return errors.Trace(mgr.mgr.CommentIssue(repo.GetOwner(), repo.GetRepo(), issue.GetNumber(), fmt.Sprintf("Update score success, the task will rewarded %d after merged.", points)))
}

//func (mgr *Manager)canPickUp(repo *types.Repo, level, login string) (int, error) {
//	score, err := mgr.mgr.GetUserSeasonScoreByRepo(login, repo.GetOwner(), repo.GetRepo(), mgr.mgr.Config.Season)
//	if err != nil {
//		return 0, errors.Trace(err)
//	}
//	return score, nil
//}

func (mgr *Manager) mayNotChallengeIssue(repo *types.Repo, issue *github.Issue, login string) error {
	var (
		seasonLabel       = fmt.Sprintf("challenge-program-%d", mgr.mgr.Config.Season)
		hasChallengeLabel = false
		comment           string
	)
	for _, label := range issue.Labels {
		if strings.ToLower(label.GetName()) == seasonLabel {
			hasChallengeLabel = true
		}
	}
	if !hasChallengeLabel {
		comment = fmt.Sprintf("@%s this is not a challenge issue.", login)
	} else {
		comment = fmt.Sprintf("@%s sorry for this accident, we will fix it ASAP.\n@you06 something wrong with this challenge issue, PTAL.", login)
	}
	return errors.Trace(mgr.mgr.CommentIssue(repo.GetOwner(), repo.GetRepo(),
		issue.GetNumber(), comment))
}

func (mgr *Manager) watchExpiredPick() {
	for range time.Tick(10 * time.Minute) {
		mgr.watchExpiredPickOnce()
	}
}

func (mgr *Manager) watchExpiredPickOnce() {
	picks, err := mgr.mgr.GetDoingPicks()
	if err != nil {
		log.Error(err)
		return
	}

	ruleStartTime, _ := time.Parse("2006-01-02", "2020-03-30")
	ruleCloseTime, _ := time.Parse("2006-01-02", "2020-04-06")
	startTime := time.Now().AddDate(0, 0, -7)
	for _, pick := range picks {
		if (pick.GetCreatedAt().Before(startTime) && pick.GetCreatedAt().After(ruleStartTime)) ||
			(pick.GetCreatedAt().Before(ruleStartTime) && time.Now().After(ruleCloseTime)) {
			if err := mgr.closeExpiredPick(pick); err != nil {
				log.Error(err)
			} else {
				log.Infof("close pick%d, %s/%s due to expired", pick.ID, pick.Owner, pick.Repo)
			}
		}
	}
}

func (mgr *Manager) closeExpiredPick(pick *types.Pick) error {
	task, err := mgr.mgr.GetTaskById(pick.GetTaskID())
	if err != nil {
		return errors.Trace(err)
	}

	pick.Status = "failed"
	if err := mgr.mgr.UpdatePick(pick); err != nil {
		return errors.Trace(err)
	}

	comment := "This pick has been automatically canceled after more than a week."
	if err := mgr.mgr.CommentIssue(task.GetOwner(), task.GetRepo(), task.GetIssueNumber(), comment); err != nil {
		return errors.Trace(err)
	}
	return nil
}
