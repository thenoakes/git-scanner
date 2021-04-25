A simple shell script which scans a user's home directory for git repositories
and prints out the remote URL for it if there is one.

Originally written in shell script (run `./git-scanner.sh`);

partially rewritten in golang (`go build git-scanner.go` to create binary, then `./git-scanner`).