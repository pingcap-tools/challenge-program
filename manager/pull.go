package manager

import (
	"github.com/google/go-github/github"
	"github.com/jinzhu/gorm"
	"github.com/juju/errors"
	"github.com/pingcap/challenge-program/pkg/types"
	"github.com/pingcap/challenge-program/util"
)

func (mgr *Manager) GetPullByNumber(owner, repo string, number int) (*types.Pull, error) {
	var pull types.Pull
	if err := mgr.storage.FindOne(&pull, "owner=? AND repo=? AND pull_number=?", owner, repo, number); err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		} else {
			return nil, errors.Trace(err)
		}
	}
	return &pull, nil
}

func (mgr *Manager) CreatePull(repo *types.Repo, pull *github.PullRequest) error {
	mgr.Lock()
	defer mgr.Unlock()
	return mgr.CreatePullNoLock(repo, pull)
}

func (mgr *Manager) CreatePullNoLock(repo *types.Repo, pull *github.PullRequest) error {
	p, err := mgr.MakePullPatch(repo, pull)
	if err != nil {
		return errors.Trace(err)
	}
	if p != nil {
		return errors.Trace(mgr.UpdatePull(p))
	}
	return nil
}

func (mgr *Manager) UpdatePull(pull *types.Pull) error {
	return errors.Trace(mgr.storage.Save(pull))
}

func (mgr *Manager) MakePullPatch(repo *types.Repo, pull *github.PullRequest) (*types.Pull, error) {
	p, err := mgr.GetPullByNumber(repo.GetOwner(), repo.GetRepo(), pull.GetNumber())
	if err == nil && p == nil {
		return mgr.MakePull(repo, pull)
	} else if err != nil {
		return nil, errors.Trace(err)
	}

	if p.UpdatedAt.Equal(pull.GetUpdatedAt()) {
		return nil, nil
	}

	var labels []string
	for _, label := range pull.Labels {
		labels = append(labels, label.GetName())
	}
	status := pull.GetState()
	if pull.MergedAt != nil {
		status = "merged"
	}

	p.Title = pull.GetTitle()
	p.Body = pull.GetBody()
	p.Label = util.EncodeStringSlice(labels)
	p.Status = status
	p.UpdatedAt = pull.GetUpdatedAt()
	p.ClosedAt = pull.GetClosedAt()
	p.MergedAt = pull.GetMergedAt()
	return p, nil
}

func (mgr *Manager) MakePull(repo *types.Repo, pull *github.PullRequest) (*types.Pull, error) {
	isMember, err := mgr.isMember(pull.GetUser().GetLogin())
	if err != nil {
		return nil, errors.Trace(err)
	}
	relation := "member"
	if !isMember {
		relation = "not member"
	}
	var labels []string
	for _, label := range pull.Labels {
		labels = append(labels, label.GetName())
	}

	status := pull.GetState()
	if pull.MergedAt != nil {
		status = "merged"
	}

	p := types.Pull{
		Owner:       repo.GetOwner(),
		Repo:        repo.GetRepo(),
		Number:      pull.GetNumber(),
		Title:       pull.GetTitle(),
		Body:        pull.GetBody(),
		User:        pull.GetUser().GetLogin(),
		Association: pull.GetAuthorAssociation(),
		Relation:    relation,
		Label:       util.EncodeStringSlice(labels),
		Status:      status,
		CreatedAt:   pull.GetCreatedAt(),
		UpdatedAt:   pull.GetUpdatedAt(),
		ClosedAt:    pull.GetClosedAt(),
		MergedAt:    pull.GetMergedAt(),
	}
	return &p, nil
}
