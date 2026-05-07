# purevipsgen

[![Go Reference](https://pkg.go.dev/badge/github.com/upamune/purevipsgen/vips.svg)](https://pkg.go.dev/github.com/upamune/purevipsgen/vips)
[![CI](https://github.com/upamune/purevipsgen/actions/workflows/ci.yml/badge.svg)](https://github.com/upamune/purevipsgen/actions/workflows/ci.yml)

purevipsgen is a Go binding generator for [libvips](https://github.com/libvips/libvips) - a fast and efficient image processing library.

purevipsgen is a fork of [github.com/cshum/vipsgen](https://github.com/cshum/vipsgen). The original project established the generator architecture, generated API shape, examples, and documentation foundation this repository builds on. We are grateful for that work.

libvips is generally 4-8x [faster](https://github.com/libvips/libvips/wiki/Speed-and-memory-use) than ImageMagick with low memory usage, thanks to its [demand-driven, horizontally threaded](https://github.com/libvips/libvips/wiki/Why-is-libvips-quick) architecture.

Existing Go libvips bindings rely on manually written code that is often incomplete, error-prone, and difficult to maintain as libvips evolves.
purevipsgen solves this by generating type-safe, documented Go bindings from GObject introspection and calling libvips dynamically through [`purego`](https://github.com/ebitengine/purego).

- **Comprehensive**: Bindings for around [300 libvips operations](https://www.libvips.org/API/current/function-list.html)
- **Type-Safe**: Proper Go types for all libvips C enums and structs
- **Idiomatic**: Clean Go APIs that feel natural to use
- **No cgo for consumers**: generated packages use purego to load libvips at runtime
- **Streaming**: `VipsSource` and `VipsTarget` integration with Go `io.Reader` and `io.Writer`

You can use purevipsgen in two ways:

- **Import directly**: Use the pre-generated library `github.com/upamune/purevipsgen/vips` for the latest default installation of libvips, or see [pre-generated packages](#pre-generated-packages)
- **Generate custom bindings**: Run the purevipsgen command to create bindings for your specific libvips version and installation


## Quick Start

Use homebrew to install vips and pkg-config:
```
brew install vips pkg-config
```

Code generation uses GObject introspection and may need cgo flags on macOS:

```bash
export CGO_CFLAGS_ALLOW="-Xpreprocessor"
```

Use the package directly:

```bash
go get -u github.com/upamune/purevipsgen/vips
```

Operations support parameters and optional arguments through structs, maintaining direct equivalence with the [libvips API](https://www.libvips.org/API/current/).
Pass `nil` to use default behavior for optional arguments.
See [examples](https://github.com/upamune/purevipsgen/tree/main/examples) for common usage patterns.


```go
package main

import (
	"log"
	"net/http"

	"github.com/upamune/purevipsgen/vips"
)

func main() {
	// Fetch an image from http.Get
	resp, err := http.Get("https://raw.githubusercontent.com/cshum/imagor/master/testdata/gopher.png")
	if err != nil {
		log.Fatalf("Failed to fetch image: %v", err)
	}
	defer resp.Body.Close()

	// Create source from io.ReadCloser
	source := vips.NewSource(resp.Body)
	defer source.Close() // source needs to remain available during image lifetime

	// Shrink-on-load via creating image from thumbnail source with options
	image, err := vips.NewThumbnailSource(source, 800, &vips.ThumbnailSourceOptions{
		Height: 1000,
		FailOn: vips.FailOnError, // Fail on first error
	})
	if err != nil {
		log.Fatalf("Failed to load image: %v", err)
	}
	defer image.Close() // always close images to free memory

	// Add a yellow border using vips_embed
	border := 10
	if err := image.Embed(
		border, border,
		image.Width()+border*2,
		image.Height()+border*2,
		&vips.EmbedOptions{
			Extend:     vips.ExtendBackground,       // extend with colour from the background property
			Background: []float64{255, 255, 0, 255}, // Yellow border
		},
	); err != nil {
		log.Fatalf("Failed to add border: %v", err)
	}

	log.Printf("Processed image: %dx%d\n", image.Width(), image.Height())

	// Save the result as WebP file with options
	err = image.Webpsave("resized-gopher.webp", &vips.WebpsaveOptions{
		Q:              85,   // Quality factor (0-100)
		Effort:         4,    // Compression effort (0-6)
		SmartSubsample: true, // Better chroma subsampling
	})
	if err != nil {
		log.Fatalf("Failed to save image as WebP: %v", err)
	}
	log.Println("Successfully saved processed images")
}
```

## Pre-generated Packages

purevipsgen provides pre-generated bindings checked in under the paths below. All packages use the same `vips` package name and API - only the import path differs.

| Import Path | libvips Version | Use When |
|-------------|----------------|----------|
| `github.com/upamune/purevipsgen/vips` | 8.17.0 | Default generated package for this checkout |
| `github.com/upamune/purevipsgen/vips817` | 8.17.0 | Versioned import path for libvips 8.17.x |
| `github.com/upamune/purevipsgen/vips816` | 8.17.0 | Compatibility import path; regenerate with libvips 8.16.x before publishing 8.16 bindings |

**Important:** Only import ONE of these packages in your project. Choose based on your installed libvips version.

Check your libvips version with `vips --version`, then use the corresponding import:

```go
// For libvips 8.18.x (latest - recommended)
import "github.com/upamune/purevipsgen/vips"

// For libvips 8.17.x
import "github.com/upamune/purevipsgen/vips817"

// For libvips 8.16.x
import "github.com/upamune/purevipsgen/vips816"

func main() {
    // API is identical across all versions
    img, err := vips.NewImageFromFile("input.jpg", nil)
    if err != nil {
        log.Fatal(err)
    }
    defer img.Close()
    
    err = img.Resize(0.5, nil)
    // ...
}
```

## Image Loaders

**Generic loaders** — [`NewImageFromFile`](https://pkg.go.dev/github.com/upamune/purevipsgen/vips#NewImageFromFile), [`NewImageFromBuffer`](https://pkg.go.dev/github.com/upamune/purevipsgen/vips#NewImageFromBuffer), [`NewImageFromSource`](https://pkg.go.dev/github.com/upamune/purevipsgen/vips#NewImageFromSource) — automatically detect the image format and accept `LoadOptions`, a generic options struct covering common options across formats. Since not every format supports every option, use the **format-specific loaders** — [`NewGifload`](https://pkg.go.dev/github.com/upamune/purevipsgen/vips#NewGifload), [`NewJpegloadBuffer`](https://pkg.go.dev/github.com/upamune/purevipsgen/vips#NewJpegloadBuffer), [`NewPngloadSource`](https://pkg.go.dev/github.com/upamune/purevipsgen/vips#NewPngloadSource), etc. — for precise, type-safe control. A few common examples:

**Animated GIF** — `N: -1` loads all frames ([full example](https://github.com/upamune/purevipsgen/tree/main/examples/from_file)):

```go
image, err := vips.NewGifload("animation.gif", &vips.GifloadOptions{
    N: -1, // -1 = load all frames
})
```

**JPEG auto-rotation** — rotate by EXIF orientation on load:

```go
source := vips.NewSource(reader)
defer source.Close()

image, err := vips.NewJpegloadSource(source, &vips.JpegloadSourceOptions{
    Autorotate: true,
})
```

## Working with Animated Images

libvips represents multi-frame images (animated GIF, WebP) as a single vertically stacked image where each frame occupies one page of height `PageHeight`. purevipsgen exposes the page metadata and provides dedicated helpers for operations that must process each frame individually.

### Metadata

```go
img.Pages()                 // number of frames
img.PageHeight()            // height of a single frame in pixels
delays, _ := img.PageDelay() // per-frame delay in milliseconds
img.Loop()                  // loop count (0 = infinite)
```

### Multi-page Helpers

Some operations — rotate, crop, embed — cannot be expressed as a single libvips pipeline call across a stacked image. purevipsgen provides helpers that apply the operation across frames while preserving the public Go API:

```go
// Rotate all frames
err = img.RotMultiPage(vips.AngleD90)

// Crop all frames to the same region
err = img.ExtractAreaMultiPage(left, top, width, height)

// Embed (pad/extend) all frames to a new canvas
err = img.EmbedMultiPage(left, top, newWidth, newHeight, &vips.EmbedMultiPageOptions{
    Extend:     vips.ExtendBackground,
    Background: []float64{0, 0, 0, 0},
})
```

These methods automatically fall through to the equivalent single-frame operation when the image has only one page.

## Code Generation

Code generation requires libvips to be built with GObject introspection support.

```bash
go install github.com/upamune/purevipsgen/cmd/purevipsgen@latest
```

Generate the bindings:

```bash
purevipsgen -out ./vips
```

Use your custom-generated code:

```go
package main

import (
    "yourproject/vips"
)
```

### Command Line Options

```
Usage of purevipsgen:
  -debug
        Enable debug json output
  -extract
        Extract embedded templates to a directory
  -extract-dir string
        Directory to extract templates to (default "./templates")
  -include-test
        Include test files in generated output
  -out string
        Output directory (default "./vips")
  -templates string
        Template directory (uses embedded templates if not specified)
```

### How Code Generation Works

The generation process has three main layers:

1. **Introspection analysis**: purevipsgen uses GObject introspection to analyze the installed libvips API, extracting operation metadata, argument types, enum definitions, and defaults.

2. **Purego binding layer**: generated `vips.go` functions register libvips symbols with purego and call them directly, converting Go scalars, strings, arrays, and opaque libvips pointers as needed.

3. **Go method layer**: generated `image.go` methods expose an idiomatic API on `*Image`, including option structs for libvips optional arguments.

For example, generated image methods select the required-argument call when options are nil and the option-aware call when an options struct is provided:

```go
// ResizeOptions optional arguments for vips_resize
type ResizeOptions struct {
    // Kernel Resampling kernel
    Kernel Kernel
    // Gap Reducing gap
    Gap float64
    // Vscale Vertical scale image by this factor
    Vscale float64
}

// DefaultResizeOptions creates default value for vips_resize optional arguments
func DefaultResizeOptions() *ResizeOptions {
    return &ResizeOptions{
        Kernel: Kernel(5),
        Gap: 2,
    }
}

// Resize vips_resize resize an image
func (r *Image) Resize(scale float64, options *ResizeOptions) error {
    if options != nil {
        out, err := purevipsgenResizeWithOptions(r.image, scale,
                                           options.Kernel, options.Gap, options.Vscale)
        if err != nil {
            return err
        }
        r.setImage(out)
        return nil
    }
    out, err := purevipsgenResize(r.image, scale)
    if err != nil {
        return err
    }
    r.setImage(out)
    return nil
}
```

This layer provides idiomatic Go methods, options structs for optional parameters, Go type system integration.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

### Development Workflow

Do not commit generated code in `vips*/` directories. Generated bindings are created automatically by CI.

**For contributors:** Submit PRs with source code changes only. Maintainers will regenerate bindings after merge.

**For maintainers:** After merging a fork PR, manually run the CI workflow from the Actions tab to regenerate bindings.

## License

MIT
