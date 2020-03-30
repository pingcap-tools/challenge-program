package api

import (
	"github.com/kataras/iris"
	"github.com/pingcap/challenge-program/pcp/manager"
)

// ManagerHandler is manager api handler
type ManagerHandler struct {
	mgr *manager.Manager
}

func newManagerHandler(mgr *manager.Manager) *ManagerHandler {
	return &ManagerHandler{
		mgr: mgr,
	}
}

func (hdl *ManagerHandler) Ping(ctx iris.Context) {
	ctx.WriteString("pong")
}
