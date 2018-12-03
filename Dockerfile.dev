FROM golang:1.11
WORKDIR /go/src/github.com/camptocamp/terraboard
COPY . .
RUN go get -u github.com/golang/dep/cmd/dep && dep ensure
RUN go install .
ENTRYPOINT ["terraboard"]
CMD [""]
