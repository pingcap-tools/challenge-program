package manager

import (
	"github.com/jinzhu/gorm"
	"github.com/juju/errors"
	"github.com/pingcap/challenge-program/pkg/types"
)

func (mgr *Manager) GetTaskgroupsBySeason(season int) ([]*types.Taskgroup, error) {
	var taskgroups []*types.Taskgroup

	if err := mgr.storage.Find(&taskgroups, "season=?", season); err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return []*types.Taskgroup{}, nil
		} else {
			return []*types.Taskgroup{}, errors.Trace(err)
		}
	}
	return taskgroups, nil
}
