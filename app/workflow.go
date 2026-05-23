package app

import (
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func ProcessClaim(ctx workflow.Context, req ClaimRequest) (ClaimResult, error) {
	options := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts: 3,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, options)

	// Step 1: Verify Policy
	var status string
	err := workflow.ExecuteActivity(ctx, VerifyPolicy, req).Get(ctx, &status)
	if err != nil {
		return ClaimResult{}, err
	}

	if status == "expired" {
		result := ClaimResult{
			ClaimID: req.ClaimID,
			Status:  "rejected",
			Reason:  "Policy expired",
		}
		workflow.ExecuteActivity(ctx, NotifyCustomer, result).Get(ctx, nil)
		return result, nil
	}

	// Step 2: Approve Claim
	var approved bool
	err = workflow.ExecuteActivity(ctx, ApproveClaim, req).Get(ctx, &approved)
	if err != nil {
		return ClaimResult{}, err
	}

	// Step 3: Build result and notify
	result := ClaimResult{
		ClaimID: req.ClaimID,
		Status:  "approved",
		Reason:  "Policy active, claim approved",
	}
	if !approved {
		result.Status = "rejected"
		result.Reason = "Claim not approved"
	}

	err = workflow.ExecuteActivity(ctx, NotifyCustomer, result).Get(ctx, nil)
	if err != nil {
		return ClaimResult{}, err
	}

	return result, nil
}
