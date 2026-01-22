#!/bin/bash
# demo.sh - Demo script to showcase gastop with live data
#
# This script creates test beads, convoys, and spawns polecats
# so you can watch gastop display real-time activity.
#
# Usage: ./scripts/demo.sh [rig]
#   rig: The rig to use (default: gastop)

set -e

RIG="${1:-gastop}"
TOWN_ROOT="${GT_TOWN_ROOT:-$HOME/gt/gastop}"

echo "=========================================="
echo "  gastop Demo Script"
echo "=========================================="
echo ""
echo "Rig: $RIG"
echo "Town: $TOWN_ROOT"
echo ""

# Check if gt and bd are available
if ! command -v gt &> /dev/null; then
    echo "Error: 'gt' command not found. Please install Gas Town."
    exit 1
fi

if ! command -v bd &> /dev/null; then
    echo "Error: 'bd' command not found. Please install Gas Town."
    exit 1
fi

# Change to rig directory
cd "$TOWN_ROOT/$RIG" 2>/dev/null || {
    echo "Error: Rig '$RIG' not found at $TOWN_ROOT/$RIG"
    echo "Available rigs:"
    gt rig list
    exit 1
}

echo "Step 1: Creating test beads..."
echo "----------------------------------------"

# Create beads with different statuses and priorities
bd new -t bug -p 1 --title "Critical: Dashboard not loading" --body "The main dashboard fails to load on startup. High priority fix needed." 2>/dev/null || true
bd new -t feature -p 2 --title "Add dark mode support" --body "Implement a dark mode toggle in the settings panel." 2>/dev/null || true
bd new -t task -p 2 --title "Refactor event handler" --body "Clean up the event handler code for better maintainability." 2>/dev/null || true
bd new -t bug -p 3 --title "Memory leak in refresh loop" --body "Memory usage increases over time during long sessions." 2>/dev/null || true
bd new -t feature -p 3 --title "Add keyboard shortcuts help" --body "Show a help overlay with all available keyboard shortcuts." 2>/dev/null || true
bd new -t chore -p 4 --title "Update dependencies" --body "Update all Go dependencies to latest versions." 2>/dev/null || true

echo "Created test beads."
echo ""

# Get the IDs of beads we just created (or existing ones)
echo "Step 2: Listing available beads..."
echo "----------------------------------------"
bd list --status=open --limit=10

echo ""
echo "Step 3: Creating a convoy to track work..."
echo "----------------------------------------"

# Get some bead IDs to track
BEADS=$(bd list --status=open --json --limit=4 2>/dev/null | grep '"id"' | head -4 | sed 's/.*"id": "\([^"]*\)".*/\1/' | tr '\n' ' ')

if [ -n "$BEADS" ]; then
    # Create a convoy
    gt convoy create "Demo Sprint" $BEADS 2>/dev/null || echo "Convoy may already exist or failed to create"
    echo "Created convoy tracking: $BEADS"
else
    echo "No beads found to track in convoy"
fi

echo ""
echo "Step 4: Spawning polecats to work on issues..."
echo "----------------------------------------"

# Get ready beads (no blockers)
READY_BEADS=$(bd ready --json --limit=3 2>/dev/null | grep '"id"' | head -3 | sed 's/.*"id": "\([^"]*\)".*/\1/')

count=0
for bead in $READY_BEADS; do
    if [ $count -lt 3 ]; then
        echo "Slinging $bead to $RIG..."
        gt sling "$bead" "$RIG" 2>/dev/null || echo "  (may already be assigned or polecat limit reached)"
        count=$((count + 1))
        sleep 2  # Give time between spawns
    fi
done

echo ""
echo "Step 5: Current status..."
echo "----------------------------------------"
echo ""
echo "Polecats:"
gt polecat list "$RIG" 2>/dev/null || echo "No polecats running"

echo ""
echo "Convoys:"
gt convoy list 2>/dev/null || echo "No convoys"

echo ""
echo "=========================================="
echo "  Demo Ready!"
echo "=========================================="
echo ""
echo "Now run gastop to see everything in action:"
echo ""
echo "  cd $TOWN_ROOT/gastop-standalone && ./gastop"
echo ""
echo "Or if gastop is in your PATH:"
echo ""
echo "  gastop"
echo ""
echo "Controls:"
echo "  j/k     - Navigate up/down"
echo "  h/l     - Switch panels"
echo "  x       - Kill polecat / close bead"
echo "  /       - Search beads"
echo "  f       - Filter by status"
echo "  +/-     - Adjust refresh speed"
echo "  ?       - Help"
echo "  q       - Quit"
echo ""
echo "The polecats will work autonomously. Watch them in gastop!"
echo ""
