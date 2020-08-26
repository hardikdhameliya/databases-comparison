package nutsdb

import (
	"github.com/philippgille/gokv/encoding"
	"github.com/philippgille/gokv/util"
	"github.com/xujiajun/nutsdb"
)

// Store is a gokv.Store implementation for nutsdb.
type Store struct {
	Db     *nutsdb.DB
	Bucket string
	Codec  encoding.Codec
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

	err = s.Db.Update(func(tx *nutsdb.Tx) error {
		return tx.Put(s.Bucket, []byte(k), data, nutsdb.Persistent)
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
	err = s.Db.View(func(tx *nutsdb.Tx) error {
		item, err := tx.Get(s.Bucket, []byte(k))
		if err != nil {
			return err
		}
		data = item.Value
		return nil
	})
	// If no value was found return false
	if err == nutsdb.ErrKeyNotFound {
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

	return s.Db.Update(func(tx *nutsdb.Tx) error {
		return tx.Delete(s.Bucket, []byte(k))
	})
}

// Close closes the store.
func (s Store) Close() error {
	s.Db.Close()
	return nil
}

// Options are the options for the nutsdb store.
type Options struct {
	// Config records params for creating DB object.
	Config nutsdb.Options
	// Dir represents Open the database located in which dir.
	Dir string
	// A bucket is a collection of unique keys that are associated with values.
	Bucket string
	// Encoding format.
	Codec encoding.Codec
}

// DefaultOptions is an Options object with default values.
var DefaultOptions = Options{
	Config: nutsdb.DefaultOptions,
	Dir:    ".",
	Bucket: "gigamon",
	Codec:  encoding.JSON,
}

// NewStore creates a nutsdb store.
func NewStore(options *Options) (Store, error) {
	if options == nil {
		options = &DefaultOptions
	}
	options.Config.Dir = options.Dir
	db, err := nutsdb.Open(options.Config)
	if err != nil {
		return Store{}, err
	}

	result := Store{
		Db:     db,
		Bucket: options.Bucket,
		Codec:  options.Codec,
	}

	return result, nil
}
