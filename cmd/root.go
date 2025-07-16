/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"runtime"

	"github.com/danlafeir/devctl/pkg/plugin"
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
		osName := runtime.GOOS
		arch := runtime.GOARCH
		if osName == "darwin" {
			osName = "darwin"
		} else if osName == "linux" {
			osName = "linux"
		} else {
			cmd.PrintErrln("Unsupported OS for update")
			return
		}
		if arch == "amd64" || arch == "x86_64" {
			arch = "amd64"
		} else if arch == "arm64" || arch == "aarch64" {
			arch = "arm64"
		} else {
			cmd.PrintErrln("Unsupported architecture for update")
			return
		}

		// Get the latest hash from GitHub
		latestHash, err := getLatestHash(osName, arch)
		if err != nil {
			cmd.PrintErrf("Failed to get latest hash: %v\n", err)
			return
		}

		cmd.Printf("Current hash: %s\n", BuildGitHash)
		cmd.Printf("Latest hash: %s\n", latestHash)

		if BuildGitHash == latestHash {
			cmd.Println("Already up to date.")
			return
		}

		filename := "devctl-" + osName + "-" + arch + "-" + latestHash
		url := "https://raw.githubusercontent.com/danlafeir/devctl/main/bin/release/" + filename
		cmd.Printf("Downloading %s...\n", url)
		resp, err := http.Get(url)
		if err != nil {
			cmd.PrintErrf("Failed to download binary: %v\n", err)
			return
		}
		if resp.StatusCode != 200 {
			cmd.PrintErrf("Failed to download binary: HTTP %d\n", resp.StatusCode)
			return
		}
		defer resp.Body.Close()
		tmpFile, err := os.CreateTemp("", "devctl-update-*")
		if err != nil {
			cmd.PrintErrf("Failed to create temp file: %v\n", err)
			return
		}
		defer os.Remove(tmpFile.Name())
		_, err = io.Copy(tmpFile, resp.Body)
		if err != nil {
			cmd.PrintErrf("Failed to write binary: %v\n", err)
			return
		}
		tmpFile.Close()
		if err := os.Chmod(tmpFile.Name(), 0o755); err != nil {
			cmd.PrintErrf("Failed to set permissions on new binary: %v\n", err)
			return
		}
		self, err := os.Executable()
		if err != nil {
			cmd.PrintErrln("Could not determine current executable path.")
			return
		}
		err = os.Rename(tmpFile.Name(), self)
		if err != nil {
			cmd.PrintErrf("Failed to replace binary: %v\n", err)
			return
		}
		cmd.Println("devctl updated to latest version.")
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
	plugin.RegisterPlugins(rootCmd)
}
