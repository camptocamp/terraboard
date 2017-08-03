FROM golang:1.8.3 as builder
#FROM tb-builder as builder
RUN go get github.com/aws/aws-sdk-go \
           github.com/Sirupsen/logrus \
		   github.com/hashicorp/terraform \
		   github.com/jinzhu/gorm github.com/lib/pq \
		   github.com/jessevdk/go-flags \
		   github.com/pmezard/go-difflib || echo
WORKDIR /go/src/github.com/camptocamp/terraboard
COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build \
    -ldflags "-linkmode external -extldflags -static -X main.version=`git describe --always --dirty`" \
	-o terraboard main.go \
	&& strip terraboard

FROM scratch
WORKDIR /
COPY static /static
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/src/github.com/camptocamp/terraboard/terraboard /
ENTRYPOINT ["/terraboard"]
CMD [""]
