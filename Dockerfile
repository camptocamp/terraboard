FROM golang:1.15 as builder
WORKDIR /opt/build
COPY . .
RUN make build

FROM scratch
WORKDIR /
COPY static /static
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /opt/build/terraboard /
ENTRYPOINT ["/terraboard"]
CMD [""]
