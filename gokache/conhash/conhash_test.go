package conhash

import (
	"fmt"
	"testing"
)

func TestHashing(t *testing.T) {
	conHash := InitConHash(5, nil)

	// Given the above hashHandler function, this will give vNodeC with "hashHandlers":
	conHash.Add("127.0.0.1:7000", "127.0.0.1:7001", "127.0.0.1:7002")
	fmt.Println(conHash.HashMap)
	fmt.Println(conHash.Get("name"))

	conHash.Add("127.0.0.1:7003")
	fmt.Println(conHash.HashMap)

	fmt.Println(conHash.Get("name"))
}
