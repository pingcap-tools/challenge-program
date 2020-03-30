package syncer

import (
	"strings"

	"github.com/juju/errors"
	"github.com/ngaut/log"
	types2 "github.com/pingcap/challenge-program/pkg/types"
)

func (s *Syncer) parseRepos(repoStrs []string) error {
	log.Info(repoStrs)
	var repos []*types2.Repo
	for _, repoStr := range repoStrs {
		r := strings.Split(repoStr, "/")
		if len(r) != 2 {
			return errors.Errorf("repo name %s invalid, example pingcap/tidb", repoStr)
		}
		repos = append(repos, &types2.Repo{
			Owner: r[0],
			Repo:  r[1],
		})
	}

	s.repos = repos
	return nil
}
