package basic

// Driver defines storage model
type Driver interface {
	Close() error
	Save(model interface{}) error
	FindOne(model interface{}, stmt string, cond ...interface{}) error
	Find(model interface{}, stmt string, cond ...interface{}) error
	Update(model interface{}, stmt string, cond []interface{}, update map[string]interface{}) error
	UpdateAll(model interface{}, update map[string]interface{}) error
	Scan(model interface{}, scanModel interface{}, selectStmt, stmt string, cond ...interface{}) error
	Group(model interface{}, scanModel interface{}, selectStmt, groupStmt string, stmt string, cond ...interface{}) error
}
