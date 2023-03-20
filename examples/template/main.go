package main

import (
	"io"

	inspector "github.com/RocketChat/gin-inspector"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	err := inspector.LoadHtml(r)
	if err != nil {
		panic(err)
	}
	debug := true

	if debug {
		r.Use(inspector.InspectorStats())

		r.GET("/_inspector", inspector.Frontend)
		r.POST("/test", func(c *gin.Context) {
			body, err := io.ReadAll(c.Request.Body)
			if err != nil {
				c.String(500, err.Error())
				return
			}
			c.String(200, string(body))
		})
	}

	r.Run(":8080")
}
