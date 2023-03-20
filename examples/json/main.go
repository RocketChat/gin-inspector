package main

import (
	"github.com/RocketChat/gin-inspector"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	debug := true

	if debug {
		r.Use(inspector.InspectorStats())
		r.GET("/_inspector", inspector.JsonFrontend)
	}

	r.Run()
}
