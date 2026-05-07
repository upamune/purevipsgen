# Repository Guidelines

## Project Structure & Module Organization

This Go module generates and publishes libvips bindings. The CLI entrypoint is `cmd/purevipsgen/main.go`. Generator code lives in `internal/generator`, introspection helpers in `internal/introspection`, and embedded templates in `internal/templates`. Support utilities are in `pointer`. Generated binding packages are checked in under `vips`, `vips816`, and `vips817`. Usage examples are in `examples/*`.

## Build, Test, and Development Commands

Use `nix develop` or run `direnv allow` to enter the pinned shell. It provides Go, libvips, pkg-config, gobject-introspection, make, and golangci-lint. Without Nix, install `libvips`, `pkg-config`, and Go 1.24 or newer.

- `make check`: verifies `pkg-config` can find libvips.
- `make build`: builds `bin/purevipsgen` with version metadata.
- `make test`: clears the Go test cache and runs package tests serially with coverage.
- `make generate`: builds the CLI and regenerates the default `vips` package.
- `make generate-custom`: regenerates using templates from `internal/templates`.
- `make lint`: runs `golangci-lint run ./...` when the linter is installed.

For direct iteration, use `CGO_CFLAGS_ALLOW=-Xpreprocessor go test ./...` when exercising the generator. Generated binding packages should also pass with `CGO_ENABLED=0 go test ./...`.

## Coding Style & Naming Conventions

Use standard Go formatting: run `gofmt` on changed `.go` files and keep imports organized with `go fmt` / `go test`. Package names are short and lowercase (`generator`, `introspection`, `pointer`). Generated binding packages use the `vips` package name even when the import path is versioned. Keep C helper files and headers beside the generated Go package they support, matching existing names such as `connection.c`, `util.h`, and `vips.go`.

## Testing Guidelines

Tests use Go’s standard `testing` package plus `stretchr/testify`. Put tests beside the package under test with `_test.go` suffixes, as in `pointer/pointer_test.go` and `vips/image_test.go`. Because generated bindings dynamically load libvips through purego, confirm the installed libvips version matches the package under test when running version-specific directories.

## Commit & Pull Request Guidelines

Recent commits use concise Conventional Commit-style prefixes such as `fix:`, `docs:`, `build:`, and `chore:`. Keep generated binding updates separate when practical; CI may auto-commit `vips*` outputs for supported libvips versions. PRs should describe the user-visible change, note any regenerated packages, and include the exact test command run. Link issues when applicable and add example code or screenshots only for documentation/example changes that benefit from them.

## Agent-Specific Instructions

Do not hand-edit generated bindings unless the task is explicitly about generated output. Prefer changing `internal/templates` or `internal/generator`, then regenerate and inspect the resulting `vips*` diff.
