# FROM golang:1.8.3 as builder
FROM tb-builder as builder
RUN go get github.com/aws/aws-sdk-go github.com/Sirupsen/logrus github.com/hashicorp/terraform github.com/mattn/go-sqlite3 github.com/jinzhu/gorm
WORKDIR /go/src/github.com/camptocamp/terraboard
COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags "-linkmode external -extldflags -static" \
  -o terraboard main.go

FROM scratch
WORKDIR /
COPY static /static
COPY index.html /
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/src/github.com/camptocamp/terraboard/terraboard /
ENTRYPOINT ["/terraboard"]
CMD [""]
