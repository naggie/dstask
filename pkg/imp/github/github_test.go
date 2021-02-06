package github

import (
	"testing"
	"time"

	"github.com/naggie/dstask"
	"gopkg.in/yaml.v2"
)

// NOTE: not sure yet what is the best way to put newlines in actual task yaml files
const tpl1 = `summary: "GH/{{.RepoOwner}}/{{.RepoName}}/{{.Number}}: {{.Title}}"
tags: ["{{.Milestone}}", "extraTag"]
project: "some-project"
priority: P2
notes: "state: {{.State}}\nurl: {{.URL}}\n opened on {{.CreatedAt}} by {{.Author}}\n{{.Body}}"`

func TestToTask(t *testing.T) {
	type testCase struct {
		tpl     string
		owner   string
		repo    string
		issue   Issue
		expTask dstask.Task
	}
	cases := []testCase{
		{
			tpl:   tpl1,
			owner: "my_user",
			repo:  "my_repo",
			issue: Issue{
				Number: 1234,
				Body:   "body content of issue",
				Title:  "ISSUE-TITLE",
				Author: Author{
					Name: "author-name",
				},
				URL:       "http://github.com/my_user/my_repo/1234",
				CreatedAt: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
				Milestone: Milestone{
					Description: "some milestone description", // we ignore this, can't use this in templates
					Number:      42,                           // we ignore this, can't use this in templates
					Title:       "my-milestone",               // this is what we allow to expand
				},
				State:  "OPEN",
				Closed: false,
			},
			expTask: dstask.Task{
				UUID:        "e09f6975-8a79-133d-78f8-837b85a1754c",
				Status:      "pending",
				Summary:     "GH/my_user/my_repo/1234: ISSUE-TITLE",
				Notes:       "state: OPEN\nurl: http://github.com/my_user/my_repo/1234\n opened on 2009-11-10 23:00:00 +0000 UTC by author-name\nbody content of issue",
				Tags:        []string{"my-milestone", "extraTag"},
				Project:     "some-project",
				Priority:    "P2",
				DelegatedTo: "",
				Created:     time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
			},
		},
		{
			tpl:   tpl1,
			owner: "my_user",
			repo:  "my_repo",
			issue: Issue{
				Number: 1234,
				Body:   "body content of issue",
				Title:  "ISSUE-TITLE",
				Author: Author{
					Name: "author-name",
				},
				URL:       "http://github.com/my_user/my_repo/1234",
				CreatedAt: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
				Milestone: Milestone{},
				State:     "CLOSED",
				Closed:    true,
				ClosedAt:  time.Date(2020, time.January, 10, 23, 0, 0, 0, time.UTC),
			},
			expTask: dstask.Task{
				UUID:        "e09f6975-8a79-133d-78f8-837b85a1754c",
				Status:      "resolved",
				Summary:     "GH/my_user/my_repo/1234: ISSUE-TITLE",
				Notes:       "state: CLOSED\nurl: http://github.com/my_user/my_repo/1234\n opened on 2009-11-10 23:00:00 +0000 UTC by author-name\nbody content of issue",
				Tags:        []string{"extraTag"},
				Project:     "some-project",
				Priority:    "P2",
				DelegatedTo: "",
				Created:     time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
				Resolved:    time.Date(2020, time.January, 10, 23, 0, 0, 0, time.UTC),
			},
		},
	}

	issueData := NewIssueData()

	for _, c := range cases {

		var tplTask dstask.Task
		err := yaml.Unmarshal([]byte(c.tpl), &tplTask)
		if err != nil {
			t.Fatalf("Failed to unmarshal template: %s", err.Error())
		}

		issueData.Init(c.owner, c.repo, c.issue)
		templates := ParseTemplates(tplTask)
		task, _ := issueData.ToTask(templates)
		if !task.Equals(c.expTask) {
			t.Errorf("ToTask() mismatch.\nwant:\n%#v\ngot:\n%#v\n", c.expTask, task)
		}

	}
}
