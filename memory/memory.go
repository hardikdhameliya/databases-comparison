package memory

import (
	"fmt"
	"sync"

	"github.com/philippgille/gokv/encoding"
	"github.com/philippgille/gokv/util"
)

// Store is a gokv.Store implementation for In-memory storage.
type Store struct {
	Db    map[string][]byte
	Codec encoding.Codec
	mu    *sync.RWMutex
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

	s.mu.Lock()
	s.Db[k] = data
	s.mu.Unlock()

	return nil

}

// Get retrieves the stored value for the given key.
func (s Store) Get(k string, v interface{}) (found bool, err error) {
	if err := util.CheckKeyAndValue(k, v); err != nil {
		return false, err
	}

	s.mu.RLock()
	data, ok := s.Db[k]
	s.mu.Unlock()
	if !ok {
		return false, nil
	}

	return true, s.Codec.Unmarshal(data, v)
}

// Delete deletes the stored value for the given key.
func (s Store) Delete(k string) error {
	if err := util.CheckKey(k); err != nil {
		return err
	}
	if _, ok := s.Db[k]; ok {
		s.mu.Lock()
		delete(s.Db, k)
		s.mu.Unlock()
		return nil

	}
	return fmt.Errorf("Key not found: %v", k)

}

// Close closes the store.
func (s Store) Close() error {
	s.mu.Lock()
	s.Db = nil
	s.mu.Unlock()
	return nil
}

// Options are the options for the In-memory store.
type Options struct {
	// Encoding format.
	Codec encoding.Codec
}

// DefaultOptions is an Options object with default values.
var DefaultOptions = Options{
	Codec: encoding.JSON,
}

// NewStore creates a In-memory store.
func NewStore(options *Options) (Store, error) {
	if options == nil {
		options = &DefaultOptions
	}
	db := make(map[string][]byte)
	result := Store{
		Db:    db,
		Codec: options.Codec,
		mu:    &sync.RWMutex{},
	}
	return result, nil

}
