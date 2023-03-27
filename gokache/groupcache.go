package gokache

import (
	"fmt"
	"sync"
	"time"
)

// A Group is a cache namespace and associated data loaded spread over
type Group struct {
	name         string
	mu           sync.Mutex
	onBackSource SourceBacker
	coreCache    cache
	backSourceCd map[string]int64
	peers        PeerPicker
}

// A SourceBacker loads data for a key.
type SourceBacker interface {
	Request(key string) ([]byte, error)
}

// A BsFunc implements SourceBacker with a function.
type BsFunc func(key string) ([]byte, error)

// Request :implements SourceBacker interface function
// cn: 该 method 属于 BsFunc 类型对象中的方法
// cn: 因为该 method 实现了 GET 方法, 它自动属于 SourceBacker 类型
func (fun BsFunc) Request(key string) ([]byte, error) {
	// cn: 此时还未定义
	// cn: 将 key 转换为 BsFunc 类型
	return fun(key)
}

var (
	mu            sync.RWMutex
	groups        = make(map[string]*Group)
	backSourceCdS = 3
)

// NewGroup :create a new instance of Group
func NewGroup(name string, c int64, sourceBacker SourceBacker) *Group {
	if sourceBacker == nil {
		panic("nil SourceBacker")
	}

	mu.Lock()
	defer mu.Unlock()

	g := &Group{
		name:         name,
		onBackSource: sourceBacker,
		coreCache:    cache{capacity: c},
		backSourceCd: make(map[string]int64),
	}

	groups[name] = g

	return g
}

// RegisterPeers :registers a PeerPicker for choosing remote peer
func (g *Group) RegisterPeers(peers PeerPicker) {
	if g.peers != nil {
		panic("RegisterPeerPicker called more than once")
	}

	g.peers = peers
}

// GetGroup :returns the named group previously created with NewGroup, or nil if there's no such group.
func GetGroup(name string) *Group {
	mu.RLock()
	defer mu.RUnlock()

	g := groups[name]

	return g
}

// Get value for a key from cache
func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key is required")
	}

	if v, ok := g.coreCache.get(key); ok {
		return v, nil
	}

	return g.backSource(key)
}

func (g *Group) Set(key string, val string) error {
	if key == "" {
		return fmt.Errorf("key is required")
	}

	value := ByteView{b: cloneBytes([]byte(val))}

	g.coreCache.set(key, value)

	return nil
}

func (g *Group) backSource(key string) (value ByteView, err error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	// 复检数据是否存在 / recheck if val exists
	if v, ok := g.coreCache.get(key); ok {
		return v, nil
	}

	// 回源cd / back-source cool down
	if g.backSourceCd[key] > 0 && g.backSourceCd[key] > time.Now().Unix() {
		print("backSourceCd")
		return ByteView{}, nil
	}

	fmt.Println("run load")

	// cn: 从数据源获取数据设置到缓存中, 依赖回调函数 onBackSource
	bytes, err := g.onBackSource.Request(key)
	if err != nil {
		return ByteView{}, err
	}

	value = ByteView{b: cloneBytes(bytes)}

	g.coreCache.set(key, value)
	g.backSourceCd[key] = time.Now().Unix() + int64(backSourceCdS)

	return value, nil
}
