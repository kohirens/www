# Kohirens World Wide Web (WWW) Package

Provides utilities for www development.


## [![CircleCI](https://dl.circleci.com/status-badge/img/gh/kohirens/www/tree/main.svg?style=svg)](https://dl.circleci.com/status-badge/redirect/gh/kohirens/www/tree/main)

When using any of the ResponseXXX() function, it is a great idea to set
`www.FooterText` to a value that is correct for your site. It defaults to
copyright symbol and current year.

```go
package main

import (
	"io"
	"net/http"
	"os"

	"github.com/kohirens/stdlib/logger"
	"github.com/kohirens/www"
)

var log = &logger.Standard{}

func main() {
	var mainErr error

	defer func() {
		if mainErr != nil {
			log.Fatf("fatal could not start the web server: %v", mainErr.Error())
			os.Exit(1)
		}
		os.Exit(0)
	}()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, e1 := io.ReadAll(r.Body)
		if e1 != nil {
			www.Respond500(w)
			return
		}

		www.Respond200(w, []byte(""), www.ContentTypeHtml)
	})
	
	// run the web server
	mainErr = http.ListenAndServeTLS(
		":443",
		"/home/app/pki/certs/server.crt",
		"/home/app/pki/private/server.key",
		nil,
	)

	log.Infof("handler returned")
}
```