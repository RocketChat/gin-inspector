package main

import (
	"html/template"
	"io"
	"net/http"
	"time"

	"github.com/RocketChat/gin-inspector"
	"github.com/gin-gonic/gin"
)

func formatDate(t time.Time) string {
	return t.Format(time.RFC822)
}

func main() {
	r := gin.Default()
	r.Delims("{{", "}}")

	r.SetFuncMap(template.FuncMap{
		"formatDate": formatDate,
	})

	r.LoadHTMLFiles("inspector.html")
	debug := true

	if debug {
		r.Use(inspector.InspectorStats())

		r.GET("/_inspector", func(c *gin.Context) {
			c.HTML(http.StatusOK, "inspector.html", map[string]interface{}{
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
