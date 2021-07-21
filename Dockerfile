FROM golang:1.15 as builder
WORKDIR /opt/build
COPY . .
RUN make build

FROM node:16.5-alpine3.14 as node-builder
WORKDIR /opt/build
COPY static/terraboard-vuejs ./terraboard-vuejs
WORKDIR /opt/build/terraboard-vuejs
RUN npm install
RUN npm run build

FROM scratch
WORKDIR /
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /opt/build/terraboard /
COPY --from=node-builder /opt/build/terraboard-vuejs/dist /static
ENTRYPOINT ["/terraboard"]
CMD [""]
