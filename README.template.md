# ${APP}

Checks for updates of external terraform modules referenced in given Terraform (0.12.x) modules. Outputs JSONL (one JSON object per line), or Markdown tables (`-pretty, -p`).

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
$ ${APP} -updates -pretty examples
```

```markdown
${EXAMPLE_PRETTY}
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
${EXAMPLES_MAIN_TF}
```

### List modules with their current versions

```sh
# default operation: list modules with their current versions and version constraints (if specified)
$ ${APP} examples
```

```json
${EXAMPLE_LIST}
```

with `-pretty`:

${EXAMPLE_LIST_PRETTY}

### Check for module updates

```sh
# -update: check for module updates from (usually) remote sources
$ ${APP} -updates examples
```

```json
${EXAMPLE_UPDATES}
```

with `-pretty`:

${EXAMPLE_UPDATES_PRETTY}

### Check for updates of specific modules

```sh
# -update and -module: check for updates of specific modules
$ ${APP} -updates -module=consul_github_https -module=consul_github_ssh examples
```

```json
${EXAMPLE_UPDATES_SINGLE}
```

with `-pretty`:

${EXAMPLE_UPDATES_SINGLE_PRETTY}

## Get it

Using go get:

```bash
go get -u github.com/keilerkonzept/${APP}
```

Or [download the binary for your platform](https://github.com/keilerkonzept/${APP}/releases/latest) from the releases page.

## Usage

```text
${APP} [PATHS...]

$USAGE
```
