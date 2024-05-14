FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod go.sum .

RUN go mod download

COPY . .

RUN go build -o /go/bin/main ./cmd/api/main.go

FROM gcr.io/distroless/static-debian12

COPY --from=builder /go/bin/main /usr/bin/

CMD ["main"]
