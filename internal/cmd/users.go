package cmd

import (
	"context"

	hackeronecli "github.com/scottbrown/hackerone-cli"
	"github.com/spf13/cobra"
)

func NewUsersCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	usersCmd := &cobra.Command{
		Use:   "users",
		Short: "Manage HackerOne users",
	}

	usersCmd.AddCommand(newUsersGetCmd(clientFactory))
	usersCmd.AddCommand(newUsersGetByIDCmd(clientFactory))

	return usersCmd
}

func newUsersGetCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	return &cobra.Command{
		Use:   "get <username>",
		Short: "Get a user by username",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			user, err := client.GetUser(context.Background(), args[0])
			if err != nil {
				return err
			}
			return writeOutput(cmd,user)
		},
	}
}

func newUsersGetByIDCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	return &cobra.Command{
		Use:   "get-by-id <id>",
		Short: "Get a user by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			user, err := client.GetUserByID(context.Background(), args[0])
			if err != nil {
				return err
			}
			return writeOutput(cmd,user)
		},
	}
}
