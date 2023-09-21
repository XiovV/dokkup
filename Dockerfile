FROM golang:1.21-alpine as build

WORKDIR /go/src/agent
COPY . .
EXPOSE 8080

RUN cd /go/src/agent/cmd/agent && go get
RUN cd /go/src/agent/cmd/agent && CGO_ENABLED=0 go build -ldflags="-s" -o /go/bin/agent

FROM alpine
COPY --from=build /go/bin/agent / 

CMD ["/agent"]
