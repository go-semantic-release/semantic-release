/*
go-travis is a Go client library for accessing the Travis CI API.

Installation

go-travis requires Go version 1.1 or greater.

    $ go get github.com/Ableton/go-travis


Usage

Interaction with the Travis CI API is done through a Client instance.

    import travis "github.com/AbletonAppDev/go-travis"
    client := travis.NewClient(travis.TRAVIS_API_DEFAULT_URL, "AQFvXR7r88s2Db5-dMYo3g")

Constructing it with the NewClient helper requires two arguments:

First, the Travis CI API URL you wish to communicate with. Different Travis CI plans are accessed through different URLs.
go-travis exposes constants for these URLs:

    TRAVIS_API_DEFAULT_URL -> default api.travis-ci.org endpoint for the free Travis "Open Source" plan.
    TRAVIS_API_PRO_URL -> the api.travis-ci.com endpoint for the paid Travis pro plans.

Second, a Travis CI token with which to authenticate.
If you wish to run requests unauthenticated, pass an empty string.
It is possible at any time to authenticate the Client instance with a Travis token or a Github token.
For more information see [Authentication]().

Services

The Client instance's Service attributes provide access to Travis CI API resources.

    opt := &travis.BuildListOptions{EventType: "pull request"}
    builds, response, err := client.Builds.ListFromRepository("mygithubuser/mygithubrepo", opt)
    if err != nil {
                log.Fatal(err)
    }

Service methods will often take an Option (sub-)type instance as input. These types, like
BuildListOptions allow narrowing and filtering your requests.


Authentication

The Client instance supports both authenticated and unauthenticated interactions with the Travis CI
API. Note that both Pro and Enterprise will require almost all API calls to be authenticated.

It is possible to use the client unauthenticated. However some resources won't be accesible.

    unauthClient := travis.NewClient(travis.TRAVIS_API_DEFAULT_URL, "")
    builds, _, _, resp, err :=
    unauthClient.Builds.ListFromRepository("mygithubuser/myopensourceproject", nil)
    // Do something with your builds

    _, err := unauthClient.Jobs.Cancel(12345)
    if err != nil {
                // This operation is unavailable in unauthenticated mode and will
                        // throw an error at you.
    }

The Client instance supports authentication with both Travis token and Github token.

    authClient := travis.NewClient(travis.TRAVIS_API_DEFAULT_URL, "mytravistoken")
    builds, _, _, resp, err := authClient.Builds.ListFromRepository("mygithubuser/myopensourceproject",
    nil)
    // Do something with your builds

    _, err := unauthClient.Jobs.Cancel(12345)
    // Your job is succesfully canceled

However, authentication with a Github token will require and extra step (and request).

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

Pagination

The services support resource pagination through the ListOption type. Every services `Option` type
implements the ListOption type.

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
*/
package travis
