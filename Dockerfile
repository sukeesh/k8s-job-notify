FROM golang:alpine as builder

WORKDIR /go/src/github.com/sukeesh/k8s-job-notify
ADD . /go/src/github.com/sukeesh/k8s-job-notify
RUN go build -o /app .

FROM alpine

COPY --from=builder /app /app
ENTRYPOINT ["/app"]