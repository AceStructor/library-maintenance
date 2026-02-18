FROM golang:1.26-alpine AS builder

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

COPY library-maintenance-api/. .
RUN go mod download
RUN go build -o /app .

FROM alpine:3.21 AS final

COPY --from=builder /app /bin/app

EXPOSE 6000
CMD ["bin/app"]
