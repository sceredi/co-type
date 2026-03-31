# Build stage
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Copy workspace files
COPY go.work go.work.sum ./

# Copy all modules
COPY broker/go.mod broker/go.sum ./broker/
COPY client/go.mod client/go.sum ./client/
COPY common/go.mod common/go.sum ./common/
COPY server/go.mod server/go.sum ./server/

# Download deps
RUN go work sync

# Copy source
COPY broker/ ./broker/
COPY client/ ./client/
COPY common/ ./common/
COPY server/ ./server/

# Build the specific binary
ARG MODULE
RUN go build -o /bin/app ./${MODULE}/cmd/app/main.go

# Final stage
FROM alpine:3.21
ARG MODULE
COPY --from=builder /bin/app /bin/app
ENTRYPOINT ["/bin/app"]

# ======== HOW TO RUN ========
# docker build --build-arg MODULE=broker -t co-type/broker .
