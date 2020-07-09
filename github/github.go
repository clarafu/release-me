package github

import (
	"context"
	"fmt"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

type GitHub struct {
	client *githubv4.Client
}

type PullRequest struct {
	ID     string
	Number int
	Title  string
	Body   string
	Author string
	Labels []string
	Merged bool
	Url    string
}

func (pr PullRequest) HasLabel(label string) bool {
	for _, lbl := range pr.Labels {
		if lbl == label {
			return true
		}
	}
	return false
}

func New(token string) GitHub {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	httpClient := oauth2.NewClient(context.Background(), src)
	return GitHub{
		client: githubv4.NewClient(httpClient),
	}
}

func (g GitHub) FetchLabelsForPullRequest(owner, repo string, pullRequestNumber int) ([]string, error) {
	var PullRequestlabelsQuery struct {
		Repository struct {
			PullRequest struct {
				Labels struct {
					Nodes []struct {
						Name string
					}
				} `graphql:"labels(first: 10)"`
			} `graphql:"pullRequest(number: $prNumber)"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}

	PRVariables := map[string]interface{}{
		"owner":    githubv4.String(owner),
		"name":     githubv4.String(repo),
		"prNumber": githubv4.Int(pullRequestNumber),
	}

	err := g.client.Query(context.Background(), &PullRequestlabelsQuery, PRVariables)
	if err != nil {
		return nil, err
	}

	var labels []string
	for _, node := range PullRequestlabelsQuery.Repository.PullRequest.Labels.Nodes {
		labels = append(labels, node.Name)
	}

	return labels, nil
}

func (g GitHub) FetchLatestReleaseCommitSHA(owner, repo string) (string, error) {
	var releaseSHAQuery struct {
		Repository struct {
			Releases struct {
				Nodes []struct {
					Tag struct {
						Target struct {
							Oid string
						}
					}
				}
			} `graphql:"releases(last: 1, orderBy:{direction: ASC, field: NAME})"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}

	releaseSHAVariables := map[string]interface{}{
		"owner": githubv4.String(owner),
		"name":  githubv4.String(repo),
	}

	err := g.client.Query(context.Background(), &releaseSHAQuery, releaseSHAVariables)
	if err != nil {
		return "", err
	}

	// If the repository does not have a release, return an empty string
	var releaseSHA string
	if len(releaseSHAQuery.Repository.Releases.Nodes) > 0 {
		releaseSHA = releaseSHAQuery.Repository.Releases.Nodes[0].Tag.Target.Oid
	}

	return releaseSHA, nil
}

func (g GitHub) FetchPullRequestsAfterCommit(owner, repo, branch, commitSHA string) ([]PullRequest, error) {
	var pullRequestsQuery struct {
		Repository struct {
			Ref struct {
				Target struct {
					Commit struct {
						History struct {
							Nodes []struct {
								Oid                    string
								AssociatedPullRequests struct {
									Nodes []struct {
										ID     string
										Title  string
										Body   string
										Author struct {
											Login string
										}
										Labels struct {
											Nodes []struct {
												Name string
											}
										} `graphql:"labels(first: 10)"`
										Number int
										Merged bool
										Url    githubv4.URI
									}
								} `graphql:"associatedPullRequests(first: 5)"`
							}
							PageInfo struct {
								EndCursor   githubv4.String
								HasNextPage bool
							}
						} `graphql:"history(first: 100, after: $commitCursor)"`
					} `graphql:"... on Commit"`
				}
			} `graphql:"ref(qualifiedName: $branch)"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}

	pullRequestsVariables := map[string]interface{}{
		"owner":        githubv4.String(owner),
		"name":         githubv4.String(repo),
		"branch":       githubv4.String(branch),
		"commitCursor": (*githubv4.String)(nil),
	}

	pullRequests := []PullRequest{}
	seen := map[string]struct{}{}

	for {
		err := g.client.Query(context.Background(), &pullRequestsQuery, pullRequestsVariables)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch pull requests from github: %w", err)
		}

		for _, commit := range pullRequestsQuery.Repository.Ref.Target.Commit.History.Nodes {
			if commit.Oid == commitSHA {
				return pullRequests, nil
			}

			for _, pr := range commit.AssociatedPullRequests.Nodes {
				if !pr.Merged {
					continue
				}
				if _, exists := seen[pr.ID]; exists {
					continue
				}
				seen[pr.ID] = struct{}{}

				labels := make([]string, len(pr.Labels.Nodes))
				for i, l := range pr.Labels.Nodes {
					labels[i] = l.Name
				}

				url, err := pr.Url.MarshalJSON()
				if err != nil {
					return nil, fmt.Errorf("failed to format url: %w", err)
				}

				pullRequests = append(pullRequests, PullRequest{
					ID:     pr.ID,
					Number: pr.Number,
					Title:  pr.Title,
					Body:   pr.Body,
					Author: pr.Author.Login,
					Labels: labels,
					Merged: pr.Merged,
					Url:    string(url),
				})
			}
		}

		if !pullRequestsQuery.Repository.Ref.Target.Commit.History.PageInfo.HasNextPage {
			return pullRequests, nil
		}

		pullRequestsVariables["commitCursor"] = pullRequestsQuery.Repository.Ref.Target.Commit.History.PageInfo.EndCursor
	}
	return pullRequests, nil
}
