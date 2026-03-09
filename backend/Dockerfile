# Build Stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# We need to copy the entire backend structure because of local imports (albion/common, etc)
COPY . .

# Build the specific service passed via build-arg or context
ARG SERVICE_NAME
WORKDIR /app/${SERVICE_NAME}

# Initialize module if not already done (it is, but for safety)
# RUN go mod download 

# Build binary
RUN go build -o /app/binary .

# Run Stage
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/binary .
# Copy metadata.json if it's the metadata service (optional, or volume mount)
# COPY --from=builder /app/backend/metadata/metadata.json . 

EXPOSE 8080 8081 8082

CMD ["./binary"]
