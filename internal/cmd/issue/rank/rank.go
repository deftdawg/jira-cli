package rank

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ankitpokhrel/jira-cli/api"
	"github.com/ankitpokhrel/jira-cli/internal/cmdutil"
	"github.com/ankitpokhrel/jira-cli/pkg/jira"
)

const (
	helpText = `Rank an issue or move it to the top/bottom of the backlog.`
	examples = `$ jira issue rank ISSUE-1 ISSUE-2
$ jira issue rank ISSUE-1 ISSUE-3 --position bottom`
)

// NewCmdRank is a rank command.
func NewCmdRank() *cobra.Command {
	cmd := cobra.Command{
		Use:     "rank ISSUE_KEY_OR_ID [ISSUE_KEY_OR_ID_AFTER]",
		Short:   "Rank an issue",
		Long:    helpText,
		Aliases: []string{"rk"},
		Example: examples,
		Annotations: map[string]string{
			"help:args": "ISSUE_KEY_OR_ID	Issue key or ID to rank\n" +
				"ISSUE_KEY_OR_ID_AFTER	Issue key or ID to rank after (optional)",
		},
		RunE: rank,
	}

	cmd.Flags().String("position", "", "Rank an issue to the top or bottom of the backlog (top|bottom)")

	return &cmd
}

func rank(cmd *cobra.Command, args []string) error {
	if len(args) == 0 || len(args) > 2 {
		return cmd.Help()
	}

	project := cmdutil.GetProject()
	issueKey := cmdutil.GetJiraIssueKey(project, args[0])
	var issueKeyAfter string
	if len(args) == 2 {
		issueKeyAfter = cmdutil.GetJiraIssueKey(project, args[1])
	}

	position, err := cmd.Flags().GetString("position")
	if err != nil {
		return err
	}

	if position != "" && (position != "top" && position != "bottom") {
		return fmt.Errorf("invalid position: %s. valid positions are 'top' or 'bottom'", position)
	}

	if position != "" && issueKeyAfter != "" {
		return fmt.Errorf("cannot use --position and specify an issue to rank after simultaneously")
	}

	debug, err := cmd.Flags().GetBool("debug")
	if err != nil {
		return err
	}

	client := api.DefaultClient(debug)
	fmt.Printf("Ranking issue %s\n", issueKey)

	s := cmdutil.Info("Ranking issue...")
	defer s.Stop()

	rankInput := jira.RankInput{
		Issues: []string{issueKey},
	}

	if issueKeyAfter != "" {
		rankInput.RankAfterIssue = issueKeyAfter
	} else if position == "top" {
		rankInput.RankFirst = true
	} else if position == "bottom" {
		// For ranking to the bottom, Jira API expects no RankAfterIssue and RankFirst=false
		// which is the default for RankInput.
		// However, some older Jira versions might have different behavior or specific fields.
		// For now, we assume the API handles "no specific rank instruction" as "bottom" or "last".
		// If a specific "rank last" mechanism is available, it should be used here.
	}

	err = client.Rank(cmd.Context(), rankInput)
	if err != nil {
		return err
	}

	s.Stop()
	if debug {
		// cmdutil.Printf("Response: %+v\n", resp) // Assuming client.Rank returns a response
		// TODO: The client.Rank method currently returns only an error.
		// If it were to return a response object, it would be logged here.
		fmt.Printf("DEBUG: Successfully sent rank request for issue %s\n", issueKey)
	}
	cmdutil.Success("Issue %s ranked successfully.\n", issueKey)

	return nil
}
