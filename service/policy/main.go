package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/verify-policy", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			PolicyID string `json:"policy_id"`
		}
		json.NewDecoder(r.Body).Decode(&req)

		status := "active"
		if req.PolicyID == "EXPIRED" {
			status = "expired"
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status":"%s"}`, status)
	})
	fmt.Println("Policy Service running on :9998")
	http.ListenAndServe(":9998", nil)
}
