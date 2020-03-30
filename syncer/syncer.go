package syncer

import (
	"github.com/juju/errors"
	"github.com/pingcap/challenge-program/config"
	"github.com/pingcap/challenge-program/manager"
	types2 "github.com/pingcap/challenge-program/pkg/types"
)

type Syncer struct {
	cfg   *config.Config
	mgr   *manager.Manager
	repos []*types2.Repo
}

func New(cfg *config.Config, mgr *manager.Manager) (*Syncer, error) {
	s := Syncer{
		cfg: cfg,
		mgr: mgr,
	}

	if err := s.parseRepos(cfg.Repos); err != nil {
		return nil, errors.Trace(err)
	}

	return &s, nil
}
