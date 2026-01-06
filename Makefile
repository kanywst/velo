.PHONY: lint lint-backend lint-frontend audit audit-backend audit-frontend check-all

# Run all checks
check-all: lint audit

# Linting
lint: lint-backend lint-frontend

lint-backend:
	@echo "==> Linting Backend (Go)..."
	go vet ./...

lint-frontend:
	@echo "==> Linting Frontend (Vue)..."
	cd frontend && npm run lint

# Vulnerability Checks
audit: audit-backend audit-frontend

audit-backend:
	@echo "==> Auditing Backend (Go)..."
	@echo "Running go mod verify..."
	go mod verify
	@if command -v govulncheck >/dev/null 2>&1; then \
		echo "Running govulncheck..."; \
		govulncheck ./...; \
	else \
		echo "govulncheck not found. Skipping vulnerability scan."; \
		echo "Install with: go install golang.org/x/vuln/cmd/govulncheck@latest"; \
	fi

audit-frontend:
	@echo "==> Auditing Frontend (Node)..."
	cd frontend && npm audit --audit-level=high
