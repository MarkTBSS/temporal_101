package app

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func VerifyPolicy(ctx context.Context, req ClaimRequest) (string, error) {
	resp, err := http.Post("http://localhost:9998/verify-policy",
		"application/json",
		strings.NewReader(fmt.Sprintf(`{"policy_id":"%s"}`, req.PolicyID)))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result map[string]string
	json.NewDecoder(resp.Body).Decode(&result)
	return result["status"], nil
}

func ApproveClaim(ctx context.Context, req ClaimRequest) (bool, error) {
	resp, err := http.Post("http://localhost:9999/approve-claim",
		"application/json",
		strings.NewReader(fmt.Sprintf(`{"claim_id":"%s","amount":%f}`, req.ClaimID, req.Amount)))
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	var result map[string]bool
	json.NewDecoder(resp.Body).Decode(&result)
	return result["approved"], nil
}

func NotifyCustomer(ctx context.Context, result ClaimResult) error {
	fmt.Printf("Notified customer: ClaimID=%s Status=%s Reason=%s\n", result.ClaimID, result.Status, result.Reason)
	return nil
}
