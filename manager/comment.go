package manager

import (
	"github.com/google/go-github/github"
	"github.com/jinzhu/gorm"
	"github.com/juju/errors"
	"github.com/pingcap/community/pkg/types"
)

type CommentPatchBase interface {
	GetID() int64
	GetBody() string
	GetUser() *github.User
	GetHTMLURL() string
}

func (m *Manager)GetCommentByID(id int64) (*types.Comment, error) {
	var comment types.Comment
	if err := m.storage.FindOne(&comment, "comment_id=?", id); err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		} else {
			return nil, errors.Trace(err)
		}
	}
	return &comment, nil
}

func (mgr *Manager)UpdateComment(comment *types.Comment) error {
	return errors.Trace(mgr.storage.Save(comment))
}

func (mgr *Manager)MakeCommentPatch(repo *types.Repo, patchBase CommentPatchBase, attach *types.CommentAttach) (*types.Comment, error) {
	c, err := mgr.GetCommentByID(patchBase.GetID())
	if err == nil && c == nil {
		return mgr.MakeComment(repo, patchBase, attach)
	} else if err != nil {
		return nil, errors.Trace(err)
	}

	if c.GetUpdatedAt().Equal(attach.GetUpdatedAt()) {
		return nil, nil
	}

	c.Body = patchBase.GetBody()
	c.UpdatedAt = attach.GetUpdatedAt()
	c.Association = attach.GetAuthorAssociation()
	return c, nil
}

func (mgr *Manager)MakeComment(repo *types.Repo, patchBase CommentPatchBase, attach *types.CommentAttach) (*types.Comment, error) {
	isMember, err := mgr.isMember(patchBase.GetUser().GetLogin())
	if err != nil {
		return nil, errors.Trace(err)
	}
	relation := "member"
	if !isMember {
		relation = "not member"
	}
	return &types.Comment{
		Owner: repo.GetOwner(),
		Repo: repo.GetRepo(),
		CommentID: patchBase.GetID(),
		CommentType: attach.GetCommentType(),
		Number: attach.GetNumber(),
		Body: patchBase.GetBody(),
		User: patchBase.GetUser().GetLogin(),
		Url: patchBase.GetHTMLURL(),
		Association: attach.GetAuthorAssociation(),
		Relation: relation,
		CreatedAt: attach.GetCreatedAt(),
		UpdatedAt: attach.GetUpdatedAt(),
	}, nil
}
