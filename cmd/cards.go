package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var cardsCmd = &cobra.Command{
	Use:   "cards",
	Short: "Card commands",
}

var cardsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List cards",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, _ := getClient()
		filter, _ := cmd.Flags().GetString("filter")
		collID, _ := cmd.Flags().GetString("collection-id")

		params := map[string]string{"f": filter}
		if collID != "" {
			params["collection_id"] = collID
		}

		data := must(client.Get("/card", params))

		var items []map[string]interface{}
		json.Unmarshal(data, &items)

		if jsonFlag {
			printOut(mustMarshal(items), true)
		} else {
			for _, card := range items {
				fmt.Printf("  %5v  %-8s  %s\n", card["id"], card["display"], card["name"])
			}
		}
		return nil
	},
}

var cardsGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get card",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, _ := getClient()
		data := must(client.Get("/card/"+args[0], nil))
		printOut(data, jsonFlag)
		return nil
	},
}

var cardsRunCmd = &cobra.Command{
	Use:   "run <id>",
	Short: "Run a card query",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, _ := getClient()
		limit, _ := cmd.Flags().GetInt("limit")
		parameters, _ := cmd.Flags().GetString("parameters")

		body := map[string]interface{}{"ignore_cache": true}
		if parameters != "" {
			var params interface{}
			json.Unmarshal([]byte(parameters), &params)
			body["parameters"] = params
		}

		data := must(client.Post("/card/"+args[0]+"/query", body))

		var result map[string]json.RawMessage
		json.Unmarshal(data, &result)

		var qdata map[string]json.RawMessage
		json.Unmarshal(result["data"], &qdata)

		var rows [][]interface{}
		var cols []map[string]interface{}
		json.Unmarshal(qdata["rows"], &rows)
		json.Unmarshal(qdata["cols"], &cols)

		colNames := make([]string, len(cols))
		for i, c := range cols {
			if n, ok := c["display_name"].(string); ok && n != "" {
				colNames[i] = n
			} else if n, ok := c["name"].(string); ok {
				colNames[i] = n
			}
		}

		if limit > 0 && limit < len(rows) {
			rows = rows[:limit]
		}

		records := make([]map[string]interface{}, len(rows))
		for i, row := range rows {
			rec := make(map[string]interface{})
			for j, val := range row {
				if j < len(colNames) {
					rec[colNames[j]] = val
				}
			}
			records[i] = rec
		}

		if jsonFlag {
			printOut(mustMarshal(records), true)
		} else {
			ts := time.Now().Format("20060102-150405")
			tmp := filepath.Join(os.TempDir(), "metabase-"+ts)
			os.MkdirAll(tmp, 0755)
			outFile := filepath.Join(tmp, fmt.Sprintf("card-%s.json", args[0]))
			b, _ := json.MarshalIndent(records, "", "  ")
			os.WriteFile(outFile, b, 0644)
			fmt.Printf("Results (%d rows) → %s\n", len(rows), outFile)
		}
		return nil
	},
}

var cardsImportCmd = &cobra.Command{
	Use:   "import",
	Short: "Create or update a card from JSON",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, _ := getClient()
		file, _ := cmd.Flags().GetString("file")
		id, _ := cmd.Flags().GetInt("id")
		collID, _ := cmd.Flags().GetString("collection-id")
		dbID, _ := cmd.Flags().GetString("database-id")

		var cardData map[string]interface{}
		if file == "-" {
			json.NewDecoder(os.Stdin).Decode(&cardData)
		} else {
			b, _ := os.ReadFile(file)
			json.Unmarshal(b, &cardData)
		}

		if collID != "" {
			cardData["collection_id"] = collID
		}
		if dbID != "" {
			if dq, ok := cardData["dataset_query"].(map[string]interface{}); ok {
				dq["database"] = dbID
			}
		}

		var result json.RawMessage
		if id > 0 {
			result = must(client.Put(fmt.Sprintf("/card/%d", id), cardData))
			var r map[string]interface{}
			json.Unmarshal(result, &r)
			fmt.Printf("Updated card %d: %v\n", id, r["name"])
		} else {
			result = must(client.Post("/card", cardData))
			var r map[string]interface{}
			json.Unmarshal(result, &r)
			fmt.Printf("Created card %v: %v\n", r["id"], r["name"])
		}
		if jsonFlag {
			printOut(result, true)
		}
		return nil
	},
}

var cardsArchiveCmd = &cobra.Command{
	Use:   "archive <id>",
	Short: "Archive a card",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, _ := getClient()
		must(client.Put("/card/"+args[0], map[string]interface{}{"archived": true}))
		fmt.Printf("Archived card %s\n", args[0])
		return nil
	},
}

var cardsDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Permanently delete a card",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, _ := getClient()
		force, _ := cmd.Flags().GetBool("force")

		if !force {
			fmt.Printf("Permanently delete card %s? [y/N] ", args[0])
			r := bufio.NewReader(os.Stdin)
			resp, _ := r.ReadString('\n')
			if strings.ToLower(strings.TrimSpace(resp)) != "y" {
				fmt.Println("Cancelled.")
				return nil
			}
		}
		must(client.Delete("/card/" + args[0]))
		fmt.Printf("Deleted card %s\n", args[0])
		return nil
	},
}

func init() {
	rootCmd.AddCommand(cardsCmd)
	cardsCmd.AddCommand(cardsListCmd, cardsGetCmd, cardsRunCmd, cardsImportCmd, cardsArchiveCmd, cardsDeleteCmd)

	cardsListCmd.Flags().String("filter", "mine", "Filter: mine, all, bookmarked, etc.")
	cardsListCmd.Flags().String("collection-id", "", "Filter by collection ID")
	cardsListCmd.Flags().String("database-id", "", "Filter by database ID")

	cardsRunCmd.Flags().Int("limit", 0, "Row limit (0 = all)")
	cardsRunCmd.Flags().String("parameters", "", "JSON parameters array")

	cardsImportCmd.Flags().String("file", "", "JSON file path (or - for stdin)")
	cardsImportCmd.MarkFlagRequired("file")
	cardsImportCmd.Flags().Int("id", 0, "Update existing card with this ID")
	cardsImportCmd.Flags().String("collection-id", "", "Override collection ID")
	cardsImportCmd.Flags().String("database-id", "", "Override database ID")

	cardsDeleteCmd.Flags().Bool("force", false, "Skip confirmation prompt")
}
