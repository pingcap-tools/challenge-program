package api

import (
	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris"
	"github.com/pingcap/challenge-program/pcp/manager"
)

func CreateRouter(app *iris.Application, mgr *manager.Manager) {
	hdl := newManagerHandler(mgr)
	crs := cors.New(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:8080",
			"http://localhost:3005",
			"http://127.0.0.1:3005",
			"https://pingcap.com",
			"https://www.pingcap.com",
			"https://pingcap.netlify.com",
		},
		AllowCredentials: true,
	})
	party := app.Party("/api", crs).AllowMethods(iris.MethodOptions)
	mockParty := app.Party("/mock/api", crs).AllowMethods(iris.MethodOptions)

	party.Get("/ping", hdl.Ping)

	// github webhook
	party.Post("/webhook", hdl.Webhook)

	// github invite
	party.Post("/invite/github/{github:string}", hdl.InviteByGithubID)
	party.Post("/invite/email/{email:string}", hdl.InviteByEmail)

	// rank
	party.Get("/rank", hdl.GetRank)
	party.Get("/rank/all", hdl.GetRankAll)
	party.Get("/rank/season/{season:int}", hdl.GetRankBySeason)
	party.Get("/task", hdl.GetMockTasksAll)
	party.Get("/task/level/{level:string}", hdl.GetMockTasksLevel)
	party.Get("/task/tasks/owner/{owner:string}/repo/{repo:string}", hdl.GetMockTasksRepo)
	party.Get("/taskgroup", hdl.GetTaskgroups)

	mockParty.Get("/rank", hdl.GetMockRank)
	mockParty.Get("/rank/all", hdl.GetMockRank)
	mockParty.Get("/rank/season/{season:int}", hdl.GetRankBySeason)
	mockParty.Get("/task", hdl.GetMockTasksAll)
	mockParty.Get("/task/level/{level:string}", hdl.GetMockTasksLevel)
	mockParty.Get("/task/tasks/owner/{owner:string}/repo/{repo:string}", hdl.GetMockTasksRepo)
	mockParty.Get("/taskgroup", hdl.GetMockTaskgroups)
}
