package main

import (
	"log"
	"net/http"
	"simple_gin/gee"
	"time"
)

func my_mid() gee.HandleFunc {
	return func(c *gee.Context) {
		t := time.Now()
		c.JSON(200, "success")
		log.Printf("[%d] %s in %v", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}

func main() {
	r := gee.New()

	r.Use(gee.Logger())

	r.GET("/index", func(c *gee.Context) {
		c.HTML(http.StatusOK, "<h1>Hi, I'm Jiyeon<h1>")
	})

	v1 := r.Group("/v1")
	v1.Use(my_mid())
	{
		v1.GET("/", func(c *gee.Context) {
			c.HTML(http.StatusOK, "<h1>Hello Gee<h1>")
		})
		v1.GET("/hello", func(c *gee.Context) {
			c.String(http.StatusOK, "hello %s, there is %s\n", c.Query("name"), c.Path)
		})
	}

	v2 := r.Group("/v2")
	{
		v2.GET("/hello/:name", func(c *gee.Context) {
			c.String(http.StatusOK, "hello %s, you are at %s", c.Param("name"), c.Path)
		})
		v2.POST("/login", func(c *gee.Context) {
			c.JSON(http.StatusOK, gee.H{
				"username": c.PostForm("username"),
				"password": c.PostForm("password"),
			})
		})
	}

	r.GET("/assets/*filepath", func(c *gee.Context) {
		c.JSON(http.StatusOK, gee.H{
			"filepath": c.Param("filepath"),
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
