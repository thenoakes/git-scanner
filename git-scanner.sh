#!/bin/bash

# Get the working directory as for storing the list
pwd=$(pwd)

# Switch to home
cd ~

# Construct the command for processing and adding each repo to the list:
# Print the repo name, check for existence of remotes then print the URI if it exits
processor=$(cat <<- CMD
    echo "\033[94m\${0:2:\${#0}-7}\033[0m"
    output=\$(git --git-dir="\${0:2:\${#0}-7}/.git" remote -v)
    printf "    "
    if ((\${#output} > 0))
    then
        git --git-dir="\${0:2:\${#0}-7}/.git" ls-remote --get-url
    else
        echo "\033[90mNo remotes\033[0m"
    fi
CMD
)

# Run find
find . -path "./.*/*" -prune -o -path "./Library/*" -prune -o -name ".git" -exec sh -c "$processor" {} \;

# Switch back to working directory
cd - > /dev/null
