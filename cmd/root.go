/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/danlafeir/dev/cmd/jwt"
	"github.com/danlafeir/cli-go/pkg/plugin"
	"github.com/danlafeir/cli-go/pkg/update"
	"github.com/spf13/cobra"
)

// These are provided by main.go
var BuildGitHash string
var BuildLatestHash string

var rootCmd = &cobra.Command{
	Use:   "dev",
	Short: "A pluggable cli to reduce developer friction",
	Long:  `This is a tool to avoid 'magic spellbooks' that people accumulate with text files and copied commands`,
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update dev to the latest version",
	Run: func(cmd *cobra.Command, args []string) {
		update.RunUpdateWithConfig(update.Config{
			AppName: "dev",
			Repo:    "danlafeir/dev",
		}, BuildGitHash, cmd)
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

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.dev.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().MarkHidden("toggle")

	// Hide the help flag
	rootCmd.PersistentFlags().MarkHidden("help")

	// Disable the help command
	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(jwt.GetJWTCommand())
	plugin.RegisterPlugins(rootCmd, "dev-")
}
