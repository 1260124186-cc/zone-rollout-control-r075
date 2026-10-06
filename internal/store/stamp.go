package store

import (
	"context"
	"sync"

	"solo-0014-zone-rollout-control/internal/domain"
)

// StampRecorder collects the snapshot stamp that store operations observe
// while serving a single request. The last recorded stamp wins, so a
// request that commits a write reports the revision produced by that
// commit, while a read-only request reports the revision it read from.
type StampRecorder struct {
	mu    sync.Mutex
	stamp domain.SnapshotStamp
	valid bool
}

func NewStampRecorder() *StampRecorder {
	return &StampRecorder{}
}

func (r *StampRecorder) Record(stamp domain.SnapshotStamp) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.stamp = stamp
	r.valid = true
}

func (r *StampRecorder) Stamp() (domain.SnapshotStamp, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.stamp, r.valid
}

type stampRecorderContextKey struct{}

func WithStampRecorder(ctx context.Context, recorder *StampRecorder) context.Context {
	return context.WithValue(ctx, stampRecorderContextKey{}, recorder)
}

func StampRecorderFrom(ctx context.Context) *StampRecorder {
	recorder, _ := ctx.Value(stampRecorderContextKey{}).(*StampRecorder)
	return recorder
}

func stampOf(state *domain.State) domain.SnapshotStamp {
	return domain.SnapshotStamp{
		Revision:    state.Revision,
		Plans:       len(state.Plans),
		Revisions:   len(state.Revisions),
		Nodes:       len(state.Nodes),
		Leases:      len(state.Leases),
		LeaseGroups: len(state.LeaseGroups),
		Probes:      len(state.Probes),
		Assessments: len(state.Assessments),
		AuditEvents: len(state.AuditEvents),
	}
}

func recordStamp(ctx context.Context, state *domain.State) {
	if recorder := StampRecorderFrom(ctx); recorder != nil {
		recorder.Record(stampOf(state))
	}
}
