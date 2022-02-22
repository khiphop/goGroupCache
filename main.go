package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"gokache"
	"strconv"
)

type Source struct {
	Data string `json:"data"`
}

var (
	selfPort = 7002
	selfIp   = "127.0.0.1"
)

func main() {
	gokache.HttpGet("http://127.0.0.1:8013")

	gokache.NewGroup("user", 10000, gokache.BsFunc(func(key string) ([]byte, error) {
		// return []byte(key), nil
		rs := gokache.HttpGet("http://127.0.0.1:8013" + "?key=" + key)
		r := fetchData(rs)
		return r, nil
	}))
	gokache.NewGroup("club", 10000, gokache.BsFunc(func(key string) ([]byte, error) {
		// return []byte(key), nil
		rs := gokache.HttpGet("http://127.0.0.1:8013" + "?key=" + key)
		r := fetchData(rs)
		return r, nil
	}))

	nd := gokache.InitNode(selfIp + ":" + string(selfPort))

	router := gin.Default()
	router.SetTrustedProxies([]string{"127.0.0.1"})

	// api handler
	router.GET("/cache/:group/:key", func(ginC *gin.Context) {
		getHandler(nd, ginC.Param("group"), ginC.Param("key"), ginC, false)
	})

	router.GET("/inner/:group/:key", func(ginC *gin.Context) {
		getHandler(nd, ginC.Param("group"), ginC.Param("key"), ginC, true)
	})

	router.POST("/cache", func(ginC *gin.Context) {
		group := ginC.PostForm("group")
		key := ginC.PostForm("key")
		val := ginC.PostForm("val")
		setHandler(nd, group, key, val, ginC, false)
	})

	router.POST("/inner", func(ginC *gin.Context) {
		group := ginC.PostForm("group")
		key := ginC.PostForm("key")
		val := ginC.PostForm("val")
		setHandler(nd, group, key, val, ginC, true)
	})

	router.POST("/node", func(ginC *gin.Context) {
		ip := ginC.PostForm("ip")
		port := ginC.PostForm("port")

		nd.RegNode(ip + ":" + port)

		ginC.JSON(200, gin.H{
			"code": 200,
			"msg":  "success",
			"data": string(nd.DisplayNode()),
		})
	})

	router.DELETE("/node", func(ginC *gin.Context) {
		ip := ginC.PostForm("ip")
		port := ginC.PostForm("port")

		nd.RemoveNode(ip + ":" + port)

		ginC.JSON(200, gin.H{
			"code": 200,
			"msg":  "success",
			"data": string(nd.DisplayNode()),
		})
	})

	router.PUT("/data", func(ginC *gin.Context) {
		_ = gokache.BackupData()

		ginC.JSON(200, gin.H{
			"code": 200,
			"msg":  "success",
			"data": true,
		})
	})

	err := router.Run(":" + strconv.Itoa(selfPort))
	if err != nil {
		return
	}
}

func getHandler(nd *gokache.NodeDispatch, group string, key string, ginC *gin.Context, inner bool) {
	d, _ := nd.GetHandler(group, key, inner)

	ginC.JSON(200, gin.H{
		"code": 200,
		"msg":  "success",
		"data": string(d),
	})
}

func setHandler(nd *gokache.NodeDispatch, group string, key string, val string, ginC *gin.Context, inner bool) {
	_ = nd.SetHandler(group, key, val, inner)

	ginC.JSON(200, gin.H{
		"code": 200,
		"msg":  "success",
		"data": true,
	})
}

func fetchData(res []byte) []byte {
	var tt Source
	err := json.Unmarshal(res, &tt)
	if err != nil {
		fmt.Println("json err:", err)
	}

	fmt.Println(tt.Data)

	return []byte(tt.Data)
}
