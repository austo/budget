#!/bin/bash

set -o nounset
set -o errexit

EXENAME='../out/'$(echo $1 | sed 's/\.go//')

printf "$(tput bold)$(tput setaf 2)Building target: \"%s\"...$(tput sgr0)\n\n" "$EXENAME"

go build -o $EXENAME $1