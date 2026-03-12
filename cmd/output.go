package cmd

import (
	"encoding/json"
	"fmt"
	"os"
)

func printJSON(v interface{}) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	_ = enc.Encode(v)
}

// printOut prints data. If asJSON, wraps in {"success":true,"data":...}.
func printOut(data json.RawMessage, asJSON bool) {
	var v interface{}
	json.Unmarshal(data, &v)
	if asJSON {
		printJSON(map[string]interface{}{"success": true, "data": v})
	} else {
		printJSON(v)
	}
}

func mustMarshal(v interface{}) json.RawMessage {
	b, _ := json.Marshal(v)
	return b
}

func die(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
