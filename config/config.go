package config

import (
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/juju/errors"
)

// Config contains configuration options.
type Config struct {
	Host          string     `toml:"host"`
	Port          int        `toml:"port"`
	Season        int        `toml:"season"`
	GithubToken   string     `toml:"github"`
	WenhookSecret string     `toml:"webhook-secret"`
	Org           string     `toml:"org"`
	Database      *Database  `toml:"database"`
	Projects      []*Project `toml:"projects"`
	Repos         []string   `toml:"repos"`
}

// Slack configuration
// type Slack struct {
// 	DMUser string `toml:"dm-user"`
// 	Token string `toml:"token"`
// }

// Database defines db configuration
type Database struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	Username string `toml:"username"`
	Password string `toml:"password"`
	Database string `toml:"database"`
}

type Project struct {
	Repo               string `toml:"repo"`
	ProjectID          int64  `toml:"project_id"`
	EasyColumnID       int64  `toml:"easy_column_id"`
	MediumColumnID     int64  `toml:"medium_column_id"`
	HardColumnID       int64  `toml:"hard_column_id"`
	InProgressColumnID int64  `toml:"in_progress_column_id"`
	FinishedColumnID   int64  `toml:"finished_column_id"`
}

var globalConf = Config{
	Host: "0.0.0.0",
	Port: 30000,
	Org:  "",
	Database: &Database{
		Host:     "127.0.0.1",
		Port:     3306,
		Username: "root",
		Password: "",
		Database: "chaos_commander",
	},
}

// GetGlobalConfig returns the global configuration for this server.
func GetGlobalConfig() *Config {
	return &globalConf
}

// Load loads config options from a toml file_logger.
func (c *Config) Load(confFile string) error {
	_, err := toml.DecodeFile(confFile, c)
	return errors.Trace(err)
}

// Init do some prepare works
func (c *Config) Init() error {
	return nil
}

// FindProject return project if exist
func (c *Config) FindProject(owner, repo string) *Project {
	fullname := fmt.Sprintf("%s/%s", owner, repo)
	for _, project := range c.Projects {
		if project.Repo == fullname {
			return project
		}
	}
	return nil
}
