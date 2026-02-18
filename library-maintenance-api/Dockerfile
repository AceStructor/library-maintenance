FROM golang:1.22-alpine
WORKDIR /app
COPY library-maintenance-api/. .
RUN go build -o app
EXPOSE 6000
CMD ["./app"]
