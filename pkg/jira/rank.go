package jira

// RankInput represents the data required to rank an issue.
type RankInput struct {
	Issues         []string `json:"issues"`
	RankAfterIssue string   `json:"rankAfterIssue,omitempty"`
	RankFirst      bool     `json:"rankFirst,omitempty"`
}
