FROM golang:1.10.1-alpine3.7 as builder
COPY app.go .
RUN go build -o /app .

FROM alpine:3.7
CMD ["./app"]
COPY --from=builder /app .
