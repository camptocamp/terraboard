FROM golang:1.8.3 as builder
WORKDIR /go/src/github.com/camptocamp/terraboard
COPY . .
RUN go get -u github.com/golang/dep/cmd/dep && dep ensure
RUN CGO_ENABLED=1 GOOS=linux go build \
    -ldflags "-linkmode external -extldflags -static -X main.version=`git describe --always`" \
	-o terraboard main.go \
	&& strip terraboard

FROM scratch
WORKDIR /
COPY static /static
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/src/github.com/camptocamp/terraboard/terraboard /
ENTRYPOINT ["/terraboard"]
CMD [""]
