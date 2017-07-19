FROM golang:1.7.3 as builder
WORKDIR /go/src/github.com/camptocamp/terraboard
RUN go get github.com/aws/aws-sdk-go
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o terraboard main.go

FROM scratch
WORKDIR /
COPY static /static
COPY index.html /
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/src/github.com/camptocamp/terraboard/terraboard /
ENTRYPOINT ["/terraboard"]
CMD [""]
