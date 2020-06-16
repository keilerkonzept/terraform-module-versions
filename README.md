# terraform-module-versions

![Build](https://github.com/keilerkonzept/terraform-module-versions/workflows/Build/badge.svg)

Checks for updates of external terraform modules referenced in given Terraform (0.10.x - 0.12.x) modules. Outputs JSONL (one JSON object per line), or Markdown tables (`-pretty, -p`).

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
$ terraform-module-versions -updates -pretty examples
```

```markdown
| UPDATE? |        NAME         |   CONSTRAINT    | VERSION | LATEST MATCHING | LATEST |
|---------|---------------------|-----------------|---------|-----------------|--------|
| ?       | consul_aws          | >=0.5.0,<=1.0.0 |         | 0.7.4           | 0.7.4  |
| ?       | consul              | > 0.1.0         |         | 0.7.4           | 0.7.4  |
| (Y)     | consul_github_https | 0.7.3           |         | 0.7.3           | 0.7.4  |
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
$ cat examples/main.tf examples/0.12.x.tf
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
  source = "git::ssh://git@github.com/keilerkonzept/terraform-module-versions?ref=0.10.0"
  version = "~> 0.10"
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
# default operation: list modules with their current versions and version constraints (if specified)
$ terraform-module-versions examples
```

```json
{
  "path": "examples/0.12.x.tf",
  "name": "consul_aws",
  "type": "registry",
  "source": "hashicorp/consul/aws",
  "constraint": ">=0.5.0,<=1.0.0"
}
{
  "path": "examples/main.tf",
  "name": "consul",
  "type": "registry",
  "source": "hashicorp/consul/aws",
  "constraint": "> 0.1.0"
}
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
  "source": "git::ssh://git@github.com/keilerkonzept/terraform-module-versions?ref=0.10.0",
  "constraint": "~> 0.10",
  "version": "0.10.0"
}
```

with `-pretty`:

|   TYPE   |        NAME         |   CONSTRAINT    | VERSION |                                    SOURCE                                    |
|----------|---------------------|-----------------|---------|------------------------------------------------------------------------------|
| registry | consul_aws          | >=0.5.0,<=1.0.0 |         | hashicorp/consul/aws                                                         |
| registry | consul              | > 0.1.0         |         | hashicorp/consul/aws                                                         |
| git      | example_git_ssh     | ~> 0.10         | 0.10.0  | git::ssh://git@github.com/keilerkonzept/terraform-module-versions?ref=0.10.0 |
| git      | consul_github_ssh   | 0.1.0           | 0.1.0   | git@github.com:hashicorp/terraform-aws-consul?ref=0.1.0                      |
| git      | consul_github_https | 0.7.3           |         | github.com/hashicorp/terraform-aws-consul                                    |

### Check for module updates

```sh
# -update: check for module updates from (usually) remote sources
$ terraform-module-versions -updates examples
```

```json
{
  "path": "examples/main.tf",
  "name": "consul",
  "constraint": "> 0.1.0",
  "constraintUpdate": true,
  "latestMatching": "0.7.4",
  "latestOverall": "0.7.4"
}
{
  "path": "examples/0.12.x.tf",
  "name": "consul_aws",
  "constraint": ">=0.5.0,<=1.0.0",
  "constraintUpdate": true,
  "latestMatching": "0.7.4",
  "latestOverall": "0.7.4"
}
{
  "path": "examples/main.tf",
  "name": "consul_github_https",
  "constraint": "0.7.3",
  "latestMatching": "0.7.3",
  "latestOverall": "0.7.4",
  "nonMatchingUpdate": true
}
```

with `-pretty`:

| UPDATE? |        NAME         |   CONSTRAINT    | VERSION | LATEST MATCHING | LATEST |
|---------|---------------------|-----------------|---------|-----------------|--------|
| ?       | consul_aws          | >=0.5.0,<=1.0.0 |         | 0.7.4           | 0.7.4  |
| ?       | consul              | > 0.1.0         |         | 0.7.4           | 0.7.4  |
| (Y)     | consul_github_https | 0.7.3           |         | 0.7.3           | 0.7.4  |

### Check for updates of specific modules

```sh
# -update and -module: check for updates of specific modules
$ terraform-module-versions -updates -module=consul_github_https -module=consul_github_ssh examples
```

```json
{
  "path": "examples/main.tf",
  "name": "consul_github_https",
  "constraint": "0.7.3",
  "latestMatching": "0.7.3",
  "latestOverall": "0.7.4",
  "nonMatchingUpdate": true
}
```

with `-pretty`:

| UPDATE? |        NAME         | CONSTRAINT | VERSION | LATEST MATCHING | LATEST |
|---------|---------------------|------------|---------|-----------------|--------|
| (Y)     | consul_github_https | 0.7.3      |         | 0.7.3           | 0.7.4  |

## Get it

Using go get:

```bash
go get -u github.com/keilerkonzept/terraform-module-versions
```

Or [download the binary for your platform](https://github.com/keilerkonzept/terraform-module-versions/releases/latest) from the releases page.

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
