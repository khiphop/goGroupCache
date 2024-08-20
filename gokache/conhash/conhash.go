package conhash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

// Hash :maps bytes to uint32
type Hash func(data []byte) uint32

// Map :contains all hashed hashRing
type Map struct {
	hashHandler Hash
	// virtual node's count
	virtualNodeC int
	// hash ring
	hashRing []int // Sorted
	// this is map if node and his hashVal. key: node's hashVal, val: node's name
	HashMap map[int]string
}

// InitConHash :creates a Map instance
// 构造函数 InitConHash() 允许自定义虚拟节点倍数和 Hash 函数
func InitConHash(vnCount int, fn Hash) *Map {
	// cn: 此时还未定义哈希环, 长度为0, 需要通过Add函数添加keys
	m := &Map{
		virtualNodeC: vnCount,
		hashHandler:  fn,
		HashMap:      make(map[int]string),
	}

	if m.hashHandler == nil {
		m.hashHandler = crc32.ChecksumIEEE
	}

	return m
}

// Add :adds some hashRing to the hashHandler.
func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.virtualNodeC; i++ {
			hashVal := int(m.hashHandler([]byte(key + "-" + strconv.Itoa(i))))
			m.hashRing = append(m.hashRing, hashVal)
			m.HashMap[hashVal] = key
		}
	}

	// slice ASC sort
	sort.Ints(m.hashRing)
}

func (m *Map) Del(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.virtualNodeC; i++ {
			hashVal := int(m.hashHandler([]byte(key + "-" + strconv.Itoa(i))))

			// del element from slice
			for i := 0; i < len(m.hashRing); {
				if m.hashRing[i] == hashVal {
					m.hashRing = append(m.hashRing[:i], m.hashRing[i+1:]...)
				} else {
					i++
				}
			}

			delete(m.HashMap, hashVal)
		}
	}

	// slice ASC sort
	sort.Ints(m.hashRing)
}

// Get :gets the closest item in the hashHandler to the provided key.
func (m *Map) Get(key string) string {
	if len(m.hashRing) == 0 {
		return ""
	}

	hash := int(m.hashHandler([]byte(key)))

	// Binary search for appropriate replica.
	// return the minimum i which make f(i)=true
	idx := sort.Search(len(m.hashRing), func(i int) bool {
		return m.hashRing[i] >= hash
	})

	return m.HashMap[m.hashRing[idx%len(m.hashRing)]]
}
