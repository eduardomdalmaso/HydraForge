---
name: validate-project
description: >-
  Procedure for running full project validation, DDD layer compliance audits,
  modular frontend file size checks, Go unit tests, build checks, HTTP endpoint verification, and documentation synchronization on HydraStream.
  Use when validating changes, performing code review, or before completing tasks.
---

# 🛡️ HydraStream Project Validation Skill

Follow this multi-step runbook to validate the integrity, DDD compliance, modular frontend rules (<100 lines), build health, runtime status, and documentation sync of HydraStream.

## Step 1: DDD Domain Purity Audit

Verify that `internal/domain` does not contain forbidden infrastructure imports (`net/http`, `database/sql`, etc.):

```bash
grep -rn "net/http" internal/domain/ && echo "❌ Violations found!" || echo "✅ Domain layer is pure!"
```

## Step 2: Modular Frontend Line Count Audit

Verify that no web file under `web/` exceeds 100 lines:

```bash
wc -l web/css/*.css web/css/*/*.css web/js/*.js web/js/*/*.js web/views/*.html web/index.html
```

*Expected result:* All web files have fewer than 100 lines.

## Step 3: Run Unit & Concurrency Tests

Execute fast Go unit tests across all packages:

```bash
go test -v ./...
```

*Expected result:* All domain and application tests pass cleanly.

## Step 4: Compile Production Binary

Ensure the application builds without syntax or dependency errors:

```bash
go build -o bin/hydrastream ./cmd/hydrastream
```

## Step 5: Endpoint & Runtime Smoke Test

Launch the binary in the background, make test API calls, and verify responses:

```bash
./bin/hydrastream &
SERVER_PID=$!
sleep 1

curl -s http://localhost:8080/healthz
curl -s http://localhost:8080/api/v1/info
curl -s http://localhost:8080/api/v1/streams

kill $SERVER_PID
```

## Step 6: Architecture Documentation Audit

Confirm that mandatory architectural docs are synchronized with codebase changes:

- Ensure [`.gemini/GEMINI.md`](file:///home/hades/Documents/HydraStream/.gemini/GEMINI.md) reflects current rules and package layout.
- Ensure [`.gemini/docs/architecture.md`](file:///home/hades/Documents/HydraStream/.gemini/docs/architecture.md) reflects current component diagrams and subsystem specs.

## Step 7: Final Sanity Check

Check that git status is clean and all new components follow the Ports & Adapters structure:

- `internal/domain` -> Entities & Domain Logic
- `internal/ports` -> Interfaces (`StreamUseCase`, `StreamRepository`)
- `internal/application` -> Services
- `internal/adapters` -> HTTP & Store implementations
