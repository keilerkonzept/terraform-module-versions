# terraform-module-versions

Checks for updates of external terraform modules referenced in given `.tf` files. Outputs JSONL (one JSON object per line), or Markdown tables (`-pretty, -p`).

Supported module sources:
- **Git** (`git::`, `github.com/...`, or `git@github.com:...` values for `source`) with SemVer tags.
- **Terraform Registry** (public `<NAMESPACE>/<NAME>/<PROVIDER>` and private `<HOSTNAME>/<NAMESPACE>/<NAME>/<PROVIDER>`)

```sh
$ terraform-module-versions -updates -pretty examples/main.tf
```

```markdown
| UPDATE? |        NAME         | CONSTRAINT | VERSION | LATEST MATCHING | LATEST |
|---------|---------------------|------------|---------|-----------------|--------|
| Y       | example_git_ssh     | ~> 0.10    | 0.10.0  | 0.11.2          | 0.11.2 |
| ?       | consul              | > 0.1.0    |         | 0.7.3           | 0.7.3  |
| (Y)     | consul_github_ssh   | 0.1.0      | 0.1.0   |                 | 0.7.3  |
|         | consul_github_https | 0.7.3      |         | 0.7.3           | 0.7.3  |
```

## Contents

- [Examples](#examples)
  - [List modules with their current versions](#list-modules-with-their-current-versions)
  - [Check for module updates](#check-for-module-updates)
  - [Check for updates of specific modules](#check-for-updates-of-specific-modules)
- [Get it](#get-it)
- [Usage](#usage)

## Examples

```sh
$ cat examples/main.tf
```

```terraform
module "consul" {
  source = "hashicorp/consul/aws"
  version = "> 0.1.0"
}

module "consul_github_https" {
  source = "github.com/hashicorp/terraform-aws-consul"
  version = "0.7.3"
}

module "consul_github_ssh" {
  source = "git@github.com:hashicorp/terraform-aws-consul?ref=0.1.0"
  version = "0.1.0"
}

module "example_git_ssh" {
  source = "git::ssh://git@github.com/sgreben/terraform-module-versions?ref=0.10.0"
  version = "~> 0.10"
}
```

### List modules with their current versions

```sh
# default operation: list modules with their current versions and version constraints (if specified)
$ terraform-module-versions examples/main.tf
```

```json
{
  "path": "examples/main.tf",
  "name": "consul_github_https",
  "type": "git",
  "source": "github.com/hashicorp/terraform-aws-consul",
  "constraint": "0.7.3"
}
{
  "path": "examples/main.tf",
  "name": "consul_github_ssh",
  "type": "git",
  "source": "git@github.com:hashicorp/terraform-aws-consul?ref=0.1.0",
  "constraint": "0.1.0",
  "version": "0.1.0"
}
{
  "path": "examples/main.tf",
  "name": "example_git_ssh",
  "type": "git",
  "source": "git::ssh://git@github.com/sgreben/terraform-module-versions?ref=0.10.0",
  "constraint": "~> 0.10",
  "version": "0.10.0"
}
{
  "path": "examples/main.tf",
  "name": "consul",
  "type": "registry",
  "source": "hashicorp/consul/aws",
  "constraint": "> 0.1.0"
}
```

### Check for module updates

```sh
# -update: check for module updates from (usually) remote sources
$ terraform-module-versions -updates examples/main.tf
```

```json
{
  "path": "examples/main.tf",
  "name": "consul",
  "constraint": "> 0.1.0",
  "constraintUpdate": true,
  "latestMatching": "0.7.3",
  "latestOverall": "0.7.3"
}
{
  "path": "examples/main.tf",
  "name": "consul_github_https",
  "constraint": "0.7.3",
  "latestMatching": "0.7.3",
  "latestOverall": "0.7.3"
}
{
  "path": "examples/main.tf",
  "name": "consul_github_ssh",
  "constraint": "0.1.0",
  "version": "0.1.0",
  "latestOverall": "0.7.3",
  "nonMatchingUpdate": true
}
{
  "path": "examples/main.tf",
  "name": "example_git_ssh",
  "constraint": "~> 0.10",
  "version": "0.10.0",
  "constraintUpdate": true,
  "latestMatching": "0.11.2",
  "matchingUpdate": true,
  "latestOverall": "0.11.2"
}
```

### Check for updates of specific modules

```sh
# -update and -module: check for updates of specific modules
$ terraform-module-versions -updates -module=consul_github_https examples/main.tf
```

```json
{
  "path": "examples/main.tf",
  "name": "consul_github_https",
  "constraint": "0.7.3",
  "latestMatching": "0.7.3",
  "latestOverall": "0.7.3"
}
{
  "path": "examples/main.tf",
  "name": "consul_github_ssh",
  "constraint": "0.1.0",
  "version": "0.1.0",
  "latestOverall": "0.7.3",
  "nonMatchingUpdate": true
}
```

## Get it

Using go get:

```bash
go get -u github.com/sgreben/terraform-module-versions
```

Or [download the binary for your platform](https://github.com/sgreben/terraform-module-versions/releases/latest) from the releases page.

## Usage

```text
terraform-module-versions [PATHS...]

Usage of terraform-module-versions:
  -module value
    	include this module (may be specified repeatedly. by default, all modules are included)
  -p	(alias for -pretty)
  -pretty
    	human-readable output
  -q	(alias for -quiet)
  -quiet
    	suppress log output (stderr)
  -u	(alias for -updates)
  -update
    	(alias for -updates)
  -updates
    	check for updates
  -version
    	print version and exit
```
