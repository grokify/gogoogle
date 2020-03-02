# Installing Prerequisites

## Google Go Client: Git Clone vs. Go Get

The Google Go Client can be installed using the following:

Get the source here: https://code.googlesource.com/google-api-go-client

```
$ cd $GOPATH/src/google.golang.org
$ git clone https://code.googlesource.com/google-api-go-client api
$ cd api
$ go get ./...
```

If you use `go get` per the docs, you may end up with the following error:

`$ go get google.golang.org/apis/slides/v1
package google.golang.org/apis/slides/v1: unrecognized import path "google.golang.org/apis/slides/v1" (parse https://google.golang.org/apis/slides/v1?go-get=1: no go-import meta tags ())`

> [Christopher Berger wrote](https://forum.golangbridge.org/t/no-go-import-meta-tags-error-on-go-get/6312/2):
> The message “no go-import meta tags” indicates that go get fails because it thinks this URL is a vanity import path (scroll down to "For code hosted on other servers) 600 and tries to get the meta information to find the real repository.

### Dependencies

The following dependencies may not be installed.

```
./google.golang.org/api/internal/gensupport/resumable.go:16:2: cannot find package "github.com/googleapis/gax-go/v2"
./google.golang.org/api/transport/http/dial.go:18:2: cannot find package "go.opencensus.io/plugin/ochttp"
./google.golang.org/api/transport/http/internal/propagation/http.go:19:2: cannot find package "go.opencensus.io/trace"
./google.golang.org/api/transport/http/internal/propagation/http.go:20:2: cannot find package "go.opencensus.io/trace/propagation"
./google.golang.org/api/internal/conn_pool.go:8:2: cannot find package "google.golang.org/grpc"
./google.golang.org/api/internal/pool.go:10:2: cannot find package "google.golang.org/grpc/naming"
```

## GRPC

Similar to Google Client which uses `google.golang.org`, instead of `$ go get -u google.golang.org/grpc`

Use:

```
$ cd $GOPATH/src/google.golang.org
$ git clone https://github.com/grpc/grpc-go grpc
$ cd grpc
$ go get ./...
```

## Google API Extensions

Repo: https://github.com/googleapis/gax-go

```
$ go get -u github.com/googleapis/gax-go/v2
```

Error:

```
`./google.golang.org/api/internal/gensupport/resumable.go:16:2: cannot find package "github.com/googleapis/gax-go/v2"`
```

## Install OpenCensus

https://opencensus.io/quickstart/go/tracing/

Use: 

`$ go get go.opencensus.io/... && go get contrib.go.opencensus.io/exporter/zipkin`

Docs result in an error:

```
$ go get go.opencensus.io/* && go get contrib.go.opencensus.io/exporter/zipkin
package go.opencensus.io/*: go.opencensus.io/*: invalid import path: malformed import path "go.opencensus.io/*": invalid char '*'
```