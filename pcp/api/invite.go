package api

import (
	"github.com/juju/errors"
	"github.com/kataras/iris"
	"github.com/ngaut/log"
)

func (hdl *ManagerHandler)InviteByGithubID(ctx iris.Context) {
	login := ctx.Params().Get("github")
	if login == "" {
		log.Error("empty login")
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.WriteString("empty login")
		return
	}
	if err := hdl.mgr.InviteByGithubID(login); err != nil {
		log.Errorf("invalid login %v", errors.ErrorStack(err))
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.WriteString(err.Error())
	} else {
		ctx.WriteString("ok")
	}
}

func (hdl *ManagerHandler)InviteByEmail(ctx iris.Context) {
	email := ctx.Params().Get("email")
	if email == "" {
		log.Error("empty email")
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.WriteString("empty email")
		return
	}
	if err := hdl.mgr.InviteByEmail(email); err != nil {
		log.Errorf("invalid login %v", errors.ErrorStack(err))
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.WriteString(err.Error())
	} else {
		ctx.WriteString("ok")
	}
}
