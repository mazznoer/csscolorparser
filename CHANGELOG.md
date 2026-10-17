# Changelog

## v0.2.0

### Added

- Support for:
    - `color(srgb ...)`
    - `color(srgb-linear ...)`
    - `color(xyz ...)`
    - `color(xyz-d50 ...)`
    - `color(xyz-d65 ...)`
- Support for relative color formats.
- New `ParseColors()` for parsing multiple colors separated by comma.

## v0.1.8

### Fixed

- Lab color conversion.

## v0.1.7

### Added

- Support parsing `lab()` and `lch()` color format.
- Make `NamedColors` public.
- `Color.RGBString()` returns modern CSS RGB format.

## v0.1.6

- Doc and code improvements.

## v0.1.5

### Added

- `Color.Clamp()`

## v0.1.4

### Added

- Support parsing `oklab()` and `oklch()` color format.
- `FromOklab()`, `FromOklch()`, `FromLinearRgb()`, `FromHsl()`, `FromHsv()`, `FromHwb()`

### Fixed

- Update `oklab` formula.

