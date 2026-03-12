package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var collCmd = &cobra.Command{
	Use:   "collections",
	Short: "Collection commands",
}

var collTreeCmd = &cobra.Command{
	Use:   "tree",
	Short: "Show collection tree",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, _ := getClient()
		search, _ := cmd.Flags().GetString("search")
		depth, _ := cmd.Flags().GetInt("L")

		data := must(client.Get("/collection/tree", nil))

		if jsonFlag {
			printOut(data, true)
			return nil
		}

		var rawItems []interface{}
		json.Unmarshal(data, &rawItems)

		var printTree func(items []interface{}, indent int)
		printTree = func(items []interface{}, indent int) {
			for _, raw := range items {
				item, ok := raw.(map[string]interface{})
				if !ok {
					continue
				}
				name, _ := item["name"].(string)
				if search != "" && !strings.Contains(strings.ToLower(name), strings.ToLower(search)) {
					continue
				}
				fmt.Printf("%s[%5v]  %s\n", strings.Repeat("  ", indent), item["id"], name)
				if children, ok := item["children"].([]interface{}); ok && indent < depth-1 {
					printTree(children, indent+1)
				}
			}
		}
		printTree(rawItems, 0)
		return nil
	},
}

var collGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get collection",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, _ := getClient()
		data := must(client.Get("/collection/"+args[0], nil))
		printOut(data, jsonFlag)
		return nil
	},
}

var collItemsCmd = &cobra.Command{
	Use:   "items <id>",
	Short: "List collection items",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, _ := getClient()
		models, _ := cmd.Flags().GetString("models")

		params := map[string]string{}
		if models != "" {
			params["models"] = models
		}

		data := must(client.Get("/collection/"+args[0]+"/items", params))

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
			for _, item := range items {
				fmt.Printf("  [%-12s] %5v  %s\n", item["model"], item["id"], item["name"])
			}
		}
		return nil
	},
}

var collCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a collection",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, _ := getClient()
		name, _ := cmd.Flags().GetString("name")
		parentID, _ := cmd.Flags().GetString("parent-id")
		description, _ := cmd.Flags().GetString("description")

		body := map[string]interface{}{"name": name}
		if parentID != "" {
			body["parent_id"] = parentID
		}
		if description != "" {
			body["description"] = description
		}

		data := must(client.Post("/collection", body))
		printOut(data, jsonFlag)
		return nil
	},
}

var collUpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update a collection",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, _ := getClient()
		name, _ := cmd.Flags().GetString("name")
		parentID, _ := cmd.Flags().GetString("parent-id")

		body := map[string]interface{}{}
		if name != "" {
			body["name"] = name
		}
		if parentID != "" {
			body["parent_id"] = parentID
		}

		data := must(client.Put("/collection/"+args[0], body))
		printOut(data, jsonFlag)
		return nil
	},
}

var collArchiveCmd = &cobra.Command{
	Use:   "archive <id>",
	Short: "Archive a collection",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, _ := getClient()
		must(client.Put("/collection/"+args[0], map[string]interface{}{"archived": true}))
		fmt.Printf("Archived collection %s\n", args[0])
		return nil
	},
}

func init() {
	rootCmd.AddCommand(collCmd)
	collCmd.AddCommand(collTreeCmd, collGetCmd, collItemsCmd, collCreateCmd, collUpdateCmd, collArchiveCmd)

	collTreeCmd.Flags().String("search", "", "Filter by name")
	collTreeCmd.Flags().Int("L", 1, "Max depth")
	collTreeCmd.Flags().Bool("include-archived", false, "Include archived")

	collItemsCmd.Flags().String("models", "", "Filter by model type")
	collItemsCmd.Flags().Bool("archived", false, "Include archived")
	collItemsCmd.Flags().String("sort-by", "", "Sort field")
	collItemsCmd.Flags().String("sort-dir", "", "Sort direction (asc/desc)")

	collCreateCmd.Flags().String("name", "", "Collection name")
	collCreateCmd.MarkFlagRequired("name")
	collCreateCmd.Flags().String("parent-id", "", "Parent collection ID")
	collCreateCmd.Flags().String("description", "", "Description")

	collUpdateCmd.Flags().String("name", "", "New name")
	collUpdateCmd.Flags().String("parent-id", "", "New parent ID")
}
