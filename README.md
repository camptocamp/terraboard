# Terraboard

Web Dashboard to inspect Terraform States


[![Docker Pulls](https://img.shields.io/docker/pulls/camptocamp/terraboard.svg)](https://hub.docker.com/r/camptocamp/terraboard/)
[![Go Report Card](https://goreportcard.com/badge/github.com/camptocamp/terraboard)](https://goreportcard.com/report/github.com/camptocamp/terraboard)
[![Gitter](https://img.shields.io/gitter/room/camptocamp/terraboard.svg)](https://gitter.im/camptocamp/terraboard)
[![By Camptocamp](https://img.shields.io/badge/by-camptocamp-fb7047.svg)](http://www.camptocamp.com)


## What is it?

Terraboard is a web dashboard to visualize and query
[Terraform](https://terraform.io) states.

![Screenshot](screenshot.png)


### Requirements

Terraboard currently supports getting the Terraform states from AWS S3. It
requires:

* A **versioned** S3 bucket name with one or more Terraform states,
  named with a `.tfstate` suffix
* AWS credentials with the following rights on the bucket:
   - `s3:GetObject`
   - `s3:ListBucket`
   - `s3:ListBucketVersions`
   - `s3:GetObjectVersion`
* A running PostgreSQL database


## Use with Docker

```shell
$ docker run -d -p 8080:80 \
   -e AWS_REGION=<AWS_DEFAULT_REGION> \
   -e AWS_ACCESS_KEY_ID=<AWS_ACCESS_KEY_ID> \
   -e AWS_SECRET_ACCESS_KEY=<AWS_SECRET_ACCESS_KEY> \
   -e AWS_BUCKET=<terraform-bucket> \
   -e DB_PASSWORD="mygreatpasswd" \
   --link postgres:db \
   camptocamp/terraboard:latest
```

and point your browser to http://localhost:8080


### Authentication and base URL

Terraboard does not implement authentication. Instead, it is recommended to use
an authentication proxy such as [oauth2_proxy](https://github.com/bitly/oauth2_proxy).

If you need to set a route path for Terraboard, you can set a base URL by
passing it as the `BASE_URL` environment variable.



## Install from source

```shell
$ go get github.com/camptocamp/terraboard
```

## Development


### Testing

```shell
$ docker-compose build && docker-compose up -d
# Point your browser to http://localhost
```

### Contributing

Please report bugs on the [GitHub project
page](https://github.com/camptocamp/terraboard/issues).

We welcome contributions in the form of [Pull
Requests](https://github.com/camptocamp/terraboard/pulls).


