FROM golang:1.18-alpine as builder
WORKDIR /build
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o /zebra cmd/main.go

FROM alpine:3
COPY --from=builder zebra /bin/main
COPY .env ./
ENTRYPOINT ["/bin/main"]