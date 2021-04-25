#!/bin/bash

source ./mnopts/mnopts.sh "ignore-none:i show-pending:s" "" "$@"

[[ "$opt_i" = "true" ]] && echo "Ignoring repositories with no remote specified (-i option)"
[[ "$opt_s" = "true" ]] && echo "Showing pending changes where relevant (-s option)"

# Get the working directory as for storing the list
pwd=$(pwd)

# Switch to home
cd ~

# Construct the command for processing and adding each repo to the list:
# Print the repo name, check for existence of remotes then print the URI if it exits
processor=$(cat <<- CMD
    # Slice the path ${0} to remove the ./ characters at the start and the /.git at the end
    repoPath="\${0:2:\${#0}-7}"
    # Get the git remotes by nameif there are any
    output=\$(git --git-dir="\${repoPath}/.git" remote)
    if ((\${#output} > 0))
    then
        # Print that a repo was found and give the name and URL of each remote from git if there are any
        echo "\033[94m\${repoPath}\033[0m"
        for remote in \$output
        do
            printf "    \$remote: "
            git --git-dir="\${repoPath}/.git" remote get-url \$remote
        done
        # Get the file status output from git
        status="\$(git -c color.status=always --git-dir="\${repoPath}/.git" --work-tree="\${repoPath}" status -s)"
        if ((\${#status} > 0)) && [ "$opt_s" == "true" ]
        then
            # If configured to print the changes appending a four-space indent
            echo "    ------ PENDING CHANGES ------"
            echo "\$status" | sed 's/^/    /'
        fi
    elif [ "$opt_i" != "true" ]
    then
        # Print that a repo was found with no remotes if configured to
        echo "\033[94m\${repoPath}\033[0m"
        echo "    \033[90mNo remotes\033[0m"
    fi
CMD
)

# Run find from the current direcrory 
#   (but skip subdirectories of hidden directories and Library - this just avoids undesired output from these directories)
#                Run the script - {} is replaced by the path
find -s . \
    -path "./.Trash" -prune -o \
    -path "./Pictures/*" -prune -o \
    -path "./.*/*" -prune -o -path "./Library/*" -prune -o \
    -name ".git" -exec sh -c "$processor" {} \;

# Switch back to working directory
cd - > /dev/null
