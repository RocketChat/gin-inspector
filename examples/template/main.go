package main

import (
	"io"
	"net/http"

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

		r.GET("/_inspector", func(c *gin.Context) {
			c.HTML(http.StatusOK, inspector.HtmlName, map[string]interface{}{
				"title":      "Gin Inspector",
				"pagination": inspector.GetPaginator(),
			})

		})
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
