package main

import (
    "log"
    "net/http"

    "admission-webhook/webhook"
)

func main() {
    http.HandleFunc("/validate", webhook.ValidatePod)
	http.HandleFunc("/hello", webhook.Hello)
	http.HandleFunc("/debug", webhook.DebugRequest)

    log.Println("Starting webhook server on :8443")
    err := http.ListenAndServeTLS(":8443", "certs/tls.crt", "certs/tls.key", nil)
    if err != nil {
        log.Fatal(err)
    }
}
