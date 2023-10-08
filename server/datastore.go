package server

import (
	"fmt"
	"sync"
	"time"
)

const TTLconst = 60 * time.Second

type Datastore struct {
	Data map[string][]byte
	TTL  map[string]time.Time
	mu   sync.Mutex
}

func NewDataStore() *Datastore {
	datastore := &Datastore{Data: make(map[string][]byte), TTL: make(map[string]time.Time), mu: sync.Mutex{}}
	return datastore
}

func (datastore *Datastore) addData(data []byte, key string) {
	datastore.mu.Lock()
	datastore.Data[key] = data
	datastore.mu.Unlock()
	datastore.resetTTL(key)
}

// for testing
func (datastore *Datastore) addDataWithTTL(data []byte, key string, TTL time.Duration) {
	datastore.mu.Lock()
	datastore.Data[key] = data
	now := time.Now()
	ttl := now.Add(TTL)
	datastore.TTL[key] = ttl
	datastore.mu.Unlock()
}

func (datastore *Datastore) getData(key string) []byte {
	return datastore.Data[key]
}

func (datastore *Datastore) resetTTL(key string) {
	datastore.mu.Lock()
	defer datastore.mu.Unlock()
	now := time.Now()
	ttl := now.Add(TTLconst)
	datastore.TTL[key] = ttl
}

func (datastore *Datastore) removeExpired() {
	datastore.mu.Lock()
	defer datastore.mu.Unlock()
	for key := range datastore.Data {
		ttl := datastore.TTL[key]
		if datastore.isExpired(ttl) {
			delete(datastore.Data, key)
			delete(datastore.TTL, key)
			fmt.Println("Deleted expired data with key: " + key)
		}
	}
}

func (datastore *Datastore) isExpired(TTL time.Time) bool {
	return time.Now().After(TTL)
}
