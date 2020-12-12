# terraform-module-versions

Checks for updates of external terraform modules referenced in given Terraform (0.10.x - 0.12.x) modules. Outputs Markdown tables by default, as well as JSONL (`-o jsonl`, one JSON object per line), JSON (`-o json`), and JUnit XML (`-o junit`).

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
| UPDATE? |              NAME               | CONSTRAINT | VERSION | LATEST MATCHING | LATEST |
|---------|---------------------------------|------------|---------|-----------------|--------|
| (Y)     | consul_github_https_missing_ref | 0.7.3      |         | v0.7.3          | v0.8.0 |
| (Y)     | consul_github_https_no_ref      |            |         |                 | v0.8.0 |
| Y       | consul_github_ssh               | ~0.1.0     | 0.1.0   | v0.1.2          | v0.8.0 |
| (Y)     | example_git_scp                 | ~> 0.12    | 0.12.0  |                 | 2.1.1  |
| (Y)     | example_git_ssh_branch          |            | master  |                 | 2.1.1  |
```

## Contents

- [Examples](#examples)
  - [List modules with their current versions](#list-modules-with-their-current-versions)
  - [Check for module updates](#check-for-module-updates)
  - [Check for updates of specific modules](#check-for-updates-of-specific-modules)
- [Get it](#get-it)
- [Usage](#usage)
  - [`list`](#list)
  - [`check`](#check)

## Examples

```sh
$ cat examples/main.tf examples/0.12.x.tf
```

```terraform
module "consul" {
  source = "hashicorp/consul/aws"
  version = "> 0.1.0"
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

module "example_git_scp" {
  source = "git::git@github.com:keilerkonzept/terraform-module-versions?ref=0.12.0"
  version = "~> 0.12"
}
//
// Example including 0.12.x syntax demo from the announcement (https://www.hashicorp.com/blog/announcing-terraform-0-12/)
//

module "consul_aws" {
  source  = "hashicorp/consul/aws"
  version = ">=0.5.0,<=1.0.0"
}

data "consul_key_prefix" "environment" {
  path = "apps/example/env"
}

resource "aws_elastic_beanstalk_environment" "example" {
  name        = "test_environment"
  application = "testing"

  setting {
    namespace = "aws:autoscaling:asg"
    name      = "MinSize"
    value     = "1"
  }

  dynamic "setting" {
    for_each = data.consul_key_prefix.environment.var
    content {
      namespace = "aws:elasticbeanstalk:application:environment"
      name      = setting.key
      value     = setting.value
    }
  }
}

output "environment" {
  value = {
    id = aws_elastic_beanstalk_environment.example.id
    vpc_settings = {
      for s in aws_elastic_beanstalk_environment.example.all_settings :
      s.name => s.value
      if s.namespace == "aws:ec2:vpc"
    }
  }
}
```

### List modules with their current versions

```sh
# list modules with their current versions and version constraints (if specified)
$ terraform-module-versions list examples
```

|   TYPE   |              NAME               |   CONSTRAINT    | VERSION |                                    SOURCE                                    |
|----------|---------------------------------|-----------------|---------|------------------------------------------------------------------------------|
| registry | consul_aws                      | >=0.5.0,<=1.0.0 |         | hashicorp/consul/aws                                                         |
| registry | consul                          | > 0.1.0         |         | hashicorp/consul/aws                                                         |
| git      | example_git_ssh_branch          |                 | master  | git::ssh://git@github.com/keilerkonzept/terraform-module-versions?ref=master |
| git      | example_git_scp                 | ~> 0.12         | 0.12.0  | git::git@github.com:keilerkonzept/terraform-module-versions?ref=0.12.0       |
| git      | consul_github_ssh               | ~0.1.0          | 0.1.0   | git@github.com:hashicorp/terraform-aws-consul?ref=0.1.0                      |
| git      | consul_github_https_no_ref      |                 |         | github.com/hashicorp/terraform-aws-consul                                    |
| git      | consul_github_https_missing_ref | 0.7.3           |         | github.com/hashicorp/terraform-aws-consul                                    |
| git      | consul_github_https             | 0.8.0           | v0.8.0  | github.com/hashicorp/terraform-aws-consul?ref=v0.8.0                         |

with `-o json`:

```json
[
  {
    "path": "examples/0.12.x.tf",
    "name": "consul_aws",
    "type": "registry",
    "source": "hashicorp/consul/aws",
    "constraint": ">=0.5.0,<=1.0.0"
  },
  {
    "path": "examples/main.tf",
    "name": "consul",
    "type": "registry",
    "source": "hashicorp/consul/aws",
    "constraint": "> 0.1.0"
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
  }
]
```

### Check for module updates

```sh
# check: check for module updates from (usually) remote sources
$ terraform-module-versions check examples
```

| UPDATE? |              NAME               | CONSTRAINT | VERSION | LATEST MATCHING | LATEST |
|---------|---------------------------------|------------|---------|-----------------|--------|
| (Y)     | consul_github_https_missing_ref | 0.7.3      |         | v0.7.3          | v0.8.0 |
| (Y)     | consul_github_https_no_ref      |            |         |                 | v0.8.0 |
| Y       | consul_github_ssh               | ~0.1.0     | 0.1.0   | v0.1.2          | v0.8.0 |
| (Y)     | example_git_scp                 | ~> 0.12    | 0.12.0  |                 | 2.1.1  |
| (Y)     | example_git_ssh_branch          |            | master  |                 | 2.1.1  |

with `-o json`:

```json
[
  {
    "path": "examples/main.tf",
    "name": "consul_github_https_missing_ref",
    "constraint": "0.7.3",
    "latestMatching": "v0.7.3",
    "latestOverall": "v0.8.0",
    "nonMatchingUpdate": true
  },
  {
    "path": "examples/main.tf",
    "name": "consul_github_https_no_ref",
    "latestOverall": "v0.8.0",
    "nonMatchingUpdate": true
  },
  {
    "path": "examples/main.tf",
    "name": "consul_github_ssh",
    "constraint": "~0.1.0",
    "version": "0.1.0",
    "latestMatching": "v0.1.2",
    "latestOverall": "v0.8.0",
    "matchingUpdate": true,
    "nonMatchingUpdate": true
  },
  {
    "path": "examples/main.tf",
    "name": "example_git_scp",
    "constraint": "~> 0.12",
    "version": "0.12.0",
    "latestOverall": "2.1.1",
    "nonMatchingUpdate": true
  },
  {
    "path": "examples/main.tf",
    "name": "example_git_ssh_branch",
    "version": "master",
    "latestOverall": "2.1.1",
    "nonMatchingUpdate": true
  }
]
```

```sh
# check -all: check for updates, include up-to-date-modules in output
$ terraform-module-versions check -all examples
```

| UPDATE? |              NAME               |   CONSTRAINT    | VERSION | LATEST MATCHING | LATEST |
|---------|---------------------------------|-----------------|---------|-----------------|--------|
| ?       | consul_aws                      | >=0.5.0,<=1.0.0 |         | 0.8.0           | 0.8.0  |
| ?       | consul                          | > 0.1.0         |         | 0.8.0           | 0.8.0  |
|         | consul_github_https             | 0.8.0           | v0.8.0  |                 | v0.8.0 |
| (Y)     | consul_github_https_missing_ref | 0.7.3           |         | v0.7.3          | v0.8.0 |
| (Y)     | consul_github_https_no_ref      |                 |         |                 | v0.8.0 |
| Y       | consul_github_ssh               | ~0.1.0          | 0.1.0   | v0.1.2          | v0.8.0 |
| (Y)     | example_git_scp                 | ~> 0.12         | 0.12.0  |                 | 2.1.1  |
| (Y)     | example_git_ssh_branch          |                 | master  |                 | 2.1.1  |

### Check for updates of specific modules

```sh
# check -module: check for updates of specific modules
$ terraform-module-versions check -all -module=consul_github_https -module=consul_github_ssh examples
```

| UPDATE? |        NAME         | CONSTRAINT | VERSION | LATEST MATCHING | LATEST |
|---------|---------------------|------------|---------|-----------------|--------|
|         | consul_github_https | 0.8.0      | v0.8.0  |                 | v0.8.0 |
| Y       | consul_github_ssh   | ~0.1.0     | 0.1.0   | v0.1.2          | v0.8.0 |

```sh
# check -module: check for updates of specific modules
$ terraform-module-versions check -module=consul_github_https -module=consul_github_ssh examples
```

| UPDATE? |       NAME        | CONSTRAINT | VERSION | LATEST MATCHING | LATEST |
|---------|-------------------|------------|---------|-----------------|--------|
| Y       | consul_github_ssh | ~0.1.0     | 0.1.0   | v0.1.2          | v0.8.0 |

with `-o json`:

```json
[
  {
    "path": "examples/main.tf",
    "name": "consul_github_ssh",
    "constraint": "~0.1.0",
    "version": "0.1.0",
    "latestMatching": "v0.1.2",
    "latestOverall": "v0.8.0",
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
  -output markdown  output format, one of [markdown-wide junit json jsonl markdown]
  -q false          (alias for -quiet)
  -quiet false      suppress log output (stderr)
```

### `list`

```text
USAGE
  terraform-module-versions list [options] [<path> ...]

List referenced terraform modules with their detected versions

FLAGS
  -module ...       include this module (may be specified repeatedly. by default, all modules are included)
  -o markdown       (alias for -output)
  -output markdown  output format, one of [markdown markdown-wide junit json jsonl]
```

### `check`

```text
USAGE
  terraform-module-versions check [options] [<path> ...]

Check referenced terraform modules' sources for newer versions

FLAGS
  -a false                           (alias for -all)
  -all false                         include modules without updates
  -e false                           (alias for -updates-found-nonzero-exit)
  -module ...                        include this module (may be specified repeatedly. by default, all modules are included)
  -o markdown                        (alias for -output)
  -output markdown                   output format, one of [markdown-wide junit json jsonl markdown]
  -updates-found-nonzero-exit false  exit with a nonzero code when modules with updates are found
```
