# Terraboard

Website: [https://camptocamp.github.io/terraboard](https://camptocamp.github.io/terraboard)

![Terraboard Logo](logo/terraboard_logo.png)

🌍 📋 A web dashboard to inspect Terraform States


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

It currently supports S3 as a remote state backend, and dynamoDB for
retrieving lock informations. Also supports Terraform Cloud (in past Terraform Enterprise [more](https://www.terraform.io/docs/cloud/index.html#note-about-product-names))


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

Terraboard currently supports getting the Terraform states from AWS S3 and Terraform Cloud (in past Terraform Enterprise [more](https://www.terraform.io/docs/cloud/index.html#note-about-product-names)). It
requires:
- Terraform states from AWS S3:
  * A **versioned** S3 bucket name with one or more Terraform states,
    named with a `.tfstate` suffix
  * AWS credentials with the following rights on the bucket:
    - `s3:GetObject`
    - `s3:ListBucket`
    - `s3:ListBucketVersions`
    - `s3:GetObjectVersion`
  * If you want to retrieve lock states
  [from a dynamoDB table](https://www.terraform.io/docs/backends/types/s3.html#dynamodb_table),
  you need to make sure the provided AWS credentials have `dynamodb:Scan` access to that
  table.
- Terraform states from Terraform Cloud:
  * Account on [Terraform Cloud](https://app.terraform.io/)
  * Existing organization
  * Token assigned to an organization
- A running PostgreSQL database


## Configuration

Terraboard currently supports configuration in three different ways:

1. Environment variables
2. CLI parameters
3. Configuration file (YAML). A configuration file example can be found in the root directory of this repository.

The precedence of configurations is as described below.

### Available parameters

|CLI|ENV|YAML|Description|Default|
|---|---|----|-----------|-------|
|`-V` or `--version`| - | - | Prints version | - |
|`-p` or `--port`|`TERRABOARD_PORT`|`web.port`|Port to listen on| 8080 |
|`-c` or `--config-file`|`CONFIG_FILE`|-|Config File path| - |
|`-l` or `--log-level` | `TERRABOARD_LOG_LEVEL` | `log.level` | Set log level (debug, info, warn, error, fatal, panic) | info |
|`--log-format` | `TERRABOARD_LOG_FORMAT` | `log.format` | Set log format (plain, json) | plain |
|`--db-host` | `DB_HOST` | `db.host` | Database host | db |
|`--db-port` | `DB_PORT` | `db.port` | Database port | 5432 |
|`--db-user` | `DB_USER` | `db.user` | Database user | gorm |
|`--db-password` | `DB_PASSWORD` | `db.password` | Database password | - |
|`--db-name` | `DB_NAME` | `db.name` | Database name | gorm |
|`--db-sslmode` | `DB_SSLMODE` | `db.sslmode` | SSL mode enforced for database access (require, verify-full, verify-ca, disable) | require |
|`--no-sync` | - | `db.no-sync` | Do not sync database | false |
|`--sync-interval` | - | `db.sync-interval` | DB sync interval (in minutes) | 1 |
|`--dynamodb-table` | `AWS_DYNAMODB_TABLE` | `aws.dynamodb-table` | AWS DynamoDB table for locks | - |
|`--s3-bucket` | `AWS_BUCKET` | `aws.bucket` | AWS S3 bucket | - |
|`--app-role-arn` | `APPRoleArn` | `aws.app-role-arn` | Role ARN to Assume | - |
|`--key-prefix` | `AWS_KEY_PREFIX` | `aws.key-prefix` | AWS Key Prefix | - |
|`--file-extension` | `AWS_FILE_EXTENSION` | `aws.file-extension` | File extension of state files | .tfstate |
|`--base-url` | `TERRABOARD_BASE_URL` | `web.base-url` | Base URL | / |
|`--logout-url` | `TERRABOARD_LOGOUT_URL` | `web.logout-url` | Logout URL | - |
|`--tfe-address` | `TFE_ADDRESS` | `tfe.tfe-address` | Terraform Enterprise address for states access | - |
|`--tfe-token` | `TFE_TOKEN` | `tfe.tfe-token` | Terraform Enterprise token for states access | - |
|`--tfe-organization` | `TFE_ORGANIZATION` | `tfe.tfe-organization` | Terraform Enterprise organization for states access | - |
|`--gcs-bucket` | `N/A` | `gcp.gcs-buckets` | Google Cloud Storage buckets to access | - |
|`--gcp-sa-key-path` | `GCP_SA_KEY_PATH` | `gcp.gcp-sa-key-path` | Path to the service account key to use for Google Cloud Storage | - |

## Use with Docker

### Docker-compose

To use the included compose file, you will need to configure an OAuth application from Github, more information in next section [Authentication and base URL](#Authentication-and-base-URL)

Configuration file can be provided to the container using a [volume](https://docs.docker.com/compose/compose-file/#volumes) or a [configuration](https://docs.docker.com/compose/compose-file/#configs).

```shell
# Set oauth information:
export OAUTH_CLIENT_ID=<>
export OAUTH_CLIENT_SECRET=<>
export OAUTH_COOKIE_SECRET=<>

# Set AWS credentials as environment variables:
export AWS_ACCESS_KEY_ID=<access_key>
export AWS_SECRET_ACCESS_KEY=<access_secret>

# Set AWS configuration as environment variables:
export AWS_DEFAULT_REGION=<AWS default region>
export AWS_BUCKET=<S3 Bucket name>
export AWS_DYNAMODB_TABLE=<Aws DynamoDB Table>
export AWS_KEY_PREFIX=<AWS Key Prefix>

# Set basic Terraboard configuration
export TERRABOARD_LOG_LEVEL=<Set log level (debug, info, warn, error, fatal, panic)>

docker-compose up
```

Then point your browser to http://localhost.

### Docker command line

```shell
# Set AWS credentials as environment variables:
export AWS_ACCESS_KEY_ID=<access_key>
export AWS_SECRET_ACCESS_KEY=<access_secret>

# Set AWS configuration as environment variables:
export AWS_DEFAULT_REGION=<AWS default region>
export AWS_BUCKET=<S3 Bucket name>
export AWS_DYNAMODB_TABLE=<AWS_DYNAMODB_TABLE>

# Spin up the two containers and a network for them to communciate on:
docker network create terraboard
docker run --name db \
  -e POSTGRES_USER=gorm \
  -e POSTGRES_DB=gorm \
  -e POSTGRES_PASSWORD="<mypassword>" \
  --net terraboard \
  --detach \
  --restart=always \
  postgres:9.5

docker run -p 8080:8080 \
  -e AWS_ACCESS_KEY_ID="${AWS_ACCESS_KEY_ID}" \
  -e AWS_SECRET_ACCESS_KEY="${AWS_SECRET_ACCESS_KEY}" \
  -e AWS_REGION="${AWS_DEFAULT_REGION}" \
  -e AWS_BUCKET="${AWS_BUCKET}" \
  -e WS_DYNAMODB_TABLE="${AWS_DYNAMODB_TABLE}" \
  -e DB_PASSWORD="<mypassword>" \
  -e DB_SSLMODE="disable" \
  --net terraboard \
  camptocamp/terraboard:latest
```

Then point your browser to http://localhost:8080.

## Use with Rancher

[Camptocamp's Rancher Catalog](https://github.com/camptocamp/camptocamp-rancher-catalog)
contains a Terraboard template to automate its installation in Cattle.


## Authentication and base URL

Terraboard does not implement authentication. Instead, it is recommended to use
an authentication proxy such as [oauth2_proxy](https://github.com/camptocamp/oauth2_proxy).

If you need to set a route path for Terraboard, you can set a base URL by
passing it as the `BASE_URL` environment variable.

When using an authentication proxy, Terraboard will retrieve the logged in
user and email from the headers passed by the proxy.
You can also pass a `TERRABOARD_LOGOUT_URL` parameter to allow users to
sign out of the proxy.


## Install from source

```shell
$ go get github.com/camptocamp/terraboard
```

## Compatibility Matrix

|Terraboard|Max Terraform version|
|----------|---------------------|
| 0.15.0   |  0.12.7             |
| 0.16.0   |  0.12.7             |
| 0.17.0   |  0.12.18            |
| 0.18.0   |  0.12.18            |
| 0.19.0   |  0.12.20            |
| 0.20.0   |  0.12.26            |
| 0.21.0   |  0.12.28            |
| 0.22.0   |  0.13.0             |

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
