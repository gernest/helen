# helen [![Build Status](https://travis-ci.org/gernest/helen.svg)](https://travis-ci.org/gernest/helen) [![GoDoc](https://godoc.org/github.com/gernest/helen?status.svg)](https://godoc.org/github.com/gernest/helen) [![Coverage Status](https://coveralls.io/repos/gernest/helen/badge.svg?branch=master&service=github)](https://coveralls.io/github/gernest/helen?branch=master)

Helen is the static assets handler for golang based web applications. Helen simply handles serving your static files that is javascript, stylesheets and images.

# features

* Fast
* Javascript and Stylesheets minifications
* support
 - [gorilla mux](https://github.com/gorilla/mux)
 - [echo](https://github.com/labstack/echo)
 - [http.ServeMux](https://godoc.org/net/http#ServeMux)
 - [httrouter](https://github.com/julienschmidt/httprouter)[coming soon]
* middlewares. You can chain any alice compatible middlwares.


# Motivation

After using different golang frameworks and libraries for building web applications, I usually came across the need to serve my static files and heck It is a bit complicated.

So helen, is a simple aproach to handle static assets. Helen supports different golang routers, you can use this as a way to remind youself how you can handle static assets for your favorite router.

# Installation

	go get github.com/gernest/helen


# How to use

```go
package main

import (
	"log"
	"net/http"

	"github.com/gernest/helen"
	"github.com/gorilla/mux"
)

func main() {

	// Create the instance of your router
	server := mux.NewRouter()

	// Create a new helen.Static instance. We are passing "static" as the directory
	// we want to serve static content from.
	static := helen.NewStatic("fixtures")

	// We bind anything matching /static/ route to our static handler
	static.Bind("/static/", server)

	// You can register other handlers to your server or whatever you want to do with it.

	log.Fatal(http.ListenAndServe(":8000", server))
}

```

In this example we will be sriving contents of the fixtures directory as static files.

This example uses `gorilla mux` router. Note that you can bind the helen handler at any point of your application. This example also works for all the supported routers.


# Middlewares

You can add any alice compatible middlewares to the `*Static` instance. If you  want gzip compression and you have a wonderful implementation called `gzipMe`.

Then,
```go
	static:=helen.Static("static")
	static.Use(gzipMe)
```

Will register your middleware. Not that you can pass whatever number of middlewares you want to the `Static.Use` method.


# Documentation
[![GoDoc](https://godoc.org/github.com/gernest/helen?status.svg)](https://godoc.org/github.com/gernest/helen)

# TODO
- Optimize memory usage
- Write benchmarks

# Contributing

Start with clicking the star button to make the author and his neighbors happy. Then fork the repository and submit a pull request for whatever change you want to be added to this project.

If you have any questions, just open an issue.

# Author
Geofrey Ernest <geofreyernest@live.com>

Twitter  : [@gernesti](https://twitter.com/gernesti)

Facebook : [Geofrey Ernest](https://www.facebook.com/geofrey.ernest.35)

# Licence

This project is released under the MIT licence. See [LICENCE](LICENCE) for more details.