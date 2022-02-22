package gokache

import (
	"encoding/json"
	"fmt"
	"gokache/conhash"
	"os"
	"sync"
)

type NodeDispatch struct {
	self  string
	mu    sync.Mutex
	Peers *conhash.Map
}

var (
	vNodeC = 5
)

// InitNode Set updates the pool's list of peers.
func InitNode(peer string) *NodeDispatch {
	nd := &NodeDispatch{
		self:  peer,
		Peers: conhash.InitConHash(vNodeC, nil),
	}

	nd.Peers.Add(peer)
	nd.cmd()

	return nd
}

func (nd *NodeDispatch) cmd() {
	for k, v := range os.Args {
		if k == 0 {
			continue
		}

		if v == "-dr" {
			DataRestore(nd)
		}
	}
}

// RegNode Set updates the pool's list of peers.
func (nd *NodeDispatch) RegNode(peer ...string) {
	nd.mu.Lock()
	defer nd.mu.Unlock()

	nd.Peers.Add(peer...)
}

// RemoveNode Set updates the pool's list of peers.
func (nd *NodeDispatch) RemoveNode(peer ...string) {
	nd.mu.Lock()
	defer nd.mu.Unlock()

	nd.Peers.Del(peer...)
}

// ChooseNode picks a peer according to key
func (nd *NodeDispatch) ChooseNode(key string) (string, bool) {
	nd.mu.Lock()
	defer nd.mu.Unlock()

	if peer := nd.Peers.Get(key); peer != "" && peer != nd.self {
		return peer, false
	}

	return nd.self, true
}

// DisplayNode picks a peer according to key
func (nd *NodeDispatch) DisplayNode() []byte {
	nd.mu.Lock()
	defer nd.mu.Unlock()

	mJson, _ := json.Marshal(nd.Peers.HashMap)
	return mJson
}

func (nd *NodeDispatch) GetHandler(group string, key string, inner bool) ([]byte, error) {
	var b []byte

	if node, ok := nd.ChooseNode(key); !ok && !inner {
		fmt.Println("Trigger InnerGet")

		b = InnerGet(node, group, key)
	} else {
		fmt.Println("Trigger LocalGet")

		gc := GetGroup(group)
		rs, _ := gc.Get(key)
		b = rs.ByteSlice()
	}

	return b, nil
}

func (nd *NodeDispatch) SetHandler(group string, key string, val string, inner bool) error {
	if node, ok := nd.ChooseNode(key); !ok && !inner {
		fmt.Println("Trigger InnerGet")
		InnerSet(node, group, key, val)
	} else {
		fmt.Println("Trigger LocalSet")

		gc := GetGroup(group)
		_ = gc.Set(key, val)
	}

	return nil
}
