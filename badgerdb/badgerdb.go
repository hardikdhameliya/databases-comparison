package badgerdb

import (
	"github.com/dgraph-io/badger"
	"github.com/philippgille/gokv/encoding"
	"github.com/philippgille/gokv/util"
)

// Store is a gokv.Store implementation for BadgerDB.
type Store struct {
	Db    *badger.DB
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

	err = s.Db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(k), data)
	})
	if err != nil {
		return err
	}
	return nil
}

// Get retrieves the stored value for the given key.
func (s Store) Get(k string, v interface{}) (found bool, err error) {
	if err := util.CheckKeyAndValue(k, v); err != nil {
		return false, err
	}

	var data []byte
	err = s.Db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(k))
		if err != nil {
			return err
		}
		data, err = item.ValueCopy(nil)
		if err != nil {
			return err
		}
		return nil
	})
	// If no value was found return false
	if err == badger.ErrKeyNotFound {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, s.Codec.Unmarshal(data, v)
}

// Delete deletes the stored value for the given key.
func (s Store) Delete(k string) error {
	if err := util.CheckKey(k); err != nil {
		return err
	}

	return s.Db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(k))
	})
}

// Close closes the store.
func (s Store) Close() error {
	return s.Db.Close()
}

// Options are the options for the BadgerDB store.
type Options struct {
	// Database file directory.
	Dir string
	// Encoding format.
	Codec encoding.Codec
}

// DefaultOptions is an Options object with default values.
var DefaultOptions = Options{
	Dir:   "BadgerDB",
	Codec: encoding.JSON,
}

// NewStore creates a new BadgerDB store.
func NewStore(options *Options) (Store, error) {

	if options == nil {
		options = &DefaultOptions
	}

	opts := badger.DefaultOptions(options.Dir)
	db, err := badger.Open(opts)
	if err != nil {
		return Store{}, err
	}
	result := Store{
		Db:    db,
		Codec: options.Codec,
	}

	return result, nil
}
