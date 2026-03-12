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

var dashCmd = &cobra.Command{
	Use:   "dashboards",
	Short: "Dashboard commands",
}

var dashListCmd = &cobra.Command{
	Use:   "list",
	Short: "List dashboards",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, _ := getClient()
		collID, _ := cmd.Flags().GetString("collection-id")

		params := map[string]string{}
		if collID != "" {
			params["collection_id"] = collID
		}

		data := must(client.Get("/dashboard", params))

		var items []map[string]interface{}
		json.Unmarshal(data, &items)

		if jsonFlag {
			printOut(mustMarshal(items), true)
		} else {
			for _, d := range items {
				fmt.Printf("  %5v  %s\n", d["id"], d["name"])
			}
		}
		return nil
	},
}

var dashGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get dashboard",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, _ := getClient()
		includeCards, _ := cmd.Flags().GetBool("include-cards")

		data := must(client.Get("/dashboard/"+args[0], nil))

		if includeCards {
			var dash map[string]interface{}
			json.Unmarshal(data, &dash)

			var dashcards []map[string]interface{}
			if dc, ok := dash["dashcards"]; ok {
				b, _ := json.Marshal(dc)
				json.Unmarshal(b, &dashcards)
			}

			cards := map[string]interface{}{}
			for _, dc := range dashcards {
				if id, ok := dc["card_id"]; ok && id != nil {
					cid := fmt.Sprintf("%v", id)
					if cdata, err := client.Get("/card/"+cid, nil); err == nil {
						var card interface{}
						json.Unmarshal(cdata, &card)
						cards[cid] = card
					}
				}
			}
			dash["_cards"] = cards
			printOut(mustMarshal(dash), jsonFlag)
		} else {
			printOut(data, jsonFlag)
		}
		return nil
	},
}

var dashExportCmd = &cobra.Command{
	Use:   "export <id>",
	Short: "Export dashboard and all cards to temp files",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, _ := getClient()

		data := must(client.Get("/dashboard/"+args[0], nil))

		var dash map[string]interface{}
		json.Unmarshal(data, &dash)

		var dashcards []map[string]interface{}
		if dc, ok := dash["dashcards"]; ok {
			b, _ := json.Marshal(dc)
			json.Unmarshal(b, &dashcards)
		}

		ts := time.Now().Format("20060102-150405")
		tmp := filepath.Join(os.TempDir(), "metabase-"+ts)
		os.MkdirAll(tmp, 0755)

		envelope := map[string]interface{}{
			"export_version": "1.0",
			"type":           "dashboard",
			"dashboard":      dash,
		}
		b, _ := json.MarshalIndent(envelope, "", "  ")
		dashFile := filepath.Join(tmp, fmt.Sprintf("dashboard-%s.json", args[0]))
		os.WriteFile(dashFile, b, 0644)
		fmt.Printf("Dashboard → %s\n", dashFile)

		for _, dc := range dashcards {
			if id, ok := dc["card_id"]; ok && id != nil {
				cid := fmt.Sprintf("%v", id)
				if cdata, err := client.Get("/card/"+cid, nil); err == nil {
					var card interface{}
					json.Unmarshal(cdata, &card)
					cb, _ := json.MarshalIndent(card, "", "  ")
					cf := filepath.Join(tmp, fmt.Sprintf("card-%s.json", cid))
					os.WriteFile(cf, cb, 0644)
					fmt.Printf("Card %s → %s\n", cid, cf)
				}
			}
		}
		return nil
	},
}

var dashImportCmd = &cobra.Command{
	Use:   "import",
	Short: "Create or update a dashboard from JSON",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, _ := getClient()
		file, _ := cmd.Flags().GetString("file")
		id, _ := cmd.Flags().GetInt("id")
		collID, _ := cmd.Flags().GetString("collection-id")

		var dashData map[string]interface{}
		if file == "-" {
			json.NewDecoder(os.Stdin).Decode(&dashData)
		} else {
			b, _ := os.ReadFile(file)
			json.Unmarshal(b, &dashData)
		}

		// Unwrap export envelope if present
		if _, ok := dashData["export_version"]; ok {
			if d, ok := dashData["dashboard"].(map[string]interface{}); ok {
				dashData = d
			}
		}

		if collID != "" {
			dashData["collection_id"] = collID
		}

		var result json.RawMessage
		if id > 0 {
			result = must(client.Put(fmt.Sprintf("/dashboard/%d", id), dashData))
			var r map[string]interface{}
			json.Unmarshal(result, &r)
			fmt.Printf("Updated dashboard %d: %v\n", id, r["name"])
		} else {
			result = must(client.Post("/dashboard", dashData))
			var r map[string]interface{}
			json.Unmarshal(result, &r)
			fmt.Printf("Created dashboard %v: %v\n", r["id"], r["name"])
		}
		if jsonFlag {
			printOut(result, true)
		}
		return nil
	},
}

var dashRevisionsCmd = &cobra.Command{
	Use:   "revisions <id>",
	Short: "List dashboard revisions",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, _ := getClient()
		data := must(client.Get("/revision", map[string]string{
			"entity": "dashboard",
			"id":     args[0],
		}))
		printOut(data, jsonFlag)
		return nil
	},
}

var dashRevertCmd = &cobra.Command{
	Use:   "revert <id> <revision_id>",
	Short: "Revert dashboard to a revision",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, _ := getClient()
		data := must(client.Post("/revision/revert", map[string]string{
			"entity":      "dashboard",
			"id":          args[0],
			"revision_id": args[1],
		}))
		printOut(data, jsonFlag)
		return nil
	},
}

var dashArchiveCmd = &cobra.Command{
	Use:   "archive <id>",
	Short: "Archive a dashboard",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, _ := getClient()
		must(client.Put("/dashboard/"+args[0], map[string]interface{}{"archived": true}))
		fmt.Printf("Archived dashboard %s\n", args[0])
		return nil
	},
}

var dashDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Permanently delete a dashboard",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, _ := getClient()
		force, _ := cmd.Flags().GetBool("force")

		if !force {
			fmt.Printf("Permanently delete dashboard %s? [y/N] ", args[0])
			r := bufio.NewReader(os.Stdin)
			resp, _ := r.ReadString('\n')
			if strings.ToLower(strings.TrimSpace(resp)) != "y" {
				fmt.Println("Cancelled.")
				return nil
			}
		}
		must(client.Delete("/dashboard/" + args[0]))
		fmt.Printf("Deleted dashboard %s\n", args[0])
		return nil
	},
}

func init() {
	rootCmd.AddCommand(dashCmd)
	dashCmd.AddCommand(dashListCmd, dashGetCmd, dashExportCmd, dashImportCmd,
		dashRevisionsCmd, dashRevertCmd, dashArchiveCmd, dashDeleteCmd)

	dashListCmd.Flags().String("collection-id", "", "Filter by collection ID")

	dashGetCmd.Flags().Bool("include-cards", false, "Fetch and embed card data")

	dashImportCmd.Flags().String("file", "", "JSON file path (or - for stdin)")
	dashImportCmd.MarkFlagRequired("file")
	dashImportCmd.Flags().Int("id", 0, "Update existing dashboard with this ID")
	dashImportCmd.Flags().String("collection-id", "", "Override collection ID")

	dashDeleteCmd.Flags().Bool("force", false, "Skip confirmation prompt")
}
