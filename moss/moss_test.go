package moss

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/couchbase/moss"
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
	buf, _ := json.Marshal(NS)
	writeOptions := moss.WriteOptions{}
	m, _ := moss.NewCollection(moss.CollectionOptions{})
	m.Start()
	for i := 0; i < b.N; i++ {
		batch, _ := m.NewBatch(0, 0)
		key := fmt.Sprintf("sen%d", i)
		batch.Set([]byte(key), buf)
		m.ExecuteBatch(batch, writeOptions)
		batch.Close()
	}

	m.Close()

}

// Due to limitation of benchmark capabilities in moss implementation, couldn't get the benchmark result of
// Get operation for in-memory datastore.
func BenchmarkGet(b *testing.B) {
	tmpDir, _ := ioutil.TempDir("", "benchStore")
	defer os.RemoveAll(tmpDir)

	store, coll, keys := createStoreAndWriteNItems(tmpDir, 10000, 100)
	defer store.Close()
	defer coll.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ss, err := coll.Snapshot()
		if err != nil || ss == nil {
			panic("Snapshot() failed!")
		}
		_, err = ss.Get(keys[i%len(keys)], moss.ReadOptions{})
		if err != nil {
			panic("Snapshot-Get() failed!")
		}

		ss.Close()
	}
}

// Due to limitation of benchmark capabilities in moss implementation, couldn't get the benchmark result of
// Delete operation for in-memory datastore.
func BenchmarkDelete(b *testing.B) {

	tmpDir, _ := ioutil.TempDir("", "benchStore")
	defer os.RemoveAll(tmpDir)

	store, coll, keys := createStoreAndWriteNItems(tmpDir, 10000, 100)
	defer store.Close()
	defer coll.Close()

	b.ResetTimer()
	batch, err := coll.NewBatch(0, 0)
	if err != nil {
		panic("NewBatch() failed!")
	}
	for i := 0; i < b.N; i++ {
		err := batch.Del(keys[i%len(keys)])
		if err != nil {
			panic("Snapshot-Get() failed!")
		}
	}

}

func createStoreAndWriteNItems(tmpDir string, items int,
	batches int) (s *moss.Store, c moss.Collection, ks [][]byte) {

	store, coll, err := moss.OpenStoreCollection(tmpDir,
		moss.StoreOptions{},
		moss.StorePersistOptions{})
	if err != nil || store == nil {
		panic("OpenStoreCollection() failed!")
	}

	keys := make([][]byte, items)

	if batches > items {
		batches = 1
	}
	itemsPerBatch := items / batches
	itemCount := 0
	val := NS
	for i := 0; i < batches; i++ {
		if itemsPerBatch > items-itemCount {
			itemsPerBatch = items - itemCount
		}

		if itemsPerBatch <= 0 {
			break
		}

		batch, err := coll.NewBatch(itemsPerBatch, itemsPerBatch*20)
		if err != nil {
			panic("NewBatch() failed!")
		}

		for j := 0; j < itemsPerBatch; j++ {
			k := []byte(fmt.Sprintf("sen%d", i))
			sid := fmt.Sprintf("sen%d", i)
			val.SensorID = sid

			v, _ := json.Marshal(val)
			itemCount++

			batch.Set(k, v)
			keys[j] = k
		}

		err = coll.ExecuteBatch(batch, moss.WriteOptions{})
		if err != nil {
			panic("ExecuteBatch() failed!")
		}
	}

	return store, coll, keys
}
