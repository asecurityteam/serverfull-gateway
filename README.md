<a id="markdown-serverfull-gateway---an-aws-api-gateway-stand-in-for-serverfull" name="serverfull-gateway---an-aws-api-gateway-stand-in-for-serverfull"></a>
# Serverfull-Gateway - An AWS API Gateway Stand-In For Serverfull
[![GoDoc](https://godoc.org/github.com/asecurityteam/serverfull-gateway?status.svg)](https://godoc.org/github.com/asecurityteam/serverfull-gateway)
[![Build Status](https://travis-ci.com/asecurityteam/serverfull-gateway.png?branch=master)](https://travis-ci.com/asecurityteam/serverfull-gateway)
[![codecov.io](https://codecov.io/github/asecurityteam/serverfull-gateway/coverage.svg?branch=master)](https://codecov.io/github/asecurityteam/serverfull-gateway?branch=master)

*Status: Incubation*

<!-- TOC -->

- [Serverfull-Gateway - An AWS API Gateway Stand-In For Serverfull](#serverfull-gateway---an-aws-api-gateway-stand-in-for-serverfull)
    - [Overview](#overview)
    - [Quick Start](#quick-start)
    - [Using The Docker Image](#using-the-docker-image)
    - [Configuration](#configuration)
    - [Templates](#templates)
        - [Request Templates](#request-templates)
        - [Response Templates](#response-templates)
    - [Contributing](#contributing)
        - [License](#license)
        - [Contributing Agreement](#contributing-agreement)

<!-- /TOC -->

<a id="markdown-overview" name="overview"></a>
## Overview

**Deprecation Notice:** This package will be archived and made read-only on January 30th, 2024. After January 30th this repo will cease to be maintained on Github.

This project is a compliment to
[Serverfull](https://github.com/asecurityteam/serverfull) and implements a subset of
the features offered by the AWS API Gateway Proxy for AWS Lambda. It is built on our
other project, [transportd](https://github.com/asecurityteam/transportd), which
bundles our HTTP performance and resiliency tooling into a service.

Generally speaking, if you have access to AWS Lambda then it is probably best to use
AWS Lambda and the AWS API Gateway Proxy event mapping features that are provided by
AWS. This tool is a bit specialized and was conceived to enable:

-   Teams who want to adopt a serverless style of development but without full access
    to a serverless/FaaS provider such as those who are developing on bare metal or
    private cloud infrastructures.

-   Teams who are ready to migrate away from AWS Lambda but aren't ready to rewrite
    large portions of existing code such as those who are either changing cloud
    providers or moving to EC2 in order to fine-tune the performance characteristics
    of their runtimes.

We hope to keep the scope of this project small. Our goals for functionality are:

-   Have the gateway, all routes, and input/output validation driven by valid OpenAPI
    specifications.

-   Adapt between HTTP callers and Lambda functions.

-   Provide reasonable compatibility with AWS API Gateway features as they relate to
    managing request/response content and linking to AWS Lambda.

We don't intend to add features outside of these goals except as they might be included
from `transportd`.

<a id="markdown-quick-start" name="quick-start"></a>
## Quick Start

Everything is driven by an OpenAPI specification. We use extensions to configure
the server:

```yaml
openapi: 3.0.0
# General runtime configuration for the server. For details,
# see the `Runtime Settings` in in the transpord documentation.
x-runtime:
  httpserver:
    address: ":9090"
  logger:
    level: "INFO"
    output: "STDOUT"
  stats:
    output: "DATADOG"
    datadog:
      address: "statsd:8126"
      flushinterval: "10s"
  signals:
    installed:
      - "OS"
    os:
      signals:
        - 2 # SIGINT
        - 15 # SIGTERM
x-transportd:
  backends:
    - app
  app:
    host: "http://app:8080" # Location of the serverfull instance.
    pool:
      ttl: "24h"
      count: 1
info:
  version: 1.0.0
  title: Sample specification
  description: Used for demonstration
  contact:
    name: n/a
    email: na@localhost.com
  license:
    name: Apache 2.0
    url: 'https://www.apache.org/licenses/LICENSE-2.0.html'
paths:
  /:
    post:
      x-transportd:
        backend: app
        enabled:
          - lambda # lambda is the transportd plugin that this project provides.
        lambda:
          arn: "sample" # The name of the lambda as configured in serverfull
          async: false # Fire and forget or wait for a response
          # Go template strings are used to map gateway input to lambda input
          request: '{"value": "#!.Request.Body.inputValue!#"}'
          # Template strings are also used to map lambda output to gateway output
          success: '{"status": 200, "body": {"v":"#!.Response.Body.someValue!#"}}'
          error: '{"status": 500, "bodyPassthrough": true}'
          # See "Templates" for more details
          authenticate: false # Toggle to enable AWS authentication for requests.
          session: # Optional AWS session authentication configuration.
            region: "us-west-2" # Required AWS region for requests (if auth enabled).
            sharedprofile: # Optional profile configuration.
              file: "" # Blank for AWS default location
              profile: "" # Name of the AWS profile defined in the file.
            static: # Optional static credentials
              id: "" # AWS access key ID.
              secret: "" # AWS secret key.
              token: "" # Optional temporary access token.
            assumerole: # Optional assume role credentials.
              role: "" # ARN of the role to assume.
              externalid: "" # Optional cross-account role identifier.
      description: Sample endpoint.
      requestBody:
        required: true
        description: The request body.
        application/json:
          schema:
            type: string
      responses:
        "200":
          description: "Success"
```

With a configured OpenAPI document:

```bash
TRANSPORTD_OPENAPI_SPECIFICATION_FILE="api.yaml" go run main.go
```

<a id="markdown-using-the-docker-image" name="using-the-docker-image"></a>
## Using The Docker Image

We provide a docker image of the latest build that you can add your
specification to:

```docker
FROM atlassian/severfull-gateway
COPY api.yaml .
ENV TRANSPORTD_OPENAPI_SPECIFICATION_FILE="api.yaml"
```

The base container is defined in this project's Dockerfile and uses the
`scratch` base internally.

<a id="markdown-configuration" name="configuration"></a>
## Configuration

There are a suite of options provided by the `transportd` service that are
documented with that project. This project adds a plugin for `transportd`
that adds the following option block to each path:

```yaml
lambda:
          arn: "sample" # The name of the lambda as configured in serverfull
          async: false # Fire and forget or wait for a response
          # Go template strings are used to map gateway input to lambda input
          request: '{"value": "#!.Request.Body.inputValue!#"}'
          # Template strings are also used to map lambda output to gateway output
          success: '{"status": 200, "body": {"v":"#!.Response.Body.someValue!#"}}'
          error: '{"status": 500, "bodyPassthrough": true}'
```

<a id="markdown-templates" name="templates"></a>
## Templates

<a id="markdown-request-templates" name="request-templates"></a>
### Request Templates

Lambda functions only receive the `POST` body of a request to the Invoke API. A good
consumer facing API, however, is going to use the breadth of features offered by
the OpenAPI specification to define inputs such as putting values in query parameters
or URL segments. To bridge the gap we offer a request template configuration that maps
incoming HTTP requests to Lambda input payloads. We use the Go
[template syntax](https://golang.org/pkg/text/template/) which has its own rich
documentation around the format. The two elements of the template feature that we
customize in this project are 1) the template delimiters and 2) the content injected
into the templates.

Because the `{` and `}` characters are so common in JSON we have adjusted them to be
the `#!` pair for open and `!#` pair for close. Anything within these pairs of characters
will be interpreted as template text. When generating a request template, the purpose is
to generate a shape that can be serialized to JSON. The shape we inject is:

```json
{
    "Request": {
        "URL": {
            "string": "map[string]string"
        },
        "Query": {
            "string": ["string"]
        },
        "Header": {
            "string": ["string"]
        },
        "Body":{}
    }
}
```

The shape for errors is the same except that the `Body` attribute is actually
an AWS error body in the form of:

```json
{
    "errorMessage": "string",
    "errorType": "string",
    "stackTrace": ["string"]
}
```

rather than the response from the Lambda.

<a id="markdown-response-templates" name="response-templates"></a>
### Response Templates

Lambda functions have a limited set of response codes they can return and don't come
with a rich set of headers. To help with these, we provide the ability to template
both success and error responses. The shape injected into success responses is:

```json
{
    "Request": {
        "URL": {
            "string": "string"
        },
        "Query": {
            "string": ["string"]
        },
        "Header": {
            "string": ["string"]
        },
        "Body":{}
    },
    "Response": {
        "Status": 200,
        "Header": {
            "string": ["string"]
        },
        "Body": {}
    }
}
```

<a id="markdown-contributing" name="contributing"></a>
## Contributing

<a id="markdown-license" name="license"></a>
### License

This project is licensed under Apache 2.0. See LICENSE.txt for details.

<a id="markdown-contributing-agreement" name="contributing-agreement"></a>
### Contributing Agreement

Atlassian requires signing a contributor's agreement before we can accept a patch. If
you are an individual you can fill out the [individual
CLA](https://na2.docusign.net/Member/PowerFormSigning.aspx?PowerFormId=3f94fbdc-2fbe-46ac-b14c-5d152700ae5d).
If you are contributing on behalf of your company then please fill out the [corporate
CLA](https://na2.docusign.net/Member/PowerFormSigning.aspx?PowerFormId=e1c17c66-ca4d-4aab-a953-2c231af4a20b).
