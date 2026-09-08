# IIoT Smart Analytics & Predictive Maintenance Platform

Monorepo for the platform: a simulated industrial line (Factory I/O) feeds
telemetry to a Go backend, which persists it, sends it to a Python AI
service for anomaly detection / RUL prediction, streams live results to a
React dashboard, and supports operator + automatic control commands.

## Repo layout

```
.
├── backend/        # Go + Gin + sqlc + PostgreSQL
├── frontend/        # React operator dashboard
├── ai-service/       # Python anomaly detection + RUL microservice
├── docs/             # Architecture
└── docker-compose.yml  # (coming soon — full-stack local dev, all 3 services)
```