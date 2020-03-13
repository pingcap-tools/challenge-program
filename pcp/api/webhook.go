package api

import (
	"bytes"
	"io/ioutil"
	"strings"

	"github.com/google/go-github/github"
	"github.com/juju/errors"
	"github.com/kataras/iris"
	"github.com/ngaut/log"
	"github.com/pingcap/community/pkg/types"
)

// HookBody for parsing webhook
type HookBody struct {
	Repository struct {
		FullName string `json:"full_name"`
	}
}

func (hdl *ManagerHandler) Webhook(ctx iris.Context) {
	r := ctx.Request()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorf("body read error %v", errors.ErrorStack(err))
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.WriteString(err.Error())
		return
	}
	// restore body for iris ReadJSON use
	r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	hookBody := HookBody{}
	if err := ctx.ReadJSON(&hookBody); err != nil {
		// body parse error
		log.Errorf("read json error %v", errors.ErrorStack(err))
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.WriteString(err.Error())
		return
	}
	repoInfo := strings.Split(hookBody.Repository.FullName, "/")
	if len(repoInfo) != 2 {
		// invalid repo name
		log.Errorf("invalid repo name")
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.WriteString(err.Error())
		return
	}
	repo := &types.Repo{
		Owner: repoInfo[0],
		Repo:  repoInfo[1],
	}

	// restore body for github ValidatePayload use
	r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	payload, err := github.ValidatePayload(r, []byte(hdl.mgr.GetConfig().WenhookSecret))
	if err != nil {
		// invalid payload
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.WriteString(err.Error())
		log.Errorf("invalid payload %v", errors.ErrorStack(err))
		return
	}
	event, err := github.ParseWebHook(github.WebHookType(r), payload)
	if err != nil {
		// event parse err
		log.Errorf("webhook parse error %v", errors.ErrorStack(err))
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.WriteString(err.Error())
		return
	}
	ctx.WriteString("ok")
	hdl.mgr.Webhook(repo, event)
}
