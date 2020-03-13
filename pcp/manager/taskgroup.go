package manager

import (
	"fmt"
	"sort"

	"github.com/juju/errors"
	"github.com/pingcap/community/pkg/types"
)

func (mgr *Manager) GetMockTaskgroups(season int) []types.Taskgroup {
	return []types.Taskgroup{
		{
			Season:      2,
			Owner:       "pingcap",
			Repo:        "tidb",
			Title:       "task group 1",
			IssueNumber: 14486,
			Issue:       "https://github.com/pingcap/tidb/issues/14486",
			Progress:    60,
			DoingUsers: []*types.User{
				{
					User:   "you06",
					Avatar: "https://avatars3.githubusercontent.com/u/9587680?s=460&v=4",
					GitHub: "https://github.com/you06/",
				},
				{
					User:   "illyrix",
					Avatar: "https://avatars3.githubusercontent.com/u/12008675?s=460&v=4",
					GitHub: "https://github.com/illyrix/",
				},
			},
		},
		{
			Season:      2,
			Owner:       "tikv",
			Repo:        "tikv",
			Title:       "task group 2",
			IssueNumber: 6519,
			Issue:       "https://github.com/tikv/tikv/issues/6519",
			Progress:    30,
			DoingUsers: []*types.User{
				{
					User:   "you06",
					Avatar: "https://avatars3.githubusercontent.com/u/9587680?s=460&v=4",
					GitHub: "https://github.com/you06/",
				},
			},
		},
	}
}

func (mgr *Manager) GetTaskgroups(season int) ([]*types.Taskgroup, error) {
	if season == -1 {
		season = mgr.mgr.Config.Season
	}
	taskgroups, err := mgr.mgr.GetTaskgroupsBySeason(season)
	if err != nil {
		return nil, errors.Trace(err)
	}

	for _, taskgroup := range taskgroups {
		mgr.ComposeTaskgroup(season, taskgroup)
	}

	sort.Sort(types.SortByVote(taskgroups))

	return taskgroups, nil
}

func (mgr *Manager) ComposeTaskgroup(season int, taskgroup *types.Taskgroup) error {
	var (
		tasks         []*types.Task
		taskIDs       []int
		picks         []*types.Pick
		completeTasks int
		topUsers      []*types.User
		err           error
	)

	taskgroup.Issue = fmt.Sprintf("https://github.com/%s/%s/issues/%d", taskgroup.Owner, taskgroup.Repo, taskgroup.IssueNumber)
	tasks, err = mgr.mgr.GetTasksByTaskgroup(taskgroup.ID)
	if err != nil {
		return errors.Trace(err)
	}
	for _, task := range tasks {
		taskIDs = append(taskIDs, task.ID)
		if task.Status == "success" {
			completeTasks++
		}
	}

	picks, err = mgr.mgr.GetPicksByTasks(taskIDs)
	if err != nil {
		return errors.Trace(err)
	}
	for i := 0; i < 5 && i < len(picks); i++ {
		pick := picks[i]
		switch pick.Status {
		case "doing", "review":
			{
				if pick.TeamID == 0 {
					user, err := mgr.mgr.GetUserByLogin(pick.User, season)
					if err != nil {
						return errors.Trace(err)
					}
					topUsers = append(topUsers, user)
				} else {
					team, err := mgr.mgr.GetTeamByID(pick.TeamID)
					if err != nil {
						return errors.Trace(err)
					}
					topUsers = append(topUsers, team.Users[0])
				}
			}
		}
	}

	for i, user := range topUsers {
		githubUser, err := mgr.mgr.GetGithubUser(user.User)
		if err != nil {
			return errors.Trace(err)
		}
		user.GitHub = githubUser.GetHTMLURL()
		user.Avatar = githubUser.GetAvatarURL()
		topUsers[i] = user
	}

	totalTasks := len(tasks)
	if totalTasks == 0 {
		totalTasks = 1
	}
	taskgroup.DoingUsers = topUsers
	taskgroup.Progress = (100 * completeTasks) / totalTasks
	taskgroup.Tasks = tasks

	return nil
}
