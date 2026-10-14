FROM golang:1.26-bookworm AS builder

RUN apt-get update && apt-get install -y --no-install-recommends \
    gcc \
    librdkafka-dev \
    pkg-config \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o /app/exe ./cmd/

FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    librdkafka1 \
    && rm -rf /var/lib/apt/lists/*

RUN groupadd -r app && useradd -r -g app app

WORKDIR /app

COPY --from=builder /app/exe /app/exe
COPY --from=builder /app/docs ./docs

RUN chown -R app:app /app
USER app

EXPOSE 5051
CMD ["/app/exe"]