#!/bin/bash
# Wrapper to run gastop with the correct town root
TOWN_ROOT="${GT_TOWN_ROOT:-/Users/davidsenack/gt/gastop}"
exec "$(dirname "$0")/gastop" -town "$TOWN_ROOT" "$@"
