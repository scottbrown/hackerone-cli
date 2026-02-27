package cmd

import (
	hackeronecli "github.com/scottbrown/hackerone-cli"
	"github.com/spf13/cobra"
)

func NewAutomationsCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	automationsCmd := &cobra.Command{
		Use:   "automations",
		Short: "Manage automations",
	}

	automationsCmd.AddCommand(newAutomationsListCmd(clientFactory))
	automationsCmd.AddCommand(newAutomationsGetCmd(clientFactory))
	automationsCmd.AddCommand(newAutomationsTriggerCmd(clientFactory))
	automationsCmd.AddCommand(newAutomationsRunsCmd(clientFactory))

	return automationsCmd
}

func newAutomationsListCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	var page, pageSize int
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List automations",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			automations, err := client.ListAutomations(cmd.Context(), hackeronecli.PageParams{Number: page, Size: pageSize})
			if err != nil {
				return err
			}
			return writeOutput(cmd,automations)
		},
	}
	cmd.Flags().IntVar(&page, "page", 0, "Page number")
	cmd.Flags().IntVar(&pageSize, "page-size", 0, "Page size")
	return cmd
}

func newAutomationsGetCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get an automation by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			auto, err := client.GetAutomation(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			return writeOutput(cmd,auto)
		},
	}
}

func newAutomationsTriggerCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	return &cobra.Command{
		Use:   "trigger <id>",
		Short: "Trigger an automation",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			if err := client.TriggerAutomation(cmd.Context(), args[0]); err != nil {
				return err
			}
			return writeMessage(cmd, "Automation triggered successfully")
		},
	}
}

func newAutomationsRunsCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	runsCmd := &cobra.Command{
		Use:   "runs",
		Short: "Manage automation runs",
	}

	runsCmd.AddCommand(newAutomationsRunsListCmd(clientFactory))
	runsCmd.AddCommand(newAutomationsRunsGetCmd(clientFactory))
	runsCmd.AddCommand(newAutomationsRunsLogsCmd(clientFactory))

	return runsCmd
}

func newAutomationsRunsListCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	var page, pageSize int
	cmd := &cobra.Command{
		Use:   "list <automation-id>",
		Short: "List runs for an automation",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			runs, err := client.ListAutomationRuns(cmd.Context(), args[0], hackeronecli.PageParams{Number: page, Size: pageSize})
			if err != nil {
				return err
			}
			return writeOutput(cmd,runs)
		},
	}
	cmd.Flags().IntVar(&page, "page", 0, "Page number")
	cmd.Flags().IntVar(&pageSize, "page-size", 0, "Page size")
	return cmd
}

func newAutomationsRunsGetCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	return &cobra.Command{
		Use:   "get <automation-id> <run-id>",
		Short: "Get a specific automation run",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			run, err := client.GetAutomationRun(cmd.Context(), args[0], args[1])
			if err != nil {
				return err
			}
			return writeOutput(cmd,run)
		},
	}
}

func newAutomationsRunsLogsCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	return &cobra.Command{
		Use:   "logs <automation-id> <run-id>",
		Short: "Get logs for an automation run",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			logs, err := client.GetAutomationRunLogs(cmd.Context(), args[0], args[1])
			if err != nil {
				return err
			}
			return writeOutput(cmd,logs)
		},
	}
}
