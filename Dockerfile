FROM golang:1.22 as builder

COPY . /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux go build .

FROM alpine

COPY --from=builder /app/media-check /app/media-check
RUN /sbin/apk add --no-cache findutils

CMD ["/app/media-check"]
