package platform

import (
	"context"
	"sync"
	"time"
)

type ConfigVersion struct {
	Version   int64
	Payload   string
	CreatedAt time.Time
}
type ConfigStore interface {
	Save(context.Context, ConfigVersion) error
	Latest(context.Context) (ConfigVersion, error)
	History(context.Context) []ConfigVersion
}
type MemoryConfigStore struct {
	mu    sync.RWMutex
	items []ConfigVersion
}

func NewMemoryConfigStore() *MemoryConfigStore { return &MemoryConfigStore{} }
func (s *MemoryConfigStore) Save(ctx context.Context, v ConfigVersion) error {
	if e := ctx.Err(); e != nil {
		return e
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items = append(s.items, v)
	return nil
}
func (s *MemoryConfigStore) Latest(ctx context.Context) (ConfigVersion, error) {
	if e := ctx.Err(); e != nil {
		return ConfigVersion{}, e
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	if len(s.items) == 0 {
		return ConfigVersion{}, nil
	}
	return s.items[len(s.items)-1], nil
}
func (s *MemoryConfigStore) History(ctx context.Context) []ConfigVersion {
	_ = ctx
	s.mu.RLock()
	defer s.mu.RUnlock()
	o := make([]ConfigVersion, len(s.items))
	copy(o, s.items)
	return o
}
