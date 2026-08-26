# syntax=docker/dockerfile:1

# ---------------------------------------------------------------------------
# Stage 1 — build the React/Vite frontend into static assets
# ---------------------------------------------------------------------------
FROM node:22-alpine AS frontend-builder
WORKDIR /app/frontend

COPY frontend/package.json frontend/package-lock.json ./
RUN npm ci

COPY frontend/ ./
RUN npm run build
# -> produces /app/frontend/dist

# ---------------------------------------------------------------------------
# Stage 2 — build the Go backend binary
# ---------------------------------------------------------------------------
FROM golang:1.25-alpine AS backend-builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server .

# ---------------------------------------------------------------------------
# Stage 3 — minimal runtime image
# ---------------------------------------------------------------------------
FROM alpine:3.20
RUN apk add --no-cache ca-certificates
WORKDIR /app

# Go server binary
COPY --from=backend-builder /app/server ./server

# Built frontend, served as static files by the Go server (see main.go)
COPY --from=frontend-builder /app/frontend/dist ./frontend/dist

EXPOSE 3000
ENTRYPOINT ["./server"]
