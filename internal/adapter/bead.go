package adapter

import (
	"context"
	"strconv"
	"time"

	"github.com/davidsenack/gastop/internal/model"
)

// BeadListOpts specifies options for listing beads.
type BeadListOpts struct {
	Status   string // Filter by status
	Limit    int    // Max results
	Assignee string // Filter by assignee
	Type     string // Filter by type
	Ready    bool   // Only ready work
	Blocked  bool   // Only blocked work
}

// ListBeads returns beads matching the options.
func (a *Adapter) ListBeads(ctx context.Context, opts BeadListOpts) ([]model.Bead, error) {
	cacheKey := "beads"

	var args []string
	if opts.Ready {
		args = []string{"ready", "--json"}
		cacheKey += ":ready"
	} else if opts.Blocked {
		args = []string{"blocked", "--json"}
		cacheKey += ":blocked"
	} else {
		args = []string{"list", "--json"}
	}

	if opts.Status != "" && !opts.Ready && !opts.Blocked {
		args = append(args, "--status="+opts.Status)
		cacheKey += ":status=" + opts.Status
	}
	if opts.Limit > 0 {
		args = append(args, "--limit", strconv.Itoa(opts.Limit))
		cacheKey += ":limit=" + strconv.Itoa(opts.Limit)
	} else {
		args = append(args, "--limit", "100")
	}
	if opts.Assignee != "" {
		args = append(args, "--assignee="+opts.Assignee)
		cacheKey += ":assignee=" + opts.Assignee
	}
	if opts.Type != "" {
		args = append(args, "--type="+opts.Type)
		cacheKey += ":type=" + opts.Type
	}

	out, err := a.execBD(ctx, args...)
	if err != nil {
		if cached, ok := a.getCache(cacheKey, 5*time.Minute); ok {
			beads := cached.([]model.Bead)
			for i := range beads {
				beads[i].Stuck = true
				beads[i].StuckReason = "stale data"
			}
			return beads, nil
		}
		return nil, err
	}

	var beads []model.Bead
	if err := parseJSON(out, &beads); err != nil {
		return nil, err
	}

	// Compute age for each bead
	for i := range beads {
		beads[i].ComputeAge()
	}

	a.setCache(cacheKey, beads)
	return beads, nil
}

// GetBead returns detailed info for a specific bead.
func (a *Adapter) GetBead(ctx context.Context, id string) (*model.Bead, error) {
	cacheKey := "bead:" + id

	out, err := a.execBD(ctx, "show", id, "--json")
	if err != nil {
		if cached, ok := a.getCache(cacheKey, 5*time.Minute); ok {
			bead := cached.(*model.Bead)
			bead.Stuck = true
			bead.StuckReason = "stale data"
			return bead, nil
		}
		return nil, err
	}

	var bead model.Bead
	if err := parseJSON(out, &bead); err != nil {
		return nil, err
	}

	bead.ComputeAge()
	a.setCache(cacheKey, &bead)
	return &bead, nil
}

// ListReadyBeads returns beads ready for work (no blockers).
func (a *Adapter) ListReadyBeads(ctx context.Context, limit int) ([]model.Bead, error) {
	return a.ListBeads(ctx, BeadListOpts{Ready: true, Limit: limit})
}

// ListBlockedBeads returns blocked beads.
func (a *Adapter) ListBlockedBeads(ctx context.Context) ([]model.Bead, error) {
	return a.ListBeads(ctx, BeadListOpts{Blocked: true})
}
