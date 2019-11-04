package utils

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// IsOrgExisting checks whether the specified org exists
func IsOrgExisting(accessToken string, orgName string) (bool, error) {

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	orgs, resp, err := client.Organizations.List(ctx, "", nil)

	if resp.StatusCode == http.StatusUnauthorized {
		fmt.Println("GitHub Personal Access Token is unauthorized")
		return false, nil
	}
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

// HasTokenRepoScope checks whether the provided access token has repo rights
func HasTokenRepoScope(accessToken string) (bool, error) {

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	_, resp, err := client.Organizations.List(ctx, "", nil)
	if err != nil {
		fmt.Println(err)
		return false, err
	}

	if resp.StatusCode == http.StatusUnauthorized {
		fmt.Println("GitHub Personal Access Token is unauthorized")
		return false, nil
	}

	return strings.Contains(resp.Header.Get("X-OAuth-Scopes"), "repo"), nil
}
