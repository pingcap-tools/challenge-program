package api

import (
	"github.com/kataras/iris"
)

func (hdl *ManagerHandler) GetMockTasksAll(ctx iris.Context) {
	rank := hdl.mgr.GetMockTasksAll(-1)
	ctx.JSON(rank)
}

func (hdl *ManagerHandler) GetMockTasksLevel(ctx iris.Context) {
	rank := hdl.mgr.GetMockTasksLevel(-1)
	ctx.JSON(rank)
}

func (hdl *ManagerHandler) GetMockTasksRepo(ctx iris.Context) {
	rank := hdl.mgr.GetMockTasksRepo(-1)
	ctx.JSON(rank)
}
