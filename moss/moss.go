package moss

import (
	"github.com/couchbase/moss"
	"github.com/philippgille/gokv/encoding"
	"github.com/philippgille/gokv/util"
)

// Store is a gokv.Store implementation for moss.
type Store struct {
	Collection moss.Collection
	Batch      moss.Batch
	Codec      encoding.Codec
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
	if err := s.Batch.Set([]byte(k), data); err != nil {
		return err
	}
	if err := s.Collection.ExecuteBatch(s.Batch, moss.WriteOptions{}); err != nil {
		return err
	}
	return nil
}

// Get retrieves the stored value for the given key.
func (s Store) Get(k string, v interface{}) (found bool, err error) {
	if err := util.CheckKeyAndValue(k, v); err != nil {
		return false, err
	}
	ropts := moss.ReadOptions{}
	ss, err := s.Collection.Snapshot()
	defer ss.Close()
	data, err := ss.Get([]byte(k), ropts)
	if err != nil || data == nil {
		return false, err
	}
	return true, s.Codec.Unmarshal(data, v)
}

// Delete deletes the stored value for the given key.
func (s Store) Delete(k string) error {
	if err := util.CheckKey(k); err != nil {
		return err
	}
	if err := s.Batch.Del([]byte(k)); err != nil {
		return err
	}
	if err := s.Collection.ExecuteBatch(s.Batch, moss.WriteOptions{}); err != nil {
		return err
	}

	return nil
}

// Close closes the store.
func (s Store) Close() error {
	if err := s.Batch.Close(); err != nil {
		return err
	}
	return nil
}

// Options are the options for the moss store.
type Options struct {
	// Collection represents an ordered mapping of key-val entries.
	Collection moss.CollectionOptions
	// Encoding format.
	Codec encoding.Codec
}

// DefaultOptions is an Options object with default values.
var DefaultOptions = Options{
	Collection: moss.CollectionOptions{},
	Codec:      encoding.JSON,
}

// NewStore creates a moss store.
func NewStore(options *Options) (Store, error) {
	if options == nil {
		options = &DefaultOptions
	}

	col, err := moss.NewCollection(options.Collection)
	if err != nil {
		return Store{}, err
	}

	batch, err := col.NewBatch(0, 0)
	if err != nil {
		return Store{}, err
	}

	result := Store{
		Collection: col,
		Batch:      batch,
		Codec:      options.Codec,
	}

	return result, nil
}
