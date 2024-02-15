# :package::rocket: semantic-release
[![CI](https://github.com/go-semantic-release/semantic-release/workflows/CI/badge.svg?branch=master)](https://github.com/go-semantic-release/semantic-release/actions?query=workflow%3ACI+branch%3Amaster)
[![pipeline status](https://gitlab.com/go-semantic-release/semantic-release/badges/master/pipeline.svg)](https://gitlab.com/go-semantic-release/semantic-release/pipelines)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-semantic-release/semantic-release)](https://goreportcard.com/report/github.com/go-semantic-release/semantic-release)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/go-semantic-release/semantic-release/v2)](https://pkg.go.dev/github.com/go-semantic-release/semantic-release/v2)

> fully automated package/module/image publishing

This project aims to be an alternative to the original [semantic-release](https://github.com/semantic-release/semantic-release) implementation. Using Go, `semantic-release` can be installed by downloading a single binary and is, therefore, easier to install and does not require Node.js and npm. Furthermore, `semantic-release` has a built-in plugin system that allows to extend and customize its functionality.

### Features

- Automated version and release management
- No external dependencies required
- Runs on Linux, macOS and Windows
- Fully extensible via plugins
- Automated changelog generation
- Supports GitHub, GitLab and git
- Support for maintaining multiple major version releases

## How does it work?
Instead of writing [meaningless commit messages](http://whatthecommit.com/), we can take our time to think about the changes in the codebase and write them down. Following the [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) specification it is then possible to generate a helpful changelog and to derive the next semantic version number from them.

When `semantic-release` is setup it will do that after every successful continuous integration build of your default branch and publish the new version for you. This way no human is directly involved in the release process and your releases are guaranteed to be [unromantic and unsentimental](http://sentimentalversioning.org/).

_Source: [semantic-release/semantic-release#how-does-it-work](https://github.com/semantic-release/semantic-release#how-does-it-work)_

You can enforce semantic commit messages using [a git hook](https://github.com/hazcod/semantic-commit-hook).

## Installation


### Option 1: Use the go-semantic-release GitHub Action ([go-semantic-release/action](https://github.com/go-semantic-release/action))

### Option 2: Install manually

```bash
curl -SL https://get-release.xyz/semantic-release/linux/amd64 -o ./semantic-release && chmod +x ./semantic-release
```

### Option 3: Install via npm

```bash
npm install -g go-semantic-release
```

## Examples

### Releasing a Go application with GitHub Actions
Full example can be found at [go-semantic-release/example-go-application](https://github.com/go-semantic-release/example-go-application).

Example [.github/workflows/ci.yml](https://github.com/go-semantic-release/example-go-application/blob/main/.github/workflows/ci.yml) config:
```yaml
name: CI
on:
  push:
    branches:
      - '**'
  pull_request:
    branches:
      - '**'
jobs:
  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v3
        with:
          go-version: 1.19
      - uses: golangci/golangci-lint-action@v3
  test:
    runs-on: ubuntu-latest
    needs: lint
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v3
        with:
          go-version: 1.19
      - run: go test -v ./...
  release:
    runs-on: ubuntu-latest
    needs: test
    permissions:
      contents: write
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v3
        with:
          go-version: 1.19
      - uses: go-semantic-release/action@v1
        with:
          hooks: goreleaser
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

### Example GitLab CI Config

#### Job Token
If you do not provide a PAT the [job token](https://docs.gitlab.com/ee/ci/jobs/ci_job_token.html) will be used.
This restricted token can create releases but not read commits. The [git strategy](https://docs.gitlab.com/ee/ci/runners/configure_runners.html#git-strategy) must be set to clone so that we can read the commits from the repository. See example below

#### Personal Access Token
Create a new Gitlab personal access token with the `api` scope [here](https://gitlab.com/-/user_settings/personal_access_tokens).
Ensure the CI variable is protected and masked as the `GITLAB_TOKEN` has a lot of rights. There is an open issue for project specific [tokens](https://gitlab.com/gitlab-org/gitlab/issues/756)
You can set the GitLab token via the `GITLAB_TOKEN` environment variable or the `-token` flag.

.gitlab-ci.yml
```yml
  variables:
    # Only needs if using job token
    GIT_STRATEGY: clone
 stages:
  # other stages
  - release

release:
  image:
    name: registry.gitlab.com/go-semantic-release/semantic-release:latest
    entrypoint: [""]
  stage: release
  # when: manual # Add this if you want to manually create releases
  only:
    - master
  script:
    - semantic-release # Add --allow-no-changes if you want to create a release for each push
```

### Releasing a Go application with GitLab CI
The full example can be found at https://gitlab.com/go-semantic-release/example-go-application.

Example [.gitlab-ci.yml](https://gitlab.com/go-semantic-release/example-go-application/-/blob/main/.gitlab-ci.yml) config:
```yaml
image: golang:1.19

stages:
  - test
  - release

test:
  stage: test
  except:
    - tags
  script:
    - go test -v ./...
    - go build ./
    - ./example-go-application

release:
  stage: release
  only:
    - main
  script:
    - curl -SL https://get-release.xyz/semantic-release/linux/amd64 -o ./semantic-release && chmod +x ./semantic-release
    - ./semantic-release --hooks goreleaser
```

## Plugin System

Since v2, semantic-release is equipped with a plugin system. The plugins are standalone binaries that use [hashicorp/go-plugin](https://github.com/hashicorp/go-plugin) as a plugin library. `semantic-release` automatically downloads the necessary plugins if they don't exist locally. The plugins are stored in the `.semrel` directory of the current working directory in the following format: `.semrel/<os>_<arch>/<plugin name>/<version>/`. The go-semantic-release plugins registry (https://registry.go-semantic-release.xyz/) is used to resolve plugins and to download the correct binary. With the new [plugin-registry](https://github.com/go-semantic-release/plugin-registry) service the API also supports batch requests to resolve multiple plugins at once and caching of the plugins.

### Running semantic-release in an air-gapped environment
As plugins are only downloaded if they do not exist in the `.semrel` folder, it is possible to download the plugins and archive the `.semrel` folder. This way it is possible to run `semantic-release` in an air-gapped environment.

```bash
# specify all required plugins and download them
./semantic-release --download-plugins --show-progress --provider github --ci-condition github --hooks goreleaser
# archive the .semrel folder
tar -czvf ./semrel-plugins.tgz .semrel/

# copy the archive to the air-gapped environment

# extract the archive
tar -xzvf ./semrel-plugins.tgz
# run semantic-release
./semantic-release --provider github --condition github --hooks goreleaser
```

### Plugins

* Provider ([Docs](https://pkg.go.dev/github.com/go-semantic-release/semantic-release/v2/pkg/provider?tab=doc#Provider))
  * [GitHub](https://github.com/go-semantic-release/provider-github)
  * [GitLab](https://github.com/go-semantic-release/provider-gitlab)
  * [Git](https://github.com/go-semantic-release/provider-git)
* CI Condition ([Docs](https://pkg.go.dev/github.com/go-semantic-release/semantic-release/v2/pkg/condition?tab=doc#CICondition))
  * [GitHub Actions](https://github.com/go-semantic-release/condition-github)
  * [GitLab CI](https://github.com/go-semantic-release/condition-gitlab)
  * [Bitbucket Pipelines](https://github.com/go-semantic-release/condition-bitbucket)
  * [Default](https://github.com/go-semantic-release/condition-default)
* Commit Analyzer ([Docs](https://pkg.go.dev/github.com/go-semantic-release/semantic-release/v2/pkg/analyzer?tab=doc#CommitAnalyzer))
  * [Conventional Commits](https://github.com/go-semantic-release/commit-analyzer-cz)
* Changelog Generator ([Docs](https://pkg.go.dev/github.com/go-semantic-release/semantic-release/v2/pkg/generator?tab=doc#ChangelogGenerator))
  * [Default](https://github.com/go-semantic-release/changelog-generator-default)
* Hooks ([Docs](https://pkg.go.dev/github.com/go-semantic-release/semantic-release/v2/pkg/hooks?tab=doc#Hooks))
  * [GoReleaser](https://github.com/go-semantic-release/hooks-goreleaser)
  * [npm-binary-releaser](https://github.com/go-semantic-release/hooks-npm-binary-releaser)
  * [plugin-registry-update](https://github.com/go-semantic-release/hooks-plugin-registry-update)
  * [exec](https://github.com/go-semantic-release/hooks-exec)
* Files Updater ([Docs](https://pkg.go.dev/github.com/go-semantic-release/semantic-release/v2/pkg/updater?tab=doc#FilesUpdater))
  * [npm](https://github.com/go-semantic-release/files-updater-npm)
  * [helm](https://github.com/go-semantic-release/files-updater-helm)

### Configuration

Plugins can be configured using CLI flags or the `.semrelrc` config file. By using a `@` sign after the plugin name, the required version of the plugin can be specified. Otherwise, any locally installed version will be used. If the plugin does not exist locally, the latest version will be downloaded. This is an example of the `.semrelrc` config file:

```json
{
  "plugins": {
    "commit-analyzer": {
      "name": "default@^1.0.0"
    },
    "ci-condition": {
      "name": "default"
    },
    "changelog-generator": {
      "name": "default",
      "options": {
        "emojis": "true"
      }
    },
    "provider": {
      "name": "gitlab",
      "options": {
        "gitlab_projectid": "123456"
      }
    },
    "files-updater": {
      "names": ["npm"]
    }
  }
}
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

Copyright © 2023 [Christoph Witzko](https://twitter.com/christophwitzko)
