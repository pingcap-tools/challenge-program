package api

import (
	"github.com/juju/errors"
	"github.com/kataras/iris"
	"github.com/ngaut/log"
)

func (hdl *ManagerHandler) GetMockTaskgroups(ctx iris.Context) {
	rank := hdl.mgr.GetMockTaskgroups(-1)
	ctx.JSON(rank)
}

func (hdl *ManagerHandler) GetTaskgroups(ctx iris.Context) {
	rank, err := hdl.mgr.GetTaskgroups(-1)
	if err != nil {
		log.Errorf("get rank error %v", errors.ErrorStack(err))
	}
	ctx.JSON(rank)
}
