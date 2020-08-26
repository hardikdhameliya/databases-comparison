package ristretto

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
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
	s, err := NewStore(nil)
	if err != nil {
		panic(err)
	}
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("sen%d", i)
		s.Set(key, NS)
	}
}

func BenchmarkGet(b *testing.B) {
	d := 1000
	s, err := createStoreAndWriteNItems(d)
	if err != nil {
		panic(err)
	}
	defer s.Close()

	newdata := new(NetworkStats)
	for i := 0; i < b.N; i++ {
		j := rand.Intn(d)
		k := fmt.Sprintf("sen%d", j)
		if f, _ := s.Get(k, newdata); f != true {
			fmt.Printf("Can not read data for the key:%v\n", k)
		}
	}
}
func BenchmarkDelete(b *testing.B) {
	d := 1000
	s, err := createStoreAndWriteNItems(d)
	if err != nil {
		panic(err)
	}
	defer s.Close()

	for i := 0; i < b.N; i++ {
		j := rand.Intn(d)
		k := fmt.Sprintf("sen%d", j)
		s.Delete(k)
	}
}

func createStoreAndWriteNItems(items int) (Store, error) {
	s, err := NewStore(nil)
	if err != nil {
		return s, err
	}

	for i := 0; i < items; i++ {
		k := fmt.Sprintf("sen%d", i)
		s.Set(k, NS)
	}
	return s, nil
}
