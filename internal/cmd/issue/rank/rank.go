package rank

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ankitpokhrel/jira-cli/internal/cmdutil"
	"github.com/ankitpokhrel/jira-cli/pkg/jira"
)

const (
	helpText = `Rank an issue or issues relative to another issue.
You must specify the target issue(s) and a reference issue to rank before or after.`
	examples = `$ jira issue rank ISSUE-1 --after ISSUE-2
$ jira issue rank ISSUE-1,ISSUE-3 --before ISSUE-4`
)

// NewCmdRank is a rank command.
func NewCmdRank() *cobra.Command {
	cmd := cobra.Command{
		Use:     "rank <ISSUE_KEY_OR_KEYS>",
		Short:   "Rank an issue or issues relative to another issue",
		Long:    helpText,
		Example: examples,
		Aliases: []string{},
		Run:     rank,
	}

	cmd.Flags().String("after", "", "Reference issue key to rank target issue(s) after")
	cmd.Flags().String("before", "", "Reference issue key to rank target issue(s) before")

	return &cmd
}

func rank(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		cmdutil.Failed("Missing required argument: <ISSUE_KEY_OR_KEYS>")
		return
	}
	targetIssuesStr := args[0]
	if targetIssuesStr == "" {
		cmdutil.Failed("Target issue key(s) cannot be empty.")
		return
	}
	targetIssueKeys := strings.Split(targetIssuesStr, ",")
	if len(targetIssueKeys) == 0 {
		cmdutil.Failed("No target issue keys provided.")
		return
	}
	for _, key := range targetIssueKeys {
		if strings.TrimSpace(key) == "" {
			cmdutil.Failed("One of the target issue keys is empty.")
			return
		}
	}

	beforeKey, _ := cmd.Flags().GetString("before")
	afterKey, _ := cmd.Flags().GetString("after")

	if beforeKey == "" && afterKey == "" {
		cmdutil.Failed("You must specify either --before or --after a reference issue.")
		return
	}
	if beforeKey != "" && afterKey != "" {
		cmdutil.Failed("You cannot specify both --before and --after.")
		return
	}
	if (beforeKey != "" && strings.TrimSpace(beforeKey) == "") || (afterKey != "" && strings.TrimSpace(afterKey) == "") {
		cmdutil.Failed("Reference issue key for --before or --after cannot be empty.")
		return
	}
	
	debug, err := cmd.Flags().GetBool("debug")
	if err != nil {
		cmdutil.Printf("Failed to get debug mode: %s\n", err)
	}

	client, err := cmdutil.Client(debug)
	if err != nil {
		cmdutil.Failed("Failed to get Jira client: %v", err)
		return
	}

	payload := jira.IssueRankPayload{
		Issues:          targetIssueKeys,
		RankBeforeIssue: strings.TrimSpace(beforeKey),
		RankAfterIssue:  strings.TrimSpace(afterKey),
	}

	err = client.RankIssues(payload)
	if err != nil {
		cmdutil.Failed("Failed to rank issues: %v", err)
		return
	}

	cmdutil.Success("Issue(s) ranked successfully.")
}
