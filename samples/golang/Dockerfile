FROM golang:1.13 as builder

WORKDIR /go/src/github.com/bookun/skyway-webrtc-gateway/samples/golang

ENV GO111MODULE=on
COPY . .

RUN go build -o /callee ./cmd/callee/main.go
RUN go build -o /caller ./cmd/caller/main.go

FROM hmlatapie/gstreamer
WORKDIR /gw-client
COPY --from=builder /callee /gw-client/callee
COPY --from=builder /caller /gw-client/caller
