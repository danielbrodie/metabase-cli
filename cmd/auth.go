package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"

	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authentication commands",
}

var authStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show current auth status",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, _ := getClient()
		data, err := client.Get("/user/current", nil)
		if err != nil {
			printOut(mustMarshal(map[string]interface{}{
				"authenticated": false,
				"url":           client.Profile.URL,
			}), jsonFlag)
			return nil
		}
		var user map[string]interface{}
		json.Unmarshal(data, &user)
		info := map[string]interface{}{
			"authenticated": true,
			"url":           client.Profile.URL,
			"email":         user["email"],
			"name":          strings.TrimSpace(fmt.Sprintf("%v %v", user["first_name"], user["last_name"])),
		}
		printOut(mustMarshal(info), jsonFlag)
		return nil
	},
}

var authLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log in to Metabase",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, cfg := getClient()

		urlVal, _ := cmd.Flags().GetString("url")
		emailVal, _ := cmd.Flags().GetString("email")
		passVal, _ := cmd.Flags().GetString("password")

		if urlVal == "" {
			if client.Profile.URL != "" {
				urlVal = client.Profile.URL
			} else {
				urlVal = promptLine("Metabase URL: ")
			}
		}
		if emailVal == "" {
			emailVal = promptLine("Email: ")
		}
		if passVal == "" {
			passVal = promptPassword("Password: ")
		}

		client.Profile.URL = strings.TrimRight(urlVal, "/")
		result := must(client.Post("/session", map[string]string{
			"username": emailVal,
			"password": passVal,
		}))

		var r map[string]interface{}
		json.Unmarshal(result, &r)
		token, _ := r["id"].(string)

		client.Profile.Token = token
		cfg.Profiles[profileFlag] = client.Profile
		if err := cfg.Save(); err != nil {
			return fmt.Errorf("save config: %w", err)
		}
		fmt.Printf("✓ Logged in as %s → %s\n", emailVal, urlVal)
		return nil
	},
}

var authTokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Print current session token",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, _ := getClient()
		fmt.Println(client.Profile.Token)
		return nil
	},
}

var authLogoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Log out of Metabase",
	RunE: func(cmd *cobra.Command, args []string) error {
		_, cfg := getClient()
		if p, ok := cfg.Profiles[profileFlag]; ok {
			p.Token = ""
			cfg.Save()
		}
		fmt.Println("Logged out.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(authCmd)
	authCmd.AddCommand(authStatusCmd, authLoginCmd, authTokenCmd, authLogoutCmd)

	authLoginCmd.Flags().String("url", "", "Metabase URL")
	authLoginCmd.Flags().String("email", "", "Email address")
	authLoginCmd.Flags().String("password", "", "Password")
}

func promptLine(p string) string {
	fmt.Print(p)
	r := bufio.NewReader(os.Stdin)
	line, _ := r.ReadString('\n')
	return strings.TrimSpace(line)
}

func promptPassword(p string) string {
	fmt.Print(p)
	b, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()
	if err != nil {
		return promptLine("")
	}
	return string(b)
}
