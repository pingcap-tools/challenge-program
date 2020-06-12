package manager

import (
	"github.com/pingcap/challenge-program/config"
	"github.com/pingcap/challenge-program/manager"
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
	m := &Manager{
		mgr: mgr,
	}
	go m.watchExpiredPick()
	return m
}

func (mgr *Manager) GetConfig() *config.Config {
	return mgr.mgr.Config
}
