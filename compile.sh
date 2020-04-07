#!/bin/bash

shc -T -f git-scanner.sh
chmod a+rx git-scanner.sh.x
sudo cp git-scanner.sh.x /usr/local/bin/git-scanner
