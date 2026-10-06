package store

import (
	"context"
	"sync"

	"solo-0014-zone-rollout-control/internal/domain"
)

// SnapshotRecorder records the snapshot stamp observed by the store operation
// that ran for a given request. It lets a successful response carry the stamp
// of the state visible after its handler completed instead of a stamp read
// before the request was processed.
type SnapshotRecorder struct {
	mu    sync.Mutex
	stamp domain.SnapshotStamp
	seen  bool
}

func NewSnapshotRecorder() *SnapshotRecorder {
	return &SnapshotRecorder{}
}

func (r *SnapshotRecorder) observe(stamp domain.SnapshotStamp) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.stamp = stamp
	r.seen = true
}

// Observe records a stamp observed outside a store operation (for example a
// fallback stamp before the handler runs). Store operations overwrite it with
// the stamp they actually observed or committed.
func (r *SnapshotRecorder) Observe(stamp domain.SnapshotStamp) {
	r.observe(stamp)
}

func (r *SnapshotRecorder) Stamp() (domain.SnapshotStamp, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.stamp, r.seen
}

type snapshotRecorderContextKey struct{}

func ContextWithSnapshotRecorder(ctx context.Context, recorder *SnapshotRecorder) context.Context {
	return context.WithValue(ctx, snapshotRecorderContextKey{}, recorder)
}

func SnapshotRecorderFromContext(ctx context.Context) *SnapshotRecorder {
	recorder, _ := ctx.Value(snapshotRecorderContextKey{}).(*SnapshotRecorder)
	return recorder
}

func snapshotRecorder(ctx context.Context) *SnapshotRecorder {
	return SnapshotRecorderFromContext(ctx)
}
