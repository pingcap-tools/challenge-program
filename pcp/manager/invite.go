package manager

import (
	"context"
	"github.com/google/go-github/github"
	"github.com/juju/errors"
	"github.com/ngaut/log"
)

const org = "tidb-perf-challenge"

func (mgr *Manager)InviteByGithubID(login string) error {
	m, _, _ := mgr.mgr.Github.Organizations.GetOrgMembership(context.Background(), "sre-bot", org)
	log.Info(m)
	user, _, err := mgr.mgr.Github.Users.Get(context.Background(), login)
	if err != nil {
		return errors.Trace(err)
	}
	userID := user.GetID()
	log.Infof("user %s, id %d", user.GetLogin(), user.GetID())
	if _, _, err := mgr.mgr.Github.Organizations.CreateOrgInvitation(context.Background(),
		org, &github.CreateOrgInvitationOptions{
			InviteeID: &userID,
		}); err != nil {
		return errors.Trace(err)
	}
	return nil
}

func (mgr *Manager)InviteByEmail(email string) error {
	if _, _, err := mgr.mgr.Github.Organizations.CreateOrgInvitation(context.Background(),
		org, &github.CreateOrgInvitationOptions{
			Email: &email,
		}); err != nil {
		return errors.Trace(err)
	}
	return nil
}
