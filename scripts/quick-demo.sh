#!/bin/bash
# quick-demo.sh - Quick demo that spawns polecats on existing beads
#
# This is a lighter version that doesn't create new beads,
# just spawns polecats on existing ready work.

set -e

RIG="${1:-gastop}"
TOWN_ROOT="${GT_TOWN_ROOT:-$HOME/gt/gastop}"

echo "Quick Demo - Spawning polecats on ready work in $RIG"
echo ""

cd "$TOWN_ROOT/$RIG" 2>/dev/null || {
    echo "Rig not found. Trying gastown..."
    RIG="gastown"
    cd "$TOWN_ROOT/$RIG"
}

# Show current state
echo "Current beads:"
bd list --status=open --limit=5 2>/dev/null || echo "No open beads"
echo ""

echo "Ready to work (no blockers):"
bd ready --limit=5 2>/dev/null || echo "No ready beads"
echo ""

# Get ready beads and sling them
READY=$(bd ready --json --limit=2 2>/dev/null | grep '"id"' | head -2 | sed 's/.*"id": "\([^"]*\)".*/\1/')

for bead in $READY; do
    echo "Slinging $bead..."
    gt sling "$bead" "$RIG" 2>&1 | head -3
    sleep 1
done

echo ""
echo "Polecats now running:"
gt polecat list "$RIG"

echo ""
echo "Run gastop to watch: cd $TOWN_ROOT/gastop-standalone && ./gastop"
