FROM golang:1.25.0-bookworm AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO=ENABLED=0 go build -o /app/cmd/exe /app/cmd

FROM alpine:3.23
WORKDIR /app
COPY --from=builder /app/cmd/exe /app
CMD ["/app/exe"]