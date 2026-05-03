FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o contacts-stats main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/contacts-stats /app/contacts-stats
COPY index.html /app/index.html

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

EXPOSE 8080

ENTRYPOINT ["/app/contacts-stats"]
CMD ["serve"]
