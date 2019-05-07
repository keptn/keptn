package utils

import (
	"context"
	"strings"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// OrgExists checks whether the specified org exists
func OrgExists(accessToken string, orgName string) (bool, error) {

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	orgs, _, err := client.Organizations.List(ctx, "", nil)
	if err != nil {
		return false, err
	}

	for _, org := range orgs {
		if org.Name != nil && *org.Name == orgName ||
			org.URL != nil && strings.Contains(*org.URL, orgName) {
			return true, nil
		}
	}
	return false, nil
}

// CheckRepoScopeOfToken checks whether the provided access token has repo rights
func CheckRepoScopeOfToken(accessToken string) (bool, error) {

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	_, resp, err := client.Organizations.List(ctx, "", nil)
	if err != nil {
		return false, err
	}

	return strings.Contains(resp.Header.Get("X-OAuth-Scopes"), "repo"), nil
}
