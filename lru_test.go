package exlru

import (
	"testing"
	"time"
)

func TestBaseExpireGet(t *testing.T) {
	lru := NewExLru(10, 10 * time.Second)
	lru.Add("aa", "ab")
	_, ok := lru.Get("aa")
	if !ok {
		t.Errorf("get key not existed")
	}
	time.Sleep(5 * time.Second)
	_, ok = lru.Get("aa")
	if !ok {
		t.Errorf("key expired error")
	}
	time.Sleep(5 * time.Second)
	_, ok = lru.Get("aa")
	if ok {
		t.Errorf("key not expired")
	}
}


func TestLengthLimit(t *testing.T) {
	lru := NewExLru(2, 10 * time.Second)
	lru.Add("key1", 1)
	lru.Add("key2", "b")
	v, ok := lru.Get("key1")
	if !ok || v != 1 {
		t.Errorf("get key error")
	}
	v, ok = lru.Get("key2")
	if !ok || v != "b" {
		t.Errorf("get key error")
	}
	l := []string{"1", "3"}
	lru.Add("key3", l)
	_, ok = lru.Get("key1")
	if ok {
		t.Errorf("lru not work")
	}
	v, ok = lru.Get("key3")
	if !ok {
		t.Errorf("lru not work")
	}
	if v.([]string)[0] != l[0] {
		t.Errorf("lru not work")
	}
}

