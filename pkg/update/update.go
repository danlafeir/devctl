package update

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"runtime"

	"github.com/spf13/cobra"
)

func getLatestHash(apiURL, osName, arch string) (string, error) {
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

func RunUpdate(currentHash string, cmd *cobra.Command) {
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
	apiURL := "https://api.github.com/repos/danlafeir/devctl/contents/bin/release"
	latestHash, err := getLatestHash(apiURL, osName, arch)
	if err != nil {
		cmd.PrintErrf("Failed to get latest hash: %v\n", err)
		return
	}
	cmd.Printf("Current hash: %s\n", currentHash)
	cmd.Printf("Latest hash: %s\n", latestHash)
	if currentHash == latestHash {
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
		if os.IsPermission(err) {
			cmd.Println("Permission denied. Retrying with sudo...")
			mvCmd := exec.Command("sudo", "mv", tmpFile.Name(), self)
			mvCmd.Stdin = os.Stdin
			mvCmd.Stdout = os.Stdout
			mvCmd.Stderr = os.Stderr
			if err := mvCmd.Run(); err != nil {
				cmd.PrintErrf("Failed to replace binary with sudo: %v\n", err)
				return
			}
		} else {
			cmd.PrintErrf("Failed to replace binary: %v\n", err)
			return
		}
	}
	cmd.Println("devctl updated to latest version.")
}
