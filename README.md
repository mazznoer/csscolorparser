# CSS Color Parser

[![PkgGoDev](https://pkg.go.dev/badge/github.com/mazznoer/csscolorparser)](https://pkg.go.dev/github.com/mazznoer/csscolorparser)
[![Build Status](https://travis-ci.org/mazznoer/csscolorparser.svg?branch=master)](https://travis-ci.org/mazznoer/csscolorparser)
[![Build Status](https://github.com/mazznoer/csscolorparser/workflows/Go/badge.svg)](https://github.com/mazznoer/csscolorparser/actions)
[![go report](https://goreportcard.com/badge/github.com/mazznoer/csscolorparser)](https://goreportcard.com/report/github.com/mazznoer/csscolorparser)
[![codecov](https://codecov.io/gh/mazznoer/csscolorparser/branch/master/graph/badge.svg)](https://codecov.io/gh/mazznoer/csscolorparser)

Go (Golang) CSS color parser.

It support W3C's CSS color module level 4.

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

## Supported Format

It support named colors, hexadecimal (`#rgb`, `#rgba`, `#rrggbb`, `#rrggbbaa`), `rgb()`, `rgba()`, `hsl()`, `hsla()`, `hwb()`, and `hsv()`.

```
lime
#0f0
#0f0f
#00ff00
#00ff00ff
rgb(0,255,0)
rgb(0% 100% 0%)
rgb(0 255 0 / 100%)
rgba(0,255,0,1)
hsl(120,100%,50%)
hsl(120deg 100% 50%)
hsl(-240 100% 50%)
hsl(-240deg 100% 50%)
hsl(0.3333turn 100% 50%)
hsl(133.333grad 100% 50%)
hsl(2.0944rad 100% 50%)
hsla(120,100%,50%,100%)
hwb(120 0% 0%)
hwb(480deg 0% 0% / 100%)
hsv(120,100%,100%)
hsv(120deg 100% 100% / 100%)
```

## Try It Online

[Go Playground](https://play.golang.org/p/7kb62KSARwa)
