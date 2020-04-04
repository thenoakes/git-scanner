#!/bin/bash

# Get the working directory as for storing the list
pwd=$(pwd)

# Switch to home
cd ~

# Construct the command for processing and adding each repo to the list
processor=$(cat <<- CMD
    echo "\033[94m\${0:2:\${#0}-7}\033[0m";
    printf "    ";
    git --git-dir=\${0:2:\${#0}-7}/.git ls-remote --get-url
CMD
)

# Run find
find . \( -path "./.*/*" -prune -o -path "./Library/*" -prune \) -o -name ".git" -exec sh -c "$processor" {} \;

# Switch back to working directory
cd "$pwd" > /dev/null
