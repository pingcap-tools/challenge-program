package manager

import (
	"github.com/pingcap/community/config"
	"github.com/pingcap/community/manager"
)

type Manager struct {
	mgr *manager.Manager
}

func New(mgr *manager.Manager) *Manager {
	//mgr.SyncOrgTeams()
	//go func(mgr *manager.Manager) {
	//	ticker := time.NewTicker(5*time.Minute)
	//	for {
	//		select {
	//		case <- ticker.C:
	//			mgr.SyncOrgTeams()
	//		}
	//	}
	//}(mgr)
	return &Manager{
		mgr: mgr,
	}
}

func (mgr *Manager) GetConfig() *config.Config {
	return mgr.mgr.Config
}
