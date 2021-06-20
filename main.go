package main

import (
	"os"

	"github.com/yuzuy/happi/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
