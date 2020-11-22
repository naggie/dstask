package github

import (
	"context"
	"fmt"
	"strconv"
	"text/template"

	"github.com/naggie/dstask"
	"github.com/naggie/dstask/pkg/sync/config"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

// Github is a client to fetch issues as tasks from GitHub
type Github struct {
	client *githubv4.Client
	cfg    config.Github

	// options derived from config
	milestone int // for filtering
	templates Templates

	// iterator state
	done   bool
	cursor *githubv4.String
}

type Templates struct {
	Summary  *template.Template
	Project  *template.Template
	Priority *template.Template
	Notes    *template.Template
	Tags     []*template.Template
}

func ParseTemplates(task dstask.Task) Templates {
	t := Templates{
		Summary:  template.Must(template.New("summary").Parse(task.Summary)),
		Project:  template.Must(template.New("project").Parse(task.Project)),
		Priority: template.Must(template.New("priority").Parse(task.Priority)),
		Notes:    template.Must(template.New("notes").Parse(task.Notes)),
	}

	for i, tag := range task.Tags {
		t.Tags = append(t.Tags, template.Must(template.New("tag"+strconv.Itoa(i)).Parse(tag)))
	}

	return t
}

// NewClient creates a new Github client
func NewClient(cfg config.Github) (*Github, error) {

	httpClient := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: cfg.Token},
	))

	g := Github{
		client: githubv4.NewClient(httpClient),
		cfg:    cfg,
	}
	if cfg.Milestone != "" {
		// we first must figure out the id of the milestone

		var mq MilestoneQuery
		variables := map[string]interface{}{
			"owner": githubv4.String(cfg.User),
			"name":  githubv4.String(cfg.Repo),
			"query": githubv4.String(cfg.Milestone),
		}
		err := g.client.Query(context.Background(), &mq, variables)
		if err != nil {
			return nil, fmt.Errorf("could execute lookup query for milestone %q: %s", cfg.Milestone, err.Error())
		}
		if len(mq.Repository.Milestones.Edges) != 1 {
			return nil, fmt.Errorf("could not look up milestone %q: got %d results", cfg.Milestone, len(mq.Repository.Milestones.Edges))
		}
		g.milestone = mq.Repository.Milestones.Edges[0].Node.Number
	}

	g.templates = ParseTemplates(cfg.TemplateTask)
	return &g, nil
}

// Next returns the next batch of issues from a given repository.
func (gh *Github) Next() ([]dstask.Task, error) {

	if gh.done {
		return nil, nil
	}

	var tasks []dstask.Task

	states := []githubv4.IssueState{githubv4.IssueStateOpen}

	if gh.cfg.GetClosed {
		states = append(states, githubv4.IssueStateClosed)
	}
	filterBy := githubv4.IssueFilters{States: &states}

	if gh.cfg.Assignee != "" {
		f := githubv4.String(gh.cfg.Assignee)
		filterBy.Assignee = &f
	}
	if len(gh.cfg.Labels) != 0 {
		var labels []githubv4.String
		for _, label := range gh.cfg.Labels {
			labels = append(labels, githubv4.String(label))
		}
		filterBy.Labels = &labels
	}

	// you would think that filtering by milestone can be done by something like this:
	//	if gh.cfg.Milestone != "" {
	//		f := githubv4.String("33")             // either by id...
	//              f := githubv4.String(gh.cfg.Milestone) // ... or by name
	//		filterBy.Milestone = &f
	//	}
	// .... but neither of these seem to work.
	// seems you need to first lookup the milestone ID and then use a different query type altogether.

	variables := map[string]interface{}{
		"owner":       githubv4.String(gh.cfg.User),
		"name":        githubv4.String(gh.cfg.Repo),
		"issueCursor": gh.cursor,
		"count":       githubv4.Int(50),
		"filterBy":    filterBy,
	}

	var err error
	var issues IssueConnection

	if gh.cfg.Milestone == "" {
		var q Query
		err = gh.client.Query(context.Background(), &q, variables)
		issues = q.Repository.IssueConnection
	} else {
		var q QueryWithMilestone
		variables["milestone"] = githubv4.Int(gh.milestone)
		err = gh.client.Query(context.Background(), &q, variables)
		issues = q.Repository.Milestone.IssueConnection
	}
	if err != nil {
		return tasks, err
	}

	if len(issues.Edges) == 0 {
		gh.done = true
		return tasks, nil
	}

	issueData := NewIssueData()

	for _, edge := range issues.Edges {

		issueData.Init(gh.cfg.User, gh.cfg.Repo, edge.Node)
		task, err := issueData.ToTask(gh.templates)
		if err != nil {
			return tasks, err
		}

		tasks = append(tasks, task)

	}

	if issues.PageInfo.HasNextPage {
		gh.cursor = &issues.PageInfo.EndCursor
	} else {
		gh.done = true
	}

	return tasks, nil
}
