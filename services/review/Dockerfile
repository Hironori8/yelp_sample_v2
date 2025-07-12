FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy shared models
COPY ./shared/models ./shared/models

# Copy service files
COPY ./services/review .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/main .

CMD ["./main"]