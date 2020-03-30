package syncer

import (
	"time"

	"github.com/google/go-github/github"
	"github.com/juju/errors"
	"github.com/ngaut/log"
	"github.com/pingcap/challenge-program/pkg/types"
)

const perPage = 100

func (s *Syncer) pollingComment(repo *types.Repo, ch chan *github.PullRequest) {

	for pull := range ch {
		issueComments, _, err := s.fetchIssueComments(repo.GetOwner(), repo.GetRepo(), pull.GetNumber())
		if err != nil {
			log.Error("fetch issue comment failed %v", errors.Trace(err))
		}
		pullComments, _, err := s.fetchPullComments(repo.GetOwner(), repo.GetRepo(), pull.GetNumber())
		if err != nil {
			log.Error("fetch issue comment failed %v", errors.Trace(err))
		}
		pullReviews, rateSafe, err := s.fetchPullReviews(repo.GetOwner(), repo.GetRepo(), pull.GetNumber())

		for _, issueComment := range issueComments {
			patch, err := s.mgr.MakeCommentPatch(repo, issueComment, &types.CommentAttach{
				CommentType: "common comment",
				Number:      pull.GetNumber(),
				CreatedAt:   issueComment.GetCreatedAt(),
				UpdatedAt:   issueComment.GetUpdatedAt(),
				Association: issueComment.GetAuthorAssociation(),
			})
			if err != nil {
				log.Errorf("create patch failed %v", err)
				continue
			} else if patch == nil {
				continue
			}
			if err := s.mgr.UpdateComment(patch); err != nil {
				log.Errorf("update patch failed %v", err)
			}
		}

		for _, pullComment := range pullComments {
			patch, err := s.mgr.MakeCommentPatch(repo, pullComment, &types.CommentAttach{
				CommentType: "review comment",
				Number:      pull.GetNumber(),
				CreatedAt:   pullComment.GetCreatedAt(),
				UpdatedAt:   pullComment.GetUpdatedAt(),
				Association: pullComment.GetAuthorAssociation(),
			})
			if err != nil {
				log.Errorf("create patch failed %v", err)
				continue
			} else if patch == nil {
				continue
			}
			if err := s.mgr.UpdateComment(patch); err != nil {
				log.Errorf("update patch failed %v", err)
			}
		}

		for _, pullReview := range pullReviews {
			patch, err := s.mgr.MakeCommentPatch(repo, pullReview, &types.CommentAttach{
				CommentType: "review",
				Number:      pull.GetNumber(),
				CreatedAt:   pullReview.GetSubmittedAt(),
				UpdatedAt:   pullReview.GetSubmittedAt(),
				Association: "",
			})
			if err != nil {
				log.Errorf("create patch failed %v", err)
				continue
			} else if patch == nil {
				continue
			}
			if err := s.mgr.UpdateComment(patch); err != nil {
				log.Errorf("update patch failed %v", err)
			}
		}

		if !rateSafe {
			time.Sleep(time.Hour)
		}
	}
}

func (s *Syncer) fetchIssueComments(owner, repo string, issueNumber int) ([]*github.IssueComment, bool, error) {
	finish := false
	page := 1
	rateSafe := true
	var r []*github.IssueComment
	for !finish {
		comments, res, err := s.mgr.ListIssueComments(owner, repo, issueNumber, page)
		// rate
		if res == nil {
			rateSafe = false
		} else {
			rateSafe = res.Rate.Remaining > 1000
		}
		// error
		if err != nil {
			return []*github.IssueComment{}, rateSafe, errors.Trace(err)
		}
		// if finish
		if len(comments) < perPage {
			finish = true
		} else {
			comment := comments[len(comments)-1]
			c, err := s.mgr.GetCommentByID(comment.GetID())
			if err != nil {
				return []*github.IssueComment{}, rateSafe, errors.Trace(err)
			} else if c != nil {
				if c.GetUpdatedAt().Equal(comment.GetUpdatedAt()) {
					finish = true
				}
			}
		}
		r = append(r, comments...)
		page++
	}
	return r, rateSafe, nil
}

func (s *Syncer) fetchPullComments(owner, repo string, pullNumber int) ([]*github.PullRequestComment, bool, error) {
	finish := false
	page := 1
	rateSafe := true
	var r []*github.PullRequestComment
	for !finish {
		comments, res, err := s.mgr.ListPullComments(owner, repo, pullNumber, page)
		// rate
		if res == nil {
			rateSafe = false
		} else {
			rateSafe = res.Rate.Remaining > 1000
		}
		// error
		if err != nil {
			return []*github.PullRequestComment{}, rateSafe, errors.Trace(err)
		}
		// if finish
		if len(comments) < perPage {
			finish = true
		} else {
			comment := comments[len(comments)-1]
			c, err := s.mgr.GetCommentByID(comment.GetID())
			if err != nil {
				return []*github.PullRequestComment{}, rateSafe, errors.Trace(err)
			} else if c != nil {
				if c.GetUpdatedAt().Equal(comment.GetUpdatedAt()) {
					finish = true
				}
			}
		}
		r = append(r, comments...)
		page++
	}
	return r, rateSafe, nil
}

func (s *Syncer) fetchPullReviews(owner, repo string, pullNumber int) ([]*github.PullRequestReview, bool, error) {
	finish := false
	page := 1
	rateSafe := true
	var r []*github.PullRequestReview
	for !finish {
		reviews, res, err := s.mgr.ListPullReviews(owner, repo, pullNumber, page)
		// rate
		if res == nil {
			rateSafe = false
		} else {
			rateSafe = res.Rate.Remaining > 1000
		}
		// error
		if err != nil {
			return []*github.PullRequestReview{}, rateSafe, errors.Trace(err)
		}
		// if finish
		if len(reviews) < perPage {
			finish = true
		} else {
			for _, review := range reviews {
				c, err := s.mgr.GetCommentByID(review.GetID())
				if err != nil {
					return []*github.PullRequestReview{}, rateSafe, errors.Trace(err)
				} else if c != nil {
					if c.GetCreatedAt().Equal(review.GetSubmittedAt()) {
						finish = true
					}
				}
			}
		}
		r = append(r, reviews...)
		page++
	}
	return r, rateSafe, nil
}
