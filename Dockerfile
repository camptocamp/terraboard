FROM golang:1.7.3 as builder
WORKDIR /go/src/github.com/camptocamp/terraboard
RUN go get github.com/aws/aws-sdk-go
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o terraboard-server main.go

FROM scratch
COPY --from=builder /go/src/github.com/camptocamp/terraboard/terraboard-server /
ENTRYPOINT ["/terraboard-server"]
CMD [""]
