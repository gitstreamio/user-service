package client

import (
	"golang.org/x/oauth2"
	"github.com/google/go-github/github"
	"context"
)

func CreateGitHubUserClient(ctx context.Context, accessToken string) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	return github.NewClient(tc)
}
