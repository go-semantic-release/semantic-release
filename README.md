# :package::rocket: semantic-release
[![CI](https://github.com/go-semantic-release/semantic-release/workflows/CI/badge.svg?branch=master)](https://github.com/go-semantic-release/semantic-release/actions?query=workflow%3ACI+branch%3Amaster)
[![pipeline status](https://gitlab.com/go-semantic-release/semantic-release/badges/master/pipeline.svg)](https://gitlab.com/go-semantic-release/semantic-release/pipelines)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-semantic-release/semantic-release)](https://goreportcard.com/report/github.com/go-semantic-release/semantic-release)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/go-semantic-release/semantic-release/v2)](https://pkg.go.dev/github.com/go-semantic-release/semantic-release/v2)

> fully automated package/module/image publishing

A more lightweight and standalone version of [semantic-release](https://github.com/semantic-release/semantic-release).

## How does it work?
Instead of writing [meaningless commit messages](http://whatthecommit.com/), we can take our time to think about the changes in the codebase and write them down. Following the [AngularJS Commit Message Conventions](https://docs.google.com/document/d/1QrDFcIiPjSLDn3EL15IJygNPiHORgU1_OOAqWjiDU5Y/edit) it is then possible to generate a helpful changelog and to derive the next semantic version number from them.

When `semantic-release` is setup it will do that after every successful continuous integration build of your master branch (or any other branch you specify) and publish the new version for you. This way no human is directly involved in the release process and your releases are guaranteed to be [unromantic and unsentimental](http://sentimentalversioning.org/).

_Source: [semantic-release/semantic-release#how-does-it-work](https://github.com/semantic-release/semantic-release#how-does-it-work)_

You can enforce semantic commit messages using [a git hook](https://github.com/hazcod/semantic-commit-hook).

## Installation
__Install the latest version of semantic-release__
```bash
curl -SL https://get-release.xyz/semantic-release/linux/amd64 -o ./semantic-release && chmod +x ./semantic-release
```

## Example GitHub Actions

For examples, look at the [go-semantic-release GitHub Action](https://github.com/go-semantic-release/action).

## Example GitLab CI Config

### GitLab token
It is necessary to create a new Gitlab personal access token with the `api` scope [here](https://gitlab.com/profile/personal_access_tokens).
Ensure the CI variable is protected and masked as the `GITLAB_TOKEN` has a lot of rights. There is an open issue for project specific [tokens](https://gitlab.com/gitlab-org/gitlab/issues/756)
You can set the GitLab token via the `GITLAB_TOKEN` environment variable or the `-token` flag.

.gitlab-ci.yml
```yml
 stages:
  # other stages
  - release

release:
  image: registry.gitlab.com/go-semantic-release/semantic-release:latest # Replace this with the current release
  stage: release
  # Remove this if you want a release created for each push to master
  when: manual
  only:
    - master
  script:
    - release
```


## Beta release support
Beta release support empowers you to release beta, rc, etc. versions with `semantic-release` (e.g. v2.0.0-beta.1). To enable this feature you need to create a new branch (e.g. beta/v2) and check in a `.semrelrc` file with the following content:
```
{
  "maintainedVersion": "2-beta"
}
```
If you commit to this branch a new incremental pre-release is created everytime you push. (2.0.0-beta.1, 2.0.0-beta.2, ...)

## Licence

The [MIT License (MIT)](http://opensource.org/licenses/MIT)

Copyright © 2020 [Christoph Witzko](https://twitter.com/christophwitzko)
