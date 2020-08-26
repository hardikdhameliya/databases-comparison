package pudge

import (
	"encoding/json"

	"github.com/philippgille/gokv/encoding"
	"github.com/philippgille/gokv/util"
	"github.com/recoilme/pudge"
)

// Store is a gokv.Store implementation for pudge.
type Store struct {
	Db    *pudge.Db
	Codec encoding.Codec
}

// Set stores the given value for the given key.
func (s Store) Set(k string, v interface{}) error {
	if err := util.CheckKeyAndValue(k, v); err != nil {
		return err
	}
	var data []byte
	var err error
	if data, err = json.Marshal(v); err != nil {
		return err
	}
	err = s.Db.Set(k, data)
	return err

}

// Get retrieves the stored value for the given key.
func (s Store) Get(k string, v interface{}) (found bool, err error) {
	if err := util.CheckKeyAndValue(k, v); err != nil {
		return false, err
	}

	var data []byte
	if err := s.Db.Get(k, &data); err == pudge.ErrKeyNotFound {
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

	return s.Db.Delete(k)
}

// Close closes the store.
func (s Store) Close() error {
	return s.Db.Close()
}

// Options are the options for the pudge store.
type Options struct {
	// Config records params for creating DB object.
	Config *pudge.Config
	// File represents Open the database.
	File string
	// Encoding format.
	Codec encoding.Codec
}

// DefaultOptions is an Options object with default values.
var DefaultOptions = Options{
	Config: &pudge.Config{
		SyncInterval: 0,
		FileMode:     0666,
		DirMode:      0777,
	},
	File:  "./db",
	Codec: encoding.JSON,
}

// NewStore creates a pudge store.
func NewStore(options *Options) (Store, error) {
	if options == nil {
		options = &DefaultOptions
	}
	db, err := pudge.Open(options.File, options.Config)
	if err != nil {
		return Store{}, err
	}
	result := Store{
		Db:    db,
		Codec: options.Codec,
	}
	return result, nil
}
