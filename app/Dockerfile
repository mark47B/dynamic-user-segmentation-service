FROM golang:1.18-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:latest

RUN apk update && apk upgrade && \
    apk add --no-cache make

WORKDIR /app
COPY --from=builder /app/app /app/app
EXPOSE 8080

CMD ["./app"]