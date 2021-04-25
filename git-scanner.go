package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	const (
		blueColour    = "\033[94m"
		greyColour    = "\033[90m"
		defaultColour = "\033[0m"
	)

	homeDir, err := os.UserHomeDir()

	// isHiddenRoot returns true for a hidden directory which is a direct descendent of ~
	isHiddenRoot := func(path string, info os.FileInfo) bool {
		return path == filepath.Join(homeDir, info.Name()) && strings.HasPrefix(info.Name(), ".")
	}

	getRepoName := func(path string) string {
		shortPath := strings.Replace(path, homeDir, "~", 1)
		// Strip the ~ prefix and .git suffix and print to console
		return shortPath[2 : len(shortPath)-5]
	}

	if err != nil {
		fmt.Println("Fatal error: unable to determine home directory")
		os.Exit(1)
	}

	// Walk all filepaths from the user's home directory...
	err = filepath.Walk(homeDir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				//return err
				// TODO: Basically we're skipping because we don't want to deal with 'access denied'
				// but maybe we can do better
				return filepath.SkipDir
			}

			// Skip hidden directories which are direct descendents
			if info.IsDir() && isHiddenRoot(path, info) {
				return filepath.SkipDir
			}

			// Log all .git directories
			if info.IsDir() && info.Name() == ".git" {
				// Print the path relative to ~
				fmt.Print(blueColour, getRepoName(path), defaultColour, "\n")

				// Get the git remotes by name if there are any
				gitCmd := exec.Command("git", "--git-dir="+path, "remote")
				result, _ := gitCmd.Output()

				// Parse the results
				remoteOutput := strings.Split(cleanse(result), "\n")

				// Filter out empty lines
				remotes := cleanseLines(remoteOutput)

				// Obtain URLs and print alongside remote names
				if len(remotes) > 0 {
					for _, remote := range remotes {
						remoteCmd := exec.Command("git", "--git-dir="+path, "remote", "get-url", remote)
						remoteUrl, _ := remoteCmd.Output()
						fmt.Printf("    %s: %s\n", remote, cleanse(remoteUrl))
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

}

func cleanse(output []byte) string {
	return strings.TrimSpace(string(output))
}

func cleanseLines(output []string) []string {
	results := []string{}
	for _, line := range output {
		if len(line) > 0 {
			results = append(results, line)
		}
	}
	return results
}
