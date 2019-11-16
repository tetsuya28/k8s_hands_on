FROM golang:alpine AS builder
RUN apk update && apk add --no-cache git
WORKDIR /app
COPY . .
RUN go get -d -v
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-w -s" -o app
RUN chmod +x /app/app

FROM alpine:3.10
COPY --from=builder /app/app .
COPY .env .
COPY .env.docker .
ENV GO_ENV docker
EXPOSE 8000
ENTRYPOINT ["./app"]
