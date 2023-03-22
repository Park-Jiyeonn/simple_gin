package main

import (
	"net/http"
	"simple_gin/gee"
)

func main() {
	r := gee.New()
	r.GET("/", func(c *gee.Context) {
		c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
	})
	r.GET("/hello", func(c *gee.Context) {
		c.String(http.StatusOK, "hello %s, there is %s\n", c.Query("name"), c.Path)
	})
	r.POST("/login", func(c *gee.Context) {
		c.JSON(http.StatusOK, gee.H{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		})
	})
	r.Run(":8081")
}

// git 代理：
//git config --global http.proxy http://127.0.0.1:1080
//git config --global https.proxy http://127.0.0.1:1080
// 取消代理：
//git config --global --unset http.proxy
//git config --global --unset https.proxy
