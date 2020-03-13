package manager
//
//import (
//	"context"
//	"github.com/google/go-github/github"
//	"github.com/ngaut/log"
//	"github.com/pingcap/errors"
//)
//
//type Team struct {
//	Origin  *github.Team
//	Members []*github.User
//}
//
//func (mgr *Manager)SyncOrgTeams() {
//	githubTeams, _, err := mgr.Github.Teams.ListTeams(context.Background(), mgr.Config.Org, nil)
//	if err != nil {
//		log.Error("sync teams error %v", errors.ErrorStack(err))
//		return
//	}
//
//	var teams []*Team
//	for _, githubTeam := range githubTeams {
//		teams = append(teams, &Team{
//			Origin: githubTeam,
//		})
//	}
//	mgr.OrgTeams = teams
//	mgr.SyncOrgTeamMember()
//}
//
//func (mgr *Manager)SyncOrgTeamMember() {
//	for _, team := range mgr.OrgTeams {
//		members, _, err := mgr.Github.Teams.ListTeamMembers(context.Background(), team.Origin.GetID(), nil)
//		if err != nil {
//			log.Error("sync team memberes error %v", errors.ErrorStack(err))
//			continue
//		}
//		team.Members = members
//	}
//}
//
//func (mgr *Manager)GetTeamByUser(login string) *Team {
//	var t *Team
//	for _, team := range mgr.OrgTeams {
//		for _, member := range team.Members {
//			if member.GetLogin() == login {
//				if t == nil {
//					t = team
//				} else {
//					// this user has 2 teams
//					// maybe send a slack message
//					log.Error("user %s has more than one team %s, %s", t.Origin.GetName(), team.Origin.GetName())
//				}
//			}
//		}
//	}
//	return t
//}
//
//func (mgr *Manager)GetTeamById(id int64) *Team {
//	for _, team := range mgr.OrgTeams {
//		if team.Origin.GetID() == id {
//			return team
//		}
//	}
//	return nil
//}
