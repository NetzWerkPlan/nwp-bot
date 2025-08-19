# Build stage
FROM golang:1.25.0 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY go/ ./go/

RUN CGO_ENABLED=0 GOOS=linux go build -o /bot ./go

# Final Image Stage
FROM gcr.io/distroless/static:nonroot
WORKDIR /app

COPY --from=builder /bot /app/bot
COPY config/ /app/config/

USER 65532:65532

EXPOSE 8080

ENTRYPOINT ["/app/bot"]
