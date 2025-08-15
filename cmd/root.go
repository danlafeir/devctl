/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"

	"github.com/danlafeir/devctl/cmd/jwt"
	"github.com/danlafeir/devctl/pkg/plugin"
	"github.com/danlafeir/devctl/pkg/update"
	"github.com/spf13/cobra"
)

// These are provided by main.go
var BuildGitHash string
var BuildLatestHash string

var rootCmd = &cobra.Command{
	Use:   "devctl",
	Short: "A pluggable cli to reduce developer friction",
	Long:  `This is a tool to avoid 'magic spellbooks' that people accumulate with text files and copied commands`,
}

// getLatestHash fetches the latest available hash from GitHub
func getLatestHash(osName, arch string) (string, error) {
	apiURL := "https://api.github.com/repos/danlafeir/devctl/contents/bin/release"

	resp, err := http.Get(apiURL)
	if err != nil {
		return "", fmt.Errorf("failed to fetch from GitHub API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var contents []struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&contents); err != nil {
		return "", fmt.Errorf("failed to decode GitHub API response: %v", err)
	}

	// Find the latest binary for this OS/ARCH
	pattern := fmt.Sprintf("devctl-%s-%s-([a-zA-Z0-9]+)", osName, arch)
	re := regexp.MustCompile(pattern)

	var latestHash string
	for _, item := range contents {
		matches := re.FindStringSubmatch(item.Name)
		if len(matches) > 1 {
			hash := matches[1]
			if hash > latestHash {
				latestHash = hash
			}
		}
	}

	if latestHash == "" {
		return "", fmt.Errorf("no binary found for %s/%s", osName, arch)
	}

	return latestHash, nil
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update devctl to the latest version",
	Run: func(cmd *cobra.Command, args []string) {
		update.RunUpdate(BuildGitHash, cmd)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.devctl.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().MarkHidden("toggle")

	// Hide the help flag
	rootCmd.PersistentFlags().MarkHidden("help")

	// Disable the help command
	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(jwt.GetJWTCommand())
	plugin.RegisterPlugins(rootCmd)
}
