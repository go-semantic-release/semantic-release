## 🚨 Upgrade to semantic-release v2 🚨

`semantic-release` v2 is now available. If you run into any problems, please create a [GitHub issue](https://github.com/go-semantic-release/semantic-release/issues/new). You can always downgrade to v1 with:

```
curl -SL https://get-release.xyz/semantic-release/linux/amd64/1.22.1 -o ./semantic-release && chmod +x ./semantic-release
```

### Breaking changes

* It is now necessary to use **double dashes** for CLI flags (e.g. `--dry`)
* **Travis CI** support has been **removed**
* Some CLI flags have changed:

|             v1             |                        v2                        |
|:--------------------------:|:------------------------------------------------:|
| `-vf`                      | `-f`                                             |
| `--noci`                   | `--no-ci`                                        |
| `--ghe-host <host>`        | `--provider-opt "github_enterprise_host=<host>"` |
| `--travis-com`             | _removed_                                        |
| `--gitlab`                 | `--provider gitlab`                              |
| `--gitlab-base-url <url>`  | `--provider-opt "gitlab_baseurl=<url>"`          |
| `--gitlab-project-id <id>` | `--provider-opt "gitlab_projectid=<id>"`         |
| `--slug`                   | `--provider-opt "slug=<url>"`                    |
