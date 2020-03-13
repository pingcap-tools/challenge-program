package github

import (
	"context"
	"github.com/google/go-github/github"
	"github.com/ngaut/log"
	"github.com/pingcap/errors"
	"golang.org/x/oauth2"
	"os"
)

// GetGithubClient return client with auth
func GetGithubClient(token string) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(context.Background(), ts)

	client := github.NewClient(tc)
	//projectCardOptions := github.ProjectCardOptions{
	//	Note: "noteaaa",
	//}
	//c, _ , err := client.Projects.CreateProjectCard(context.Background(), 6873026, &projectCardOptions)
	//log.Info(c, err)
	// get project snippet
	// p, _, err := client.Repositories.ListProjects(context.Background(), "pingcap", "pd", nil)
	// log.Info(p, err)
	// os.Exit(1)
	// cs, _, err := client.Projects.ListProjectColumns(context.Background(), 3420390, nil)
	// log.Info(cs, err)
	// for _, c := range cs {
	// 	log.Info(c.GetName(), c.GetID())
	// }
	// os.Exit(1)
	user, _, err := client.Users.Get(context.Background(), "")
	if err != nil {
		log.Errorf("get user info failed %v", errors.ErrorStack(err))
		os.Exit(1)
	}
	log.Infof("token user %s", user.GetLogin())
	return client
}
