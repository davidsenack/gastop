package adapter

import (
	"context"
	"time"

	"github.com/davidsenack/gastop/internal/model"
)

// ConvoyListOpts specifies options for listing convoys.
type ConvoyListOpts struct {
	All    bool   // Include closed convoys
	Status string // Filter by status (open, closed)
}

// ListConvoys returns all convoys matching the options.
// It queries both gt convoy list (town-level) and bd list -t convoy (rig-level).
func (a *Adapter) ListConvoys(ctx context.Context, opts ConvoyListOpts) ([]model.Convoy, error) {
	cacheKey := "convoys"

	args := []string{"convoy", "list", "--json"}
	if opts.All {
		args = append(args, "--all")
		cacheKey += ":all"
	}
	if opts.Status != "" {
		args = append(args, "--status="+opts.Status)
		cacheKey += ":status=" + opts.Status
	}

	var convoys []model.Convoy

	// Try gt convoy list first (town-level convoys)
	out, err := a.execGT(ctx, args...)
	if err == nil {
		_ = parseJSON(out, &convoys)
	}

	// Also query rig-level convoys via bd list -t convoy
	bdArgs := []string{"list", "-t", "convoy", "--json", "--limit", "50"}
	if opts.All {
		bdArgs = append(bdArgs, "--all")
	}
	bdOut, bdErr := a.execBD(ctx, bdArgs...)
	if bdErr == nil {
		var rigConvoys []model.Convoy
		if parseJSON(bdOut, &rigConvoys) == nil {
			convoys = append(convoys, rigConvoys...)
		}
	}

	// If both failed, try cache
	if err != nil && bdErr != nil {
		if cached, ok := a.getCache(cacheKey, 5*time.Minute); ok {
			convoys := cached.([]model.Convoy)
			for i := range convoys {
				convoys[i].Stuck = true
				convoys[i].StuckReason = "stale data"
			}
			return convoys, nil
		}
		return nil, err
	}

	// Compute progress for each convoy
	for i := range convoys {
		convoys[i].ComputeProgress()
	}

	a.setCache(cacheKey, convoys)
	return convoys, nil
}

// GetConvoyStatus returns detailed status for a convoy.
func (a *Adapter) GetConvoyStatus(ctx context.Context, id string) (*model.Convoy, error) {
	cacheKey := "convoy:" + id

	out, err := a.execGT(ctx, "convoy", "status", id, "--json")
	if err != nil {
		if cached, ok := a.getCache(cacheKey, 5*time.Minute); ok {
			convoy := cached.(*model.Convoy)
			convoy.Stuck = true
			convoy.StuckReason = "stale data"
			return convoy, nil
		}
		return nil, err
	}

	var convoy model.Convoy
	if err := parseJSON(out, &convoy); err != nil {
		return nil, err
	}

	convoy.ComputeProgress()
	a.setCache(cacheKey, &convoy)
	return &convoy, nil
}
