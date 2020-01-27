package cacher

// NilCache noop cache which does nothing useful
type NilCache struct{}

// NewNilCache creates a new noop cache Instance
func NewNilCache() Cacher {
	return &NilCache{}
}

// Load return nil
func (c NilCache) Load(v interface{}) error {
	return nil
}

// Store return nil
func (c NilCache) Store(v interface{}) error {
	return nil
}

// Clear return nil
func (c NilCache) Clear() error {
	return nil
}

// Expired return true that means cache is always expired
func (c NilCache) Expired() bool {
	return true
}
