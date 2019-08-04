# ${APP}

Checks for updates of external terraform modules referenced in given `.tf` files. Outputs JSONL (one JSON object per line). Pretty-printed output TBD.

Supported module sources:
- **Git** (`git::`, `github.com/...`, or `git@github.com:...` values for `source`) with SemVer tags.
- **Terraform Registry** (public `<NAMESPACE>/<NAME>/<PROVIDER>` and private `<HOSTNAME>/<NAMESPACE>/<NAME>/<PROVIDER>`)

## Contents

- [Example](#example)
- [Get it](#get-it)
- [Usage](#usage)

## Example

```sh
$ cat main.tf
```

```terraform
module "consul" {
  source = "hashicorp/consul/aws"
  version = "> 0.1.0"
}

module "example" {
  source = "git::ssh://git@github.com/sgreben/${APP}?ref=0.10.0"
  version = "~> 0.10"
}
```

```sh
# default operation: list current modules with their versions and version constraints (if specified)
$ ${APP} main.tf
```

```json
{"path":"main.tf","name":"consul","source":"hashicorp/consul/aws","versionConstraint":"> 0.1.0","type":"terraform-registry"}
{"path":"main.tf","name":"example","source":"git::ssh://git@github.com/sgreben/${APP}?ref=0.10.0","version":"0.10.0","versionConstraint":"~> 0.10","type":"git"}
```

```sh
# -update: check for module updates from (usually) remote sources
$ ${APP} -updates main.tf
```

```json
{"path":"main.tf","name":"consul","source":"hashicorp/consul/aws","versionConstraint":"> 0.1.0","type":"terraform-registry","latestMatchingUpdate":"0.7.3","updates":["0.1.1","0.1.2","0.2.0","0.2.1","0.2.2","0.3.0","0.3.1","0.3.2","0.3.3","0.3.4","0.3.5","0.3.6","0.3.7","0.3.8","0.3.9","0.3.10","0.4.0","0.4.1","0.4.2","0.4.3","0.4.4","0.4.5","0.5.0","0.6.0","0.6.1","0.7.0","0.7.1","0.7.2","0.7.3"],"hasMinorUpdate":true,"hasPatchUpdate":true}
{"path":"main.tf","name":"example","source":"git::ssh://git@github.com/sgreben/${APP}?ref=0.10.0", "version":"0.10.0","versionConstraint":"~> 0.10","type":"git","latestMatchingUpdate":"0.11.0","updates":["0.11.0"],"hasMinorUpdate":true}

```

## Get it

Using go get:

```bash
go get -u github.com/sgreben/${APP}
```

Or [download the binary](https://github.com/sgreben/${APP}/releases/latest) from the releases page.

```bash
# Linux
curl -LO https://github.com/sgreben/${APP}/releases/download/${VERSION}/${APP}_${VERSION}_linux_x86_64.zip
unzip ${APP}_${VERSION}_linux_x86_64.zip

# OS X
curl -LO https://github.com/sgreben/${APP}/releases/download/${VERSION}/${APP}_${VERSION}_osx_x86_64.zip
unzip ${APP}_${VERSION}_osx_x86_64.zip

# Windows
curl -LO https://github.com/sgreben/${APP}/releases/download/${VERSION}/${APP}_${VERSION}_windows_x86_64.zip
unzip ${APP}_${VERSION}_windows_x86_64.zip
```

## Usage

```text
${APP} [PATHS...]

$USAGE
```
