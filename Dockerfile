# Build phase
FROM golang:1.19-alpine AS builder

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GO111MODULE=on

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go mod download

WORKDIR /app/src

RUN go build -o transacion-api

# Copy phase
FROM alpine

WORKDIR /app

COPY --from=builder /app/src/transacion-api .

COPY --from=builder /app/src/docs .

EXPOSE 8098

CMD ["./transacion-api"]