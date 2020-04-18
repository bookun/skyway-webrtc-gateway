FROM golang:1.14 as builder

WORKDIR /go/src/github.com/bookun/skyway-webrtc-gateway/samples/golang

ENV GO111MODULE=on

RUN go build ./cmd/callee/main.go
RUN go build ./cmd/caller/main.go

#FROM hmlatapie/gstreamer
#WORKDIR /gw-client
#COPY --from=builder /go/src/github.com/bookun/skyway-webrtc-gateway/samples/golang/callee /gw-client
#COPY --from=builder /go/src/github.com/bookun/skyway-webrtc-gateway/samples/golang/caller /gw-client
