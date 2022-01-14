package lru

import (
	"fmt"
	"testing"
)

type String string

func (d String) Len() int {
	return len(d)
}

func TestOnRemove(t *testing.T) {
	keys := make([]string, 0)
	remover := func(key string, value Value) {
		fmt.Printf("remove %s: %s \n", key, value)
		keys = append(keys, key)
	}
	lru := InitLru(3, 0, remover)

	fmt.Println("k1")
	lru.Set("k1", String("123")) // 8

	fmt.Println("k2")
	lru.Set("k2", String("222")) // 8+5=13

	fmt.Println("k3")
	lru.Set("k3", String("333")) // 13+7=20, 20-8=12

	fmt.Println("k4")
	lru.Set("k4", String("444")) // 12+5=17, 17-5=12

	//expect := []string{"key1", "k2"}

	// only return type's size
	//fmt.Println(unsafe.Sizeof(lru.cache))

	//if !reflect.DeepEqual(expect, keys) {
	//	fmt.Println(keys)
	//	t.Fatalf("Call OnRemove failed, expect keys equals to %s", expect)
	//}
}

//func TestSet(t *testing.T) {
//	remover := func(key string, value Value) {
//		fmt.Printf("remove %s: %s \n", key, value)
//	}
//
//	lru := InitLru(1<<9, 1, remover)
//
//	// val 必需是一个interface
//	lru.Set("key", String("111"))
//	lru.Set("key1", String("111"))
//
//	//time.Sleep(2 * time.Second)
//	r, _ := lru.Request("key")
//	fmt.Println(r)
//}
