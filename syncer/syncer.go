package syncer

import (
	"github.com/juju/errors"
	"github.com/pingcap/community/config"
	"github.com/pingcap/community/manager"
	types2 "github.com/pingcap/community/pkg/types"
)


type Syncer struct {
	cfg *config.Config
	mgr *manager.Manager
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
