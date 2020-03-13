package manager

import "github.com/pingcap/community/pkg/types"

func (mgr *Manager) GetMockTasksAll(season int) []types.Task {
	return append(mgr.GetMockTasksLevel(-1), mgr.GetMockTasksRepo(-1)...)
}

func (mgr *Manager) GetMockTasksLevel(season int) []types.Task {
	return []types.Task{
		{
			Season: 2,
			APICompleteUser: &types.User{
				User:   "you06",
				Avatar: "https://avatars3.githubusercontent.com/u/9587680?s=460&v=4",
				GitHub: "https://github.com/you06/",
			},
			APICompleteTeam: nil,
			Owner:           "pingcap",
			Repo:            "tidb",
			Title:           "task 1",
			Issue:           "https://github.com/pingcap/tidb/issues/10467",
			Level:           "medium",
			Score:           1000,
			Status:          "success",
		},
	}
}

func (mgr *Manager) GetMockTasksRepo(season int) []types.Task {
	return []types.Task{
		{
			Season: 2,
			DoingUsers: []types.User{
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
			Owner:  "pingcap",
			Repo:   "tidb",
			Title:  "task 2",
			Issue:  "https://github.com/pingcap/tidb/issues/7546",
			Level:  "hard",
			Score:  6000,
			Status: "success",
		},
	}
}
