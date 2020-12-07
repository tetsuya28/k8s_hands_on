FROM golang:alpine AS builder
WORKDIR /app
COPY . .
RUN go get -d -v
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-w -s" -o app
RUN chmod +x /app/app

FROM alpine:3.10
COPY --from=builder /app/app .
CMD [ "./app" ]
