FROM golang:1.21 AS builder
ENV CGO_ENABLED 0
ARG VERSION
WORKDIR /go/src/app
ADD . .
RUN go build -o /kube-exec-all

FROM busybox
COPY --from=builder /kube-exec-all /kube-exec-all
WORKDIR /data/kube-exec-all
CMD ["/kube-exec-all"]
