package gokache

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type InnerResp struct {
	Data string `json:"data"`
}

var (
	tempUrl = "http://127.0.0.1:8011/"
)

func InnerGet(node string, group string, key string) []byte {
	u := fmt.Sprintf(
		"%v/%v/%v",
		"http://"+node+"/cache",
		url.QueryEscape(group),
		url.QueryEscape(key),
	)

	// temporary
	//u = tempUrl
	retJson := HttpGet(u)

	var ir InnerResp
	err := json.Unmarshal(retJson, &ir)
	if err != nil {
		fmt.Println("json err:", err)
	}

	return []byte(ir.Data)
}

func InnerSet(node string, group string, key string, val string) []byte {
	u := "http://" + node + "/cache"

	args := url.Values{}
	args.Add("group", group)
	args.Add("key", key)
	args.Add("val", val)

	retJson := HttpPostForm(u, args)

	var ir InnerResp
	err := json.Unmarshal(retJson, &ir)
	if err != nil {
		fmt.Println("json err:", err)
	}

	return []byte(ir.Data)
}

func HttpGet(url string) []byte {
	fmt.Println(url)
	client := http.Client{
		Timeout: 500 * time.Millisecond,
	}

	resp, err := client.Get(url)

	if err != nil {
		return []byte("")
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	fmt.Println("body" + string(body))
	return body
}

func HttpPostForm(api string, args url.Values) []byte {
	resp, err := http.PostForm(api, args)
	if err != nil {
		fmt.Println("error")
		return []byte("")
	}
	defer resp.Body.Close()

	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error")
		return []byte("")
	}

	return bs
}
