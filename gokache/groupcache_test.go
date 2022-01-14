package gokache

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

type GroupRule struct {
	name string
	bsr  BackSourceRule
}

type BackSourceRule struct {
	url   string
	field string
}

type IT struct {
	Ret int64 `json:"ret"`
}

func TestGetter(t *testing.T) {
	baseUrl := "http://127.0.0.1:8001"
	// SourceBacker type:interface arg:string return:[]byte
	// cn: 实际使用中的Getter应该从数据库重新获取数据, 此处为了测试, 直接输出 []byte 类型
	var f SourceBacker = BsFunc(func(key string) ([]byte, error) {
		//return []byte(key), nil
		return HttpGet(baseUrl + "?key=" + key), nil
	})

	// cn: 将变量var转换成[]byte类型，并赋值给value
	// [107 101 121]

	v, _ := f.Request("key")

	fmt.Println("v:")
	fmt.Println(string(v))

	fmt.Println(fetchCode(v))
}

func fetchCode(res []byte) int64 {
	var tt IT
	err := json.Unmarshal(res, &tt)
	if err != nil {
		fmt.Println("json err:", err)
	}

	return tt.Ret
}

func TestGet(t *testing.T) {
	//bsr := BackSourceRule{
	//	url:   "http://127.0.0.1:8001",
	//	field: "data.version",
	//}
	//
	//gr := GroupRule{
	//	name: "user",
	//	bsr: bsr,
	//}

	// 创建一个map[string:int] 长度3
	// loadCounts 是一个map, 并不是一个数字
	// db 为模拟数据库
	// loadCounts 的作用记录从数据库重新加载数据的次数
	loadCounts := make(map[string]int, len(db))

	// NewGroup(name string, capacity int64, onBackSource SourceBacker) *Group
	// 创建一个命名空间为 scores 的, cache group; 空间大小为 2048字节; 并设置回调函数
	gc := NewGroup("scores", 10000, BsFunc(
		// cn: 回调函数(callback)，在缓存不存在时，调用这个函数，得到源数据
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)

			if v, ok := db[key]; ok {
				if _, ok := loadCounts[key]; !ok {
					loadCounts[key] = 0
				}
				loadCounts[key]++
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))

	for k, v := range db {
		if view, err := gc.Get(k); err != nil || view.String() != v {
			t.Fatal("failed to get value of Tom")
		}
		if _, err := gc.Get(k); err != nil || loadCounts[k] > 1 {
			t.Fatalf("cache %s miss", k)
		}
	}

	if view, err := gc.Get("unknown"); err == nil {
		t.Fatalf("the value of unknow should be empty, but %s got", view)
	}

	fmt.Println(loadCounts)
}

//func TestGetGroup(t *testing.T) {
//	groupName := "scores"
//
//	// 初始化 scores 组, 定义的回调函数什么也没返回 ([], nil)
//	NewGroup(groupName, 1<<20, BsFunc(
//		func(key string) (bytes []byte, err error) { return }))
//
//	if group := GetGroup(groupName); group == nil || group.name != groupName {
//		t.Fatalf("group %s not exist", groupName)
//	}
//
//	if group := GetGroup(groupName + "111"); group != nil {
//		t.Fatalf("expect nil, but %s got", group.name)
//	}
//}
