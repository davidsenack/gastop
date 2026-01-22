package adapter

import (
	"context"
	"time"

	"github.com/davidsenack/gastop/internal/model"
)

// ListPolecats returns all polecats in a rig (or all rigs if rig is empty).
func (a *Adapter) ListPolecats(ctx context.Context, rig string) ([]model.Polecat, error) {
	cacheKey := "polecats"
	if rig != "" {
		cacheKey += ":" + rig
	} else {
		cacheKey += ":all"
	}

	args := []string{"polecat", "list", "--json"}
	if rig != "" {
		args = append(args, rig)
	} else {
		args = append(args, "--all")
	}

	out, err := a.execGT(ctx, args...)
	if err != nil {
		if cached, ok := a.getCache(cacheKey, 5*time.Minute); ok {
			polecats := cached.([]model.Polecat)
			for i := range polecats {
				polecats[i].Stuck = true
				polecats[i].StuckReason = "stale data"
			}
			return polecats, nil
		}
		return nil, err
	}

	var polecats []model.Polecat
	if err := parseJSON(out, &polecats); err != nil {
		return nil, err
	}

	a.setCache(cacheKey, polecats)
	return polecats, nil
}

// GetPolecatStatus returns detailed status for a polecat.
func (a *Adapter) GetPolecatStatus(ctx context.Context, rigPolecat string) (*model.Polecat, error) {
	cacheKey := "polecat:" + rigPolecat

	out, err := a.execGT(ctx, "polecat", "status", rigPolecat, "--json")
	if err != nil {
		if cached, ok := a.getCache(cacheKey, 5*time.Minute); ok {
			polecat := cached.(*model.Polecat)
			polecat.Stuck = true
			polecat.StuckReason = "stale data"
			return polecat, nil
		}
		return nil, err
	}

	var polecat model.Polecat
	if err := parseJSON(out, &polecat); err != nil {
		return nil, err
	}

	a.setCache(cacheKey, &polecat)
	return &polecat, nil
}

// ListStalePolecats returns polecats that may need cleanup.
func (a *Adapter) ListStalePolecats(ctx context.Context, rig string) ([]model.Polecat, error) {
	args := []string{"polecat", "stale"}
	if rig != "" {
		args = append(args, rig)
	}

	out, err := a.execGT(ctx, args...)
	if err != nil {
		return nil, err
	}

	// Note: stale command may not have JSON output yet
	// For now, return empty if we can't parse
	var polecats []model.Polecat
	_ = parseJSON(out, &polecats) // Ignore parse errors
	return polecats, nil
}

// HookStatus represents the response from gt hook show.
type HookStatus struct {
	Agent  string `json:"agent"`
	Status string `json:"status"` // "empty" or "hooked"
	Bead   string `json:"bead,omitempty"`
	Title  string `json:"title,omitempty"`
}

// GetHookedBead returns the bead hooked to a polecat (if any).
func (a *Adapter) GetHookedBead(ctx context.Context, rigPolecat string) (*HookStatus, error) {
	out, err := a.execGT(ctx, "hook", "show", rigPolecat, "--json")
	if err != nil {
		return nil, err
	}

	var status HookStatus
	if err := parseJSON(out, &status); err != nil {
		return nil, err
	}

	return &status, nil
}

// EnrichPolecatsWithHooks fetches hooked bead info for polecats that are working.
// This is a separate call to avoid slowing down the main list.
func (a *Adapter) EnrichPolecatsWithHooks(ctx context.Context, polecats []model.Polecat) {
	for i := range polecats {
		pc := &polecats[i]
		// Only fetch hooks for working polecats to minimize API calls
		if pc.State == "working" || pc.State == "done" {
			status, err := a.GetHookedBead(ctx, pc.FullName())
			if err == nil && status.Status == "hooked" {
				pc.HookedBead = status.Bead
				pc.HookedTitle = status.Title
			}
		}
	}
}

// EnrichPolecatWithDetails fetches detailed status for a single polecat.
func (a *Adapter) EnrichPolecatWithDetails(ctx context.Context, pc *model.Polecat) error {
	detailed, err := a.GetPolecatStatus(ctx, pc.FullName())
	if err != nil {
		return err
	}

	// Copy over detailed fields
	pc.Branch = detailed.Branch
	pc.ClonePath = detailed.ClonePath
	pc.Windows = detailed.Windows
	pc.LastActivity = detailed.LastActivity
	pc.SessionID = detailed.SessionID

	return nil
}
