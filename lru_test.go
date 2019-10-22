package exlru

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestBaseExpireGet(t *testing.T) {
	lru := NewExLru(10)
	lru.Add("bb", "zouzuz")
	lru.AddWithExpire("aa", "ab", 10*time.Second)
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
	lru.AddWithExpire("key1", 1, 10*time.Second)
	lru.AddWithExpire("key2", "b", 10*time.Second)
	v, ok := lru.Get("key1")
	if !ok || v != 1 {
		t.Errorf("get key error")
	}
	v, ok = lru.Get("key2")
	if !ok || v != "b" {
		t.Errorf("get key error")
	}
	l := []string{"1", "3"}
	lru.AddWithExpire("key3", l, 10*time.Second)
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

func TestMemCache(t *testing.T) {
	m := NewMemCache()
	a := func(ctx context.Context) (interface{}, error) {
		return time.Now(), nil
	}
	expire := 3 * time.Second
	c, _ := m.Execute(context.TODO(), "zou", "a", a, 10, &expire)
	fmt.Println(c)
	time.Sleep(time.Second)
	c2, _ := m.Execute(context.TODO(), "zou", "a", a, 10, &expire)
	fmt.Println(c2)
	if c != c2 {
		t.Errorf("lru not work %v", c)
	}
	time.Sleep(3 * time.Second)
	c3, _ := m.Execute(context.TODO(), "zou", "a", a, 10, &expire)
	if c == c3 {
		t.Errorf("lru not work %v", c)
	}

	type person struct {
		name string
		d    time.Time
	}

	b := func(ctx context.Context) (interface{}, error) {
		return &person{name: "kaka", d: time.Now()}, nil
	}
	d, _ := m.Execute(context.TODO(), "stuct", "ab", b, 10, &expire)
	dv1 := d.(*person).d
	time.Sleep(time.Second)
	d1, _ := m.Execute(context.TODO(), "stuct", "ab", b, 10, &expire)
	dv2 := d1.(*person).d
	if dv1 != dv2 {
		t.Errorf("lru not work %v, %v", d, d1)
	}
	time.Sleep(3 * time.Second)
	d2, _ := m.Execute(context.TODO(), "stuct", "ab", b, 10, &expire)
	dv3 := d2.(*person).d
	if dv1 == dv3 {
		t.Errorf("lru not work %v, %v", d, d2)
	}
}

func TestMultiExCache(t *testing.T) {
	lru := NewExLru(10)

	go func() {
		for x := 1; x < 100; x++ {
			lru.Add("aa", "bb")
		}
	}()

	go func() {
		for y := 1; y <= 200; y++ {
			lru.Add("bb", "cc")
		}
	}()

	go func() {
		for y := 1; y <= 200; y++ {
			lru.Get("aa")
		}
	}()
	time.Sleep(2 * time.Second)
}

func TestMultiMemCache(t *testing.T) {
	m := NewMemCache()
	a := func(ctx context.Context) (interface{}, error) { fmt.Printf("hello world"); return nil, nil }
	go func() {
		for x := 1; x < 100; x++ {
			_, _ = m.Execute(context.TODO(), "aa", "bb", a, 10, nil)
		}
	}()
	go func() {
		for y := 1; y <= 200; y++ {
			_, _ = m.Execute(context.TODO(), "bb", "cc", a, 10, nil)
		}
	}()

	go func() {
		for y := 1; y <= 200; y++ {
			_, _ = m.Execute(context.TODO(), "bb", "dd", a, 10, nil)
		}
	}()
	time.Sleep(2 * time.Second)
}
