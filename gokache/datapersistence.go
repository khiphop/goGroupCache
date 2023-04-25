package gokache

import (
	"bufio"
	"encoding/json"
	"fmt"
	"gokache/lru"
	"io"
	"os"
)

type Storage struct {
	Group string `json:"g"`
	Key   string `json:"k"`
	Value string `json:"v"`
}

var (
	dataFile = "./runtime/data.txt"
)

// DataRestore :set cache from file
func DataRestore(nd *NodeDispatch) error {
	fi, err := os.Open(dataFile)

	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return nil
	}
	defer fi.Close()

	br := bufio.NewReader(fi)
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}

		if string(a) == "" {
			continue
		}

		obj := readObj(a)
		err := nd.SetHandler(obj.Group, obj.Key, obj.Value, true)
		if err != nil {
			return err
		}

		fmt.Println(string(a))
	}

	// clear data storage file
	err = os.Truncate(dataFile, 0)
	if err != nil {
		return err
	}

	return nil
}

func readObj(res []byte) Storage {
	var d Storage

	err := json.Unmarshal(res, &d)
	if err != nil {
		fmt.Println("json err:", err)
	}

	return d
}

func BackupData() error {
	f, err := os.OpenFile(dataFile, os.O_WRONLY|os.O_TRUNC, 0600)
	defer f.Close()

	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	g := groups
	for groupName, gMap := range g {
		fmt.Println(groupName)

		for e := gMap.coreCache.lru.Ll.Front(); e != nil; e = e.Next() {
			kv := e.Value.(*lru.KvMap)
			fmt.Println(kv.Key)
			fmt.Println(kv.Value.(ByteView).String())

			m := map[string]string{"g": groupName, "k": kv.Key, "v": kv.Value.(ByteView).String()}
			mJson, _ := json.Marshal(m)
			mString := string(mJson) + "\n"
			_, err = f.Write([]byte(mString))
		}
	}

	return nil
}
