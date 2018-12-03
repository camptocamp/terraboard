FROM golang:1.11 as builder
WORKDIR /go/src/github.com/camptocamp/terraboard
COPY . .
RUN go get -u github.com/golang/dep/cmd/dep && dep ensure
RUN make terraboard

FROM scratch
WORKDIR /
COPY static /static
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/src/github.com/camptocamp/terraboard/terraboard /
ENTRYPOINT ["/terraboard"]
CMD [""]
