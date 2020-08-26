package bigcache

import (
	"time"

	"github.com/allegro/bigcache/v2"
	"github.com/philippgille/gokv/encoding"
	"github.com/philippgille/gokv/util"
)

// Store is a gokv.Store implementation for BigCache.
type Store struct {
	Db    *bigcache.BigCache
	Codec encoding.Codec
}

// Set stores the given value for the given key.
func (s Store) Set(k string, v interface{}) error {
	if err := util.CheckKeyAndValue(k, v); err != nil {
		return err
	}
	data, err := s.Codec.Marshal(v)
	if err != nil {
		return err
	}
	return s.Db.Set(k, data)
}

// Get retrieves the stored value for the given key.
func (s Store) Get(k string, v interface{}) (found bool, err error) {
	if err := util.CheckKeyAndValue(k, v); err != nil {
		return false, err
	}

	data, err := s.Db.Get(k)
	if err != nil {
		if err == bigcache.ErrEntryNotFound {
			return false, nil
		}
		return false, err
	}

	return true, s.Codec.Unmarshal(data, v)
}

// Delete deletes the stored value for the given key.
func (s Store) Delete(k string) error {
	if err := util.CheckKey(k); err != nil {
		return err
	}

	err := s.Db.Delete(k)
	if err != nil {
		if err == bigcache.ErrEntryNotFound {
			return nil
		}
		return err
	}
	return nil
}

// Close closes the store.
func (s Store) Close() error {
	return s.Db.Close()
}

// Options are the options for the BigCache store.
type Options struct {
	// The maximum size of the cache in MiB.
	// 0 means no limit.
	HardMaxCacheSize int
	// Time after which an entry can be evicted.
	// 0 means no eviction.
	Eviction time.Duration
	// Encoding format.
	Codec encoding.Codec
}

// DefaultOptions is an Options object with default values.
var DefaultOptions = Options{
	Codec: encoding.JSON,
}

// NewStore creates a BigCache store.
func NewStore(options *Options) (Store, error) {
	if options == nil {
		options = &DefaultOptions
	}

	config := bigcache.DefaultConfig(options.Eviction)
	cache, err := bigcache.NewBigCache(config)
	if err != nil {
		return Store{}, err
	}
	result := Store{
		Db:    cache,
		Codec: options.Codec,
	}

	return result, nil
}
