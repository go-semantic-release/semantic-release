# :package::rocket: semantic-release
[![Build Status](https://travis-ci.org/go-semantic-release/semantic-release.svg?branch=master)](https://travis-ci.org/go-semantic-release/semantic-release)

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

## Example GitHub Release

### GitHub token
It is necessary to create a new GitHub token with the `repo` or `public_repo` scope [here](https://github.com/settings/tokens/new).
You can set the GitHub token via the `GITHUB_TOKEN` environment variable or the `-token` flag.

__.travis.yml__
```yml
language: go
go:
  - 1.x
install:
  - curl -SL https://get-release.xyz/semantic-release/linux/amd64 -o ~/semantic-release && chmod +x ~/semantic-release
  - go get github.com/mitchellh/gox
  - go get github.com/tcnksm/ghr
after_success:
  - ./release
notifications:
  email: false
```

__release__
```bash
#!/bin/bash
set -e

~/semantic-release -ghr -vf
export VERSION=$(cat .version)
gox -ldflags="-s -w" -output="bin/{{.Dir}}_v"$VERSION"_{{.OS}}_{{.Arch}}"
ghr $(cat .ghr) bin/

```

## Example Docker Hub

The environment variables GITHUB_TOKEN, DOCKER_USERNAME and DOCKER_PASSWORD must be set.

__.travis.yml__
```yml
language: go
services:
  - docker
go:
  - 1.x
install:
  - curl -SL https://get-release.xyz/semantic-release/linux/amd64 -o ~/semantic-release && chmod +x ~/semantic-release
after_success:
  - ./release
notifications:
  email: false
```
__release__
```bash
#!/bin/bash

set -e

# run semantic-release
~/semantic-release -vf
export VERSION=$(cat .version)

# docker build
export IMAGE_NAME="user/imagename"
export IMAGE_NAME_VERSION="$IMAGE_NAME:$VERSION"

docker build -t $IMAGE_NAME_VERSION .
docker tag $IMAGE_NAME_VERSION $IMAGE_NAME

# push to docker hub
docker login -u $DOCKER_USERNAME -p $DOCKER_PASSWORD
docker push $IMAGE_NAME_VERSION
docker push $IMAGE_NAME

```

## Example npm

The environment variables GITHUB_TOKEN and NPM_TOKEN must be set.

__.travis.yml__
```yml
language: node_js
cache:
  directories:
    - node_modules
notifications:
  email: false
node_js:
  - '7'
  - '6'
  - '4'
after_success:
  - curl -SL https://get-release.xyz/semantic-release/linux/amd64 -o ~/semantic-release && chmod +x ~/semantic-release
  - ~/semantic-release -update package.json && npm publish
branches:
  except:
    - /^v\d+\.\d+\.\d+$/
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

Copyright © 2017 [Christoph Witzko](https://twitter.com/christophwitzko)
