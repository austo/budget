#!/bin/bash

set -o nounset
set -o errexit

EXENAME='../bin/app'
SRCFILES=$(ls src/*.go)

printf "Building target: $(tput bold)$(tput setaf 2)\"%s\"$(tput sgr0) from:\n" "$EXENAME"

for var in $SRCFILES; do
	printf "  %s\n" "$var"
done

go build -o $EXENAME $SRCFILES