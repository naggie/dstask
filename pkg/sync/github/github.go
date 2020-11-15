package github

import (
	"context"
	"crypto/md5"
	"fmt"
	"io"

	"github.com/gofrs/uuid"
	"github.com/naggie/dstask"
	"github.com/naggie/dstask/pkg/sync/config"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

// Github is a client to fetch issues as tasks from GitHub
type Github struct {
	client *githubv4.Client
	cfg    config.Github

	cursor *githubv4.String
	done   bool
}

// NewClient creates a new Github client
func NewClient(cfg config.Github) *Github {

	httpClient := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: cfg.Token},
	))

	return &Github{
		client: githubv4.NewClient(httpClient),
		cfg:    cfg,
	}
}

// Next returns the next batch of issues from a given repository.
func (gh *Github) Next() ([]dstask.Task, error) {

	if gh.done {
		return nil, nil
	}

	var tasks []dstask.Task
	var q Query

	states := []githubv4.IssueState{githubv4.IssueStateOpen}

	if gh.cfg.GetClosed {
		states = append(states, githubv4.IssueStateClosed)
	}
	filterBy := githubv4.IssueFilters{States: &states}

	if gh.cfg.Assignee != "" {
		f := githubv4.String(gh.cfg.Assignee)
		filterBy.Assignee = &f
	}

	variables := map[string]interface{}{
		"owner":       githubv4.String(gh.cfg.User),
		"name":        githubv4.String(gh.cfg.Repo),
		"issueCursor": gh.cursor,
		"count":       githubv4.Int(50),
		"filterBy":    filterBy,
	}

	hash := md5.New() // to write key issue features into, to generate the UUID

	err := gh.client.Query(context.Background(), &q, variables)
	if err != nil {
		return tasks, err
	}

	if len(q.Repository.IssueConnection.Edges) == 0 {
		gh.done = true
		return tasks, nil
	}

	for _, edge := range q.Repository.IssueConnection.Edges {

		io.WriteString(hash, "GH")
		io.WriteString(hash, "\x00")
		io.WriteString(hash, gh.cfg.User)
		io.WriteString(hash, "\x00")
		io.WriteString(hash, gh.cfg.Repo)
		io.WriteString(hash, "\x00")
		io.WriteString(hash, fmt.Sprintf("%d", edge.Node.Number))

		var id uuid.UUID

		task := dstask.Task{
			Summary: fmt.Sprintf("GH/%s/%s/%d: %s", gh.cfg.User, gh.cfg.Repo, edge.Node.Number, edge.Node.Title),
			Status:  dstask.STATUS_PENDING,
			Created: edge.Node.CreatedAt,
		}
		hash.Sum(id[:0])
		hash.Reset()
		task.UUID = id.String()

		if edge.Node.Closed {
			task.Status = dstask.STATUS_RESOLVED
			task.Resolved = edge.Node.ClosedAt
		}
		tasks = append(tasks, task)
	}

	if q.Repository.IssueConnection.PageInfo.HasNextPage {
		gh.cursor = &q.Repository.IssueConnection.PageInfo.EndCursor
	} else {
		gh.done = true
	}

	return tasks, nil
}
