package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

var dbCmd = &cobra.Command{
	Use:   "databases",
	Short: "Database commands",
}

var dbListCmd = &cobra.Command{
	Use:   "list",
	Short: "List databases",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, _ := getClient()
		data := must(client.Get("/database", nil))

		var items []map[string]interface{}
		var result map[string]json.RawMessage
		if err := json.Unmarshal(data, &result); err == nil {
			if d, ok := result["data"]; ok {
				json.Unmarshal(d, &items)
			}
		}
		if items == nil {
			json.Unmarshal(data, &items)
		}

		if jsonFlag {
			printOut(mustMarshal(items), true)
		} else {
			for _, db := range items {
				fmt.Printf("  %3v  %s  (%s)\n", db["id"], db["name"], db["engine"])
			}
		}
		return nil
	},
}

var dbGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get database details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, _ := getClient()
		includeTables, _ := cmd.Flags().GetBool("include-tables")
		includeFields, _ := cmd.Flags().GetBool("include-fields")

		params := map[string]string{}
		if includeFields {
			params["include"] = "tables.fields"
		} else if includeTables {
			params["include"] = "tables"
		}

		data := must(client.Get("/database/"+args[0], params))
		printOut(data, jsonFlag)
		return nil
	},
}

var dbMetadataCmd = &cobra.Command{
	Use:   "metadata <id>",
	Short: "Get database metadata",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, _ := getClient()
		data := must(client.Get("/database/"+args[0]+"/metadata", nil))
		printOut(data, jsonFlag)
		return nil
	},
}

var dbSchemasCmd = &cobra.Command{
	Use:   "schemas <id>",
	Short: "List database schemas",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, _ := getClient()
		data := must(client.Get("/database/"+args[0]+"/schemas", nil))
		printOut(data, jsonFlag)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(dbCmd)
	dbCmd.AddCommand(dbListCmd, dbGetCmd, dbMetadataCmd, dbSchemasCmd)

	dbListCmd.Flags().Bool("include-tables", false, "Include tables")
	dbGetCmd.Flags().Bool("include-tables", false, "Include tables")
	dbGetCmd.Flags().Bool("include-fields", false, "Include fields (implies --include-tables)")
	dbMetadataCmd.Flags().Bool("include-hidden", false, "Include hidden fields")
}
