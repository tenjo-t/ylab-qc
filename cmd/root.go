package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	root = &cobra.Command{
		Use:   "qc",
		Short: "Quasicrystal analysis CLI tools",
		Long:  "Quasicrystal analysis CLI tools",
	}
)

func Execute() {
	root.SetOutput(os.Stdout)
	if err := root.Execute(); err != nil {
		root.SetOutput(os.Stderr)
		root.Println(err)
		os.Exit(1)
	}
}

func init() {
}
