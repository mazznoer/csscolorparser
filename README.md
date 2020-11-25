# CSS Color Parser

[![PkgGoDev](https://pkg.go.dev/badge/github.com/mazznoer/csscolorparser)](https://pkg.go.dev/github.com/mazznoer/csscolorparser)
[![Build Status](https://travis-ci.org/mazznoer/csscolorparser.svg?branch=master)](https://travis-ci.org/mazznoer/csscolorparser)
[![Build Status](https://github.com/mazznoer/csscolorparser/workflows/Go/badge.svg)](https://github.com/mazznoer/csscolorparser/actions)
[![go report](https://goreportcard.com/badge/github.com/mazznoer/csscolorparser)](https://goreportcard.com/report/github.com/mazznoer/csscolorparser)
[![codecov](https://codecov.io/gh/mazznoer/csscolorparser/branch/master/graph/badge.svg)](https://codecov.io/gh/mazznoer/csscolorparser)

Go (Golang) CSS color parser.

It support W3C's CSS color module level 4

```go
import "github.com/mazznoer/csscolorparser"
```

## Usage

```go
c, err := csscolorparser.Parse("gold")

if err != nil {
	panic(err)
}
```
<!--
TODO
- all supported format
- link playground
-->
