package storage

import (
	"github.com/pingcap/challenge-program/config"
	"github.com/pingcap/challenge-program/pkg/storage/basic"
	"github.com/pingcap/challenge-program/pkg/storage/mysql"

	"github.com/juju/errors"
)

func New(config *config.Database) (basic.Driver, error) {
	driver, err := mysql.Connect(config)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return driver, nil
}
