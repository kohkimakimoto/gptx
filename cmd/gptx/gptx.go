package main

import (
	"fmt"
	"github.com/kohkimakimoto/gptx/internal"
	"os"
)

func main() {
	if err := internal.Run(os.Args); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
