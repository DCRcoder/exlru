package exlru

import (
	"fmt"
	"testing"
	"time"
)

func TestBaseGet(t *testing.T) {
	lru := NewExLru(10, 10 * time.Second)
	lru.Add("aa", "ab")
	_, ok := lru.Get("aa")
	time.Sleep(5 * time.Second)
	_, ok = lru.Get("aa")
	fmt.Println(ok)
	time.Sleep(5 * time.Second)
	_, ok = lru.Get("aa")
	fmt.Println(ok)
}