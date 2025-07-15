/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/danlafeir/devctl/cmd"
)

// BuildGitHash is set at build time via -ldflags
var BuildGitHash = "dev"

func getRemoteLatestFilename() string {
	osName := runtime.GOOS
	arch := runtime.GOARCH
	if osName == "darwin" {
		osName = "darwin"
	} else if osName == "linux" {
		osName = "linux"
	} else {
		return ""
	}
	if arch == "amd64" || arch == "x86_64" {
		arch = "amd64"
	} else if arch == "arm64" || arch == "aarch64" {
		arch = "arm64"
	} else {
		return ""
	}
	url := "https://api.github.com/repos/danlafeir/devctl/contents/bin/release"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return ""
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return ""
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	prefix := fmt.Sprintf("devctl-%s-%s-", osName, arch)
	var latest string
	for _, line := range strings.Split(string(body), "\n") {
		if idx := strings.Index(line, prefix); idx != -1 {
			start := idx
			end := strings.Index(line[start:], "\"")
			if end != -1 {
				name := line[start : start+end]
				latest = name // last match is latest (if sorted)
			}
		}
	}
	return latest
}

func checkUpgrade() {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return // fail silently
	}
	checkFile := filepath.Join(configDir, "devctl", "upgrade-check")
	os.MkdirAll(filepath.Dir(checkFile), 0o755)

	today := time.Now().Format("2006-01-02")
	var lastDate, lastHash string
	if f, err := os.Open(checkFile); err == nil {
		fmt.Fscanf(f, "%s %s", &lastDate, &lastHash)
		f.Close()
	}
	if lastDate == today {
		return // already checked today
	}

	// Check remote for latest hash
	remoteHash := getRemoteLatestFilename()
	if remoteHash != "" && remoteHash != BuildGitHash {
		fmt.Fprintf(os.Stderr, "A new version of devctl is available (hash: %s). Please upgrade.\n", remoteHash)
	}

	// Write today's check
	f, err := os.Create(checkFile)
	if err == nil {
		fmt.Fprintf(f, "%s %s", today, BuildGitHash)
		f.Close()
	}
}

func main() {
	checkUpgrade()
	cmd.Execute()
}
