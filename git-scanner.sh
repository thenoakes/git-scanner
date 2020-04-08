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
    repoPath="\${0:2:\${#0}-7}"
    output=\$(git --git-dir="\${repoPath}/.git" remote)
    if ((\${#output} > 0))
    then
        echo "\033[94m\${repoPath}\033[0m"
        for remote in \$output
        do
            printf "    \$remote: "
            git --git-dir="\${repoPath}/.git" remote get-url \$remote
        done
        status="\$(git -c color.status=always --git-dir="\${repoPath}/.git" --work-tree="\${repoPath}" status -s)"
        if ((\${#status} > 0)) && [ "$opt_s" == "true" ]
        then
            echo "    ------ PENDING CHANGES ------"
            echo "\$status" | sed 's/^/    /'
        fi
    elif [ "$opt_i" != "true" ]
    then
        echo "\033[94m\${repoPath}\033[0m"
        echo "    \033[90mNo remotes\033[0m"
    fi
CMD
)

# Run find
find . -path "./.*/*" -prune -o -path "./Library/*" -prune -o -name ".git" -exec sh -c "$processor" {} \;

# Switch back to working directory
cd - > /dev/null
