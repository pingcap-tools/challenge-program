package manager

import (
	"sync"

	"github.com/google/go-github/github"
	"github.com/juju/errors"
	"github.com/pingcap/challenge-program/config"
	githubInit "github.com/pingcap/challenge-program/pkg/github"
	"github.com/pingcap/challenge-program/pkg/storage"
	"github.com/pingcap/challenge-program/pkg/storage/basic"
)

// Manager represent schrodinger syncer
type Manager struct {
	Config  *config.Config
	storage basic.Driver
	Github  *github.Client
	Members map[string]bool
	Users   map[string]*github.User
	sync.Mutex
}

// New init manager
func New(cfg *config.Config) (*Manager, error) {
	s, err := storage.New(cfg.Database)
	if err != nil {
		return nil, errors.Trace(err)
	}

	return &Manager{
		Config:  cfg,
		storage: s,
		Github:  githubInit.GetGithubClient(cfg.GithubToken),
		Users:   make(map[string]*github.User),
		Members: make(map[string]bool),
	}, nil
}
