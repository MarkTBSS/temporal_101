# Car Insurance Claim — Temporal Hands-On

## Business Flow

```
Customer files claim → Verify policy → Approve/Reject → Notify customer
```

---

## Temporal Design

### Workflow: `ProcessClaim`
- Input: `ClaimRequest` (claimID, policyID, amount)
- Executes 3 Activities sequentially
- Returns result (approved/rejected)

### Activities (3)

| Activity | Purpose |
|----------|---------|
| `VerifyPolicy` | Call Policy Service to check if policy is valid |
| `ApproveClaim` | Call Claim Service to approve/reject |
| `NotifyCustomer` | Notify customer of the result (print log) |

### Mock Services (2 separate services)
- **Policy Service** (`service/policy/main.go`) — port 9998
  - Endpoint: `POST /verify-policy`
  - Checks policy status
- **Claim Service** (`service/claim/main.go`) — port 9999
  - Endpoint: `POST /approve-claim`
  - Approves or rejects claim

Benefit: Can shut down services independently. e.g. stop Claim Service → VerifyPolicy still passes, but ApproveClaim fails and retries.

---

## Project Structure

```
temporal-hand-on/101/
├── README.md               # Design doc (English)
├── app/
│   ├── models.go           # Structs (package app)
│   ├── workflow.go         # ProcessClaim (package app)
│   └── activities.go       # 3 Activities (package app)
├── service/
│   ├── policy/main.go     # Policy Service (port 9998)
│   └── claim/main.go      # Claim Service (port 9999)
├── worker/main.go          # Register + start Worker
├── start/main.go           # Trigger Workflow
└── go.mod
```

---

## Concepts Practiced

| Concept | How |
|---------|-----|
| Workflow + Activity | ProcessClaim calls 3 Activities |
| Error Handling | VerifyPolicy returns "expired" → reject early, skip remaining steps |
| Retry | Shut down Claim Service → Activity fails → auto retry |
| Timeout | StartToCloseTimeout: 10s, Activity times out if service is slow |
| RetryPolicy | MaximumAttempts: 3, fails after 3 retries |
| Event History | Observe in Temporal Web UI |

---

## Task Queue: `car-insurance-claims`

---

## Phases

### Phase 1: Basic Flow (Happy Path)

1. `go mod init` — create module
2. Create `app/models.go` — define structs (ClaimRequest, ClaimResult)
3. Create `service/policy/main.go` — Policy Service (port 9998)
4. Create `service/claim/main.go` — Claim Service (port 9999)
5. Create `app/activities.go` — 3 Activities (call services via HTTP)
6. Create `app/workflow.go` — ProcessClaim calls Activities sequentially
7. Create `worker/main.go` — register Workflow + Activities
8. Create `start/main.go` — trigger Workflow with sample input
9. Run 5 terminals → verify in Web UI

### Phase 2: Error Handling & Retry

1. Update Policy Service — return "expired" when policyID = "EXPIRED"
2. Update `workflow.go` — handle expired status (reject early, skip ApproveClaim)
3. Shut down Claim Service → Activity fails → observe retry in Web UI
4. Add RetryPolicy in ActivityOptions (MaximumAttempts = 3)
5. Restart Claim Service → retry succeeds → Workflow completes

### Phase 3: Timeout

1. Set StartToCloseTimeout: 10s for all Activities
2. Add 12s delay in Claim Service `/approve-claim` (simulate slow service)
3. Run again → observe timeout + retry behavior → Workflow fails after 3 attempts

---

## How to Run

```
Terminal 1: temporal server start-dev
Terminal 2: go run ./service/policy/      (port 9998)
Terminal 3: go run ./service/claim/       (port 9999)
Terminal 4: go run ./worker/
Terminal 5: go run ./start/
```

Web UI: http://localhost:8233
