package wmap

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func check(b bool) {
	if b == false {
		panic("error")
	}
}

func Test1(t *testing.T) {
	wm := New[int, int]()
	wm.Put(45, 34)
	fmt.Println("Put end")
	fmt.Println(wm.Get(45))
	fmt.Println("Get end")
	fmt.Println(wm.Del(45))
	fmt.Println("Del end")
	fmt.Println(wm.Get(45))
	fmt.Println("Get 2 end")
}

func Test2(t *testing.T) {
	seed := time.Now().UnixMilli()
	rand.Seed(seed)
	fmt.Println(seed)
	m := make(map[int]int, 0)
	wm := New[int, int]()
	for i := 0; i < 1000000; i += 1 {
		k := rand.Int()
		v := rand.Int()
		wm.Put(k, v)
		m[k] = v
	}
	for k, v := range m {
		v1, ve := wm.Get(k)
		if ve != nil {
			panic(ve)
		}
		if v1 != v {
			fmt.Println(k, v, v1)
			panic("v1!=v")
		}
	}
	dn := 0
	dm := make(map[int]int, 0)
	for k, v := range m {
		dn += 1
		de := wm.Del(k)
		if de != nil {
			panic(de)
		}
		dm[k] = v
		if dn >= len(m)/2 {
			break
		}
	}
	for k, v := range dm {
		v1, ve := wm.Get(k)
		if ve == nil {
			fmt.Println(k, v, v1)
			panic(ve)
		}
	}
	for k, v := range m {
		_, dgb := dm[k]
		if dgb {
			continue
		}
		v1, ve := wm.Get(k)
		if ve != nil {
			panic(ve)
		}
		if v1 != v {
			fmt.Println(k, v, v1)
			panic("v1!=v")
		}
	}

}
