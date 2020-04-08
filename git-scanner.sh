#!/bin/bash

source ./mnopts/mnopts.sh "ignore:i" "" "$@"

[[ "$opt_i" = "true" ]] && echo "Ignoring repositories with no remote specified (-i option)"

# Get the working directory as for storing the list
pwd=$(pwd)

# Switch to home
cd ~

# Construct the command for processing and adding each repo to the list:
# Print the repo name, check for existence of remotes then print the URI if it exits
processor=$(cat <<- CMD
    output=\$(git --git-dir="\${0:2:\${#0}-7}/.git" remote)
    if ((\${#output} > 0))
    then
        echo "\033[94m\${0:2:\${#0}-7}\033[0m"
        for remote in \$output
        do
            printf "    \$remote: "
            git --git-dir="\${0:2:\${#0}-7}/.git" remote get-url \$remote
        done
    elif [ "$opt_i" != "true" ]
    then
        echo "\033[94m\${0:2:\${#0}-7}\033[0m"
        echo "    \033[90mNo remotes\033[0m"
    fi
CMD
)

# Run find
find . -path "./.*/*" -prune -o -path "./Library/*" -prune -o -name ".git" -exec sh -c "$processor" {} \;

# Switch back to working directory
cd - > /dev/null
