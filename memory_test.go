package cache

import (
	"testing"
	"time"
)

func TestMemoryAdapter_Instance(t *testing.T) {
	adp := MemoryAdapter{}
	_, err := adp.Instance("")
	if err != nil {
		t.Fatal("impossible")
	}
}

func TestMemoryAdapterInstance_GetSetDel(t *testing.T) {
	adp := MemoryAdapter{}
	adi, err := adp.Instance("")
	if err != nil {
		t.Fatal("impossible")
	}
	var val string
	// get
	val, err = adi.Get("key1")
	if err != ErrKeyNotFound {
		t.Fatal("failed")
	}
	// set
	err = adi.Set("key1", "val1", 1)
	if err != nil {
		t.Fatal("failed")
	}
	// get
	val, err = adi.Get("key1")
	if err != nil {
		t.Fatal("failed")
	}
	if val != "val1" {
		t.Fatal("failed")
	}
	// expires
	time.Sleep(time.Second * 2)
	val, err = adi.Get("key1")
	if err != ErrKeyNotFound {
		t.Fatal("failed")
	}
	// set
	err = adi.Set("key1", "val1", 1)
	if err != nil {
		t.Fatal("failed")
	}
	// get
	val, err = adi.Get("key1")
	if err != nil {
		t.Fatal("failed")
	}
	if val != "val1" {
		t.Fatal("failed")
	}
	// del
	err = adi.Del("key1")
	if err != nil {
		t.Fatal("failed")
	}
	// get
	val, err = adi.Get("key1")
	if err != ErrKeyNotFound {
		t.Fatal("failed")
	}
}
