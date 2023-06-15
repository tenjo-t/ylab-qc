package cmd

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var xcsvCmd = &cobra.Command{
	Use:   "xcsv",
	Short: "RAS to CSV converter",
	Long:  "RAS to CSV converter",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := runXcsv(args[0]); err != nil {
			cmd.Println(err)
			os.Exit(1)
		}
	},
}

func runXcsv(path string) error {
	if !strings.HasSuffix(path, ".ras") {
		return fmt.Errorf("only RAS file")
	}

	// reader
	r, err := os.Open(path)
	if err != nil {
		return err
	}
	defer r.Close()

	s := bufio.NewScanner(r)

	// writer
	f, err := os.Create(strings.Replace(path, ".ras", ".csv", 1))
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)

	// header
	if err := w.Write([]string{"#2theta", "Intensity"}); err != nil {
		return err
	}

	for s.Scan() {
		line := s.Text()

		// skip
		if strings.HasPrefix(line, "*") || strings.HasPrefix(line, "#") {
			continue
		}

		// write
		if err := w.Write(strings.Split(line, " ")[:2]); err != nil {
			return err
		}
	}

	if err := s.Err(); err != nil {
		return err
	}

	return nil
}

func init() {
	rootCmd.AddCommand(xcsvCmd)
}
