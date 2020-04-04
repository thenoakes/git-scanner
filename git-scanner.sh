#!/bin/bash

# Get the working directory as for storing the list
pwd=$(pwd)

# Switch to home
cd ~

# Construct the command for processing and adding each repo to the list
processor="echo  "\${0:2:\${\#0}-7}" | tee -a \"${pwd}/repo-list\""

# Run find
find . \( -path "./.*/*" -prune -o -path "./Library/*" -prune \) -o -name ".git" -exec sh -c "$processor" {} \;

# Switch back to working directory
cd - > /dev/null

# echo -n "Please wait while ~ is scanned for git repositories... "

# spinWhile() {
#     cd ~
#     find . \( -path "./.*/*" -prune -o -path "./Library/*" -prune \) -o -name ".git" -print \
#         | awk '{print substr($0,3,length($0)-7)}'
#     cd - > /dev/null
# }

# . ./spin.sh
