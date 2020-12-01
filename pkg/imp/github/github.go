package github

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"text/template"

	"github.com/naggie/dstask"
	"github.com/naggie/dstask/pkg/imp"
	"github.com/naggie/dstask/pkg/imp/config"
	"github.com/shurcooL/githubv4"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

// Do runs the Github import if the user requested it
func Do(dstaskRepo string, cfg config.Config) error {
	if len(cfg.Github) == 0 {
		return nil
	}

	for i, cfgGithub := range cfg.Github {
		if cfgGithub.Token == "" {
			logrus.Infof("GitHub config section %d (%v): skipping because no token configured", i, cfgGithub.Repos)
			continue
		}
		logrus.Infof("GitHub config section %d (%v): processing", i, cfgGithub.Repos)

		gh, err := NewClient(cfgGithub)
		if err != nil {
			return err
		}
		err = gh.Run(dstaskRepo)
		if err != nil {
			return err
		}
	}
	dstask.MustGitCommit(dstaskRepo, "GitHub import")
	return nil
}

// Github is a client to process multiple repos with a given template
type Github struct {
	client *githubv4.Client
	cfg    config.Github

	// options derived from config
	templates Templates
}

// NewClient creates a new Github client
func NewClient(cfg config.Github) (*Github, error) {

	httpClient := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: cfg.Token},
	))

	g := Github{
		client:    githubv4.NewClient(httpClient),
		cfg:       cfg,
		templates: ParseTemplates(cfg.TemplateTask),
	}
	return &g, nil
}

// Run processes the issues from all requested repositories
func (gh *Github) Run(dstaskRepo string) error {
	for _, r := range gh.cfg.Repos {
		iter, err := NewRepoIter(gh.cfg, r, gh.templates, gh.client)
		if err != nil {
			return err
		}
		logrus.Infof("GitHub starting iteration for %s", r)

		for {
			tasks, err := iter.Next()
			if err != nil {
				return err
			}
			if len(tasks) == 0 {
				break
			}

			for _, t := range tasks {
				err = imp.ProcessTask(dstaskRepo, t)
				if err != nil {
					return err
				}
			}
		}

	}
	return nil
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

// RepoIter iterates all the desired issues from a given repo
type RepoIter struct {
	client *githubv4.Client
	cfg    config.Github

	// options derived from config
	repoOwner string
	repoName  string
	milestone int // for filtering
	templates Templates

	// iterator state
	done   bool
	cursor *githubv4.String
}

func NewRepoIter(cfg config.Github, repo string, templates Templates, client *githubv4.Client) (*RepoIter, error) {
	repoSplit := strings.Split(repo, "/")
	if len(repoSplit) != 2 || repoSplit[0] == "" || repoSplit[1] == "" {
		return nil, fmt.Errorf("invalid repo %q", repo)
	}

	ri := RepoIter{
		repoOwner: repoSplit[0],
		repoName:  repoSplit[1],
		client:    client,
		cfg:       cfg,
		templates: templates,
	}

	if cfg.Milestone != "" {
		// we first must figure out the id of the milestone

		var mq MilestoneQuery
		variables := map[string]interface{}{
			"owner": githubv4.String(ri.repoOwner),
			"name":  githubv4.String(ri.repoName),
			"query": githubv4.String(cfg.Milestone),
		}
		err := client.Query(context.Background(), &mq, variables)
		if err != nil {
			return nil, fmt.Errorf("could execute lookup query for milestone %q: %s", cfg.Milestone, err.Error())
		}
		if len(mq.Repository.Milestones.Edges) != 1 {
			return nil, fmt.Errorf("could not look up milestone %q: got %d results", cfg.Milestone, len(mq.Repository.Milestones.Edges))
		}
		ri.milestone = mq.Repository.Milestones.Edges[0].Node.Number
	}
	return &ri, nil
}

// Next returns the next batch of tasks, if there are any
func (ri *RepoIter) Next() ([]dstask.Task, error) {
	if ri.done {
		return nil, nil
	}
	var tasks []dstask.Task

	states := []githubv4.IssueState{githubv4.IssueStateOpen}

	if ri.cfg.GetClosed {
		states = append(states, githubv4.IssueStateClosed)
	}
	filterBy := githubv4.IssueFilters{States: &states}

	if ri.cfg.Assignee != "" {
		f := githubv4.String(ri.cfg.Assignee)
		filterBy.Assignee = &f
	}
	if len(ri.cfg.Labels) != 0 {
		var labels []githubv4.String
		for _, label := range ri.cfg.Labels {
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
		"owner":       githubv4.String(ri.repoOwner),
		"name":        githubv4.String(ri.repoName),
		"issueCursor": ri.cursor,
		"count":       githubv4.Int(50),
		"filterBy":    filterBy,
	}

	var err error
	var issues IssueConnection

	if ri.cfg.Milestone == "" {
		var q Query
		err = ri.client.Query(context.Background(), &q, variables)
		issues = q.Repository.IssueConnection
	} else {
		var q QueryWithMilestone
		variables["milestone"] = githubv4.Int(ri.milestone)
		err = ri.client.Query(context.Background(), &q, variables)
		issues = q.Repository.Milestone.IssueConnection
	}
	if err != nil {
		return tasks, err
	}

	if len(issues.Edges) == 0 {
		ri.done = true
		return tasks, nil
	}

	issueData := NewIssueData()

	for _, edge := range issues.Edges {

		issueData.Init(ri.repoOwner, ri.repoName, edge.Node)
		task, err := issueData.ToTask(ri.templates)
		if err != nil {
			return tasks, err
		}

		tasks = append(tasks, task)
	}

	if issues.PageInfo.HasNextPage {
		ri.cursor = &issues.PageInfo.EndCursor
	} else {
		ri.done = true
	}

	return tasks, nil
}
