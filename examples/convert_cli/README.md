# convert_cli

`convert_cli` is a small image conversion example built with `github.com/upamune/purevipsgen/vips`.

It loads an input image through libvips, optionally autorotates and resizes it, then writes one of the supported output formats. The implementation lives under `internal/convertcli` so it is scoped to this example and cannot be imported from outside this directory tree.

## Usage

```bash
CGO_ENABLED=0 go run ./examples/convert_cli \
  -in /path/to/input.heic \
  -out /path/to/output.webp \
  -width 1600 \
  -quality 85
```

Required flags:

- `-in`: input image path
- `-out`: output image path

Common optional flags:

- `-format`: output format. If omitted, the output extension is used.
- `-width`: maximum output width while preserving aspect ratio.
- `-height`: maximum output height while preserving aspect ratio.
- `-quality`: quality for JPEG, WebP, HEIC, AVIF, and TIFF. Default is `85`.
- `-lossless`: enable lossless mode where supported.
- `-autorotate`: apply EXIF orientation after loading. Default is `true`.

Supported output formats:

- JPEG: `.jpg` or `.jpeg`
- PNG: `.png`
- WebP: `.webp`
- HEIC/HEIF: `.heic` or `.heif`
- AVIF: `.avif`
- TIFF: `.tiff` or `.tif`
- GIF: `.gif`

## Examples

Convert HEIC to WebP:

```bash
CGO_ENABLED=0 go run ./examples/convert_cli \
  -in /Users/upamune/Downloads/sewing-threads.heic \
  -out /Users/upamune/Downloads/sewing-threads.webp \
  -width 1600 \
  -quality 85
```

Convert HEIC to JPEG:

```bash
CGO_ENABLED=0 go run ./examples/convert_cli \
  -in /Users/upamune/Downloads/sewing-threads.heic \
  -out /Users/upamune/Downloads/sewing-threads.jpg \
  -width 1600 \
  -quality 90
```

Convert HEIC to PNG:

```bash
CGO_ENABLED=0 go run ./examples/convert_cli \
  -in /Users/upamune/Downloads/sewing-threads.heic \
  -out /Users/upamune/Downloads/sewing-threads.png \
  -width 1600
```
