package syncer

import (
	"github.com/google/go-github/github"
	"github.com/ngaut/log"
	"github.com/pingcap/community/pkg/types"
	"github.com/pingcap/errors"
	"time"
)

func (s *Syncer)StartPolling() {
	for _, repo := range s.repos {
		// time.Sleep(time.Duration(time.Duration(300))*time.Second)
		log.Infof("start polling %s/%s", repo.GetOwner(), repo.GetRepo())
		pullCh := make(chan *github.PullRequest)
		go s.polling(repo, pullCh)
		// go s.pollingComment(repo, pullCh)
	}
}

func (s *Syncer)polling(repo *types.Repo, pullCh chan *github.PullRequest) {
	ch := make(chan []*github.PullRequest)
	go s.mgr.FetchPullsBatch(repo.Owner, repo.Repo, 1, &ch)
	for pulls := range ch {
		for _, pull := range pulls {
			if err := s.mgr.CreatePullNoLock(repo, pull); err != nil {}
			// pullCh <- pull
		}
	}
	log.Infof("%s/%s history pulls done", repo.GetOwner(), repo.GetRepo())
	ticker := time.NewTicker(10*time.Minute)
	for {
		select {
		case <- ticker.C:
			// should not be asynchronous
			s.fetchUpdates(repo, pullCh)
		}
	}
}

func (s *Syncer)fetchUpdates(repo *types.Repo, pullCh chan *github.PullRequest) {
	page := 1
	finish := false
	for {
		pulls, err := s.mgr.FetchPullsUpdate(repo.GetOwner(), repo.GetRepo(), page)
		if err != nil {
			log.Errorf("fetch update pulls error %v", errors.ErrorStack(err))
			break
		}

		for _, pull := range pulls {
			pullPatch, err := s.mgr.MakePullPatch(repo, pull)
			if err != nil {
				log.Errorf("make pull patch failed %v", errors.ErrorStack(err))
				continue
			}
			if pullPatch == nil {
				finish = true
				continue
			}
			log.Infof("pull  %s/%s#%d", repo.GetOwner(), repo.GetRepo(), pullPatch.GetNumber())
			if err := s.mgr.UpdatePull(pullPatch); err != nil {
				log.Errorf("update pull failed %v", errors.ErrorStack(err))
			}
			// pullCh <- pull
		}

		if !finish {
			page++
		} else {
			return
		}
	}
}
