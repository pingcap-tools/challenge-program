package manager

import (
	"github.com/google/go-github/github"
	"github.com/pingcap/challenge-program/pkg/types"
)

func (mgr *Manager) Webhook(repo *types.Repo, event interface{}) {
	switch event := event.(type) {
	case *github.PullRequestEvent:
		mgr.processPull(repo, event)
	case *github.IssueCommentEvent:
		mgr.processIssueComment(repo, event)
	case *github.IssuesEvent:
		mgr.processIssue(repo, event)
	}
}
