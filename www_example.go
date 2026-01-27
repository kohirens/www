package www

import (
	"io"
	"net/http"
)

func ExampleRespond500() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, e1 := io.ReadAll(r.Body)
		if e1 != nil {
			Respond500(w, []byte(""), ContentTypeHtml)
			return
		}

		Respond200(w, []byte(""), ContentTypeHtml)
	})

	// Run the web server.
	e1 := http.ListenAndServeTLS(
		":443",
		"/home/app/pki/certs/server.crt",
		"/home/app/pki/private/server.key",
		nil,
	)

	if e1 != nil { // Don't just throw the error away.
		panic(e1.Error())
	}
}
