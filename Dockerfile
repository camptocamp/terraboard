FROM golang:1.12 as builder
WORKDIR /go/src/github.com/camptocamp/terraboard
COPY . .
ENV GO111MODULE=on
RUN make terraboard

FROM scratch
WORKDIR /
COPY static /static
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/src/github.com/camptocamp/terraboard/terraboard /
ENTRYPOINT ["/terraboard"]
CMD [""]
