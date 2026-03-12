package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search Metabase",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, _ := getClient()

		models, _ := cmd.Flags().GetString("models")
		collID, _ := cmd.Flags().GetString("collection-id")
		limit, _ := cmd.Flags().GetInt("limit")

		params := map[string]string{
			"q":     args[0],
			"limit": fmt.Sprintf("%d", limit),
		}
		if models != "" {
			params["models"] = models
		}
		if collID != "" {
			params["collection_id"] = collID
		}

		data := must(client.Get("/search", params))

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
				model, _ := item["model"].(string)
				name, _ := item["name"].(string)
				id := item["id"]
				fmt.Printf("[%-12s] %5v  %s\n", model, id, name)
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
	searchCmd.Flags().String("models", "", "Filter by model type (card, dashboard, collection, ...)")
	searchCmd.Flags().String("collection-id", "", "Filter by collection ID")
	searchCmd.Flags().String("database-id", "", "Filter by database ID")
	searchCmd.Flags().Int("limit", 50, "Max results")
	searchCmd.Flags().Bool("archived", false, "Include archived items")
}
