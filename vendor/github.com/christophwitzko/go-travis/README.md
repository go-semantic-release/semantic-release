# go-travis

go-travis is a Go client library for accessing the [Travis CI API](http://docs.travis-ci.com/api/).

**Documentation:** [![GoDoc](https://godoc.org/github.com/Ableton/go-travis/travis?status.svg)](https://godoc.org/github.com/Ableton/go-travis)

**Build Status:** [![Build Status](https://travis-ci.org/Ableton/go-travis.svg?branch=master)](https://travis-ci.org/Ableton/go-travis)

go-travis requires Go version 1.1 or greater.

## Dive

```go
import (
    "log"
    travis "github.com/Ableton/go-travis"
)

client := travis.NewDefaultClient("")
builds, _, _, resp, err := client.Builds.ListFromRepository("Ableton/go-travis", nil)
if err != nil {
    log.Fatal(err)
}

// Now do something with the builds
```

## Installation

```bash
$ go get github.com/Ableton/go-travis
```

## Usage

Interaction with the Travis CI API is done through a `Client` instance.

```go
import travis "github.com/Ableton/go-travis"

client := travis.NewClient(travis.TRAVIS_API_DEFAULT_URL, "asuperdupertoken")
```

Constructing it with the ``NewClient`` helper requires two arguments:
* The Travis CI API URL you wish to communicate with. Different Travis CI plans are accessed through different URLs. go-travis exposes constants for these URLs:
  * ``TRAVIS_API_DEFAULT_URL``: default *api.travis-ci.org* endpoint for the free Travis "Open Source" plan.
  * ``TRAVIS_API_PRO_URL``: the *api.travis-ci.com* endpoint for the paid Travis pro plans.
* A Travis CI token with which to authenticate. If you wish to run requests unauthenticated, pass an empty string. It is possible at any time to authenticate the Client instance with a Travis token or a Github token. For more information see [Authentication]().


### Services oriented design

The ``Client`` instance's ``Service`` attributes provide access to Travis CI API resources.

```go
opt := &travis.BuildListOptions{EventType: "pull request"}
builds, response, err := client.Builds.ListFromRepository("mygithubuser/mygithubrepo", opt)
if err != nil {
        log.Fatal(err)
}
```

**Non exhaustive list of implemented services**:
+ Authentication
+ Branches
+ Builds
+ Commits
+ Jobs
+ Logs
+ Repositories
+ Requests
+ Users

(*For an up to date exhaustive list, please check out the documentation*)


**Nota**: Service methods will often take an *Option* (sub-)type instance as input. These types, like ``BuildListOptions`` allow narrowing and filtering your requests.


### Authentication

The Client instance supports both authenticated and unauthenticated interaction with the Travis CI API. **Note** that both Pro and Enterprise plans will require almost all API calls to be authenticated.


#### Unuathenticated

It is possible to use the client unauthenticated. However some resources won't be accesible.

```go
unauthClient := travis.NewClient(travis.TRAVIS_API_DEFAULT_URL, "")
builds, _, _, resp, err := unauthClient.Builds.ListFromRepository("mygithubuser/myopensourceproject", nil)
// Do something with your builds

_, err := unauthClient.Jobs.Cancel(12345)
if err != nil {
        // This operation is unavailable in unauthenticated mode and will
        // throw an error.
}
```

#### Authenticated

The Client instance supports authentication with both Travis token and Github token.

```go
authClient := travis.NewClient(travis.TRAVIS_API_DEFAULT_URL, "mytravistoken")
builds, _, _, resp, err := authClient.Builds.ListFromRepository("mygithubuser/myopensourceproject",
nil)
// Do something with your builds

_, err := unauthClient.Jobs.Cancel(12345)
// Your job is succesfully canceled
```

However, authentication with a Github token will require and extra step (and request).

```go
authWithGithubClient := travis.NewClient(travis.TRAVIS_API_DEFAULT_URL, "")
// authWithGithubClient.IsAuthenticated() will return false

err := authWithGithubClient.Authentication.UsingGithubToken("mygithubtoken")
if err != nil {
        log.Fatal(err)
}
// authWithGithubClient.IsAuthenticated()  will return true

builds, _, _, resp, err := authClient.Builds.ListFromRepository("mygithubuser/myopensourceproject",
nil)
// Do something with your builds
```


### Pagination

The services support resource pagination through the `ListOption` type. Every services `Option` type implements the `ListOption` type.

```go
client := travis.NewClient(travis.TRAVIS_API_DEFAULT_URL, "mysuperdupertoken")
opt := &travis.BuildListOptions{}

for {
        travisBuilds, _, _, _, err := tc.Builds.ListFromRepository(target, opt)
        if err != nil {
                log.Fatal(err)
        }

        // Do something with the builds

        opt.GetNextPage(travisBuilds)
        if opt.AfterNumber <= 1 {  // Travis CI resources are one-indexed (not zero-indexed)
                break
        }
}
```

## Roadmap

This library is being initially developed for internal applications at
[Ableton](http://ableton.com). Therefore API methods are implemented in the order that they are
needed by our applications. Eventually, we would like to cover the entire
Travis API, so contributions are of course [always welcome][contributing].

[contributing]: CONTRIBUTING.md

## Maintainers

* [@mst-ableton](https://github.com/mst-ableton)

## Maintainers-Emeritus

* Theo Crevon <theo.crevon@ableton.com>

## Disclaimer

This library design is heavily inspired from the amazing Google's [go-github](https://github.com/google/go-github) library. Some pieces of code have been directly extracted from there too. Therefore any obvious similarities would not be adventitious.

## License

This library is distributed under the BSD-style license found in the [LICENSE](./LICENSE)
file.
