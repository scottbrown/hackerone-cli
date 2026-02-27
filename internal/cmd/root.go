package cmd

import (
	"fmt"
	"os"

	hackeronecli "github.com/scottbrown/hackerone-cli"
	"github.com/spf13/cobra"
)

func writeOutput(cmd *cobra.Command, data interface{}) error {
	format, _ := cmd.Root().PersistentFlags().GetString("format")
	return hackeronecli.FormatOutput(cmd.OutOrStdout(), format, data)
}

func writeMessage(cmd *cobra.Command, msg string) error {
	format, _ := cmd.Root().PersistentFlags().GetString("format")
	return hackeronecli.FormatMessage(cmd.OutOrStdout(), format, msg)
}

func NewRootCmd(version string) *cobra.Command {
	var apiIdentifier string
	var apiToken string

	rootCmd := &cobra.Command{
		Use:   "h1",
		Short: "HackerOne CLI - interact with the HackerOne API",
		Long:  "A CLI tool for the HackerOne Customer Resources REST API.",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	rootCmd.PersistentFlags().StringVar(&apiIdentifier, "api-identifier", "", "API identifier (overrides HACKERONE_API_IDENTIFIER)")
	rootCmd.PersistentFlags().StringVar(&apiToken, "api-token", "", "API token (overrides HACKERONE_API_TOKEN)")
	rootCmd.PersistentFlags().String("format", "json", "Output format: json, text, or markdown")

	clientFactory := func() (*hackeronecli.Client, error) {
		id := apiIdentifier
		if id == "" {
			id = os.Getenv("HACKERONE_API_IDENTIFIER")
		}
		tok := apiToken
		if tok == "" {
			tok = os.Getenv("HACKERONE_API_TOKEN")
		}
		if id == "" || tok == "" {
			return nil, fmt.Errorf("HACKERONE_API_IDENTIFIER and HACKERONE_API_TOKEN must be set (via environment or --api-identifier/--api-token flags)")
		}
		client := hackeronecli.NewClient(id, tok)
		return client, nil
	}

	rootCmd.AddCommand(NewVersionCmd(version))
	rootCmd.AddCommand(NewActivitiesCmd(clientFactory))
	rootCmd.AddCommand(NewAnalyticsCmd(clientFactory))
	rootCmd.AddCommand(NewAssetsCmd(clientFactory))
	rootCmd.AddCommand(NewAutomationsCmd(clientFactory))
	rootCmd.AddCommand(NewCredentialsCmd(clientFactory))
	rootCmd.AddCommand(NewEmailCmd(clientFactory))
	rootCmd.AddCommand(NewOrganizationsCmd(clientFactory))
	rootCmd.AddCommand(NewProgramsCmd(clientFactory))
	rootCmd.AddCommand(NewReportsCmd(clientFactory))
	rootCmd.AddCommand(NewUsersCmd(clientFactory))

	return rootCmd
}
