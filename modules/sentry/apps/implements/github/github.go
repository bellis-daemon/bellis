package github

import (
	"context"
	"fmt"
	"github.com/bellis-daemon/bellis/modules/sentry/apps/implements"
	"github.com/bellis-daemon/bellis/modules/sentry/apps/option"
	"github.com/bellis-daemon/bellis/modules/sentry/apps/status"
	githubLib "github.com/google/go-github/v32/github"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/oauth2"
	"net/http"
	"strings"
	"time"
)

type GitHub struct {
	options         githubOptions
	githubClient    *githubLib.Client
	obfuscatedToken string
	rateRemaining   int
	rateReset       time.Time
}

func (this *GitHub) Fetch(ctx context.Context) (status.Status, error) {
	if this.rateRemaining <= 0 {
		if time.Now().After(this.rateReset) {
			this.rateRemaining = 5000
		} else {
			return &githubStatus{}, fmt.Errorf("exceeds github api rate limit until %s", this.rateReset.Format(time.DateTime))
		}
	}
	owner, repository, err := splitRepositoryName(this.options.Repository)
	if err != nil {
		return &githubStatus{}, err
	}
	repositoryInfo, response, err := this.githubClient.Repositories.Get(ctx, owner, repository)
	if err != nil {
		return &githubStatus{}, err
	}
	this.rateRemaining = response.Rate.Remaining
	this.rateReset = response.Rate.Reset.Time
	return &githubStatus{
		Stars:       repositoryInfo.GetStargazersCount(),
		Subscribers: repositoryInfo.GetSubscribersCount(),
		Watchers:    repositoryInfo.GetWatchersCount(),
		Networks:    repositoryInfo.GetNetworkCount(),
		Forks:       repositoryInfo.GetForksCount(),
		OpenIssues:  repositoryInfo.GetOpenIssuesCount(),
		Size:        repositoryInfo.GetSize(),
		Language:    repositoryInfo.GetLanguage(),
	}, nil
}

func (this *GitHub) createGitHubClient() (*githubLib.Client, error) {
	httpClient := &http.Client{}

	this.obfuscatedToken = "Unauthenticated"

	if this.options.AccessToken != "" {
		tokenSource := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: this.options.AccessToken},
		)
		oauthClient := oauth2.NewClient(context.Background(), tokenSource)
		_ = context.WithValue(context.Background(), oauth2.HTTPClient, oauthClient)

		this.obfuscatedToken = this.options.AccessToken[0:4] + "..." + this.options.AccessToken[len(this.options.AccessToken)-3:]

		return this.newGithubClient(oauthClient)
	}

	return this.newGithubClient(httpClient)
}

func (this *GitHub) newGithubClient(httpClient *http.Client) (*githubLib.Client, error) {
	if this.options.EnterpriseBaseURL != "" {
		return githubLib.NewEnterpriseClient(this.options.EnterpriseBaseURL, "", httpClient)
	}
	return githubLib.NewClient(httpClient), nil
}

type githubOptions struct {
	Repository        string
	AccessToken       string
	EnterpriseBaseURL string
}

type githubStatus struct {
	Stars       int
	Subscribers int
	Watchers    int
	Networks    int
	Forks       int
	OpenIssues  int
	Size        int
	Language    string
}

func (this *githubStatus) PullTrigger(triggerName string) *status.TriggerInfo {
	switch triggerName {

	}
	return nil
}

func splitRepositoryName(repositoryName string) (owner string, repository string, err error) {
	splits := strings.SplitN(repositoryName, "/", 2)

	if len(splits) != 2 {
		return "", "", fmt.Errorf("%v is not of format 'owner/repository'", repositoryName)
	}

	return splits[0], splits[1], nil
}

func init() {
	implements.Register("github", func(options bson.M) implements.Implement {
		ret := &GitHub{
			options:       option.ToOption[githubOptions](options),
			githubClient:  nil,
			rateRemaining: 5000,
			rateReset:     time.Now(),
		}
		var err error
		ret.githubClient, err = ret.createGitHubClient()
		if err != nil {
			panic(err)
		}
		return ret
	})
}
