# Gin Inspector

![Gin Inspector HTML Preview](https://raw.githubusercontent.com/fatihkahveci/gin-inspector/master/preview-html.png)

![Gin Inspector HTML Preview 2](https://raw.githubusercontent.com/fatihkahveci/gin-inspector/master/preview-html-2.jpg)

Gin middleware for investigating http request.

## Usage


```sh
$ go get github.com/fatihkahveci/gin-inspector
```

### JSON Response

```Go
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
```

### Html Template

```Go
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

```