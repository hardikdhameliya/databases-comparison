package ristretto

import (
	"encoding/json"
	"fmt"

	"github.com/dgraph-io/ristretto"
	"github.com/philippgille/gokv/encoding"
	"github.com/philippgille/gokv/util"
)

// Store is a gokv.Store implementation for ristretto.
type Store struct {
	Db    *ristretto.Cache
	Codec encoding.Codec
}

// Set stores the given value for the given key.
func (s Store) Set(k string, v interface{}) error {
	if err := util.CheckKeyAndValue(k, v); err != nil {
		return err
	}
	s.Db.Set(k, v, 1)
	return nil

}

// Get retrieves the stored value for the given key.
func (s Store) Get(k string, v interface{}) (found bool, err error) {
	if err := util.CheckKeyAndValue(k, v); err != nil {
		return false, err
	}
	data, found := s.Db.Get(k)
	if !found || data == nil {
		return false, fmt.Errorf("value is missing")
	}
	b, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	return true, s.Codec.Unmarshal(b, v)
}

// Delete deletes the stored value for the given key.
func (s Store) Delete(k string) error {
	if err := util.CheckKey(k); err != nil {
		return err
	}
	s.Db.Del(k)
	return nil
}

// Close closes the store.
func (s Store) Close() error {
	s.Db.Close()
	return nil
}

// Options are the options for the ristretto store.
type Options struct {
	// Config records params for creating DB object.
	Config ristretto.Config
	// Encoding format.
	Codec encoding.Codec
}

// DefaultOptions is an Options object with default values.
var DefaultOptions = Options{
	Config: ristretto.Config{
		// number of keys to track frequency of (10M).
		NumCounters: 1e7,
		// maximum cost of cache (1GB).
		MaxCost: 1 << 30,
		// number of keys per Get buffer.
		BufferItems: 1024,
	},
	Codec: encoding.JSON,
}

// NewStore creates a ristretto store.
func NewStore(options *Options) (Store, error) {
	if options == nil {
		options = &DefaultOptions
	}
	cache, err := ristretto.NewCache(&options.Config)
	if err != nil {
		return Store{}, err
	}
	result := Store{
		Codec: options.Codec,
		Db:    cache,
	}
	return result, nil
}
