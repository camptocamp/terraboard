# Terraboard

:chart_with_upwards_trend: A web dashboard to inspect Terraform States :earth_africa:


[![Docker Pulls](https://img.shields.io/docker/pulls/camptocamp/terraboard.svg)](https://hub.docker.com/r/camptocamp/terraboard/)
[![Go Report Card](https://goreportcard.com/badge/github.com/camptocamp/terraboard)](https://goreportcard.com/report/github.com/camptocamp/terraboard)
[![Gitter](https://img.shields.io/gitter/room/camptocamp/terraboard.svg)](https://gitter.im/camptocamp/terraboard)
[![Build Status](https://travis-ci.org/camptocamp/terraboard.svg?branch=master)](https://travis-ci.org/camptocamp/terraboard)
[![Coverage Status](https://coveralls.io/repos/github/camptocamp/terraboard/badge.svg?branch=master)](https://coveralls.io/github/camptocamp/terraboard?branch=master)
[![By Camptocamp](https://img.shields.io/badge/by-camptocamp-fb7047.svg)](http://www.camptocamp.com)

## What is it?

Terraboard is a web dashboard to visualize and query
[Terraform](https://terraform.io) states. It currently features:

- an overview page listing the most recently updated state files with their
  activity
- a state page with state file details, including versions and resource
  attributes
- a search interface to query resources by type, name or attributes
- a diff interface to compare state between versions

It currently only supports S3 as a remote state backend, and dynamoDB for
retrieving lock informations.


### Overview

The overview presents all the state files in the S3 bucket, by most recent
modification date.

![Screenshot Overview](screenshots/main.png)


### Search

The search view allows to find resources by various criteria.

![Screenshot Search](screenshots/search.png)


### State

The state view presents details of a Terraform state at a given version.

![Screenshot State](screenshots/state.png)


### Compare

From the state view, you can compare the current state version with another
version.

![Screenshot Compare](screenshots/compare.png)


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
* If you want to retrieve lock states
  [from a dynamoDB table](https://www.terraform.io/docs/backends/types/s3.html#dynamodb_table),
  you need to make sure the provided AWS credentials have `dynamodb:Scan` access to that
  table.


## Use with Docker

```shell
$ docker run -d -p 8080:8080 \
   -e AWS_REGION=<AWS_DEFAULT_REGION> \
   -e AWS_ACCESS_KEY_ID=<AWS_ACCESS_KEY_ID> \
   -e AWS_SECRET_ACCESS_KEY=<AWS_SECRET_ACCESS_KEY> \
   -e AWS_BUCKET=<terraform-bucket> \
   -e AWS_DYNAMODB_TABLE=<terraform-locks-table> \
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

### Architecture

Terraboard is made of two components:

#### A server process

The server is written in go and runs a web server which serves:

  - the API on known access points, taking the data from the PostgreSQL
    database
  - the index page (from [static/index.html](static/index.html)) on all other
    URLs

The server also has a routine which regularly (every 1 minute) feeds
the PostgreSQL database from the S3 bucket.

#### A web UI

The UI is an AngularJS application served from `index.html`. All the UI code
can be found in the [static/](static/) directory.



### Testing

```shell
$ docker-compose build && docker-compose up -d
# Point your browser to http://localhost
```

### Contributing


See [CONTRIBUTING.md](CONTRIBUTING.md)
