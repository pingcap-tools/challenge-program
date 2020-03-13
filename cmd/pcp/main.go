package main

import (
	"flag"
	"fmt"
	"github.com/kataras/iris"
	"os"
	"os/signal"
	"syscall"

	"github.com/juju/errors"
	"github.com/ngaut/log"
	"github.com/pingcap/community/config"
	"github.com/pingcap/community/manager"
	"github.com/pingcap/community/pcp/api"
	pcpManager "github.com/pingcap/community/pcp/manager"
)

var (
	cfg          *config.Config
	configPath   string
)

func init() {
	flag.StringVar(&configPath, "c", "", "path to PCP config")
}

func main() {
	flag.Parse()

	cfg = config.GetGlobalConfig()
	if configPath != "" {
		err := cfg.Load(configPath)
		if err != nil {
			log.Fatalf(errors.ErrorStack(err))
		}
	}
	if err := cfg.Init(); err != nil {
		log.Error(err)
	}
	// if cfg.Project == nil {
	// 	log.Fatal("project field can not be nil")
	// }

	mgr, err := manager.New(cfg)
	if err != nil {
		log.Fatalf("can't run pcp: %v", errors.ErrorStack(err))
	}
	pcp := pcpManager.New(mgr)

	go func() {
		log.Infof("begin to listen %s:%d 😄", cfg.Host, cfg.Port)
		app := iris.New()
		api.CreateRouter(app, pcp)
		if err := app.Run(iris.Addr(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port))); err != nil {
			log.Fatalf("app run error %v", err)
		}
	}()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc,
		os.Kill,
		os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	sig := <-sc
	log.Infof("Got signal %d to exit.", sig)
}
