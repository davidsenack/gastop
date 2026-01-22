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
