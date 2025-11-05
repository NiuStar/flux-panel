# Multi-stage build: frontend + backend into a single minimal image

# --- Frontend build ---
FROM node:20-alpine AS fe
WORKDIR /fe
COPY vite-frontend/package*.json ./
COPY vite-frontend/ .
# Use npm install (with legacy peer deps) to avoid lock mismatch and peer resolution failures in CI
RUN npm install --legacy-peer-deps --no-audit --no-fund && npm run build

# --- Backend build ---
FROM golang:1.25-alpine AS be
WORKDIR /app
COPY . ./
RUN apk add --no-cache build-base ca-certificates && \
    go build -mod vendor -ldflags "-w -s"  -o /app/server ./golang-backend/cmd/server && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod vendor -ldflags "-w -s" -o /app/flux-agent-linux-amd64 ./golang-backend/cmd/flux-agent && \
    CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -mod vendor -ldflags "-w -s" -o /app/flux-agent-linux-arm64 ./golang-backend/cmd/flux-agent

# --- Final runtime ---
FROM alpine:3.19
WORKDIR /app
ENV PORT=6365
COPY --from=be /app/server /app/server
COPY --from=fe /fe/dist /app/public
# publish flux-agent binaries for node installation
RUN mkdir -p /app/public/flux-agent
COPY --from=be /app/flux-agent-linux-amd64 /app/public/flux-agent/flux-agent-linux-amd64
COPY --from=be /app/flux-agent-linux-arm64 /app/public/flux-agent/flux-agent-linux-arm64
# serve install.sh from the backend container
COPY install.sh /app/install.sh
EXPOSE 6365
CMD ["/app/server"]
