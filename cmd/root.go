/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"io"
	"net/http"
	"os"
	"runtime"

	"github.com/spf13/cobra"
)

// These are provided by main.go
var BuildGitHash string
var getRemoteLatestHash func() string

var rootCmd = &cobra.Command{
	Use:   "devctl",
	Short: "A pluggable cli to reduce developer friction",
	Long:  `This is a tool to avoid 'magic spellbooks' that people accumulate with text files and copied commands`,
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
		// Get latest hash from remote
		latestHash := getRemoteLatestHash()
		if latestHash == "" {
			cmd.PrintErrln("Could not determine latest version.")
			return
		}
		if latestHash == BuildGitHash {
			cmd.Println("Already up to date.")
			return
		}
		// Download the latest binary
		filename := "devctl-" + osName + "-" + arch + "-" + latestHash
		url := "https://raw.githubusercontent.com/danlafeir/devctl/main/bin/release/" + filename
		cmd.Printf("Downloading %s...\n", url)
		resp, err := http.Get(url)
		if err != nil || resp.StatusCode != 200 {
			cmd.PrintErrln("Failed to download latest binary.")
			return
		}
		defer resp.Body.Close()
		tmpFile, err := os.CreateTemp("", "devctl-update-*")
		if err != nil {
			cmd.PrintErrln("Failed to create temp file.")
			return
		}
		defer os.Remove(tmpFile.Name())
		_, err = io.Copy(tmpFile, resp.Body)
		if err != nil {
			cmd.PrintErrln("Failed to write binary.")
			return
		}
		tmpFile.Close()
		os.Chmod(tmpFile.Name(), 0o755)
		// Replace the current binary
		self, err := os.Executable()
		if err != nil {
			cmd.PrintErrln("Could not determine current executable path.")
			return
		}
		// On Unix, we can move the file in place
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
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// Disable the help command
	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})
	rootCmd.AddCommand(updateCmd)
}
