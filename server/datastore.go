package server

import "time"

const TTLconst = 60 * time.Second

type Datastore struct {
	Data map[string][]byte
	TTL  map[string]time.Time
}

func NewDataStore() *Datastore {
	datastore := &Datastore{Data: make(map[string][]byte), TTL: make(map[string]time.Time)}
	return datastore
}

func (datastore *Datastore) addData() {

}

func (datastore *Datastore) setTTL() {

}
