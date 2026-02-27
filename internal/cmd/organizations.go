package cmd

import (
	hackeronecli "github.com/scottbrown/hackerone-cli"
	"github.com/spf13/cobra"
)

func NewOrganizationsCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	orgCmd := &cobra.Command{
		Use:   "organizations",
		Short: "Manage organizations",
	}

	orgCmd.AddCommand(newOrgsListCmd(clientFactory))
	orgCmd.AddCommand(newOrgsAuditLogCmd(clientFactory))
	orgCmd.AddCommand(newOrgsProgramsCmd(clientFactory))
	orgCmd.AddCommand(newOrgsInboxesCmd(clientFactory))
	orgCmd.AddCommand(newOrgsEligibilityCmd(clientFactory))
	orgCmd.AddCommand(newOrgsInvitationsCmd(clientFactory))
	orgCmd.AddCommand(newOrgsGroupsCmd(clientFactory))
	orgCmd.AddCommand(newOrgsMembersCmd(clientFactory))

	return orgCmd
}

func newOrgsListCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	var page, pageSize int

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List organizations",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			orgs, err := client.ListOrganizations(cmd.Context(), hackeronecli.PageParams{Number: page, Size: pageSize})
			if err != nil {
				return err
			}
			return writeOutput(cmd,orgs)
		},
	}

	cmd.Flags().IntVar(&page, "page", 0, "Page number")
	cmd.Flags().IntVar(&pageSize, "page-size", 0, "Page size")

	return cmd
}

func newOrgsAuditLogCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	var page, pageSize int

	cmd := &cobra.Command{
		Use:   "audit-log <id>",
		Short: "Get audit log for an organization",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			entries, err := client.GetOrganizationAuditLog(cmd.Context(), args[0], hackeronecli.PageParams{Number: page, Size: pageSize})
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

func newOrgsProgramsCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	var page, pageSize int

	cmd := &cobra.Command{
		Use:   "programs <id>",
		Short: "List programs for an organization",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			programs, err := client.ListOrganizationPrograms(cmd.Context(), args[0], hackeronecli.PageParams{Number: page, Size: pageSize})
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

func newOrgsInboxesCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	var page, pageSize int

	cmd := &cobra.Command{
		Use:   "inboxes <id>",
		Short: "List inboxes for an organization",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			inboxes, err := client.ListOrganizationInboxes(cmd.Context(), args[0], hackeronecli.PageParams{Number: page, Size: pageSize})
			if err != nil {
				return err
			}
			return writeOutput(cmd,inboxes)
		},
	}

	cmd.Flags().IntVar(&page, "page", 0, "Page number")
	cmd.Flags().IntVar(&pageSize, "page-size", 0, "Page size")

	return cmd
}

func newOrgsEligibilityCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	eligCmd := &cobra.Command{
		Use:   "eligibility",
		Short: "Manage eligibility settings",
	}

	eligCmd.AddCommand(newOrgsEligibilityListCmd(clientFactory))
	eligCmd.AddCommand(newOrgsEligibilityGetCmd(clientFactory))

	return eligCmd
}

func newOrgsEligibilityListCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	var page, pageSize int

	cmd := &cobra.Command{
		Use:   "list <org-id>",
		Short: "List eligibility settings",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			settings, err := client.ListEligibilitySettings(cmd.Context(), args[0], hackeronecli.PageParams{Number: page, Size: pageSize})
			if err != nil {
				return err
			}
			return writeOutput(cmd,settings)
		},
	}

	cmd.Flags().IntVar(&page, "page", 0, "Page number")
	cmd.Flags().IntVar(&pageSize, "page-size", 0, "Page size")

	return cmd
}

func newOrgsEligibilityGetCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	return &cobra.Command{
		Use:   "get <org-id> <setting-id>",
		Short: "Get an eligibility setting",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			setting, err := client.GetEligibilitySetting(cmd.Context(), args[0], args[1])
			if err != nil {
				return err
			}
			return writeOutput(cmd,setting)
		},
	}
}

func newOrgsInvitationsCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	invCmd := &cobra.Command{
		Use:   "invitations",
		Short: "Manage organization invitations",
	}

	invCmd.AddCommand(newOrgsInvitationsListCmd(clientFactory))
	invCmd.AddCommand(newOrgsInvitationsCreateCmd(clientFactory))

	return invCmd
}

func newOrgsInvitationsListCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	var page, pageSize int

	cmd := &cobra.Command{
		Use:   "list <org-id>",
		Short: "List invitations",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			invs, err := client.ListInvitations(cmd.Context(), args[0], hackeronecli.PageParams{Number: page, Size: pageSize})
			if err != nil {
				return err
			}
			return writeOutput(cmd,invs)
		},
	}

	cmd.Flags().IntVar(&page, "page", 0, "Page number")
	cmd.Flags().IntVar(&pageSize, "page-size", 0, "Page size")

	return cmd
}

func newOrgsInvitationsCreateCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	var email string
	var groupIDs []string

	cmd := &cobra.Command{
		Use:   "create <org-id>",
		Short: "Create an invitation",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			inv, err := client.CreateInvitation(cmd.Context(), args[0], hackeronecli.CreateInvitationInput{
				Email:    email,
				GroupIDs: groupIDs,
			})
			if err != nil {
				return err
			}
			return writeOutput(cmd,inv)
		},
	}

	cmd.Flags().StringVar(&email, "email", "", "Email address to invite")
	cmd.Flags().StringSliceVar(&groupIDs, "group-ids", nil, "Group IDs to assign")

	return cmd
}

func newOrgsGroupsCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	grpCmd := &cobra.Command{
		Use:   "groups",
		Short: "Manage organization groups",
	}

	grpCmd.AddCommand(newOrgsGroupsListCmd(clientFactory))
	grpCmd.AddCommand(newOrgsGroupsGetCmd(clientFactory))
	grpCmd.AddCommand(newOrgsGroupsCreateCmd(clientFactory))
	grpCmd.AddCommand(newOrgsGroupsUpdateCmd(clientFactory))

	return grpCmd
}

func newOrgsGroupsListCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	var page, pageSize int

	cmd := &cobra.Command{
		Use:   "list <org-id>",
		Short: "List groups",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			groups, err := client.ListGroups(cmd.Context(), args[0], hackeronecli.PageParams{Number: page, Size: pageSize})
			if err != nil {
				return err
			}
			return writeOutput(cmd,groups)
		},
	}

	cmd.Flags().IntVar(&page, "page", 0, "Page number")
	cmd.Flags().IntVar(&pageSize, "page-size", 0, "Page size")

	return cmd
}

func newOrgsGroupsGetCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	return &cobra.Command{
		Use:   "get <org-id> <group-id>",
		Short: "Get a group",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			group, err := client.GetGroup(cmd.Context(), args[0], args[1])
			if err != nil {
				return err
			}
			return writeOutput(cmd,group)
		},
	}
}

func newOrgsGroupsCreateCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	var name string
	var permissions []string

	cmd := &cobra.Command{
		Use:   "create <org-id>",
		Short: "Create a group",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			group, err := client.CreateGroup(cmd.Context(), args[0], hackeronecli.CreateGroupInput{
				Name:        name,
				Permissions: permissions,
			})
			if err != nil {
				return err
			}
			return writeOutput(cmd,group)
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "Group name")
	cmd.Flags().StringSliceVar(&permissions, "permissions", nil, "Permissions")

	return cmd
}

func newOrgsGroupsUpdateCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	var name string
	var permissions []string

	cmd := &cobra.Command{
		Use:   "update <org-id> <group-id>",
		Short: "Update a group",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			group, err := client.UpdateGroup(cmd.Context(), args[0], args[1], hackeronecli.UpdateGroupInput{
				Name:        name,
				Permissions: permissions,
			})
			if err != nil {
				return err
			}
			return writeOutput(cmd,group)
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "Group name")
	cmd.Flags().StringSliceVar(&permissions, "permissions", nil, "Permissions")

	return cmd
}

func newOrgsMembersCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	memCmd := &cobra.Command{
		Use:   "members",
		Short: "Manage organization members",
	}

	memCmd.AddCommand(newOrgsMembersListCmd(clientFactory))
	memCmd.AddCommand(newOrgsMembersGetCmd(clientFactory))
	memCmd.AddCommand(newOrgsMembersUpdateCmd(clientFactory))
	memCmd.AddCommand(newOrgsMembersRemoveCmd(clientFactory))

	return memCmd
}

func newOrgsMembersListCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	var page, pageSize int

	cmd := &cobra.Command{
		Use:   "list <org-id>",
		Short: "List members",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			members, err := client.ListMembers(cmd.Context(), args[0], hackeronecli.PageParams{Number: page, Size: pageSize})
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

func newOrgsMembersGetCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	return &cobra.Command{
		Use:   "get <org-id> <member-id>",
		Short: "Get a member",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			member, err := client.GetMember(cmd.Context(), args[0], args[1])
			if err != nil {
				return err
			}
			return writeOutput(cmd,member)
		},
	}
}

func newOrgsMembersUpdateCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	var groupIDs []string

	cmd := &cobra.Command{
		Use:   "update <org-id> <member-id>",
		Short: "Update a member",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			member, err := client.UpdateMember(cmd.Context(), args[0], args[1], hackeronecli.UpdateMemberInput{
				GroupIDs: groupIDs,
			})
			if err != nil {
				return err
			}
			return writeOutput(cmd,member)
		},
	}

	cmd.Flags().StringSliceVar(&groupIDs, "group-ids", nil, "Group IDs to assign")

	return cmd
}

func newOrgsMembersRemoveCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	return &cobra.Command{
		Use:   "remove <org-id> <member-id>",
		Short: "Remove a member from the organization",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			if err := client.RemoveMember(cmd.Context(), args[0], args[1]); err != nil {
				return err
			}
			return writeMessage(cmd, "Member removed.")
		},
	}
}
