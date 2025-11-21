# VPS Deployment Guide

This guide explains how to prepare the production VPS and how the GitHub Actions pipeline deploys the Go backend.

## 1. Server Requirements
- Ubuntu 22.04 or newer.
- Go runtime **not** required (binary built in CI).
- PostgreSQL accessible from VPS (managed DB or local instance).
- Systemd available to manage the backend service.

## 2. Initial VPS Setup
```bash
sudo apt update && sudo apt upgrade -y
sudo apt install -y unzip git curl

# Create deploy directory and service user
sudo mkdir -p /opt/sigap
sudo useradd -m -s /bin/bash sigap || true
sudo chown -R sigap:sigap /opt/sigap

# Create systemd service skeleton
sudo tee /etc/systemd/system/sigap-backend.service <<'EOF'
[Unit]
Description=SIGAP Backend Service
After=network.target

[Service]
WorkingDirectory=/opt/sigap
ExecStart=/opt/sigap/sigap-backend
EnvironmentFile=/opt/sigap/.env
Restart=on-failure
User=sigap
Group=sigap

[Install]
WantedBy=multi-user.target
EOF

sudo systemctl daemon-reload
sudo systemctl enable sigap-backend.service

# Open firewall port (if using ufw)
sudo ufw allow 8080/tcp
```

## 3. Environment Variables
Prepare `.env` content that matches production DB and secrets. This is injected during deployment via GitHub Secret `PRODUCTION_ENV_FILE`.

```
APP_ENV=production
SERVER_PORT=8080
DB_HOST=prod-db-host
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=supersecret
DB_NAME=sigap
JWT_SECRET=replace_me
```

## 4. GitHub Secrets
| Secret | Description |
| --- | --- |
| `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME` | DB credentials used during CI migrate/seed. |
| `PRODUCTION_ENV_FILE` | Multi-line string representation of the `.env` file above. |
| `VPS_HOST`, `VPS_PORT`, `VPS_USER` | SSH info for the deploy user. |
| `VPS_SSH_KEY` | Private key with access to the VPS user. |
| `VPS_DEPLOY_PATH` | Target directory on VPS (e.g., `/opt/sigap`). |

## 5. GitHub Actions Workflow
Stored at `.github/workflows/main.yml` and performs:
1. Checkout & cache Go modules.
2. Download dependencies.
3. Run `go test ./...`.
4. Apply migrations with `go run cmd/migrate/main.go -command up`.
5. Seed DB with `go run cmd/seed/main.go`.
6. Build Linux binary `sigap-backend`.
7. Write `backend.env` from `PRODUCTION_ENV_FILE` secret.
8. SCP binary + env to VPS (`sigap-backend`, `.env`).
9. SSH into VPS, move env → `.env`, set executable bit, restart `sigap-backend.service`.

## 6. Manual Deployment (Fallback)
If GitHub Actions is unavailable:
```bash
GOOS=linux GOARCH=amd64 go build -o sigap-backend ./cmd/main.go
scp sigap-backend .env sigap@<VPS_HOST>:/opt/sigap
ssh sigap@<VPS_HOST> "cd /opt/sigap && chmod +x sigap-backend && sudo systemctl restart sigap-backend.service"
```

## 7. Verification
```bash
ssh sigap@<VPS_HOST>
systemctl status sigap-backend.service
curl http://localhost:8080/health
```

## 8. Troubleshooting
- Check `/var/log/syslog` or `journalctl -u sigap-backend.service` for runtime issues.
- Ensure DB firewall/security group allows access from VPS IP.
- If migrations fail during CI, rerun workflow after fixing schema issues.
