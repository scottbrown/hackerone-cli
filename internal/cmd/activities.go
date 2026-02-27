package cmd

import (
	"context"

	hackeronecli "github.com/scottbrown/hackerone-cli"
	"github.com/spf13/cobra"
)

func NewActivitiesCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	activitiesCmd := &cobra.Command{
		Use:   "activities",
		Short: "Manage HackerOne activities",
	}

	activitiesCmd.AddCommand(newActivitiesGetCmd(clientFactory))
	activitiesCmd.AddCommand(newActivitiesListCmd(clientFactory))

	return activitiesCmd
}

func newActivitiesGetCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get an activity by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			activity, err := client.GetActivity(context.Background(), args[0])
			if err != nil {
				return err
			}
			return writeOutput(cmd,activity)
		},
	}
}

func newActivitiesListCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	var page int
	var pageSize int
	var updatedAtAfter string
	var updatedAtBefore string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List activities",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			params := hackeronecli.PageParams{Number: page, Size: pageSize}
			activities, err := client.ListActivities(context.Background(), params, updatedAtAfter, updatedAtBefore)
			if err != nil {
				return err
			}
			return writeOutput(cmd,activities)
		},
	}

	cmd.Flags().IntVar(&page, "page", 0, "Page number")
	cmd.Flags().IntVar(&pageSize, "page-size", 0, "Page size")
	cmd.Flags().StringVar(&updatedAtAfter, "updated-at-after", "", "Filter activities updated after this timestamp")
	cmd.Flags().StringVar(&updatedAtBefore, "updated-at-before", "", "Filter activities updated before this timestamp")

	return cmd
}
