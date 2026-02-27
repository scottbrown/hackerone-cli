package cmd

import (
	hackeronecli "github.com/scottbrown/hackerone-cli"
	"github.com/spf13/cobra"
)

func NewCredentialsCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	credCmd := &cobra.Command{
		Use:   "credentials",
		Short: "Manage credentials",
	}

	credCmd.AddCommand(newCredentialsListCmd(clientFactory))
	credCmd.AddCommand(newCredentialsCreateCmd(clientFactory))
	credCmd.AddCommand(newCredentialsUpdateCmd(clientFactory))
	credCmd.AddCommand(newCredentialsDeleteCmd(clientFactory))
	credCmd.AddCommand(newCredentialsAssignCmd(clientFactory))
	credCmd.AddCommand(newCredentialsRevokeCmd(clientFactory))
	credCmd.AddCommand(newCredentialsInquiriesCmd(clientFactory))

	return credCmd
}

func newCredentialsListCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	var page, pageSize int

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List credentials",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			creds, err := client.ListCredentials(cmd.Context(), hackeronecli.PageParams{Number: page, Size: pageSize})
			if err != nil {
				return err
			}
			return writeOutput(cmd,creds)
		},
	}

	cmd.Flags().IntVar(&page, "page", 0, "Page number")
	cmd.Flags().IntVar(&pageSize, "page-size", 0, "Page size")

	return cmd
}

func newCredentialsCreateCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	var accountName, credType, credentials, programID string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a credential",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			cred, err := client.CreateCredential(cmd.Context(), hackeronecli.CreateCredentialInput{
				AccountName: accountName,
				CredType:    credType,
				Credentials: credentials,
				ProgramID:   programID,
			})
			if err != nil {
				return err
			}
			return writeOutput(cmd,cred)
		},
	}

	cmd.Flags().StringVar(&accountName, "account-name", "", "Account name")
	cmd.Flags().StringVar(&credType, "credential-type", "", "Credential type")
	cmd.Flags().StringVar(&credentials, "credentials", "", "Credentials value")
	cmd.Flags().StringVar(&programID, "program-id", "", "Program ID")

	return cmd
}

func newCredentialsUpdateCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	var accountName, credentials string

	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update a credential",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			cred, err := client.UpdateCredential(cmd.Context(), args[0], hackeronecli.UpdateCredentialInput{
				AccountName: accountName,
				Credentials: credentials,
			})
			if err != nil {
				return err
			}
			return writeOutput(cmd,cred)
		},
	}

	cmd.Flags().StringVar(&accountName, "account-name", "", "Account name")
	cmd.Flags().StringVar(&credentials, "credentials", "", "Credentials value")

	return cmd
}

func newCredentialsDeleteCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	return &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a credential",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			if err := client.DeleteCredential(cmd.Context(), args[0]); err != nil {
				return err
			}
			return writeMessage(cmd, "Credential deleted.")
		},
	}
}

func newCredentialsAssignCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	return &cobra.Command{
		Use:   "assign <id>",
		Short: "Assign a credential",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			cred, err := client.AssignCredential(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			return writeOutput(cmd,cred)
		},
	}
}

func newCredentialsRevokeCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	return &cobra.Command{
		Use:   "revoke <id>",
		Short: "Revoke a credential",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			cred, err := client.RevokeCredential(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			return writeOutput(cmd,cred)
		},
	}
}

func newCredentialsInquiriesCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	inqCmd := &cobra.Command{
		Use:   "inquiries",
		Short: "Manage credential inquiries",
	}

	inqCmd.AddCommand(newCredInquiriesListCmd(clientFactory))
	inqCmd.AddCommand(newCredInquiriesCreateCmd(clientFactory))
	inqCmd.AddCommand(newCredInquiriesUpdateCmd(clientFactory))
	inqCmd.AddCommand(newCredInquiriesDeleteCmd(clientFactory))
	inqCmd.AddCommand(newCredInquiriesResponsesCmd(clientFactory))

	return inqCmd
}

func newCredInquiriesListCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	var page, pageSize int

	cmd := &cobra.Command{
		Use:   "list <program-id>",
		Short: "List credential inquiries for a program",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			inquiries, err := client.ListCredentialInquiries(cmd.Context(), args[0], hackeronecli.PageParams{Number: page, Size: pageSize})
			if err != nil {
				return err
			}
			return writeOutput(cmd,inquiries)
		},
	}

	cmd.Flags().IntVar(&page, "page", 0, "Page number")
	cmd.Flags().IntVar(&pageSize, "page-size", 0, "Page size")

	return cmd
}

func newCredInquiriesCreateCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	var question, fieldType string
	var required bool

	cmd := &cobra.Command{
		Use:   "create <program-id>",
		Short: "Create a credential inquiry",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			inq, err := client.CreateCredentialInquiry(cmd.Context(), args[0], hackeronecli.CreateCredentialInquiryInput{
				Question:  question,
				Required:  required,
				FieldType: fieldType,
			})
			if err != nil {
				return err
			}
			return writeOutput(cmd,inq)
		},
	}

	cmd.Flags().StringVar(&question, "question", "", "Inquiry question")
	cmd.Flags().BoolVar(&required, "required", false, "Whether the inquiry is required")
	cmd.Flags().StringVar(&fieldType, "field-type", "", "Field type")

	return cmd
}

func newCredInquiriesUpdateCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	var question, fieldType string
	var required bool

	cmd := &cobra.Command{
		Use:   "update <program-id> <inquiry-id>",
		Short: "Update a credential inquiry",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			inq, err := client.UpdateCredentialInquiry(cmd.Context(), args[0], args[1], hackeronecli.CreateCredentialInquiryInput{
				Question:  question,
				Required:  required,
				FieldType: fieldType,
			})
			if err != nil {
				return err
			}
			return writeOutput(cmd,inq)
		},
	}

	cmd.Flags().StringVar(&question, "question", "", "Inquiry question")
	cmd.Flags().BoolVar(&required, "required", false, "Whether the inquiry is required")
	cmd.Flags().StringVar(&fieldType, "field-type", "", "Field type")

	return cmd
}

func newCredInquiriesDeleteCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	return &cobra.Command{
		Use:   "delete <program-id> <inquiry-id>",
		Short: "Delete a credential inquiry",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			if err := client.DeleteCredentialInquiry(cmd.Context(), args[0], args[1]); err != nil {
				return err
			}
			return writeMessage(cmd, "Credential inquiry deleted.")
		},
	}
}

func newCredInquiriesResponsesCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	respCmd := &cobra.Command{
		Use:   "responses",
		Short: "Manage credential inquiry responses",
	}

	respCmd.AddCommand(newCredInqResponsesListCmd(clientFactory))
	respCmd.AddCommand(newCredInqResponsesDeleteCmd(clientFactory))

	return respCmd
}

func newCredInqResponsesListCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	var page, pageSize int

	cmd := &cobra.Command{
		Use:   "list <program-id> <inquiry-id>",
		Short: "List credential inquiry responses",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			responses, err := client.ListCredentialInquiryResponses(cmd.Context(), args[0], args[1], hackeronecli.PageParams{Number: page, Size: pageSize})
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

func newCredInqResponsesDeleteCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	return &cobra.Command{
		Use:   "delete <program-id> <inquiry-id> <response-id>",
		Short: "Delete a credential inquiry response",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			if err := client.DeleteCredentialInquiryResponse(cmd.Context(), args[0], args[1], args[2]); err != nil {
				return err
			}
			return writeMessage(cmd, "Credential inquiry response deleted.")
		},
	}
}
