package nutsdb

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/philippgille/gokv/encoding"
	"github.com/xujiajun/nutsdb"
)

type NetworkStats struct {
	SensorID   string           `json:"sensor_id"`
	Updated    time.Time        `json:"updated"`
	Interfaces []InterfaceStats `json:"interfaces"`
}
type InterfaceStats struct {
	Interface string `json:"interface_name"`
	TxBytes   int64  `json:"tx_bytes"`
	TxPackets int64  `json:"tx_packets"`
	TxErrors  int64  `json:"tx_errors"`
	RxBytes   int64  `json:"rx_bytes"`
	RxPackets int64  `json:"rx_packets"`
	RxErrors  int64  `json:"rx_errors"`
}

var NS = NetworkStats{
	SensorID: "ses1",
	Updated:  time.Now(),
	Interfaces: []InterfaceStats{
		{
			Interface: "eth0",
			TxBytes:   123,
			TxPackets: 345,
			TxErrors:  234,
			RxBytes:   566,
			RxPackets: 12,
			RxErrors:  12,
		},
	},
}

func BenchmarkSet(b *testing.B) {
	tmpDir, _ := ioutil.TempDir("", "store")
	defer os.RemoveAll(tmpDir)
	ops := &Options{
		Config: nutsdb.DefaultOptions,
		Dir:    tmpDir,
		Bucket: "gigamon",
		Codec:  encoding.JSON,
	}

	s, err := NewStore(ops)
	if err != nil {
		panic(err)
	}
	for i := 0; i < b.N; i++ {
		k := fmt.Sprintf("sen%d", i)
		err := s.Set(k, NS)
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkGet(b *testing.B) {
	tmpDir, _ := ioutil.TempDir("", "store")
	defer os.RemoveAll(tmpDir)
	ops := &Options{
		Config: nutsdb.DefaultOptions,
		Dir:    tmpDir,
		Bucket: "gigamon",
		Codec:  encoding.JSON,
	}

	d := 1000
	s, err := createStoreAndWriteNItems(ops, d)
	if err != nil {
		panic(err)
	}
	defer s.Close()

	newdata := new(NetworkStats)
	for i := 0; i < b.N; i++ {
		j := rand.Intn(d)
		k := fmt.Sprintf("sen%d", j)
		if f, err := s.Get(k, newdata); !f {
			fmt.Printf("Can not read data for the key:%v\n", k)
		} else if err != nil {
			panic(err)
		}
	}
}

func BenchmarkDelete(b *testing.B) {
	tmpDir, _ := ioutil.TempDir("", "store")
	defer os.RemoveAll(tmpDir)
	ops := &Options{
		Config: nutsdb.DefaultOptions,
		Dir:    tmpDir,
		Bucket: "gigamon",
		Codec:  encoding.JSON,
	}

	d := 1000
	s, err := createStoreAndWriteNItems(ops, d)
	if err != nil {
		panic(err)
	}
	defer s.Close()

	for i := 0; i < b.N; i++ {
		j := rand.Intn(d)
		k := fmt.Sprintf("sen%d", j)
		if err := s.Delete(k); err != nil {
			panic("Del() failed!")
		}
	}
}

func createStoreAndWriteNItems(option *Options, items int) (Store, error) {
	s, err := NewStore(option)
	if err != nil {
		return s, err
	}

	for i := 0; i < items; i++ {
		k := fmt.Sprintf("sen%d", i)
		s.Set(k, NS)
	}
	return s, nil
}
