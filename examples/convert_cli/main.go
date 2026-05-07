package main

import (
	"os"

	"github.com/upamune/purevipsgen/examples/convert_cli/internal/convertcli"
)

func main() {
	os.Exit(convertcli.Main(os.Args[1:]))
}
