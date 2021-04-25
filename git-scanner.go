package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	blueColour := "\033[94m"
	greyColour := "\033[90m"
	defaultColour := "\033[0m"

	homeDir, _ := os.UserHomeDir()

	// Walk all filepaths from the user's home directory...
	err := filepath.Walk(homeDir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				//return err
				// TODO: Basically we're skipping because we don't want to deal with 'access denied'
				// but maybe we can do better
				return filepath.SkipDir
			}

			// Skip hidden directories which are direct descendents
			if info.IsDir() && path == filepath.Join(homeDir, info.Name()) && strings.HasPrefix(info.Name(), ".") {
				return filepath.SkipDir
			}

			// Log all .git directories
			if info.IsDir() && info.Name() == ".git" {
				// Get a path relative to ~
				shortPath := strings.Replace(path, homeDir, "~", 1)

				// Strip the ~ prefix and .git suffix and print to console
				repoName := shortPath[2 : len(shortPath)-5]
				fmt.Print(blueColour, repoName, defaultColour, "\n")

				// Get the git remotes by name if there are any
				gitCmd := exec.Command("git", "--git-dir="+path, "remote")
				result, _ := gitCmd.Output()

				// Parse the results
				remoteOutput := strings.Split(strings.TrimSpace(string(result)), "\n")

				// Filter out empty lines
				remotes := []string{}
				for _, line := range remoteOutput {
					if len(line) > 0 {
						remotes = append(remotes, line)
					}
				}

				// Obtain URLs and print alongside remote names
				if len(remotes) > 0 {
					for _, remote := range remotes {
						remoteCmd := exec.Command("git", "--git-dir="+path, "remote", "get-url", remote)
						remoteUrl, _ := remoteCmd.Output()
						fmt.Printf("    %s: %s\n", remote, strings.TrimSpace(string(remoteUrl)))
					}
				} else {
					fmt.Print(greyColour, "    No remotes", defaultColour, "\n")
				}

				// Once we've found a .git directory, save time by not delving any further
				return filepath.SkipDir
			}

			return nil
		})

	if err != nil {
		os.Exit(0)
	}

	//fmt.Printf("%d", len(repos))

}
