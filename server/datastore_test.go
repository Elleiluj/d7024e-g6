package server

import (
	"fmt"
	"testing"
	"time"
)

func TestDatastore_print(t *testing.T) {
	fmt.Print("\n--------------------\n datastore.go\n--------------------\n")
}

func TestNewDataStore(t *testing.T) {
	datastore := NewDataStore()
	fail := false
	if datastore == nil {
		fail = true
		t.Error("NewDataStore() returned nil")
	}
	if datastore.Data == nil {
		fail = true
		t.Error("Data map is nil")
	}
	if datastore.TTL == nil {
		fail = true
		t.Error("TTL map is nil")
	}
	if !fail {
		fmt.Println("NewDatastore \tPASS")
	}
}

func TestAddData(t *testing.T) {
	datastore := NewDataStore()
	key := CreateHash("test")
	data := []byte("test")
	datastore.addData(data, key)
	retrievedData := datastore.getData(key)
	if string(retrievedData) != string(data) {
		t.Errorf("Expected data: %s, but got: %s", string(data), string(retrievedData))
	} else {
		fmt.Println("AddData \tPASS")
	}
}

func TestAddDataWithTTL(t *testing.T) {
	fail := false
	datastore := NewDataStore()
	key := "test"
	data := []byte("test")
	TTL := 2 * time.Second
	datastore.addDataWithTTL(data, key, TTL)
	retrievedData := datastore.getData(key)
	retrievedTTL := datastore.TTL[key]
	if string(retrievedData) != string(data) {
		t.Errorf("Expected data: %s, but got: %s", string(data), string(retrievedData))
		fail = true
	}
	if retrievedTTL.Before(time.Now()) {
		t.Errorf("Expected time to be after now, but was before")
		fail = true
	}
	if !fail {
		fmt.Println("AddDataWithTTL \tPASS")
	}
}

func TestGetData(t *testing.T) {
	fail := false
	datastore := NewDataStore()
	key := "test_key"
	data := []byte("test_data")
	datastore.addData(data, key)
	retrievedData := datastore.getData(key)
	if string(retrievedData) != string(data) {
		t.Errorf("Expected data: %s, but got: %s", string(data), string(retrievedData))
		fail = true
	}
	nonExistentData := datastore.getData("non_existent_key")
	if nonExistentData != nil {
		t.Errorf("Expected nil for non-existent key, but got: %s", string(nonExistentData))
		fail = true
	}
	if !fail {
		fmt.Println("GetData \tPASS")
	}
}

func TestResetTTL(t *testing.T) {
	datastore := NewDataStore()
	key := "test_key"
	data := []byte("test_data")
	datastore.addData(data, key)
	datastore.resetTTL(key)
	ttl := datastore.TTL[key]
	now := time.Now()
	if !ttl.After(now) {
		t.Errorf("Expected TTL to be after the current time, but it is not")
	} else {
		fmt.Println("ResetTTL \tPASS")
	}
}

func TestRemoveExpired(t *testing.T) {
	fail := false
	datastore := NewDataStore()
	data1 := []byte("data1")
	key1 := "key1"
	datastore.addDataWithTTL(data1, key1, 1*time.Second)

	data2 := []byte("data2")
	key2 := "key2"
	datastore.addDataWithTTL(data2, key2, 2*time.Second)

	data3 := []byte("data3")
	key3 := "key3"
	datastore.addDataWithTTL(data3, key3, 3*time.Second)

	time.Sleep(4 * time.Second)

	datastore.removeExpired()

	if datastore.getData(key1) != nil {
		fail = true
		t.Errorf("Expected data with key1 to be removed, but it is still present")
	}
	if datastore.getData(key2) != nil {
		fail = true
		t.Errorf("Expected data with key2 to be removed, but it is still present")
	}
	if datastore.getData(key3) != nil {
		fail = true
		t.Errorf("Expected data with key3 to be removed, but it is still present")
	}
	if !fail {
		fmt.Println("RemoveExpired \tPASS")
	}
}
