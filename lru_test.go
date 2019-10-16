package exlru

import (
	"testing"
	"time"
)

func TestBaseExpireGet(t *testing.T) {
	lru := NewExLru(10)
	lru.Add("bb", "zouzuz")
	lru.AddWithExpire("aa", "ab", 10 * time.Second)
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

	_, ok = lru.Get("bb")
	if !ok {
		t.Errorf("expired error")
	}
}

func TestLengthLimit(t *testing.T) {
	lru := NewExLru(2)
	lru.AddWithExpire("key1", 1, 10 * time.Second)
	lru.AddWithExpire("key2", "b", 10 * time.Second)
	v, ok := lru.Get("key1")
	if !ok || v != 1 {
		t.Errorf("get key error")
	}
	v, ok = lru.Get("key2")
	if !ok || v != "b" {
		t.Errorf("get key error")
	}
	l := []string{"1", "3"}
	lru.AddWithExpire("key3", l,  10 * time.Second)
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
