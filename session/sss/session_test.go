package sss

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/kohirens/stdlib/logger"
	"github.com/kohirens/www/session"
	"net/http"
	"os"
	"strconv"
	"time"
)

func ExampleNewStorageClient() {
	bucket := os.Getenv("S3_BUCKET_NAME")
	if bucket == "" {
		panic("missing environment variable S3_BUCKET_NAME")
	}

	log := logger.Standard{}

	// Start a new S3 storage client, it picks up its credential from the environment.
	sessionStorage := NewStorageClient(bucket, context.Background())
	// set where to store the session in the bucket.
	sessionStorage.Prefix("session/")
	// HTTP Session handler using RAM and then saving to Amazon S3 for
	// longer-term.
	sessionManager := session.NewManager(sessionStorage, time.Minute*20)

	type Counter struct {
		Visits int // Only public fields will be saved to the session.
	}

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		counter := &Counter{}

		// Load any previous session data.
		sessionManager.Load(w, req)

		// Convert the data into something your application can use.
		jsonData := sessionManager.Get("counter")
		if jsonData != nil {
			if e := json.Unmarshal(jsonData, counter); e != nil {
				panic("failed to unmarshal client info: " + e.Error())
			}
		}

		counter.Visits++
		_, e1 := w.Write([]byte(`{"count": "` + strconv.Itoa(counter.Visits) + `"}`))
		if e1 != nil {
			log.Errf(e1.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		bytes, e2 := json.Marshal(counter)
		if e2 != nil {
			log.Errf(e2.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Save data to the session.
		sessionManager.Set("counter", bytes)

		// Write the session to Amazon S3.
		if e := sessionManager.Save(); e != nil {
			log.Errf(e.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	})

	// One can use generate_cert.go in crypto/tls to generate cert.pem and key.pem.
	fmt.Printf("About to listen on 8443. Go to https://127.0.0.1:8443/")
	err := http.ListenAndServeTLS(":8443", "cert.pem", "key.pem", nil)
	if err != nil {
		log.Errf(err.Error())
	}
}
