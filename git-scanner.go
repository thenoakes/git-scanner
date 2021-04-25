package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/jessevdk/go-flags"
)

const (
	blueColour    = "\033[94m"
	greyColour    = "\033[90m"
	defaultColour = "\033[0m"
)

// opts encapsulates command-line options for git-scanner
type opts struct {
	IgnoreNone bool `short:"i" long:"ignore-none" description:"Ignore repositories with no remote specified"`
}

// parseOptsStd parses command-line options using the flag library
func parseOptsStd() *opts {
	var optI_Long = flag.Bool("ignore-none", false, "Ignore repositories with no remote specified")
	var optI_Short = flag.Bool("i", false, "")
	flag.Parse()
	return &opts{
		IgnoreNone: *optI_Short || *optI_Long,
	}
}

// parseOptsThirdParty parses command-line options using third-party library go-flags
func parseOptsThirdParty() *opts {
	var opts = &opts{}
	parser := flags.NewParser(opts, flags.Default)
	_, err := parser.Parse()

	if err != nil {	
		if !flags.WroteHelp(err) {
			parser.WriteHelp(os.Stdout)
		}
		return nil
	}
	return opts
}

func main() {
	const useNew = false

	var optI = false
	if useNew {
		opts := parseOptsThirdParty()
		if opts == nil {
			return
		}
		optI = opts.IgnoreNone
	} else {
		optI = parseOptsStd().IgnoreNone
	}

	homeDir, err := os.UserHomeDir()

	// isHiddenRoot returns true for a hidden directory which is a direct descendent of ~
	isHiddenRoot := func(path string, info os.DirEntry) bool {
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
	err = filepath.WalkDir(homeDir,
		func(path string, info os.DirEntry, err error) error {
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
				// Get the git remotes by name if there are any
				gitCmd := exec.Command("git", "--git-dir="+path, "remote")
				result, _ := gitCmd.Output()
				
				// Parse the results
				remoteOutput := strings.Split(cleanse(result), "\n")
				
				// Filter out empty lines
				remotes := cleanseLines(remoteOutput)

				if len(remotes) > 0 || !optI {
					// Print the path relative to ~
					fmt.Print(blueColour, getRepoName(path), defaultColour, "\n")
	
					// Obtain URLs and print alongside remote names
					printRemoteUrls(path, remotes)
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

func printRemoteUrls(path string, remoteNames []string) {
	if len(remoteNames) == 0 {
		fmt.Print(greyColour, "    No remotes", defaultColour, "\n")
	} else {
		for _, remoteName := range remoteNames {
			remoteCmd := exec.Command("git", "--git-dir="+path, "remote", "get-url", remoteName)
			remoteUrl, _ := remoteCmd.Output()
			fmt.Printf("    %s: %s\n", remoteName, cleanse(remoteUrl))
		}
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
