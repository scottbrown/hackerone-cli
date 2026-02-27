package cmd

import (
	"context"

	hackeronecli "github.com/scottbrown/hackerone-cli"
	"github.com/spf13/cobra"
)

func NewAnalyticsCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	analyticsCmd := &cobra.Command{
		Use:   "analytics",
		Short: "Query HackerOne analytics",
	}

	analyticsCmd.AddCommand(newAnalyticsGetCmd(clientFactory))

	return analyticsCmd
}

func newAnalyticsGetCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	var program string
	var groups string
	var startDate string
	var endDate string

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get analytics for a program",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			result, err := client.GetAnalytics(context.Background(), program, groups, startDate, endDate)
			if err != nil {
				return err
			}
			return writeOutput(cmd,result)
		},
	}

	cmd.Flags().StringVar(&program, "program", "", "Program handle (required)")
	cmd.MarkFlagRequired("program")
	cmd.Flags().StringVar(&groups, "groups", "", "Analytics groups")
	cmd.Flags().StringVar(&startDate, "start-date", "", "Start date")
	cmd.Flags().StringVar(&endDate, "end-date", "", "End date")

	return cmd
}
