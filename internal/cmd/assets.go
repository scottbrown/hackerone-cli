package cmd

import (
	"fmt"
	"strings"

	hackeronecli "github.com/scottbrown/hackerone-cli"
	"github.com/spf13/cobra"
)

func NewAssetsCmd(clientFactory func() (*hackeronecli.Client, error)) *cobra.Command {
	var orgID string

	assetsCmd := &cobra.Command{
		Use:   "assets",
		Short: "Manage assets",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if orgID == "" {
				return fmt.Errorf("--org is required for asset commands")
			}
			return nil
		},
	}

	assetsCmd.PersistentFlags().StringVar(&orgID, "org", "", "Organization ID (required)")

	assetsCmd.AddCommand(newAssetsListCmd(clientFactory, &orgID))
	assetsCmd.AddCommand(newAssetsGetCmd(clientFactory, &orgID))
	assetsCmd.AddCommand(newAssetsCreateCmd(clientFactory, &orgID))
	assetsCmd.AddCommand(newAssetsUpdateCmd(clientFactory, &orgID))
	assetsCmd.AddCommand(newAssetsArchiveCmd(clientFactory, &orgID))
	assetsCmd.AddCommand(newAssetsImportCmd(clientFactory, &orgID))
	assetsCmd.AddCommand(newAssetsImportStatusCmd(clientFactory, &orgID))
	assetsCmd.AddCommand(newAssetsScreenshotCmd(clientFactory, &orgID))
	assetsCmd.AddCommand(newAssetsPortsCmd(clientFactory, &orgID))
	assetsCmd.AddCommand(newAssetsReachabilityCmd(clientFactory, &orgID))
	assetsCmd.AddCommand(newAssetsScannerCmd(clientFactory, &orgID))
	assetsCmd.AddCommand(newAssetsScopesCmd(clientFactory, &orgID))
	assetsCmd.AddCommand(newAssetsTagsCmd(clientFactory, &orgID))

	return assetsCmd
}

func newAssetsListCmd(clientFactory func() (*hackeronecli.Client, error), orgID *string) *cobra.Command {
	var page, pageSize int
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List assets",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			assets, err := client.ListAssets(cmd.Context(), *orgID, hackeronecli.PageParams{Number: page, Size: pageSize})
			if err != nil {
				return err
			}
			return writeOutput(cmd, assets)
		},
	}
	cmd.Flags().IntVar(&page, "page", 0, "Page number")
	cmd.Flags().IntVar(&pageSize, "page-size", 0, "Page size")
	return cmd
}

func newAssetsGetCmd(clientFactory func() (*hackeronecli.Client, error), orgID *string) *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get an asset by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			asset, err := client.GetAsset(cmd.Context(), *orgID, args[0])
			if err != nil {
				return err
			}
			return writeOutput(cmd, asset)
		},
	}
}

func newAssetsCreateCmd(clientFactory func() (*hackeronecli.Client, error), orgID *string) *cobra.Command {
	var input hackeronecli.CreateAssetInput
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create an asset",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			asset, err := client.CreateAsset(cmd.Context(), *orgID, input)
			if err != nil {
				return err
			}
			return writeOutput(cmd, asset)
		},
	}
	cmd.Flags().StringVar(&input.AssetType, "type", "", "Asset type")
	cmd.Flags().StringVar(&input.Identifier, "identifier", "", "Asset identifier")
	cmd.Flags().StringVar(&input.Description, "description", "", "Asset description")
	cmd.Flags().StringVar(&input.Coverage, "coverage", "", "Coverage level")
	cmd.Flags().StringVar(&input.MaxSeverity, "max-severity", "", "Maximum severity")
	return cmd
}

func newAssetsUpdateCmd(clientFactory func() (*hackeronecli.Client, error), orgID *string) *cobra.Command {
	var input hackeronecli.UpdateAssetInput
	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update an asset",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			asset, err := client.UpdateAsset(cmd.Context(), *orgID, args[0], input)
			if err != nil {
				return err
			}
			return writeOutput(cmd, asset)
		},
	}
	cmd.Flags().StringVar(&input.AssetType, "type", "", "Asset type")
	cmd.Flags().StringVar(&input.Identifier, "identifier", "", "Asset identifier")
	cmd.Flags().StringVar(&input.Description, "description", "", "Asset description")
	cmd.Flags().StringVar(&input.Coverage, "coverage", "", "Coverage level")
	cmd.Flags().StringVar(&input.MaxSeverity, "max-severity", "", "Maximum severity")
	return cmd
}

func newAssetsArchiveCmd(clientFactory func() (*hackeronecli.Client, error), orgID *string) *cobra.Command {
	var ids string
	cmd := &cobra.Command{
		Use:   "archive",
		Short: "Archive assets by IDs",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			idList := strings.Split(ids, ",")
			if err := client.ArchiveAssets(cmd.Context(), *orgID, idList); err != nil {
				return err
			}
			return writeMessage(cmd, "Assets archived successfully")
		},
	}
	cmd.Flags().StringVar(&ids, "ids", "", "Comma-separated list of asset IDs")
	return cmd
}

func newAssetsImportCmd(clientFactory func() (*hackeronecli.Client, error), orgID *string) *cobra.Command {
	return &cobra.Command{
		Use:   "import <file>",
		Short: "Import assets from a CSV file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			result, err := client.ImportAssets(cmd.Context(), *orgID, args[0])
			if err != nil {
				return err
			}
			return writeOutput(cmd, result)
		},
	}
}

func newAssetsImportStatusCmd(clientFactory func() (*hackeronecli.Client, error), orgID *string) *cobra.Command {
	return &cobra.Command{
		Use:   "import-status <id>",
		Short: "Get the status of an asset import",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			result, err := client.GetImportStatus(cmd.Context(), *orgID, args[0])
			if err != nil {
				return err
			}
			return writeOutput(cmd, result)
		},
	}
}

func newAssetsScreenshotCmd(clientFactory func() (*hackeronecli.Client, error), orgID *string) *cobra.Command {
	return &cobra.Command{
		Use:   "screenshot <id> <file>",
		Short: "Upload a screenshot for an asset",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			if err := client.UploadAssetScreenshot(cmd.Context(), *orgID, args[0], args[1]); err != nil {
				return err
			}
			return writeMessage(cmd, "Screenshot uploaded successfully")
		},
	}
}

func newAssetsPortsCmd(clientFactory func() (*hackeronecli.Client, error), orgID *string) *cobra.Command {
	portsCmd := &cobra.Command{
		Use:   "ports",
		Short: "Manage asset ports",
	}

	portsCmd.AddCommand(newAssetsPortsListCmd(clientFactory, orgID))
	portsCmd.AddCommand(newAssetsPortsCreateCmd(clientFactory, orgID))
	portsCmd.AddCommand(newAssetsPortsDeleteCmd(clientFactory, orgID))

	return portsCmd
}

func newAssetsPortsListCmd(clientFactory func() (*hackeronecli.Client, error), orgID *string) *cobra.Command {
	var page, pageSize int
	cmd := &cobra.Command{
		Use:   "list <asset-id>",
		Short: "List ports for an asset",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			ports, err := client.ListAssetPorts(cmd.Context(), *orgID, args[0], hackeronecli.PageParams{Number: page, Size: pageSize})
			if err != nil {
				return err
			}
			return writeOutput(cmd, ports)
		},
	}
	cmd.Flags().IntVar(&page, "page", 0, "Page number")
	cmd.Flags().IntVar(&pageSize, "page-size", 0, "Page size")
	return cmd
}

func newAssetsPortsCreateCmd(clientFactory func() (*hackeronecli.Client, error), orgID *string) *cobra.Command {
	var input hackeronecli.CreatePortInput
	cmd := &cobra.Command{
		Use:   "create <asset-id>",
		Short: "Create a port for an asset",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			port, err := client.CreateAssetPort(cmd.Context(), *orgID, args[0], input)
			if err != nil {
				return err
			}
			return writeOutput(cmd, port)
		},
	}
	cmd.Flags().IntVar(&input.Port, "port", 0, "Port number")
	cmd.Flags().StringVar(&input.Protocol, "protocol", "", "Protocol (tcp/udp)")
	cmd.Flags().StringVar(&input.Service, "service", "", "Service name")
	return cmd
}

func newAssetsPortsDeleteCmd(clientFactory func() (*hackeronecli.Client, error), orgID *string) *cobra.Command {
	return &cobra.Command{
		Use:   "delete <asset-id> <port-id>",
		Short: "Delete a port from an asset",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			if err := client.DeleteAssetPort(cmd.Context(), *orgID, args[0], args[1]); err != nil {
				return err
			}
			return writeMessage(cmd, "Port deleted successfully")
		},
	}
}

func newAssetsReachabilityCmd(clientFactory func() (*hackeronecli.Client, error), orgID *string) *cobra.Command {
	reachCmd := &cobra.Command{
		Use:   "reachability",
		Short: "Asset reachability operations",
	}

	reachCmd.AddCommand(newAssetsReachabilityStatusCmd(clientFactory, orgID))
	reachCmd.AddCommand(newAssetsReachabilityCheckCmd(clientFactory, orgID))

	return reachCmd
}

func newAssetsReachabilityStatusCmd(clientFactory func() (*hackeronecli.Client, error), orgID *string) *cobra.Command {
	return &cobra.Command{
		Use:   "status <asset-id>",
		Short: "Get reachability status for an asset",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			result, err := client.GetReachabilityStatus(cmd.Context(), *orgID, args[0])
			if err != nil {
				return err
			}
			return writeOutput(cmd, result)
		},
	}
}

func newAssetsReachabilityCheckCmd(clientFactory func() (*hackeronecli.Client, error), orgID *string) *cobra.Command {
	return &cobra.Command{
		Use:   "check <asset-id>",
		Short: "Check reachability for an asset",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			result, err := client.CheckReachability(cmd.Context(), *orgID, args[0])
			if err != nil {
				return err
			}
			return writeOutput(cmd, result)
		},
	}
}

func newAssetsScannerCmd(clientFactory func() (*hackeronecli.Client, error), orgID *string) *cobra.Command {
	scannerCmd := &cobra.Command{
		Use:   "scanner",
		Short: "Scanner configuration operations",
	}

	scannerCmd.AddCommand(newAssetsScannerGetCmd(clientFactory, orgID))
	scannerCmd.AddCommand(newAssetsScannerUpdateCmd(clientFactory, orgID))

	return scannerCmd
}

func newAssetsScannerGetCmd(clientFactory func() (*hackeronecli.Client, error), orgID *string) *cobra.Command {
	return &cobra.Command{
		Use:   "get <asset-id>",
		Short: "Get scanner configuration for an asset",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			cfg, err := client.GetScannerConfig(cmd.Context(), *orgID, args[0])
			if err != nil {
				return err
			}
			return writeOutput(cmd, cfg)
		},
	}
}

func newAssetsScannerUpdateCmd(clientFactory func() (*hackeronecli.Client, error), orgID *string) *cobra.Command {
	var enabled bool
	cmd := &cobra.Command{
		Use:   "update <asset-id>",
		Short: "Update scanner configuration for an asset",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			cfg, err := client.UpdateScannerConfig(cmd.Context(), *orgID, args[0], hackeronecli.ScannerConfiguration{Enabled: enabled})
			if err != nil {
				return err
			}
			return writeOutput(cmd, cfg)
		},
	}
	cmd.Flags().BoolVar(&enabled, "enabled", false, "Enable or disable the scanner")
	return cmd
}

func newAssetsScopesCmd(clientFactory func() (*hackeronecli.Client, error), orgID *string) *cobra.Command {
	scopesCmd := &cobra.Command{
		Use:   "scopes",
		Short: "Manage asset scopes",
	}

	scopesCmd.AddCommand(newAssetsScopesAddCmd(clientFactory, orgID))
	scopesCmd.AddCommand(newAssetsScopesUpdateCmd(clientFactory, orgID))
	scopesCmd.AddCommand(newAssetsScopesArchiveCmd(clientFactory, orgID))

	return scopesCmd
}

func newAssetsScopesAddCmd(clientFactory func() (*hackeronecli.Client, error), orgID *string) *cobra.Command {
	var input hackeronecli.AssetScope
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add an asset scope",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			if err := client.AddAssetScope(cmd.Context(), *orgID, input); err != nil {
				return err
			}
			return writeMessage(cmd, "Asset scope added successfully")
		},
	}
	cmd.Flags().StringVar(&input.AssetID, "asset-id", "", "Asset ID")
	cmd.Flags().StringVar(&input.ProgramID, "program-id", "", "Program ID")
	cmd.Flags().BoolVar(&input.Eligible, "eligible", false, "Whether the asset is eligible")
	return cmd
}

func newAssetsScopesUpdateCmd(clientFactory func() (*hackeronecli.Client, error), orgID *string) *cobra.Command {
	var input hackeronecli.AssetScope
	cmd := &cobra.Command{
		Use:   "update <asset-id>",
		Short: "Update an asset scope",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			if err := client.UpdateAssetScope(cmd.Context(), *orgID, args[0], input); err != nil {
				return err
			}
			return writeMessage(cmd, "Asset scope updated successfully")
		},
	}
	cmd.Flags().StringVar(&input.ProgramID, "program-id", "", "Program ID")
	cmd.Flags().BoolVar(&input.Eligible, "eligible", false, "Whether the asset is eligible")
	return cmd
}

func newAssetsScopesArchiveCmd(clientFactory func() (*hackeronecli.Client, error), orgID *string) *cobra.Command {
	return &cobra.Command{
		Use:   "archive <asset-id>",
		Short: "Archive scopes for an asset",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			if err := client.ArchiveAssetScopes(cmd.Context(), *orgID, args[0]); err != nil {
				return err
			}
			return writeMessage(cmd, "Asset scopes archived successfully")
		},
	}
}

func newAssetsTagsCmd(clientFactory func() (*hackeronecli.Client, error), orgID *string) *cobra.Command {
	tagsCmd := &cobra.Command{
		Use:   "tags",
		Short: "Manage asset tags",
	}

	tagsCmd.AddCommand(newAssetsTagsListCmd(clientFactory, orgID))
	tagsCmd.AddCommand(newAssetsTagsCreateCmd(clientFactory, orgID))
	tagsCmd.AddCommand(newAssetsTagsCategoriesCmd(clientFactory, orgID))
	tagsCmd.AddCommand(newAssetsTagsCreateCategoryCmd(clientFactory, orgID))

	return tagsCmd
}

func newAssetsTagsListCmd(clientFactory func() (*hackeronecli.Client, error), orgID *string) *cobra.Command {
	var page, pageSize int
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List asset tags",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			tags, err := client.ListAssetTags(cmd.Context(), *orgID, hackeronecli.PageParams{Number: page, Size: pageSize})
			if err != nil {
				return err
			}
			return writeOutput(cmd, tags)
		},
	}
	cmd.Flags().IntVar(&page, "page", 0, "Page number")
	cmd.Flags().IntVar(&pageSize, "page-size", 0, "Page size")
	return cmd
}

func newAssetsTagsCreateCmd(clientFactory func() (*hackeronecli.Client, error), orgID *string) *cobra.Command {
	var name, categoryID string
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create an asset tag",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			tag, err := client.CreateAssetTag(cmd.Context(), *orgID, hackeronecli.AssetTag{
				Attributes: hackeronecli.AssetTagAttributes{Name: name, CategoryID: categoryID},
			})
			if err != nil {
				return err
			}
			return writeOutput(cmd, tag)
		},
	}
	cmd.Flags().StringVar(&name, "name", "", "Tag name")
	cmd.Flags().StringVar(&categoryID, "category-id", "", "Category ID")
	return cmd
}

func newAssetsTagsCategoriesCmd(clientFactory func() (*hackeronecli.Client, error), orgID *string) *cobra.Command {
	var page, pageSize int
	cmd := &cobra.Command{
		Use:   "categories",
		Short: "List asset tag categories",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			cats, err := client.ListAssetTagCategories(cmd.Context(), *orgID, hackeronecli.PageParams{Number: page, Size: pageSize})
			if err != nil {
				return err
			}
			return writeOutput(cmd, cats)
		},
	}
	cmd.Flags().IntVar(&page, "page", 0, "Page number")
	cmd.Flags().IntVar(&pageSize, "page-size", 0, "Page size")
	return cmd
}

func newAssetsTagsCreateCategoryCmd(clientFactory func() (*hackeronecli.Client, error), orgID *string) *cobra.Command {
	var name string
	cmd := &cobra.Command{
		Use:   "create-category",
		Short: "Create an asset tag category",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientFactory()
			if err != nil {
				return err
			}
			cat, err := client.CreateAssetTagCategory(cmd.Context(), *orgID, hackeronecli.AssetTagCategory{
				Attributes: hackeronecli.AssetTagCategoryAttributes{Name: name},
			})
			if err != nil {
				return err
			}
			return writeOutput(cmd, cat)
		},
	}
	cmd.Flags().StringVar(&name, "name", "", "Category name")
	return cmd
}
