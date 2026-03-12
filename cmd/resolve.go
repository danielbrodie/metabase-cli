package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/spf13/cobra"
)

var resolveCmd = &cobra.Command{
	Use:   "resolve <url>",
	Short: "Resolve a Metabase URL to its entity",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, _ := getClient()
		rawURL := args[0]

		patterns := []struct {
			re     *regexp.Regexp
			entity string
		}{
			{regexp.MustCompile(`/question/(\d+)`), "card"},
			{regexp.MustCompile(`/dashboard/(\d+)`), "dashboard"},
			{regexp.MustCompile(`/collection/(\d+)`), "collection"},
			{regexp.MustCompile(`/browse/databases/(\d+)`), "database"},
		}

		for _, p := range patterns {
			m := p.re.FindStringSubmatch(rawURL)
			if m == nil {
				continue
			}
			eid := m[1]
			data := must(client.Get("/"+p.entity+"/"+eid, nil))

			var entity map[string]interface{}
			json.Unmarshal(data, &entity)

			id, _ := strconv.Atoi(eid)
			result := map[string]interface{}{
				"entity": p.entity,
				"id":     id,
				"name":   entity["name"],
				"data":   entity,
			}
			printOut(mustMarshal(result), jsonFlag)
			return nil
		}

		fmt.Fprintf(os.Stderr, "{\"success\":false,\"error\":{\"code\":\"NOT_FOUND\",\"message\":\"Could not parse URL: %s\"}}\n", rawURL)
		os.Exit(1)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(resolveCmd)
}
