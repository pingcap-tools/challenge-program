package manager

import (
	"context"
	"time"

	"github.com/google/go-github/github"
	"github.com/juju/errors"
	"github.com/ngaut/log"
	"github.com/pingcap/challenge-program/config"
	"github.com/pingcap/challenge-program/util"
)

const (
	perPage      = 100
	maxRetryTime = 3
)

func (m *Manager) FetchPullsBatch(owner, repo string, startID int, ch *chan []*github.PullRequest) {
	opt := &github.PullRequestListOptions{
		State: "all",
		Sort:  "created",
		ListOptions: github.ListOptions{
			Page:    1 + int(startID/perPage),
			PerPage: perPage,
		},
	}

	if err := util.RetryOnError(context.Background(), maxRetryTime, func() error {
		pulls, _, err := m.Github.PullRequests.List(context.Background(), owner, repo, opt)
		if err != nil {
			return errors.Trace(err)
		}
		*ch <- pulls
		if len(pulls) < perPage {
			// fetch finished
			close(*ch)
		} else {
			time.Sleep(10 * time.Second)
			m.FetchPullsBatch(owner, repo, startID+perPage, ch)
		}
		return nil
	}); err != nil {
		log.Errorf("error while fetch pulls %v", errors.ErrorStack(err))
	}
}

func (m *Manager) FetchPullsUpdate(owner, repo string, page int) ([]*github.PullRequest, error) {
	var pulls []*github.PullRequest
	opt := &github.PullRequestListOptions{
		State:     "all",
		Sort:      "updated",
		Direction: "desc",
		ListOptions: github.ListOptions{
			Page:    page,
			PerPage: perPage,
		},
	}

	err := util.RetryOnError(context.Background(), maxRetryTime, func() error {
		ps, _, err := m.Github.PullRequests.List(context.Background(), owner, repo, opt)
		if err != nil {
			return errors.Trace(err)
		}
		pulls = ps
		return nil
	})

	return pulls, errors.Trace(err)
}

func (m *Manager) ListIssueComments(owner, repo string, issueNumber, page int) ([]*github.IssueComment, *github.Response, error) {
	var (
		sort      = "updated"
		direction = "desc"
	)
	return m.Github.Issues.ListComments(context.Background(), owner, repo, issueNumber, &github.IssueListCommentsOptions{
		Sort:      &sort,
		Direction: &direction,
		ListOptions: github.ListOptions{
			Page:    page,
			PerPage: perPage,
		},
	})
}

func (m *Manager) ListPullComments(owner, repo string, pullNumber, page int) ([]*github.PullRequestComment, *github.Response, error) {
	return m.Github.PullRequests.ListComments(context.Background(), owner, repo, pullNumber, &github.PullRequestListCommentsOptions{
		Sort:      "updated",
		Direction: "desc",
		ListOptions: github.ListOptions{
			Page:    page,
			PerPage: perPage,
		},
	})
}

func (m *Manager) ListPullReviews(owner, repo string, pullNumber, page int) ([]*github.PullRequestReview, *github.Response, error) {
	return m.Github.PullRequests.ListReviews(context.Background(), owner, repo, pullNumber, &github.ListOptions{
		Page:    page,
		PerPage: perPage,
	})
}

func (m *Manager) CommentIssue(owner, repo string, issueNumber int, comment string) error {
	issueComment := github.IssueComment{
		Body: &comment,
	}
	_, _, err := m.Github.Issues.CreateComment(context.Background(), owner, repo, issueNumber, &issueComment)
	return errors.Trace(err)
}

func (m *Manager) ClosePull(owner, repo string, pullNumber int) error {
	state := "closed"
	patch := &github.PullRequest{
		State: &state,
	}
	_, _, err := m.Github.PullRequests.Edit(context.Background(), owner, repo, pullNumber, patch)
	return errors.Trace(err)
}

func (m *Manager) CloseIssue(owner, repo string, issueNumber int, comment string) error {
	if comment != "" {
		if err := m.CommentIssue(owner, repo, issueNumber, comment); err != nil {
			return errors.Trace(err)
		}
	}
	state := "closed"
	_, _, err := m.Github.Issues.Edit(context.Background(), owner, repo, issueNumber, &github.IssueRequest{
		State: &state,
	})
	return errors.Trace(err)
}

func (m *Manager) AddLabel(owner, repo string, issue *github.Issue, label string) error {
	var labels []string
	hasLabel := false
	for _, l := range issue.Labels {
		labels = append(labels, l.GetName())
		if label == l.GetName() {
			hasLabel = true
		}
	}
	if hasLabel {
		return nil
	}
	labels = append(labels, label)
	_, _, err := m.Github.Issues.AddLabelsToIssue(context.Background(),
		owner, repo, issue.GetNumber(), labels)
	return errors.Trace(err)
}

func (m *Manager) AddLabels(owner, repo string, issue *github.Issue, labelAdds []string) error {
	//var labels []string
	//for _, l := range issue.Labels {
	//	labels = append(labels, l.GetName())
	//}
	//hasAllLabel := true
	//for _, l := range labelAdds {
	//	hasLabel := false
	//	for _, l1 := range labels {
	//		if l != l1 {
	//			hasLabel = true
	//		}
	//	}
	//	if hasLabel {
	//		continue
	//	}
	//	hasAllLabel = false
	//	labels = append(labels, l)
	//}
	//if hasAllLabel {
	//	return nil
	//}
	_, _, err := m.Github.Issues.AddLabelsToIssue(context.Background(),
		owner, repo, issue.GetNumber(), labelAdds)
	return errors.Trace(err)
}

func (m *Manager) ListProjectCard(id int64) ([]*github.ProjectCard, error) {
	archivedState := "not_archived"
	p, _, err := m.Github.Projects.ListProjectCards(context.Background(), id, &github.ProjectCardListOptions{
		ArchivedState: &archivedState,
		ListOptions: github.ListOptions{
			Page:    1,
			PerPage: 100,
		},
	})
	return p, err
}

func (m *Manager) GetGithubUser(login string) (*github.User, error) {
	if user, ok := m.Users[login]; ok {
		return user, nil
	}
	user, _, err := m.Github.Users.Get(context.Background(), login)
	if err != nil {
		return nil, errors.Trace(err)
	}
	m.Users[login] = user
	return user, nil
}

// func (m *Manager)CreateProjectCard(level, note string) error {
// 	var targetColumn int64
// 	var originContentID int64

// 	if level == "progress" {
// 		return m.CreateProgressCard(note)
// 	}

// 	if level == "easy" {
// 		targetColumn = m.Config.Project.EasyColumnID
// 	}
// 	if level == "medium" {
// 		targetColumn = m.Config.Project.MediumColumnID
// 	}
// 	if level == "hard" {
// 		targetColumn = m.Config.Project.HardColumnID
// 	}
// 	if targetColumn == 0 {
// 		return nil
// 	}

// 	if easyCards, err := m.ListProjectCard(m.Config.Project.EasyColumnID); err != nil {
// 		return errors.Trace(err)
// 	} else {
// 		for _, easyCard := range easyCards {
// 			if easyCard.GetNote() == note {
// 				if level == "easy" {
// 					return nil
// 				}
// 				originContentID = easyCard.GetID()
// 			}
// 		}
// 	}

// 	if mediumCards, err := m.ListProjectCard(m.Config.Project.MediumColumnID); err != nil {
// 		return errors.Trace(err)
// 	} else {
// 		for _, mediumCard := range mediumCards {
// 			if mediumCard.GetNote() == note {
// 				if level == "medium" {
// 					return nil
// 				}
// 				originContentID = mediumCard.GetID()
// 			}
// 		}
// 	}

// 	if hardCards, err := m.ListProjectCard(m.Config.Project.HardColumnID); err != nil {
// 		return errors.Trace(err)
// 	} else {
// 		for _, hardCard := range hardCards {
// 			if hardCard.GetNote() == note {
// 				if level == "hard" {
// 					return nil
// 				}
// 				originContentID = hardCard.GetID()
// 			}
// 		}
// 	}

// 	projectCardOptions := github.ProjectCardOptions{}
// 	if originContentID == 0 {
// 		projectCardOptions.Note = note
// 	} else {
// 		projectCardOptions.ContentID = originContentID
// 	}

// 	_, _, err := m.Github.Projects.CreateProjectCard(context.Background(), targetColumn, &projectCardOptions)
// 	return errors.Trace(err)
// }

func (m *Manager) CreateProgressCard(project *config.Project, note string) error {
	if progressCards, err := m.ListProjectCard(project.InProgressColumnID); err != nil {
		return errors.Trace(err)
	} else {
		for _, progressCard := range progressCards {
			if progressCard.GetNote() == note {
				return nil
			}
		}
	}

	projectCardOptions := github.ProjectCardOptions{
		Note: note,
	}

	_, _, err := m.Github.Projects.CreateProjectCard(context.Background(), project.InProgressColumnID, &projectCardOptions)
	return errors.Trace(err)
}

func (m *Manager) CreateDoneCard(project *config.Project, note string) error {
	var originContentID int64

	if progressCards, err := m.ListProjectCard(project.InProgressColumnID); err != nil {
		return errors.Trace(err)
	} else {
		for _, card := range progressCards {
			if card.GetNote() == note {
				originContentID = card.GetID()
			}
		}
	}

	projectCardOptions := github.ProjectCardOptions{}
	if originContentID == 0 {
		projectCardOptions.Note = note
	} else {
		projectCardOptions.ContentID = originContentID
	}
	_, _, err := m.Github.Projects.CreateProjectCard(context.Background(), project.FinishedColumnID, &projectCardOptions)
	return errors.Trace(err)
}
