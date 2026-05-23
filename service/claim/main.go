package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/approve-claim", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"approved":true}`)
	})
	fmt.Println("Claim Service running on :9999")
	http.ListenAndServe(":9999", nil)
}
