package main

import (
	"fmt"
	"os"

	"github.com/gabefiori/tms/internal/cli"
)

func main() {
	if err := cli.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
