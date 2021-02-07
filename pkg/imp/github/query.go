package github

import (
	"time"

	"github.com/shurcooL/githubv4"
)

type Query struct {
	Repository struct {
		IssueConnection IssueConnection `graphql:"issues(first: $count, after: $issueCursor, filterBy: $filterBy)"`
	} `graphql:"repository(owner: $owner, name: $name)"`
}

type QueryWithMilestone struct {
	Repository struct {
		Milestone struct {
			IssueConnection IssueConnection `graphql:"issues(first: $count, after: $issueCursor, filterBy: $filterBy)"`
		} `graphql:"milestone(number: $milestone)"`
	} `graphql:"repository(owner: $owner, name: $name)"`
}

type IssueConnection struct {
	Edges    []IssueEdge `graphql:"edges"`
	PageInfo PageInfo    `graphql:"pageInfo"`
}

// PageInfo helps with the paging large query responses
type PageInfo struct {
	EndCursor   githubv4.String
	HasNextPage bool
}

type IssueEdge struct {
	Cursor string `graphql:"cursor"`
	Node   Issue  `graphql:"node"`
}

type Issue struct {
	ID        string    `graphql:"id"`
	Number    int       `graphql:"number"`
	Body      string    `graphql:"body"`
	Title     string    `graphql:"title"`
	Author    Author    `graphql:"author"`
	URL       string    `graphql:"url"`
	CreatedAt time.Time `graphql:"createdAt"`
	Milestone Milestone `graphql:"milestone"`
	State     string    `graphql:"state"`
	Closed    bool      `graphql:"closed"`
	ClosedAt  time.Time `graphql:"closedAt"`
}

type Author struct {
	Name string `graphql:"login"`
}

type MilestoneQuery struct {
	Repository struct {
		Milestones MilestoneConnection `graphql:"milestones(first: 100, query: $query)"`
	} `graphql:"repository(owner: $owner, name: $name)"`
}

type MilestoneConnection struct {
	Edges    []MilestoneEdge `graphql:"edges"`
	PageInfo PageInfo        `graphql:"pageInfo"`
}

type MilestoneEdge struct {
	Cursor string    `graphql:"cursor"`
	Node   Milestone `graphql:"node"`
}

type Milestone struct {
	Description string
	Number      int
	Title       string
}
