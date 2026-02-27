package cmd

import (
	"context"
	"fmt"
	"os"

	hackeronecli "github.com/scottbrown/hackerone-cli"
	"github.com/spf13/cobra"
)

func NewEmailCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	emailCmd := &cobra.Command{
		Use:   "email",
		Short: "Send email via HackerOne",
	}

	emailCmd.AddCommand(newEmailSendCmd(clientFactory))

	return emailCmd
}

func newEmailSendCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	var to string
	var subject string
	var body string
	var bodyFile string

	cmd := &cobra.Command{
		Use:   "send",
		Short: "Send an email",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}

			emailBody := body
			if bodyFile != "" {
				content, err := os.ReadFile(bodyFile)
				if err != nil {
					return fmt.Errorf("reading body file: %w", err)
				}
				emailBody = string(content)
			}

			input := hackeronecli.SendEmailInput{
				To:      to,
				Subject: subject,
				Body:    emailBody,
			}
			if err := client.SendEmail(context.Background(), input); err != nil {
				return err
			}
			return writeMessage(cmd, "Email sent successfully.")
		},
	}

	cmd.Flags().StringVar(&to, "to", "", "Recipient email address (required)")
	cmd.MarkFlagRequired("to")
	cmd.Flags().StringVar(&subject, "subject", "", "Email subject (required)")
	cmd.MarkFlagRequired("subject")
	cmd.Flags().StringVar(&body, "body", "", "Email body")
	cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to a file containing the email body")

	return cmd
}
