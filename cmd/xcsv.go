package cmd

import (
	"bufio"
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

	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	result := [][]string{
		{"#2theta", "Intensity"},
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "*") || strings.HasPrefix(line, "#") {
			continue
		}
		result = append(result, strings.Split(line, " ")[:2])
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	if err := createCSV(&result, strings.Replace(path, ".ras", ".csv", 1)); err != nil {
		return err
	}

	return nil
}

func init() {
	rootCmd.AddCommand(xcsvCmd)
}
