# flagvar

[![](https://godoc.org/github.com/sgreben/flagvar?status.svg)](http://godoc.org/github.com/sgreben/flagvar) [![](https://goreportcard.com/badge/github.com/sgreben/flagvar/goreportcard)](https://goreportcard.com/report/github.com/sgreben/flagvar) [![coverage](http://gocover.io/_badge/github.com/sgreben/flagvar)](https://gocover.io/github.com/sgreben/flagvar) [![Build Status](https://travis-ci.org/sgreben/flagvar.svg?branch=master)](https://travis-ci.org/sgreben/flagvar)

A collection of CLI argument types for the `flag` package.

```go
import "github.com/sgreben/flagvar"
```

Each group of types (`Enum*`, `URL*`, `IP*`, ...) is also available as a separate package (also named `flagvar`), for example:

```go
import "github.com/sgreben/flagvar/enum"

var enums flagvar.EnumsCSV
```

Or just copy & paste what you need. It's public domain.

<!-- TOC -->

- [Example](#example)
- [Conventions](#conventions)
- [Types](#types)
- [Goals / design principles](#goals--design-principles)

<!-- /TOC -->

## Example

```go
package main

import (
	"flag"
	"fmt"
	"github.com/sgreben/flagvar"
)

var (
	fruit    = flagvar.Enum{Choices: []string{"apple", "banana"}}
	urls     flagvar.URLs
	settings flagvar.Assignments
)

func main() {
	flag.Var(&fruit, "fruit", fmt.Sprintf("set a fruit (%s)", fruit.Help()))
	flag.Var(&urls, "url", "add a URL")
	flag.Var(&settings, "set", fmt.Sprintf("specify a setting (%s)", settings.Help()))
	flag.Parse()
}
```

```sh
$ go run main.go -set abc=xyz -url https://github.com
# no error

$ go run main.go -set abc=xyz -url ://github.com
invalid value "://github.com" for flag -url: parse ://github.com: missing protocol scheme

$ go run main.go -fruit kiwi
invalid value "kiwi" for flag -fruit: "kiwi" must be one of [apple banana]

$ go run main.go -h
Usage:
  -fruit value
        set a fruit (one of [apple banana])
  -set value
        specify a setting (a key/value pair KEY=VALUE)
  -url value
        add a URL
```

## Conventions

- Pluralized argument types (e.g. `Strings`, `Assignments`) can be specified repeatedly, the values are collected in a slice.
- The resulting value is stored in `.Value` for singular types and in `.Values` for plural types
- The original argument string is stored in `.Text` for singular types and in `.Texts` for plural types
- -Set types (`EnumSet`, `StringSet`) de-duplicate provided values.
- -CSV types (`IntsCSV`, `EnumsCSV`) accept comma-separated values and accumulate values across flag instances if their `.Accumulate` field is set to `true`.
- Most types implement `interface{ Help() string }`, which produces a string suitable for inclusion in a help message.

## Types

Here's a compact overview:

| `flagvar` type | example CLI arg    | type of resulting Go value           |
|----------------|--------------------|--------------------------------------|
| [Alternative](https://godoc.org/github.com/sgreben/flagvar#Alternative)  |           |  |
| [Assignment](https://godoc.org/github.com/sgreben/flagvar#Assignment)  | KEY=VALUE          | struct{Key,Value} |
| [Assignments](https://godoc.org/github.com/sgreben/flagvar#Assignments) | KEY=VALUE          | []struct{Key,Value}                         |
| [AssignmentsMap](https://godoc.org/github.com/sgreben/flagvar#AssignmentsMap) | KEY=VALUE          | map[string]string                         |
| [CIDR](https://godoc.org/github.com/sgreben/flagvar#CIDR)        | 127.0.0.1/24               | struct{IPNet,IP}                              |
| [CIDRs](https://godoc.org/github.com/sgreben/flagvar#CIDRs)        | 127.0.0.1/24               | []struct{IPNet,IP}                              |
| [CIDRsCSV](https://godoc.org/github.com/sgreben/flagvar#CIDRsCSV)        | 127.0.0.1/16,10.1.2.3/8               | []struct{IPNet,IP}                              |
| [Enum](https://godoc.org/github.com/sgreben/flagvar#Enum)        | apple              | string                               |
| [Enums](https://godoc.org/github.com/sgreben/flagvar#Enums)       | apple              | []string                             |
| [EnumsCSV](https://godoc.org/github.com/sgreben/flagvar#EnumsCSV)       | apple,banana              | []string                             |
| [EnumSet](https://godoc.org/github.com/sgreben/flagvar#EnumSet)     | apple              | []string                             |
| [EnumSetCSV](https://godoc.org/github.com/sgreben/flagvar#EnumSetCSV)       | apple,banana              | []string                             |
| [File](https://godoc.org/github.com/sgreben/flagvar#File)        | ./README.md        | string                               |
| [Files](https://godoc.org/github.com/sgreben/flagvar#Files)        | ./README.md        | []string                               |
| [Floats](https://godoc.org/github.com/sgreben/flagvar#Floats)      | 1.234              | []float64                            |
| [FloatsCSV](https://godoc.org/github.com/sgreben/flagvar#FloatsCSV)      | 1.234,5.0              | []float64                            |
| [Glob](https://godoc.org/github.com/sgreben/flagvar#Glob)        | src/**.js          | glob.Glob                            |
| [Globs](https://godoc.org/github.com/sgreben/flagvar#Globs)       | src/**.js          | []glob.Glob                            |
| [Ints](https://godoc.org/github.com/sgreben/flagvar#Ints)        | 1002               | []int64                              |
| [IntsCSV](https://godoc.org/github.com/sgreben/flagvar#IntsCSV)        | 123,1002               | []int64                              |
| [IP](https://godoc.org/github.com/sgreben/flagvar#IP)        | 127.0.0.1               | net.IP                              |
| [IPs](https://godoc.org/github.com/sgreben/flagvar#IPs)        | 127.0.0.1               | []net.IP                              |
| [IPsCSV](https://godoc.org/github.com/sgreben/flagvar#IPsCSV)        | 127.0.0.1,10.1.2.3               | []net.IP                              |
| [JSON](https://godoc.org/github.com/sgreben/flagvar#JSON)        | '{"a":1}'          | interface{}                          |
| [JSONs](https://godoc.org/github.com/sgreben/flagvar#JSONs)       | '{"a":1}'          | []interface{}                        |
| [Regexp](https://godoc.org/github.com/sgreben/flagvar#Regexp)        | [a-z]+          | *regexp.Regexp                            |
| [Regexps](https://godoc.org/github.com/sgreben/flagvar#Regexps)       | [a-z]+          | []*regexp.Regexp                            |
| [Strings](https://godoc.org/github.com/sgreben/flagvar#Strings)     | "xyxy"             | []string                             |
| [StringSet](https://godoc.org/github.com/sgreben/flagvar#StringSet)  | "xyxy"             | []string                             |
| [StringSetCSV](https://godoc.org/github.com/sgreben/flagvar#StringSetCSV)       | y,x,y              | []string                             |
| [TCPAddr](https://godoc.org/github.com/sgreben/flagvar#TCPAddr)        | 127.0.0.1:10               | net.TCPAddr                              |
| [TCPAddrs](https://godoc.org/github.com/sgreben/flagvar#TCPAddrs)        | 127.0.0.1:10               | []net.TCPAddr                              |
| [TCPAddrsCSV](https://godoc.org/github.com/sgreben/flagvar#TCPAddrsCSV)        | 127.0.0.1:10,:123               | []net.TCPAddr                              |
| [Template](https://godoc.org/github.com/sgreben/flagvar#Template)    | "{{.Size}}"        | *template.Template                   |
| [Templates](https://godoc.org/github.com/sgreben/flagvar#Templates)   | "{{.Size}}"        | []*template.Template                 |
| [TemplateFile](https://godoc.org/github.com/sgreben/flagvar#TemplateFile)   | "/path/to/template.file"        | string                 |
| [Time](https://godoc.org/github.com/sgreben/flagvar#Time)        | "10:30 AM"         | time.Time                            |
| [Times](https://godoc.org/github.com/sgreben/flagvar#Times)       | "10:30 AM"         | []time.Time                          |
| [UDPAddr](https://godoc.org/github.com/sgreben/flagvar#UDPAddr)        | 127.0.0.1:10               | net.UDPAddr                              |
| [UDPAddrs](https://godoc.org/github.com/sgreben/flagvar#UDPAddrs)        | 127.0.0.1:10               | []net.UDPAddr                              |
| [UDPAddrsCSV](https://godoc.org/github.com/sgreben/flagvar#UDPAddrsCSV)        | 127.0.0.1:10,:123               | []net.UDPAddr                              |
| [UnixAddr](https://godoc.org/github.com/sgreben/flagvar#UnixAddr)        | /example.sock               | net.UnixAddr                              |
| [UnixAddrs](https://godoc.org/github.com/sgreben/flagvar#UnixAddrs)        | /example.sock               | []net.UnixAddr                              |
| [UnixAddrsCSV](https://godoc.org/github.com/sgreben/flagvar#UnixAddrsCSV)        | /example.sock,/other.sock               | []net.UnixAddr                              |
| [URL](https://godoc.org/github.com/sgreben/flagvar#URL)         | https://github.com | *url.URL                             |
| [URLs](https://godoc.org/github.com/sgreben/flagvar#URLs)        | https://github.com | []*url.URL                           |
| [Wrap](https://godoc.org/github.com/sgreben/flagvar#Wrap)        |                    |                                      |
| [WrapCSV](https://godoc.org/github.com/sgreben/flagvar#WrapCSV)        |                    |                                      |
| [WrapFunc](https://godoc.org/github.com/sgreben/flagvar#WrapFunc)    |                    |                                      |
| [WrapPointer](https://godoc.org/github.com/sgreben/flagvar#WrapPointer)    |                    |                                      |

## Goals / design principles

- Help avoid dependencies
    - Self-contained > DRY
    - Explicitly support copy & paste workflow
    - Copyable units should be easy to determine
    - Anonymous structs > shared types
- "Code-you-own" feeling, even when imported as a package
    - No private fields / methods
    - No magic
    - Simple built-in types used wherever possible
    - Avoid introducing new concepts
- Support "blind" usage
    - Zero values should be useful
    - Avoid introducing failure cases, handle any combination of parameters gracefully.
    - All "obvious things to try" should work.
