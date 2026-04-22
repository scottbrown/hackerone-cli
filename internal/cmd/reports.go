package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	hackeronecli "github.com/scottbrown/hackerone-cli"
	"github.com/spf13/cobra"
)

func NewReportsCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	reportsCmd := &cobra.Command{
		Use:   "reports",
		Short: "Manage HackerOne reports",
	}

	reportsCmd.AddCommand(newReportsListCmd(clientFactory))
	reportsCmd.AddCommand(newReportsGetCmd(clientFactory))
	reportsCmd.AddCommand(newReportsCreateCmd(clientFactory))
	reportsCmd.AddCommand(newReportsCommentCmd(clientFactory))
	reportsCmd.AddCommand(newReportsAssignCmd(clientFactory))
	reportsCmd.AddCommand(newReportsCloseCommentsCmd(clientFactory))
	reportsCmd.AddCommand(newReportsCustomFieldsCmd(clientFactory))
	reportsCmd.AddCommand(newReportsCVEsCmd(clientFactory))
	reportsCmd.AddCommand(newReportsInboxesCmd(clientFactory))
	reportsCmd.AddCommand(newReportsSeverityCmd(clientFactory))
	reportsCmd.AddCommand(newReportsStateCmd(clientFactory))
	reportsCmd.AddCommand(newReportsScopeCmd(clientFactory))
	reportsCmd.AddCommand(newReportsTitleCmd(clientFactory))
	reportsCmd.AddCommand(newReportsWeaknessCmd(clientFactory))
	reportsCmd.AddCommand(newReportsReferenceCmd(clientFactory))
	reportsCmd.AddCommand(newReportsRedactCmd(clientFactory))
	reportsCmd.AddCommand(newReportsSummaryCmd(clientFactory))
	reportsCmd.AddCommand(newReportsPDFCmd(clientFactory))
	reportsCmd.AddCommand(newReportsTransferCmd(clientFactory))
	reportsCmd.AddCommand(newReportsEscalateCmd(clientFactory))
	reportsCmd.AddCommand(newReportsDeescalateCmd(clientFactory))
	reportsCmd.AddCommand(newReportsParticipantsCmd(clientFactory))
	reportsCmd.AddCommand(newReportsAttachmentsCmd(clientFactory))
	reportsCmd.AddCommand(newReportsBountiesCmd(clientFactory))
	reportsCmd.AddCommand(newReportsBountySuggestionsCmd(clientFactory))
	reportsCmd.AddCommand(newReportsDisclosureCmd(clientFactory))
	reportsCmd.AddCommand(newReportsTagsCmd(clientFactory))
	reportsCmd.AddCommand(newReportsRetestCmd(clientFactory))
	reportsCmd.AddCommand(newReportsSwagCmd(clientFactory))

	return reportsCmd
}

func newReportsListCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	var page, pageSize int
	var programs, inboxIDs []string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List reports",
		Long:  "List reports filtered by program handle(s) or inbox ID(s). At least one of --program or --inbox-ids is required.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(programs) == 0 && len(inboxIDs) == 0 {
				return fmt.Errorf("at least one of --program or --inbox-ids is required")
			}
			client, err := clientFactory()
			if err != nil {
				return err
			}
			filter := hackeronecli.ListReportsFilter{
				Programs: programs,
				InboxIDs: inboxIDs,
			}
			reports, err := client.ListReports(cmd.Context(), hackeronecli.PageParams{Number: page, Size: pageSize}, filter)
			if err != nil {
				return err
			}
			return writeOutput(cmd,reports)
		},
	}
	cmd.Flags().IntVar(&page, "page", 0, "Page number")
	cmd.Flags().IntVar(&pageSize, "page-size", 0, "Page size")
	cmd.Flags().StringSliceVar(&programs, "program", nil, "Program handle(s) to filter by (can be specified multiple times)")
	cmd.Flags().StringSliceVar(&inboxIDs, "inbox-ids", nil, "Inbox ID(s) to filter by (can be specified multiple times)")
	return cmd
}

func newReportsGetCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get a report by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			report, err := client.GetReport(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			return writeOutput(cmd,report)
		},
	}
}

func newReportsCreateCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	var input hackeronecli.CreateReportInput

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new report",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			report, err := client.CreateReport(cmd.Context(), input)
			if err != nil {
				return err
			}
			return writeOutput(cmd,report)
		},
	}
	cmd.Flags().StringVar(&input.Title, "title", "", "Report title")
	cmd.Flags().StringVar(&input.VulnerabilityInformation, "vulnerability-info", "", "Vulnerability information")
	cmd.Flags().StringVar(&input.Impact, "impact", "", "Impact description")
	cmd.Flags().StringVar(&input.Severity, "severity", "", "Severity rating")
	cmd.Flags().StringVar(&input.ProgramID, "program-id", "", "Program ID")
	cmd.Flags().StringVar(&input.WeaknessID, "weakness-id", "", "Weakness ID")
	return cmd
}

func newReportsCommentCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	var input hackeronecli.CommentInput

	cmd := &cobra.Command{
		Use:   "comment <id>",
		Short: "Add a comment to a report",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			result, err := client.AddComment(cmd.Context(), args[0], input)
			if err != nil {
				return err
			}
			return writeOutput(cmd,result)
		},
	}
	cmd.Flags().StringVar(&input.Message, "message", "", "Comment message")
	cmd.Flags().BoolVar(&input.Internal, "internal", false, "Internal comment")
	return cmd
}

func newReportsAssignCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	var assigneeID string

	cmd := &cobra.Command{
		Use:   "assign <id>",
		Short: "Update report assignee",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			result, err := client.UpdateAssignee(cmd.Context(), args[0], assigneeID)
			if err != nil {
				return err
			}
			return writeOutput(cmd,result)
		},
	}
	cmd.Flags().StringVar(&assigneeID, "assignee-id", "", "Assignee ID")
	return cmd
}

func newReportsCloseCommentsCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	return &cobra.Command{
		Use:   "close-comments <id>",
		Short: "Close comments on a report",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			result, err := client.CloseComments(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			return writeOutput(cmd,result)
		},
	}
}

func newReportsCustomFieldsCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	var fieldsJSON string

	cmd := &cobra.Command{
		Use:   "custom-fields <id>",
		Short: "Update custom fields on a report",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			var fields map[string]interface{}
			if err := json.Unmarshal([]byte(fieldsJSON), &fields); err != nil {
				return fmt.Errorf("invalid JSON for --fields: %w", err)
			}
			result, err := client.UpdateCustomFields(cmd.Context(), args[0], fields)
			if err != nil {
				return err
			}
			return writeOutput(cmd,result)
		},
	}
	cmd.Flags().StringVar(&fieldsJSON, "fields", "{}", "Custom fields as JSON string")
	return cmd
}

func newReportsCVEsCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	var cvesStr string

	cmd := &cobra.Command{
		Use:   "cves <id>",
		Short: "Update CVEs on a report",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			cves := strings.Split(cvesStr, ",")
			result, err := client.UpdateCVEs(cmd.Context(), args[0], cves)
			if err != nil {
				return err
			}
			return writeOutput(cmd,result)
		},
	}
	cmd.Flags().StringVar(&cvesStr, "cves", "", "Comma-separated CVE IDs")
	return cmd
}

func newReportsInboxesCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	var inboxIDsStr string

	cmd := &cobra.Command{
		Use:   "inboxes <id>",
		Short: "Update inboxes on a report",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			inboxIDs := strings.Split(inboxIDsStr, ",")
			result, err := client.UpdateInboxes(cmd.Context(), args[0], inboxIDs)
			if err != nil {
				return err
			}
			return writeOutput(cmd,result)
		},
	}
	cmd.Flags().StringVar(&inboxIDsStr, "inbox-ids", "", "Comma-separated inbox IDs")
	return cmd
}

func newReportsSeverityCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	var input hackeronecli.SeverityInput

	cmd := &cobra.Command{
		Use:   "severity <id>",
		Short: "Update severity on a report",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			result, err := client.UpdateSeverity(cmd.Context(), args[0], input)
			if err != nil {
				return err
			}
			return writeOutput(cmd,result)
		},
	}
	cmd.Flags().StringVar(&input.Rating, "rating", "", "Severity rating")
	cmd.Flags().Float64Var(&input.Score, "score", 0, "Severity score")
	return cmd
}

func newReportsStateCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	var input hackeronecli.StateChangeInput

	cmd := &cobra.Command{
		Use:   "state <id>",
		Short: "Change report state",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			result, err := client.ChangeState(cmd.Context(), args[0], input)
			if err != nil {
				return err
			}
			return writeOutput(cmd,result)
		},
	}
	cmd.Flags().StringVar(&input.State, "state", "", "New state")
	cmd.Flags().StringVar(&input.Message, "message", "", "State change message")
	return cmd
}

func newReportsScopeCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	var scopeID string

	cmd := &cobra.Command{
		Use:   "scope <id>",
		Short: "Update structured scope on a report",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			result, err := client.UpdateReportScope(cmd.Context(), args[0], scopeID)
			if err != nil {
				return err
			}
			return writeOutput(cmd,result)
		},
	}
	cmd.Flags().StringVar(&scopeID, "scope-id", "", "Scope ID")
	return cmd
}

func newReportsTitleCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	var title string

	cmd := &cobra.Command{
		Use:   "title <id>",
		Short: "Update report title",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			result, err := client.UpdateTitle(cmd.Context(), args[0], title)
			if err != nil {
				return err
			}
			return writeOutput(cmd,result)
		},
	}
	cmd.Flags().StringVar(&title, "title", "", "New title")
	return cmd
}

func newReportsWeaknessCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	var weaknessID string

	cmd := &cobra.Command{
		Use:   "weakness <id>",
		Short: "Update weakness on a report",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			result, err := client.UpdateWeakness(cmd.Context(), args[0], weaknessID)
			if err != nil {
				return err
			}
			return writeOutput(cmd,result)
		},
	}
	cmd.Flags().StringVar(&weaknessID, "weakness-id", "", "Weakness ID")
	return cmd
}

func newReportsReferenceCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	var reference string

	cmd := &cobra.Command{
		Use:   "reference <id>",
		Short: "Update reference on a report",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			result, err := client.UpdateReference(cmd.Context(), args[0], reference)
			if err != nil {
				return err
			}
			return writeOutput(cmd,result)
		},
	}
	cmd.Flags().StringVar(&reference, "reference", "", "Reference string")
	return cmd
}

func newReportsRedactCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	return &cobra.Command{
		Use:   "redact <id>",
		Short: "Redact a report",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			result, err := client.RedactReport(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			return writeOutput(cmd,result)
		},
	}
}

func newReportsSummaryCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	var content string

	cmd := &cobra.Command{
		Use:   "summary <id>",
		Short: "Add summary to a report",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			result, err := client.AddSummary(cmd.Context(), args[0], content)
			if err != nil {
				return err
			}
			return writeOutput(cmd,result)
		},
	}
	cmd.Flags().StringVar(&content, "content", "", "Summary content")
	return cmd
}

func newReportsPDFCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	return &cobra.Command{
		Use:   "pdf <id>",
		Short: "Generate PDF for a report",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			result, err := client.GeneratePDF(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			return writeOutput(cmd,result)
		},
	}
}

func newReportsTransferCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	var programID string

	cmd := &cobra.Command{
		Use:   "transfer <id>",
		Short: "Transfer a report to another program",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			result, err := client.TransferReport(cmd.Context(), args[0], programID)
			if err != nil {
				return err
			}
			return writeOutput(cmd,result)
		},
	}
	cmd.Flags().StringVar(&programID, "program-id", "", "Target program ID")
	return cmd
}

func newReportsEscalateCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	var integration string

	cmd := &cobra.Command{
		Use:   "escalate <id>",
		Short: "Escalate a report",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			result, err := client.EscalateReport(cmd.Context(), args[0], integration)
			if err != nil {
				return err
			}
			return writeOutput(cmd,result)
		},
	}
	cmd.Flags().StringVar(&integration, "integration", "", "Integration name")
	return cmd
}

func newReportsDeescalateCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	return &cobra.Command{
		Use:   "de-escalate <id>",
		Short: "De-escalate a report",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			return client.DeescalateReport(cmd.Context(), args[0])
		},
	}
}

func newReportsParticipantsCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	participantsCmd := &cobra.Command{
		Use:   "participants",
		Short: "Manage report participants",
	}

	var email string
	addCmd := &cobra.Command{
		Use:   "add <id>",
		Short: "Add a participant to a report",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			result, err := client.AddParticipant(cmd.Context(), args[0], email)
			if err != nil {
				return err
			}
			return writeOutput(cmd,result)
		},
	}
	addCmd.Flags().StringVar(&email, "email", "", "Participant email")

	participantsCmd.AddCommand(addCmd)
	return participantsCmd
}

func newReportsAttachmentsCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	attachmentsCmd := &cobra.Command{
		Use:   "attachments",
		Short: "Manage report attachments",
	}

	var filePath string
	uploadCmd := &cobra.Command{
		Use:   "upload <id>",
		Short: "Upload an attachment to a report",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			result, err := client.UploadAttachment(cmd.Context(), args[0], filePath)
			if err != nil {
				return err
			}
			return writeOutput(cmd,result)
		},
	}
	uploadCmd.Flags().StringVar(&filePath, "file", "", "File path to upload")

	deleteCmd := &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete attachments from a report",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			return client.DeleteAttachment(cmd.Context(), args[0])
		},
	}

	attachmentsCmd.AddCommand(uploadCmd)
	attachmentsCmd.AddCommand(deleteCmd)
	return attachmentsCmd
}

func newReportsBountiesCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	bountiesCmd := &cobra.Command{
		Use:   "bounties",
		Short: "Manage report bounties",
	}

	var input hackeronecli.BountyInput
	awardCmd := &cobra.Command{
		Use:   "award <id>",
		Short: "Award a bounty on a report",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			result, err := client.AwardReportBounty(cmd.Context(), args[0], input)
			if err != nil {
				return err
			}
			return writeOutput(cmd,result)
		},
	}
	awardCmd.Flags().Float64Var(&input.Amount, "amount", 0, "Bounty amount")
	awardCmd.Flags().Float64Var(&input.BonusAmount, "bonus-amount", 0, "Bonus amount")
	awardCmd.Flags().StringVar(&input.Message, "message", "", "Bounty message")

	ineligibleCmd := &cobra.Command{
		Use:   "ineligible <id>",
		Short: "Mark report as ineligible for bounty",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			result, err := client.MarkIneligible(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			return writeOutput(cmd,result)
		},
	}

	bountiesCmd.AddCommand(awardCmd)
	bountiesCmd.AddCommand(ineligibleCmd)
	return bountiesCmd
}

func newReportsBountySuggestionsCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	suggestionsCmd := &cobra.Command{
		Use:   "bounty-suggestions",
		Short: "Manage bounty suggestions",
	}

	var page, pageSize int
	listCmd := &cobra.Command{
		Use:   "list <id>",
		Short: "List bounty suggestions for a report",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			suggestions, err := client.ListBountySuggestions(cmd.Context(), args[0], hackeronecli.PageParams{Number: page, Size: pageSize})
			if err != nil {
				return err
			}
			return writeOutput(cmd,suggestions)
		},
	}
	listCmd.Flags().IntVar(&page, "page", 0, "Page number")
	listCmd.Flags().IntVar(&pageSize, "page-size", 0, "Page size")

	var input hackeronecli.CreateBountySuggestionInput
	createCmd := &cobra.Command{
		Use:   "create <id>",
		Short: "Create a bounty suggestion",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			suggestion, err := client.CreateBountySuggestion(cmd.Context(), args[0], input)
			if err != nil {
				return err
			}
			return writeOutput(cmd,suggestion)
		},
	}
	createCmd.Flags().Float64Var(&input.Amount, "amount", 0, "Suggested amount")
	createCmd.Flags().StringVar(&input.Message, "message", "", "Suggestion message")

	suggestionsCmd.AddCommand(listCmd)
	suggestionsCmd.AddCommand(createCmd)
	return suggestionsCmd
}

func newReportsDisclosureCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	disclosureCmd := &cobra.Command{
		Use:   "disclosure",
		Short: "Manage report disclosure",
	}

	var state string
	updateCmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update disclosure state",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			result, err := client.UpdateDisclosure(cmd.Context(), args[0], state)
			if err != nil {
				return err
			}
			return writeOutput(cmd,result)
		},
	}
	updateCmd.Flags().StringVar(&state, "state", "", "Disclosure state")

	cancelCmd := &cobra.Command{
		Use:   "cancel <id>",
		Short: "Cancel disclosure",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			return client.CancelDisclosure(cmd.Context(), args[0])
		},
	}

	disclosureCmd.AddCommand(updateCmd)
	disclosureCmd.AddCommand(cancelCmd)
	return disclosureCmd
}

func newReportsTagsCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	var tagsStr string

	cmd := &cobra.Command{
		Use:   "tags <id>",
		Short: "Update tags on a report",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			tags := strings.Split(tagsStr, ",")
			result, err := client.UpdateTags(cmd.Context(), args[0], tags)
			if err != nil {
				return err
			}
			return writeOutput(cmd,result)
		},
	}
	cmd.Flags().StringVar(&tagsStr, "tags", "", "Comma-separated tags")
	return cmd
}

func newReportsRetestCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	retestCmd := &cobra.Command{
		Use:   "retest",
		Short: "Manage report retests",
	}

	requestCmd := &cobra.Command{
		Use:   "request <id>",
		Short: "Request a retest",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			result, err := client.RequestRetest(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			return writeOutput(cmd,result)
		},
	}

	cancelCmd := &cobra.Command{
		Use:   "cancel <id>",
		Short: "Cancel a retest",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			return client.CancelRetest(cmd.Context(), args[0])
		},
	}

	retestCmd.AddCommand(requestCmd)
	retestCmd.AddCommand(cancelCmd)
	return retestCmd
}

func newReportsSwagCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	return &cobra.Command{
		Use:   "swag <id>",
		Short: "Award swag for a report",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			result, err := client.AwardSwag(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			return writeOutput(cmd,result)
		},
	}
}
