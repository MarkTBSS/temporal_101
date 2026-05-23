package app

type ClaimRequest struct {
	ClaimID  string
	PolicyID string
	Amount   float64
}

type ClaimResult struct {
	ClaimID string
	Status  string
	Reason  string
}
