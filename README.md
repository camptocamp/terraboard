<h1 align="center">Terraboard</h1>
<p align="center">
    <img alt="Terraboard logo" height="200" src="logo/terraboard_logo.png">
</p>
<p align="center">🌍 📋 A web dashboard to inspect Terraform States</p>
<p align="center">
  <a href="https://hub.docker.com/r/camptocamp/terraboard/" target="_blank">
    <img alt="Docker Pulls" src="https://img.shields.io/docker/pulls/camptocamp/terraboard.svg" />
  </a>
  <a href="https://goreportcard.com/report/github.com/camptocamp/terraboard" target="_blank">
    <img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/camptocamp/terraboard" />
  </a>
  <a href="https://gitter.im/camptocamp/terraboard" target="_blank">
    <img alt="Gitter" src="https://img.shields.io/gitter/room/camptocamp/terraboard.svg" />
  </a>
  <a href="https://github.com/camptocamp/terraboard/actions" target="_blank">
    <img alt="Build Status" src="https://github.com/camptocamp/terraboard/workflows/Build/badge.svg" />
  </a>
  <a href="https://coveralls.io/github/camptocamp/terraboard?branch=master" target="_blank">
    <img alt="Coverage Status" src="https://coveralls.io/repos/github/camptocamp/terraboard/badge.svg?branch=master" />
  </a>
  <a href="http://www.camptocamp.com" target="_blank">
    <img alt="By Camptocamp" src="https://img.shields.io/badge/by-camptocamp-fb7047.svg" />
  </a>
  <a href="https://pkg.go.dev/github.com/camptocamp/terraboard" target="_blank">
    <img alt="Documentation" src="https://pkg.go.dev/badge/github.com/camptocamp/terraboard">
  </a>
</p>
<p align="center">Website: <a href="https://terraboard.io">https://terraboard.io</a></p>

---

<details><summary>Table of content</summary>

- [What is it?](#what-is-it)
  - [Overview](#overview)
  - [Search](#search)
  - [State](#state)
  - [Compare](#compare)
  - [Requirements](#requirements)
    - [AWS S3 (state) + DynamoDB (lock)](#aws-s3-state--dynamodb-lock)
    - [Terraform Cloud](#terraform-cloud)
- [Configuration](#configuration)
  - [Multiple buckets/providers](#multiple-bucketsproviders)
  - [Available parameters](#available-parameters)
    - [Application Options](#application-options)
    - [General Provider Options](#general-provider-options)
    - [Logging Options](#logging-options)
    - [Database Options](#database-options)
    - [AWS (and S3 compatible providers) Options](#aws-and-s3-compatible-providers-options)
    - [S3 Options](#s3-options)
    - [Terraform Enterprise Options](#terraform-enterprise-options)
    - [Google Cloud Platform Options](#google-cloud-platform-options)
    - [GitLab Options](#gitlab-options)
    - [Web](#web)
    - [Help Options](#help-options)
- [Push plans to Terraboard](#push-plans-to-terraboard)
- [Use with Docker](#use-with-docker)
  - [Docker-compose](#docker-compose)
  - [Docker command line](#docker-command-line)
- [Use with Kubernetes](#use-with-kubernetes)
- [Use with Rancher](#use-with-rancher)
- [Authentication and base URL](#authentication-and-base-url)
- [Install from source](#install-from-source)
- [Compatibility Matrix](#compatibility-matrix)
- [Development](#development)
  - [Architecture](#architecture)
    - [A server process](#a-server-process)
    - [A web UI](#a-web-ui)
  - [Testing](#testing)
  - [Contributing](#contributing)

</details>

## What is it?

Terraboard is a web dashboard to visualize and query
[Terraform](https://terraform.io) states. It currently features:

- an overview page listing the most recently updated state files with their
  activity
- a state page with state file details, including versions and resource
  attributes
- a search interface to query resources by type, name or attributes
- a diff interface to compare state between versions

It currently supports several remote state backend providers:

- [AWS S3 (state) + DynamoDB (lock)](https://www.terraform.io/docs/backends/types/s3.html)
- [S3 compatible backends (ex: MinIO)](https://min.io/)
- [Google Cloud Storage](https://www.terraform.io/docs/backends/types/gcs.html)
- [Terraform Cloud (remote)](https://www.terraform.io/docs/backends/types/remote.html)
- [GitLab](https://docs.gitlab.com/ee/user/infrastructure/terraform_state.html)

Terraboard is now able to handle multiple buckets/providers configuration! 🥳
Check *configuration* section for more details. 

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

Independently of the location of your statefiles, Terraboard needs to store an internal version of its dataset. For this purpose it requires a PostgreSQL database.
Data resiliency is not paramount though as this dataset can be rebuilt upon your statefiles at anytime.
#### AWS S3 (state) + DynamoDB (lock)

- A **versioned** S3 bucket name with one or more Terraform states, named with a `.tfstate` suffix
- AWS credentials with the following IAM permissions over the bucket:
  - `s3:GetObject`
  - `s3:ListBucket`
  - `s3:ListBucketVersions`
  - `s3:GetObjectVersion`
- If you want to retrieve lock states [from a dynamoDB table](https://www.terraform.io/docs/backends/types/s3.html#dynamodb_table), you need to make sure the provided AWS credentials have `dynamodb:Scan` access to that table.
#### Terraform Cloud

- Account on [Terraform Cloud](https://app.terraform.io/)
- Existing organization
- Token assigned to an organization

## Configuration

Terraboard currently supports configuration in three different ways:

1. Environment variables **(only usable for single provider configuration)**
2. CLI parameters **(only usable for single provider configuration)**
3. Configuration file (YAML). A configuration file example can be found in the root directory of this repository and in the `test/` subdirectory.

**Important: all flags/environment variables related to the providers settings aren't compatible with multi-provider configuration! Instead, you must use the YAML config file to be able to configure multiples buckets/providers. YAML config is able to load values from environments variables.**

The precedence of configurations is as described below.

### Multiple buckets/providers

In order for Terraboard to import states from multiples buckets or even providers, you must use the YAML configuration method:

- Set the `CONFIG_FILE` environment variable or the `-c`/`--config-file` flag to point to a valid YAML config file.
- In the YAML file, specify your desired providers configuration. For example with two MinIO buckets (using the AWS provider with compatible mode):

```yaml
# Needed since MinIO doesn't support versioning or locking
provider:
  no-locks: true
  no-versioning: true

aws:
  - endpoint: http://minio:9000/
    region: ${AWS_DEFAULT_REGION}
    s3:
      - bucket: test-bucket
        force-path-style: true
        file-extension: 
          - .tfstate

  - endpoint: http://minio:9000/
    region: eu-west-1
    s3:
      - bucket: test-bucket2
        force-path-style: true
        file-extension: 
          - .tfstate
```

In the case of AWS, don't forget to set the `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` environment variables.

That's it! Terraboard will now fetch these two buckets on DB refresh. You can also mix providers like AWS and Gitlab or anything else.
You can find a ready-to-use Docker example with two *MinIO* buckets in the `test/multiple-minio-buckets/` sub-folder. 

### Available parameters

#### Application Options

- `-V`, `--version` Display version.
- `-c`, `--config-file` <default: *$CONFIG_FILE*> Config File path
  - Env: *CONFIG_FILE*
  
#### General Provider Options

- `--no-versioning` <default: *$TERRABOARD_NO_VERSIONING*> Disable versioning support from Terraboard (useful for S3 compatible providers like MinIO)
  - Env: *TERRABOARD_NO_VERSIONING*
  - Yaml: *provider.no-versioning*
- `--no-locks` <default: *$TERRABOARD_NO_LOCKS*> Disable locks support from Terraboard (useful for S3 compatible providers like MinIO)
  - Env: *TERRABOARD_NO_LOCKS*
  - Yaml: *provider.no-locks*

#### Logging Options

- `-l`, `--log-level` <default: *"info"*> Set log level ('debug', 'info', 'warn', 'error', 'fatal', 'panic').
  - Env: *TERRABOARD_LOG_LEVEL*
  - Yaml: *log.level*
- `--log-format` <default: *"plain"*> Set log format ('plain', 'json').
  - Env: *TERRABOARD_LOG_FORMAT*
  - Yaml: *log.format*

#### Database Options

- `--db-host` <default: *"db"*> Database host.
  - Env: *DB_HOST*
  - Yaml: *database.host*
- `--db-port` <default: *"5432"*> Database port.
  - Env: *DB_PORT*
  - Yaml: *database.port*
- `--db-user` <default: *"gorm"*> Database user.
  - Env: *DB_USER*
  - Yaml: *database.user*
- `--db-password` <default: *$DB_PASSWORD*> Database password.
  - Env: *DB_PASSWORD*
  - Yaml: *database.password*
- `--db-name` <default: *"gorm"*> Database name.
  - Env: *DB_NAME*
  - Yaml: *database.name*
- `--db-sslmode` <default: *"require"*> Database SSL mode.
  - Env: *DB_SSLMODE*
  - Yaml: *database.sslmode*
- `--no-sync` Do not sync database.
  - Yaml: *database.no-sync*
- `--sync-interval` <default: *"1"*> DB sync interval (in minutes)
  - Yaml: *database.sync-interval*

#### AWS (and S3 compatible providers) Options

- `--aws-access-key` <default: *$AWS_ACCESS_KEY_ID*> AWS account access key.
  - Env: *AWS_ACCESS_KEY_ID*
  - Yaml: *aws.access-key*
- `--aws-secret-access-key` <default: *$AWS_SECRET_ACCESS_KEY*> AWS secret account access key.
  - Env: *AWS_SECRET_ACCESS_KEY*
  - Yaml: *aws.secret-access-key*
- `--aws-session-token` <default: *$AWS_SESSION_TOKEN*> AWS session token.
  - Env: *AWS_SESSION_TOKEN*
  - Yaml: *aws.session-token*
- `--dynamodb-table` <default: *$AWS_DYNAMODB_TABLE*> AWS DynamoDB table for locks.
  - Env: *AWS_DYNAMODB_TABLE*
  - Yaml: *aws.dynamodb-table*
- `--aws-endpoint` <default: *$AWS_ENDPOINT*> AWS endpoint.
  - Env: *AWS_ENDPOINT*
  - Yaml: *aws.endpoint*
- `--aws-region` <default: *$AWS_REGION*> AWS region.
  - Env: *AWS_REGION*
  - Yaml: *aws.region*
- `--aws-role-arn` <default: *$APP_ROLE_ARN*> Role ARN to Assume.
  - Env: *APP_ROLE_ARN*
  - Yaml: *aws.app-role-arn*
- `--aws-external-id` <default: *$AWS_EXTERNAL_ID*> External ID to use when assuming role.
  - Env: *AWS_EXTERNAL_ID*
  - Yaml: *aws.external-id*

#### S3 Options

- `--s3-bucket` <default: *$AWS_BUCKET*> AWS S3 bucket.
  - Env: *AWS_BUCKET*
  - Yaml: *aws.s3.bucket*
- `--key-prefix` <default: *$AWS_KEY_PREFIX*> AWS Key Prefix.
  - Env: *AWS_KEY_PREFIX*
  - Yaml: *aws.s3.key-prefix*
- `--file-extension` <default: *".tfstate"*> File extension(s) of state files.
  - Env: *AWS_FILE_EXTENSION*
  - Yaml: *aws.s3.file-extension*
- `--force-path-style` <default: *$AWS_FORCE_PATH_STYLE*> Force path style S3 bucket calls.
  - Env: *AWS_FORCE_PATH_STYLE*
  - Yaml: *aws.s3.force-path-style*

#### Terraform Enterprise Options

- `--tfe-address` <default: *$TFE_ADDRESS*> Terraform Enterprise address for states access
  - Env: *TFE_ADDRESS*
  - Yaml: *tfe.address*
- `--tfe-token` <default: *$TFE_TOKEN*> Terraform Enterprise Token for states access
  - Env: *TFE_TOKEN*
  - Yaml: *tfe.token*
- `--tfe-organization` <default: *$TFE_ORGANIZATION*> Terraform Enterprise organization for states access
  - Env: *TFE_ORGANIZATION*
  - Yaml: *tfe.organization*

#### Google Cloud Platform Options

- `--gcs-bucket` Google Cloud bucket to search
  - Yaml: *gcp.gcs-bucket*
- `--gcp-sa-key-path` <default: *$GCP_SA_KEY_PATH*> The path to the service account to use to connect to Google Cloud Platform
  - Env: *GCP_SA_KEY_PATH*
  - Yaml: *gcp.gcp-sa-key-path*

#### GitLab Options

- `--gitlab-address` <default: *"https://gitlab.com"*> GitLab address (root)
  - Env: *GITLAB_ADDRESS*
  - Yaml: *gitlab.address*
- `--gitlab-token` <default: *$GITLAB_TOKEN*> Token to authenticate upon GitLab
  - Env: *GITLAB_TOKEN*
  - Yaml: *gitlab.token*

#### Web

- `-p`, `--port` <default: *"8080"*> Port to listen on.
  - Env: *TERRABOARD_PORT*
  - Yaml: *web.port*
- `--base-url` <default: *"/"*> Base URL.
  - Env: *TERRABOARD_BASE_URL*
  - Yaml: *web.base-url*
- `--logout-url` <default: *$TERRABOARD_LOGOUT_URL*> Logout URL.
  - Env: *TERRABOARD_LOGOUT_URL*
  - Yaml: *web.logout-url*

#### Help Options

- `-h`, `--help` Show this help message

## Push plans to Terraboard

In order to send Terraform plans to Terraboard, you must wrap it in this JSON format:
```json
{
    "lineage": "<Plan's lineage>",
    "terraform_version": "<Terraform version>",
    "git_remote": "<The URL of the remote that generated this plan>",
    "git_commit": "<Commit hash>",
    "ci_url": "<The URL of the CI that sent this plan>",
    "source": "<Free field for the triggering event>",
    "plan_json": "<Terraform plan JSON export>"
}
```

And send it to `/api/plans` using **POST** method

## Use with Docker

### Docker-compose

Configuration file can be provided to the container using a [volume](https://docs.docker.com/compose/compose-file/#volumes) or a [configuration](https://docs.docker.com/compose/compose-file/#configs).

```shell
# Set AWS credentials as environment variables:
export AWS_ACCESS_KEY_ID=<access_key>
export AWS_SECRET_ACCESS_KEY=<access_secret>

# Set AWS configuration as environment variables:
export AWS_DEFAULT_REGION=<AWS default region>
export AWS_BUCKET=<S3 Bucket name>
export AWS_DYNAMODB_TABLE=<Aws DynamoDB Table>

docker-compose up
```

Then point your browser to http://localhost:8080.

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
  -e GODEBUG="netdns=go" \
  --net terraboard \
  --detach \
  --restart=always \
  postgres:9.5

docker run -p 8080:8080 \
  -e AWS_ACCESS_KEY_ID="${AWS_ACCESS_KEY_ID}" \
  -e AWS_SECRET_ACCESS_KEY="${AWS_SECRET_ACCESS_KEY}" \
  -e AWS_REGION="${AWS_DEFAULT_REGION}" \
  -e AWS_BUCKET="${AWS_BUCKET}" \
  -e AWS_DYNAMODB_TABLE="${AWS_DYNAMODB_TABLE}" \
  -e DB_PASSWORD="<mypassword>" \
  -e DB_SSLMODE="disable" \
  --net terraboard \
  camptocamp/terraboard:latest
```

Then point your browser to http://localhost:8080.


## Use with Kubernetes

A Helm chart is available on [Camptocamp's repository](https://github.com/camptocamp/charts/tree/master/terraboard).

In order to install it:

```shell
$ helm repo add c2c https://camptocamp.github.io/charts
$ helm install -v values.yaml terraboard c2c/terraboard
```


## Use with Rancher

[Camptocamp's Rancher Catalog](https://github.com/camptocamp/camptocamp-rancher-catalog)
contains a Terraboard template to automate its installation in Cattle.


## Authentication and base URL

Terraboard does not implement authentication. Instead, it is recommended to use
an authentication proxy such as [oauth2_proxy](https://github.com/bitly/oauth2_proxy).

If you need to set a route path for Terraboard, you can set a base URL by
passing it as the `BASE_URL` environment variable.

When using an authentication proxy, Terraboard will retrieve the logged in
user and email from the headers passed by the proxy.
Terraboard expects you to setup the HTTP Headers `X-Forwarded-User` and
`X-Forwarded-Email` when passing the logged in user and email. A Nginx
example can be found below:

```nginx
location / {
  ....
  auth_request_set $user   $upstream_http_x_auth_request_user;
  auth_request_set $email  $upstream_http_x_auth_request_email;
  proxy_set_header X-Forwarded-User  $user;
  proxy_set_header X-Forwarded-Email $email;
  ...
  proxy_pass http://terraboard/;
}
```

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
| 1.0.0    |  0.14.5             |
| 1.1.0    |  0.14.10            |

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
