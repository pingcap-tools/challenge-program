package manager

import (
	"github.com/juju/errors"
	"github.com/pingcap/community/pkg/types"
)

// season: -1: current season, 0: all season
func (mgr *Manager) GetRank(season int) (types.Rank, error) {
	if season == -1 {
		season = mgr.mgr.Config.Season
	}
	rank, err := mgr.mgr.GetAllRank(season)
	return rank, errors.Trace(err)
}

func (mgr *Manager) GetMockRank(season int) types.Rank {
	return []*types.RankItem{
		{
			Rank:      1,
			Type:      "team",
			Name:      "PingCAP",
			Community: false,
			Url:       "https://github.com/tidb-perf-challenge/pcp/issues/1",
			Score:     1140,
			DoingTask: "https://github.com/pingcap/tidb/issues/14486",
		},
		{
			Rank:      2,
			Type:      "team",
			Name:      "Shimokitazawa",
			Community: true,
			Url:       "https://github.com/tidb-perf-challenge/pcp/issues/1",
			Score:     514,
			DoingTask: "",
		},
		{
			Rank:      3,
			Type:      "user",
			Name:      "you06",
			Community: false,
			Url:       "https://github.com/tidb-perf-challenge/pcp/issues/1",
			Score:     500,
			DoingTask: "https://github.com/pingcap/tidb/issues/14486",
		},
		{
			Rank:      4,
			Type:      "user",
			Name:      "810s",
			Community: true,
			Url:       "https://github.com/tidb-perf-challenge/pcp/issues/1",
			Score:     350,
		},
		{
			Rank:      5,
			Type:      "user",
			Name:      "sre-bot",
			Community: false,
			Url:       "https://github.com/tidb-perf-challenge/pcp/issues/1",
			Score:     100,
		},
	}
}
