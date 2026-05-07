package convertcli

import (
	"flag"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/upamune/purevipsgen/vips"
)

type options struct {
	input      string
	output     string
	format     string
	width      int
	height     int
	quality    int
	lossless   bool
	autorotate bool
}

type result struct {
	Input        string
	Output       string
	Format       string
	OutputWidth  int
	OutputHeight int
}

func Main(args []string) int {
	if err := run(args, os.Stdout, os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	return 0
}

func run(args []string, stdout, stderr io.Writer) error {
	opts, err := parseFlags(args, stderr)
	if err != nil {
		return err
	}

	res, err := convert(opts)
	if err != nil {
		return err
	}

	fmt.Fprintf(stdout, "converted %s -> %s (%s, %dx%d)\n", res.Input, res.Output, res.Format, res.OutputWidth, res.OutputHeight)
	return nil
}

func parseFlags(args []string, output io.Writer) (options, error) {
	var opts options
	flags := flag.NewFlagSet("convert_cli", flag.ContinueOnError)
	flags.SetOutput(output)
	flags.StringVar(&opts.input, "in", "", "input image path")
	flags.StringVar(&opts.output, "out", "", "output image path")
	flags.StringVar(&opts.format, "format", "", "output format: jpg, png, webp, heic, avif, tiff, gif")
	flags.IntVar(&opts.width, "width", 0, "maximum output width; keeps aspect ratio")
	flags.IntVar(&opts.height, "height", 0, "maximum output height; keeps aspect ratio")
	flags.IntVar(&opts.quality, "quality", 85, "quality for jpeg, webp, heic, avif, and tiff")
	flags.BoolVar(&opts.lossless, "lossless", false, "use lossless mode for webp/heic/avif when supported")
	flags.BoolVar(&opts.autorotate, "autorotate", true, "apply EXIF orientation after loading")

	if err := flags.Parse(args); err != nil {
		return options{}, err
	}
	if opts.input == "" || opts.output == "" {
		flags.Usage()
		return options{}, fmt.Errorf("missing required -in or -out")
	}
	if opts.quality < 1 || opts.quality > 100 {
		return options{}, fmt.Errorf("-quality must be between 1 and 100")
	}
	if opts.width < 0 || opts.height < 0 {
		return options{}, fmt.Errorf("-width and -height must be positive")
	}

	opts.format = normalizeFormat(opts.format, opts.output)
	return opts, nil
}

func convert(opts options) (result, error) {
	vips.Startup(&vips.Config{})
	defer vips.Shutdown()

	img, err := vips.NewImageFromFile(opts.input, &vips.LoadOptions{
		FailOnError: true,
		Access:      vips.AccessSequential,
	})
	if err != nil {
		return result{}, fmt.Errorf("load %q: %w", opts.input, err)
	}
	defer img.Close()

	if opts.autorotate {
		if err := img.Autorot(nil); err != nil {
			return result{}, fmt.Errorf("autorotate: %w", err)
		}
	}

	if opts.width > 0 || opts.height > 0 {
		if err := resizeToFit(img, opts.width, opts.height); err != nil {
			return result{}, fmt.Errorf("resize: %w", err)
		}
	}

	if err := saveImage(img, opts.output, opts.format, opts.quality, opts.lossless); err != nil {
		return result{}, fmt.Errorf("save %q as %s: %w", opts.output, opts.format, err)
	}

	return result{
		Input:        opts.input,
		Output:       opts.output,
		Format:       opts.format,
		OutputWidth:  img.Width(),
		OutputHeight: img.Height(),
	}, nil
}

func resizeToFit(img *vips.Image, width, height int) error {
	if width == 0 {
		width = int(math.Round(float64(img.Width()) * float64(height) / float64(img.Height())))
	}
	if height == 0 {
		height = 1
	}

	return img.ThumbnailImage(width, &vips.ThumbnailImageOptions{
		Height: height,
		Size:   vips.SizeDown,
	})
}

func normalizeFormat(format, output string) string {
	if format == "" {
		format = strings.TrimPrefix(strings.ToLower(filepath.Ext(output)), ".")
	}
	switch strings.ToLower(format) {
	case "jpeg":
		return "jpg"
	case "tif":
		return "tiff"
	default:
		return strings.ToLower(format)
	}
}

func saveImage(img *vips.Image, output, format string, quality int, lossless bool) error {
	switch format {
	case "jpg":
		options := vips.DefaultJpegsaveOptions()
		options.Q = quality
		return img.Jpegsave(output, options)
	case "png":
		return img.Pngsave(output, nil)
	case "webp":
		options := vips.DefaultWebpsaveOptions()
		options.Q = quality
		options.Lossless = lossless
		return img.Webpsave(output, options)
	case "heic", "heif":
		options := vips.DefaultHeifsaveOptions()
		options.Q = quality
		options.Lossless = lossless
		options.Compression = vips.HeifCompressionHevc
		return img.Heifsave(output, options)
	case "avif":
		options := vips.DefaultHeifsaveOptions()
		options.Q = quality
		options.Lossless = lossless
		options.Compression = vips.HeifCompressionAv1
		return img.Heifsave(output, options)
	case "tiff":
		options := vips.DefaultTiffsaveOptions()
		options.Q = quality
		return img.Tiffsave(output, options)
	case "gif":
		return img.Gifsave(output, vips.DefaultGifsaveOptions())
	default:
		return fmt.Errorf("unsupported format %q", format)
	}
}
