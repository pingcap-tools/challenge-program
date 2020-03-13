package api

import (
	"strconv"

	"github.com/juju/errors"
	"github.com/kataras/iris"
	"github.com/ngaut/log"
)

func (hdl *ManagerHandler) GetRank(ctx iris.Context) {
	rank, err := hdl.mgr.GetRank(-1)
	if err != nil {
		log.Errorf("get rank error %v", errors.ErrorStack(err))
	}
	ctx.JSON(rank)
}

func (hdl *ManagerHandler) GetRankAll(ctx iris.Context) {
	rank, err := hdl.mgr.GetRank(0)
	if err != nil {
		log.Errorf("get rank error %v", errors.ErrorStack(err))
	}
	ctx.JSON(rank)
}

func (hdl *ManagerHandler) GetRankBySeason(ctx iris.Context) {
	season, err := strconv.Atoi(ctx.Params().Get("season"))
	if err != nil {
		log.Errorf("season parse error %v", errors.ErrorStack(err))
	}
	rank, err := hdl.mgr.GetRank(season)
	if err != nil {
		log.Errorf("get rank error %v", errors.ErrorStack(err))
	}
	ctx.JSON(rank)
}

func (hdl *ManagerHandler) GetMockRank(ctx iris.Context) {
	rank := hdl.mgr.GetMockRank(-1)
	ctx.JSON(rank)
}
