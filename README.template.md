# ${APP}

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
$ ${APP} check examples
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
  - [`list`](#list)
  - [`check`](#check)

## Examples

```sh
$ cat examples/main.tf
```

```terraform
${EXAMPLES_MAIN_TF}
```

### List modules with their current versions

```sh
# list modules with their current versions and version constraints (if specified)
$ ${APP} list examples
```

${EXAMPLE_LIST_PRETTY}

with `-o json`:

```json
${EXAMPLE_LIST}
```

### Check for module updates

```sh
# check: check for module updates from (usually) remote sources
$ ${APP} check examples
```

${EXAMPLE_UPDATES_PRETTY}

with `-o json`:

```json
${EXAMPLE_UPDATES}
```

```sh
# check -all: check for updates, include up-to-date-modules in output
$ ${APP} check -all examples
```

${EXAMPLE_UPDATES_ALL_PRETTY}

### Check for updates of specific modules

```sh
# check -module: check for updates of specific modules
$ ${APP} check -all -module=consul_github_https -module=consul_github_ssh examples
```

${EXAMPLE_UPDATES_SINGLE_ALL_PRETTY}

```sh
# check -module: check for updates of specific modules
$ ${APP} check -module=consul_github_https -module=consul_github_ssh examples
```

${EXAMPLE_UPDATES_SINGLE_PRETTY}

with `-o json`:

```json
${EXAMPLE_UPDATES_SINGLE}
```

## Get it

Using go get:

```bash
go get -u github.com/keilerkonzept/${APP}
```

Or [download the binary for your platform](https://github.com/keilerkonzept/${APP}/releases/latest) from the releases page.

## Usage

```text
$USAGE
```

### `list`

```text
${USAGE_LIST}
```

### `check`

```text
${USAGE_CHECK}
```
