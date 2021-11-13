# terraform-module-versions

Checks for updates of external terraform modules referenced in given Terraform source. Outputs Markdown tables by default, as well as JSONL (`-o jsonl`, one JSON object per line), JSON (`-o json`), and JUnit XML (`-o junit`).

Supported module sources:

- **Git** with SemVer tags
  - `git::...`
  - `github.com/...`
  - `git@github.com:...`
- **Terraform Registry**
  - public `<NAMESPACE>/<NAME>/<PROVIDER>`
  - private `<HOSTNAME>/<NAMESPACE>/<NAME>/<PROVIDER>`

## Example

```sh
$ terraform-module-versions check examples
```

```markdown
| UPDATE? |               NAME               | CONSTRAINT | VERSION | LATEST MATCHING | LATEST  |
|---------|----------------------------------|------------|---------|-----------------|---------|
| (Y)     | consul                           | ~0.7.3     |         | 0.7.11          | 0.11.0  |
| (Y)     | consul_github_https              | 0.8.0      | v0.8.0  |                 | v0.11.0 |
| (Y)     | consul_github_https_missing_ref  | 0.7.3      |         | v0.7.3          | v0.11.0 |
| (Y)     | consul_github_https_no_ref       |            |         |                 | v0.11.0 |
| Y       | consul_github_ssh                | ~0.1.0     | 0.1.0   | v0.1.2          | v0.11.0 |
| (Y)     | example_git_scp                  | ~> 0.12    | 0.12.0  |                 | 3.1.5   |
| (Y)     | example_git_ssh_branch           |            | master  |                 | 3.1.5   |
| (Y)     | example_with_prerelease_versions |            | v0.22.2 |                 | v0.22.3 |
```

## Contents

- [terraform-module-versions](#app)
  - [Example](#example)
  - [Contents](#contents)
  - [Examples](#examples)
    - [List modules with their current versions](#list-modules-with-their-current-versions)
    - [Check for module updates](#check-for-module-updates)
    - [Check for module updates using Github Token authentication](#check-for-module-updates-using-github-token-authentication)
    - [Check for updates of specific modules](#check-for-updates-of-specific-modules)
  - [Get it](#get-it)
  - [Usage](#usage)
    - [`list`](#list)
    - [`check`](#check)

## Examples

```sh
$ cat examples/main.tf
```

```terraform
module "consul" {
  source = "hashicorp/consul/aws"
  version = "~0.7.3"
}

module "consul_github_https_missing_ref" {
  source = "github.com/hashicorp/terraform-aws-consul"
  version = "0.7.3"
}

module "consul_github_https_no_ref" {
  source = "github.com/hashicorp/terraform-aws-consul"
}

module "consul_github_https" {
  source = "github.com/hashicorp/terraform-aws-consul?ref=v0.8.0"
  version = "0.8.0"
}

module "consul_github_ssh" {
  source = "git@github.com:hashicorp/terraform-aws-consul?ref=0.1.0"
  version = "~0.1.0"
}

module "example_git_ssh_branch" {
  source = "git::ssh://git@github.com/keilerkonzept/terraform-module-versions?ref=master"
}

module "example_git_scp" {
  source = "git::git@github.com:keilerkonzept/terraform-module-versions?ref=0.12.0"
  version = "~> 0.12"
}

module "example_with_prerelease_versions" {
  source = "git@github.com:kubernetes/api.git?ref=v0.22.2"
}

module "local" {
  source = "./local"
}

variable "_0_15_sensitive_example" {
  type      = string
  sensitive = true
}

output "0_15_sensitive_example" {
  value = "foo-${var._0_15_sensitive_example}"
  sensitive = true
}

output "0_15_nonsensitive_example" {
  value = nonsensitive(var._0_15_sensitive_example)
}
```

### List modules with their current versions

```sh
# list modules with their current versions and version constraints (if specified)
$ terraform-module-versions list examples
```

|   TYPE   |               NAME               | CONSTRAINT | VERSION |                                    SOURCE                                    |
|----------|----------------------------------|------------|---------|------------------------------------------------------------------------------|
| registry | consul                           | ~0.7.3     |         | hashicorp/consul/aws                                                         |
| local    | local                            |            |         | ./local                                                                      |
| git      | example_with_prerelease_versions |            | v0.22.2 | git@github.com:kubernetes/api.git?ref=v0.22.2                                |
| git      | example_git_ssh_branch           |            | master  | git::ssh://git@github.com/keilerkonzept/terraform-module-versions?ref=master |
| git      | example_git_scp                  | ~> 0.12    | 0.12.0  | git::git@github.com:keilerkonzept/terraform-module-versions?ref=0.12.0       |
| git      | consul_github_ssh                | ~0.1.0     | 0.1.0   | git@github.com:hashicorp/terraform-aws-consul?ref=0.1.0                      |
| git      | consul_github_https_no_ref       |            |         | github.com/hashicorp/terraform-aws-consul                                    |
| git      | consul_github_https_missing_ref  | 0.7.3      |         | github.com/hashicorp/terraform-aws-consul                                    |
| git      | consul_github_https              | 0.8.0      | v0.8.0  | github.com/hashicorp/terraform-aws-consul?ref=v0.8.0                         |

with `-o json`:

```json
[
  {
    "path": "examples/main.tf",
    "name": "consul",
    "type": "registry",
    "source": "hashicorp/consul/aws",
    "constraint": "~0.7.3"
  },
  {
    "path": "examples/main.tf",
    "name": "consul_github_https",
    "type": "git",
    "source": "github.com/hashicorp/terraform-aws-consul?ref=v0.8.0",
    "constraint": "0.8.0",
    "version": "v0.8.0"
  },
  {
    "path": "examples/main.tf",
    "name": "consul_github_https_missing_ref",
    "type": "git",
    "source": "github.com/hashicorp/terraform-aws-consul",
    "constraint": "0.7.3"
  },
  {
    "path": "examples/main.tf",
    "name": "consul_github_https_no_ref",
    "type": "git",
    "source": "github.com/hashicorp/terraform-aws-consul"
  },
  {
    "path": "examples/main.tf",
    "name": "consul_github_ssh",
    "type": "git",
    "source": "git@github.com:hashicorp/terraform-aws-consul?ref=0.1.0",
    "constraint": "~0.1.0",
    "version": "0.1.0"
  },
  {
    "path": "examples/main.tf",
    "name": "example_git_scp",
    "type": "git",
    "source": "git::git@github.com:keilerkonzept/terraform-module-versions?ref=0.12.0",
    "constraint": "~> 0.12",
    "version": "0.12.0"
  },
  {
    "path": "examples/main.tf",
    "name": "example_git_ssh_branch",
    "type": "git",
    "source": "git::ssh://git@github.com/keilerkonzept/terraform-module-versions?ref=master",
    "version": "master"
  },
  {
    "path": "examples/main.tf",
    "name": "example_with_prerelease_versions",
    "type": "git",
    "source": "git@github.com:kubernetes/api.git?ref=v0.22.2",
    "version": "v0.22.2"
  },
  {
    "path": "examples/main.tf",
    "name": "local",
    "type": "local",
    "source": "./local"
  }
]
```

### Check for module updates

```sh
# check: check for module updates from (usually) remote sources
$ terraform-module-versions check examples
```

| UPDATE? |               NAME               | CONSTRAINT | VERSION | LATEST MATCHING | LATEST  |
|---------|----------------------------------|------------|---------|-----------------|---------|
| (Y)     | consul                           | ~0.7.3     |         | 0.7.11          | 0.11.0  |
| (Y)     | consul_github_https              | 0.8.0      | v0.8.0  |                 | v0.11.0 |
| (Y)     | consul_github_https_missing_ref  | 0.7.3      |         | v0.7.3          | v0.11.0 |
| (Y)     | consul_github_https_no_ref       |            |         |                 | v0.11.0 |
| Y       | consul_github_ssh                | ~0.1.0     | 0.1.0   | v0.1.2          | v0.11.0 |
| (Y)     | example_git_scp                  | ~> 0.12    | 0.12.0  |                 | 3.1.5   |
| (Y)     | example_git_ssh_branch           |            | master  |                 | 3.1.5   |
| (Y)     | example_with_prerelease_versions |            | v0.22.2 |                 | v0.22.3 |

with `-o json`:

```json
[
  {
    "path": "examples/main.tf",
    "name": "consul",
    "source": "hashicorp/consul/aws",
    "constraint": "~0.7.3",
    "latestMatching": "0.7.11",
    "latestOverall": "0.11.0",
    "nonMatchingUpdate": true
  },
  {
    "path": "examples/main.tf",
    "name": "consul_github_https",
    "source": "github.com/hashicorp/terraform-aws-consul?ref=v0.8.0",
    "constraint": "0.8.0",
    "version": "v0.8.0",
    "latestOverall": "v0.11.0",
    "nonMatchingUpdate": true
  },
  {
    "path": "examples/main.tf",
    "name": "consul_github_https_missing_ref",
    "source": "github.com/hashicorp/terraform-aws-consul",
    "constraint": "0.7.3",
    "latestMatching": "v0.7.3",
    "latestOverall": "v0.11.0",
    "nonMatchingUpdate": true
  },
  {
    "path": "examples/main.tf",
    "name": "consul_github_https_no_ref",
    "source": "github.com/hashicorp/terraform-aws-consul",
    "latestOverall": "v0.11.0",
    "nonMatchingUpdate": true
  },
  {
    "path": "examples/main.tf",
    "name": "consul_github_ssh",
    "source": "git@github.com:hashicorp/terraform-aws-consul?ref=0.1.0",
    "constraint": "~0.1.0",
    "version": "0.1.0",
    "latestMatching": "v0.1.2",
    "latestOverall": "v0.11.0",
    "matchingUpdate": true,
    "nonMatchingUpdate": true
  },
  {
    "path": "examples/main.tf",
    "name": "example_git_scp",
    "source": "git::git@github.com:keilerkonzept/terraform-module-versions?ref=0.12.0",
    "constraint": "~> 0.12",
    "version": "0.12.0",
    "latestOverall": "3.1.5",
    "nonMatchingUpdate": true
  },
  {
    "path": "examples/main.tf",
    "name": "example_git_ssh_branch",
    "source": "git::ssh://git@github.com/keilerkonzept/terraform-module-versions?ref=master",
    "version": "master",
    "latestOverall": "3.1.5",
    "nonMatchingUpdate": true
  },
  {
    "path": "examples/main.tf",
    "name": "example_with_prerelease_versions",
    "source": "git@github.com:kubernetes/api.git?ref=v0.22.2",
    "version": "v0.22.2",
    "latestOverall": "v0.22.3",
    "nonMatchingUpdate": true
  }
]
```

```sh
# check -all: check for updates, include up-to-date-modules in output
$ terraform-module-versions check -all examples
```

| UPDATE? |               NAME               | CONSTRAINT | VERSION | LATEST MATCHING | LATEST  |
|---------|----------------------------------|------------|---------|-----------------|---------|
| (Y)     | consul                           | ~0.7.3     |         | 0.7.11          | 0.11.0  |
| (Y)     | consul_github_https              | 0.8.0      | v0.8.0  |                 | v0.11.0 |
| (Y)     | consul_github_https_missing_ref  | 0.7.3      |         | v0.7.3          | v0.11.0 |
| (Y)     | consul_github_https_no_ref       |            |         |                 | v0.11.0 |
| Y       | consul_github_ssh                | ~0.1.0     | 0.1.0   | v0.1.2          | v0.11.0 |
| (Y)     | example_git_scp                  | ~> 0.12    | 0.12.0  |                 | 3.1.5   |
| (Y)     | example_git_ssh_branch           |            | master  |                 | 3.1.5   |
| (Y)     | example_with_prerelease_versions |            | v0.22.2 |                 | v0.22.3 |
| ?       | local                            |            |         |                 |         |

### Check for module updates using Github Token authentication

```sh
$ export GITHUB_TOKEN="<your Github PAT>"
$ terraform-module-versions check examples
```

### Check for updates of specific modules

```sh
# check -module: check for updates of specific modules
$ terraform-module-versions check -all -module=consul_github_https -module=consul_github_ssh examples
```

| UPDATE? |        NAME         | CONSTRAINT | VERSION | LATEST MATCHING | LATEST  |
|---------|---------------------|------------|---------|-----------------|---------|
| (Y)     | consul_github_https | 0.8.0      | v0.8.0  |                 | v0.11.0 |
| Y       | consul_github_ssh   | ~0.1.0     | 0.1.0   | v0.1.2          | v0.11.0 |

```sh
# check -module: check for updates of specific modules
$ terraform-module-versions check -module=consul_github_https -module=consul_github_ssh examples
```

| UPDATE? |        NAME         | CONSTRAINT | VERSION | LATEST MATCHING | LATEST  |
|---------|---------------------|------------|---------|-----------------|---------|
| (Y)     | consul_github_https | 0.8.0      | v0.8.0  |                 | v0.11.0 |
| Y       | consul_github_ssh   | ~0.1.0     | 0.1.0   | v0.1.2          | v0.11.0 |

with `-o json`:

```json
[
  {
    "path": "examples/main.tf",
    "name": "consul_github_https",
    "source": "github.com/hashicorp/terraform-aws-consul?ref=v0.8.0",
    "constraint": "0.8.0",
    "version": "v0.8.0",
    "latestOverall": "v0.11.0",
    "nonMatchingUpdate": true
  },
  {
    "path": "examples/main.tf",
    "name": "consul_github_ssh",
    "source": "git@github.com:hashicorp/terraform-aws-consul?ref=0.1.0",
    "constraint": "~0.1.0",
    "version": "0.1.0",
    "latestMatching": "v0.1.2",
    "latestOverall": "v0.11.0",
    "matchingUpdate": true,
    "nonMatchingUpdate": true
  }
]
```

## Get it

Using go get:

```bash
go get -u github.com/keilerkonzept/terraform-module-versions
```

Or [download the binary for your platform](https://github.com/keilerkonzept/terraform-module-versions/releases/latest) from the releases page.

## Usage

```text
USAGE
  terraform-module-versions [options] <subcommand>

SUBCOMMANDS
  list     List referenced terraform modules with their detected versions
  check    Check referenced terraform modules' sources for newer versions
  version  Print version and exit

FLAGS
  -o markdown       (alias for -output)
  -output markdown  output format, one of [json jsonl junit markdown markdown-wide]
  -q=false          (alias for -quiet)
  -quiet=false      suppress log output (stderr)
```

### `list`

```text
USAGE
  terraform-module-versions list [options] [<path> ...]

List referenced terraform modules with their detected versions

FLAGS
  -module ...       include this module (may be specified repeatedly. by default, all modules are included)
  -o markdown       (alias for -output)
  -output markdown  output format, one of [json jsonl junit markdown markdown-wide]
```

### `check`

```text
USAGE
  terraform-module-versions check [options] [<path> ...]

Check referenced terraform modules' sources for newer versions

FLAGS
  -H ...                                 (alias for -registry-header)
  -a=false                               (alias for -all)
  -all=false                             include modules without updates
  -any-updates-found-nonzero-exit=false  exit with a nonzero code when modules with updates are found (ignoring version constraints)
  -e=false                               (alias for -updates-found-nonzero-exit)
  -module ...                            include this module (may be specified repeatedly. by default, all modules are included)
  -n=false                               (alias for -any-updates-found-nonzero-exit)
  -o markdown                            (alias for -output)
  -output markdown                       output format, one of [json jsonl junit markdown markdown-wide]
  -pre-release=false                     include pre-release versions
  -registry-header ...                   extra HTTP headers for requests to Terraform module registries (a key/value pair KEY:VALUE, may be specified repeatedly)
  -sed=false                             generate sed statements for upgrade
  -updates-found-nonzero-exit=false      exit with a nonzero code when modules with updates matching are found (respecting version constraints)
```
