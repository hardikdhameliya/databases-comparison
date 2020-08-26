package badgerdb

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/philippgille/gokv/encoding"
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
		Dir:   tmpDir,
		Codec: encoding.JSON,
	}

	s, err := NewStore(ops)
	defer s.Close()
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
		Dir:   tmpDir,
		Codec: encoding.JSON,
	}

	d := 1000
	s, err := createStoreAndWriteNItems(ops, d)
	defer s.Close()
	if err != nil {
		panic(err)
	}
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
		Dir:   tmpDir,
		Codec: encoding.JSON,
	}

	d := 1000
	s, err := createStoreAndWriteNItems(ops, d)
	if err != nil {
		panic(err)
	}
	for i := 0; i < b.N; i++ {
		j := rand.Intn(d)
		k := fmt.Sprintf("sen%d", j)
		err := s.Delete(k)
		if err != nil {
			panic(err)
		}
	}
	s.Close()
}

func createStoreAndWriteNItems(option *Options, items int) (Store, error) {
	s, err := NewStore(option)
	if err != nil {
		return s, err
	}

	for i := 0; i < items; i++ {
		k := fmt.Sprintf("sen%d", i)
		if err := s.Set(k, NS); err != nil {
			panic(err)
		}
	}
	return s, nil

}
