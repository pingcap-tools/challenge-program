package mysql

import (
	"fmt"

	"github.com/pingcap/challenge-program/config"
	"github.com/pingcap/challenge-program/pkg/storage/basic"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/juju/errors"
)

type driver struct {
	db *gorm.DB
}

// Connect create database connection
func Connect(config *config.Database) (basic.Driver, error) {
	fmt.Println(config)
	connect := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Username, config.Password, config.Host, config.Port, config.Database)
	db, err := gorm.Open("mysql", connect)
	db.LogMode(true)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return &driver{
		db,
	}, nil
}

// do not use errors trace in the methods
// then gorm.IsRecordNotFoundError can work in the outside
// Close connection
func (d *driver) Close() error {
	return d.db.Close()
}

// Save insert or update data
func (d *driver) Save(model interface{}) error {
	return d.db.Save(model).Error
}

// FindOne get one matched result
func (d *driver) FindOne(model interface{}, stmt string, cond ...interface{}) error {
	return d.db.Where(stmt, cond...).First(model).Error
}

// Find get all matched results
func (d *driver) Find(model interface{}, stmt string, cond ...interface{}) error {
	return d.db.Where(stmt, cond...).Find(model).Error
}

// Update data
func (d *driver) Update(model interface{}, stmt string, cond []interface{}, update map[string]interface{}) error {
	return d.db.Model(model).Where(stmt, cond...).Updates(update).Error
}

// UpdateAll update all data
func (d *driver) UpdateAll(model interface{}, update map[string]interface{}) error {
	return d.db.Model(model).Updates(update).Error
}

// Scan do scan functions
func (d *driver) Scan(model interface{}, scanModel interface{}, selectStmt, stmt string, cond ...interface{}) error {
	return d.db.Model(model).Where(stmt, cond...).Select(selectStmt).Scan(scanModel).Error
}

// Scan do scan functions
func (d *driver) Group(model interface{}, scanModel interface{}, selectStmt, groupStmt string, stmt string, cond ...interface{}) error {
	return d.db.Model(model).Where(stmt, cond...).Select(selectStmt).Group(groupStmt).Scan(scanModel).Error
}
