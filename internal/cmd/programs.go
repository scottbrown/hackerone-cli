package cmd

import (
	hackeronecli "github.com/scottbrown/hackerone-cli"
	"github.com/spf13/cobra"
)

func NewProgramsCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	programsCmd := &cobra.Command{
		Use:   "programs",
		Short: "Manage HackerOne programs",
	}

	programsCmd.AddCommand(newProgramsListCmd(clientFactory))
	programsCmd.AddCommand(newProgramsGetCmd(clientFactory))
	programsCmd.AddCommand(newProgramsAuditLogCmd(clientFactory))
	programsCmd.AddCommand(newProgramsBalanceCmd(clientFactory))
	programsCmd.AddCommand(newProgramsPaymentTransactionsCmd(clientFactory))
	programsCmd.AddCommand(newProgramsCommonResponsesCmd(clientFactory))
	programsCmd.AddCommand(newProgramsReportersCmd(clientFactory))
	programsCmd.AddCommand(newProgramsMembersCmd(clientFactory))
	programsCmd.AddCommand(newProgramsThanksCmd(clientFactory))
	programsCmd.AddCommand(newProgramsIntegrationsCmd(clientFactory))
	programsCmd.AddCommand(newProgramsTriageReviewsCmd(clientFactory))
	programsCmd.AddCommand(newProgramsWeaknessesCmd(clientFactory))
	programsCmd.AddCommand(newProgramsNotifyCmd(clientFactory))
	programsCmd.AddCommand(newProgramsMessagesCmd(clientFactory))
	programsCmd.AddCommand(newProgramsBountiesCmd(clientFactory))
	programsCmd.AddCommand(newProgramsAllowedReportersCmd(clientFactory))
	programsCmd.AddCommand(newProgramsCVERequestsCmd(clientFactory))
	programsCmd.AddCommand(newProgramsHackerInvitationsCmd(clientFactory))
	programsCmd.AddCommand(newProgramsPolicyCmd(clientFactory))
	programsCmd.AddCommand(newProgramsScopesCmd(clientFactory))
	programsCmd.AddCommand(newProgramsSwagCmd(clientFactory))

	return programsCmd
}

func newProgramsListCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	var page, pageSize int

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List programs",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			programs, err := client.ListPrograms(cmd.Context(), hackeronecli.PageParams{Number: page, Size: pageSize})
			if err != nil {
				return err
			}
			return writeOutput(cmd,programs)
		},
	}
	cmd.Flags().IntVar(&page, "page", 0, "Page number")
	cmd.Flags().IntVar(&pageSize, "page-size", 0, "Page size")
	return cmd
}

func newProgramsGetCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get a program by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			program, err := client.GetProgram(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			return writeOutput(cmd,program)
		},
	}
}

func newProgramsAuditLogCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	var page, pageSize int

	cmd := &cobra.Command{
		Use:   "audit-log <id>",
		Short: "Get audit log for a program",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			entries, err := client.GetProgramAuditLog(cmd.Context(), args[0], hackeronecli.PageParams{Number: page, Size: pageSize})
			if err != nil {
				return err
			}
			return writeOutput(cmd,entries)
		},
	}
	cmd.Flags().IntVar(&page, "page", 0, "Page number")
	cmd.Flags().IntVar(&pageSize, "page-size", 0, "Page size")
	return cmd
}

func newProgramsBalanceCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	return &cobra.Command{
		Use:   "balance <id>",
		Short: "Get program balance",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			balance, err := client.GetProgramBalance(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			return writeOutput(cmd,balance)
		},
	}
}

func newProgramsPaymentTransactionsCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	var page, pageSize int

	cmd := &cobra.Command{
		Use:   "payment-transactions <id>",
		Short: "List payment transactions for a program",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			txns, err := client.ListPaymentTransactions(cmd.Context(), args[0], hackeronecli.PageParams{Number: page, Size: pageSize})
			if err != nil {
				return err
			}
			return writeOutput(cmd,txns)
		},
	}
	cmd.Flags().IntVar(&page, "page", 0, "Page number")
	cmd.Flags().IntVar(&pageSize, "page-size", 0, "Page size")
	return cmd
}

func newProgramsCommonResponsesCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	var page, pageSize int

	cmd := &cobra.Command{
		Use:   "common-responses <id>",
		Short: "List common responses for a program",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			responses, err := client.ListCommonResponses(cmd.Context(), args[0], hackeronecli.PageParams{Number: page, Size: pageSize})
			if err != nil {
				return err
			}
			return writeOutput(cmd,responses)
		},
	}
	cmd.Flags().IntVar(&page, "page", 0, "Page number")
	cmd.Flags().IntVar(&pageSize, "page-size", 0, "Page size")
	return cmd
}

func newProgramsReportersCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	var page, pageSize int

	cmd := &cobra.Command{
		Use:   "reporters <id>",
		Short: "List reporters for a program",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			reporters, err := client.ListReporters(cmd.Context(), args[0], hackeronecli.PageParams{Number: page, Size: pageSize})
			if err != nil {
				return err
			}
			return writeOutput(cmd,reporters)
		},
	}
	cmd.Flags().IntVar(&page, "page", 0, "Page number")
	cmd.Flags().IntVar(&pageSize, "page-size", 0, "Page size")
	return cmd
}

func newProgramsMembersCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	var page, pageSize int

	cmd := &cobra.Command{
		Use:   "members <id>",
		Short: "List team members for a program",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			members, err := client.ListTeamMembers(cmd.Context(), args[0], hackeronecli.PageParams{Number: page, Size: pageSize})
			if err != nil {
				return err
			}
			return writeOutput(cmd,members)
		},
	}
	cmd.Flags().IntVar(&page, "page", 0, "Page number")
	cmd.Flags().IntVar(&pageSize, "page-size", 0, "Page size")
	return cmd
}

func newProgramsThanksCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	var page, pageSize int

	cmd := &cobra.Command{
		Use:   "thanks <id>",
		Short: "List thanks for a program",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			thanks, err := client.ListThanks(cmd.Context(), args[0], hackeronecli.PageParams{Number: page, Size: pageSize})
			if err != nil {
				return err
			}
			return writeOutput(cmd,thanks)
		},
	}
	cmd.Flags().IntVar(&page, "page", 0, "Page number")
	cmd.Flags().IntVar(&pageSize, "page-size", 0, "Page size")
	return cmd
}

func newProgramsIntegrationsCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	var page, pageSize int

	cmd := &cobra.Command{
		Use:   "integrations <id>",
		Short: "List integrations for a program",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			integrations, err := client.ListIntegrations(cmd.Context(), args[0], hackeronecli.PageParams{Number: page, Size: pageSize})
			if err != nil {
				return err
			}
			return writeOutput(cmd,integrations)
		},
	}
	cmd.Flags().IntVar(&page, "page", 0, "Page number")
	cmd.Flags().IntVar(&pageSize, "page-size", 0, "Page size")
	return cmd
}

func newProgramsTriageReviewsCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	var page, pageSize int

	cmd := &cobra.Command{
		Use:   "triage-reviews <id>",
		Short: "List triage reviews for a program",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			reviews, err := client.ListTriageReviews(cmd.Context(), args[0], hackeronecli.PageParams{Number: page, Size: pageSize})
			if err != nil {
				return err
			}
			return writeOutput(cmd,reviews)
		},
	}
	cmd.Flags().IntVar(&page, "page", 0, "Page number")
	cmd.Flags().IntVar(&pageSize, "page-size", 0, "Page size")
	return cmd
}

func newProgramsWeaknessesCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	var page, pageSize int

	cmd := &cobra.Command{
		Use:   "weaknesses <id>",
		Short: "List weaknesses for a program",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			weaknesses, err := client.ListWeaknesses(cmd.Context(), args[0], hackeronecli.PageParams{Number: page, Size: pageSize})
			if err != nil {
				return err
			}
			return writeOutput(cmd,weaknesses)
		},
	}
	cmd.Flags().IntVar(&page, "page", 0, "Page number")
	cmd.Flags().IntVar(&pageSize, "page-size", 0, "Page size")
	return cmd
}

func newProgramsNotifyCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	return &cobra.Command{
		Use:   "notify <id>",
		Short: "Notify external platform",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			result, err := client.NotifyExternalPlatform(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			return writeOutput(cmd,result)
		},
	}
}

func newProgramsMessagesCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	messagesCmd := &cobra.Command{
		Use:   "messages",
		Short: "Manage program messages",
	}

	var input hackeronecli.MessageInput
	sendCmd := &cobra.Command{
		Use:   "send <id>",
		Short: "Send a message",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			result, err := client.SendProgramMessage(cmd.Context(), args[0], input)
			if err != nil {
				return err
			}
			return writeOutput(cmd,result)
		},
	}
	sendCmd.Flags().StringVar(&input.RecipientHandle, "recipient-handle", "", "Recipient handle")
	sendCmd.Flags().StringVar(&input.Message, "message", "", "Message content")

	messagesCmd.AddCommand(sendCmd)
	return messagesCmd
}

func newProgramsBountiesCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	bountiesCmd := &cobra.Command{
		Use:   "bounties",
		Short: "Manage program bounties",
	}

	var input hackeronecli.ProgramBountyInput
	awardCmd := &cobra.Command{
		Use:   "award <id>",
		Short: "Award a bounty",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			result, err := client.AwardProgramBounty(cmd.Context(), args[0], input)
			if err != nil {
				return err
			}
			return writeOutput(cmd,result)
		},
	}
	awardCmd.Flags().StringVar(&input.ReportID, "report-id", "", "Report ID")
	awardCmd.Flags().Float64Var(&input.Amount, "amount", 0, "Bounty amount")
	awardCmd.Flags().Float64Var(&input.BonusAmount, "bonus-amount", 0, "Bonus amount")
	awardCmd.Flags().StringVar(&input.Message, "message", "", "Bounty message")

	bountiesCmd.AddCommand(awardCmd)
	return bountiesCmd
}

func newProgramsAllowedReportersCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	allowedReportersCmd := &cobra.Command{
		Use:   "allowed-reporters",
		Short: "Manage allowed reporters",
	}

	var listPage, listPageSize int
	listCmd := &cobra.Command{
		Use:   "list <id>",
		Short: "List allowed reporters",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			reporters, err := client.ListAllowedReporters(cmd.Context(), args[0], hackeronecli.PageParams{Number: listPage, Size: listPageSize})
			if err != nil {
				return err
			}
			return writeOutput(cmd,reporters)
		},
	}
	listCmd.Flags().IntVar(&listPage, "page", 0, "Page number")
	listCmd.Flags().IntVar(&listPageSize, "page-size", 0, "Page size")

	var histPage, histPageSize int
	historyCmd := &cobra.Command{
		Use:   "history <id>",
		Short: "Get allowed reporters history",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			history, err := client.GetAllowedReportersHistory(cmd.Context(), args[0], hackeronecli.PageParams{Number: histPage, Size: histPageSize})
			if err != nil {
				return err
			}
			return writeOutput(cmd,history)
		},
	}
	historyCmd.Flags().IntVar(&histPage, "page", 0, "Page number")
	historyCmd.Flags().IntVar(&histPageSize, "page-size", 0, "Page size")

	var actPage, actPageSize int
	activitiesCmd := &cobra.Command{
		Use:   "activities <id>",
		Short: "Get allowed reporter activities",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			activities, err := client.GetAllowedReporterActivities(cmd.Context(), args[0], hackeronecli.PageParams{Number: actPage, Size: actPageSize})
			if err != nil {
				return err
			}
			return writeOutput(cmd,activities)
		},
	}
	activitiesCmd.Flags().IntVar(&actPage, "page", 0, "Page number")
	activitiesCmd.Flags().IntVar(&actPageSize, "page-size", 0, "Page size")

	var unPage, unPageSize int
	usernameHistoryCmd := &cobra.Command{
		Use:   "username-history <id>",
		Short: "Get allowed reporter username history",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			history, err := client.GetAllowedReporterUsernameHistory(cmd.Context(), args[0], hackeronecli.PageParams{Number: unPage, Size: unPageSize})
			if err != nil {
				return err
			}
			return writeOutput(cmd,history)
		},
	}
	usernameHistoryCmd.Flags().IntVar(&unPage, "page", 0, "Page number")
	usernameHistoryCmd.Flags().IntVar(&unPageSize, "page-size", 0, "Page size")

	allowedReportersCmd.AddCommand(listCmd)
	allowedReportersCmd.AddCommand(historyCmd)
	allowedReportersCmd.AddCommand(activitiesCmd)
	allowedReportersCmd.AddCommand(usernameHistoryCmd)
	return allowedReportersCmd
}

func newProgramsCVERequestsCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	cveRequestsCmd := &cobra.Command{
		Use:   "cve-requests",
		Short: "Manage CVE requests",
	}

	var page, pageSize int
	listCmd := &cobra.Command{
		Use:   "list <id>",
		Short: "List CVE requests",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			requests, err := client.ListCVERequests(cmd.Context(), args[0], hackeronecli.PageParams{Number: page, Size: pageSize})
			if err != nil {
				return err
			}
			return writeOutput(cmd,requests)
		},
	}
	listCmd.Flags().IntVar(&page, "page", 0, "Page number")
	listCmd.Flags().IntVar(&pageSize, "page-size", 0, "Page size")

	var input hackeronecli.CreateCVERequestInput
	createCmd := &cobra.Command{
		Use:   "create <id>",
		Short: "Create a CVE request",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			request, err := client.CreateCVERequest(cmd.Context(), args[0], input)
			if err != nil {
				return err
			}
			return writeOutput(cmd,request)
		},
	}
	createCmd.Flags().StringVar(&input.ReportID, "report-id", "", "Report ID")
	createCmd.Flags().StringVar(&input.CveIdentifier, "cve-identifier", "", "CVE identifier")

	cveRequestsCmd.AddCommand(listCmd)
	cveRequestsCmd.AddCommand(createCmd)
	return cveRequestsCmd
}

func newProgramsHackerInvitationsCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	invitationsCmd := &cobra.Command{
		Use:   "hacker-invitations",
		Short: "Manage hacker invitations",
	}

	var page, pageSize int
	listCmd := &cobra.Command{
		Use:   "list <id>",
		Short: "List hacker invitations",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			invitations, err := client.ListHackerInvitations(cmd.Context(), args[0], hackeronecli.PageParams{Number: page, Size: pageSize})
			if err != nil {
				return err
			}
			return writeOutput(cmd,invitations)
		},
	}
	listCmd.Flags().IntVar(&page, "page", 0, "Page number")
	listCmd.Flags().IntVar(&pageSize, "page-size", 0, "Page size")

	var input hackeronecli.CreateHackerInvitationInput
	createCmd := &cobra.Command{
		Use:   "create <id>",
		Short: "Create a hacker invitation",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			invitation, err := client.CreateHackerInvitation(cmd.Context(), args[0], input)
			if err != nil {
				return err
			}
			return writeOutput(cmd,invitation)
		},
	}
	createCmd.Flags().StringVar(&input.Email, "email", "", "Hacker email")
	createCmd.Flags().StringVar(&input.Username, "username", "", "Hacker username")
	createCmd.Flags().StringVar(&input.Message, "message", "", "Invitation message")

	invitationsCmd.AddCommand(listCmd)
	invitationsCmd.AddCommand(createCmd)
	return invitationsCmd
}

func newProgramsPolicyCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	policyCmd := &cobra.Command{
		Use:   "policy",
		Short: "Manage program policy",
	}

	var policyText string
	updateCmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update program policy",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			result, err := client.UpdatePolicy(cmd.Context(), args[0], hackeronecli.PolicyInput{Policy: policyText})
			if err != nil {
				return err
			}
			return writeOutput(cmd,result)
		},
	}
	updateCmd.Flags().StringVar(&policyText, "policy", "", "Policy text")

	var filePath string
	attachCmd := &cobra.Command{
		Use:   "attach <id>",
		Short: "Attach file to program policy",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			result, err := client.AttachToPolicy(cmd.Context(), args[0], filePath)
			if err != nil {
				return err
			}
			return writeOutput(cmd,result)
		},
	}
	attachCmd.Flags().StringVar(&filePath, "file", "", "File path to attach")

	policyCmd.AddCommand(updateCmd)
	policyCmd.AddCommand(attachCmd)
	return policyCmd
}

func newProgramsScopesCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	scopesCmd := &cobra.Command{
		Use:   "scopes",
		Short: "Manage program scopes",
	}

	var listPage, listPageSize int
	listCmd := &cobra.Command{
		Use:   "list <id>",
		Short: "List structured scopes",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			scopes, err := client.ListScopes(cmd.Context(), args[0], hackeronecli.PageParams{Number: listPage, Size: listPageSize})
			if err != nil {
				return err
			}
			return writeOutput(cmd,scopes)
		},
	}
	listCmd.Flags().IntVar(&listPage, "page", 0, "Page number")
	listCmd.Flags().IntVar(&listPageSize, "page-size", 0, "Page size")

	var createInput hackeronecli.CreateScopeInput
	createCmd := &cobra.Command{
		Use:   "create <id>",
		Short: "Create a structured scope",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			scope, err := client.CreateScope(cmd.Context(), args[0], createInput)
			if err != nil {
				return err
			}
			return writeOutput(cmd,scope)
		},
	}
	createCmd.Flags().StringVar(&createInput.AssetIdentifier, "asset-identifier", "", "Asset identifier")
	createCmd.Flags().StringVar(&createInput.AssetType, "asset-type", "", "Asset type")
	createCmd.Flags().BoolVar(&createInput.EligibleForBounty, "eligible-for-bounty", false, "Eligible for bounty")
	createCmd.Flags().BoolVar(&createInput.EligibleForSubmission, "eligible-for-submission", false, "Eligible for submission")
	createCmd.Flags().StringVar(&createInput.Instruction, "instruction", "", "Instruction")

	var updateInput hackeronecli.UpdateScopeInput
	var updateBounty, updateSubmission bool
	updateCmd := &cobra.Command{
		Use:   "update <id> <scope-id>",
		Short: "Update a structured scope",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			if cmd.Flags().Changed("eligible-for-bounty") {
				updateInput.EligibleForBounty = &updateBounty
			}
			if cmd.Flags().Changed("eligible-for-submission") {
				updateInput.EligibleForSubmission = &updateSubmission
			}
			scope, err := client.UpdateProgramScope(cmd.Context(), args[0], args[1], updateInput)
			if err != nil {
				return err
			}
			return writeOutput(cmd,scope)
		},
	}
	updateCmd.Flags().BoolVar(&updateBounty, "eligible-for-bounty", false, "Eligible for bounty")
	updateCmd.Flags().BoolVar(&updateSubmission, "eligible-for-submission", false, "Eligible for submission")
	updateCmd.Flags().StringVar(&updateInput.Instruction, "instruction", "", "Instruction")

	archiveCmd := &cobra.Command{
		Use:   "archive <id> <scope-id>",
		Short: "Archive a structured scope",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			return client.ArchiveScope(cmd.Context(), args[0], args[1])
		},
	}

	scopesCmd.AddCommand(listCmd)
	scopesCmd.AddCommand(createCmd)
	scopesCmd.AddCommand(updateCmd)
	scopesCmd.AddCommand(archiveCmd)
	return scopesCmd
}

func newProgramsSwagCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	swagCmd := &cobra.Command{
		Use:   "swag",
		Short: "Manage awarded swag",
	}

	var page, pageSize int
	listCmd := &cobra.Command{
		Use:   "list <id>",
		Short: "List awarded swag",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			swag, err := client.ListAwardedSwag(cmd.Context(), args[0], hackeronecli.PageParams{Number: page, Size: pageSize})
			if err != nil {
				return err
			}
			return writeOutput(cmd,swag)
		},
	}
	listCmd.Flags().IntVar(&page, "page", 0, "Page number")
	listCmd.Flags().IntVar(&pageSize, "page-size", 0, "Page size")

	var sent bool
	updateCmd := &cobra.Command{
		Use:   "update <id> <swag-id>",
		Short: "Update awarded swag",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			swag, err := client.UpdateAwardedSwag(cmd.Context(), args[0], args[1], hackeronecli.UpdateSwagInput{Sent: sent})
			if err != nil {
				return err
			}
			return writeOutput(cmd,swag)
		},
	}
	updateCmd.Flags().BoolVar(&sent, "sent", false, "Mark swag as sent")

	swagCmd.AddCommand(listCmd)
	swagCmd.AddCommand(updateCmd)
	return swagCmd
}
