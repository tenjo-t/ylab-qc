package cmd

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var dtaCmd = &cobra.Command{
	Use:   "dta",
	Short: "Make csv and plt file from DTA ASC file",
	Long:  "Make csv and plt file from DTA ASC file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := runDta(args[0]); err != nil {
			cmd.Println(err)
			os.Exit(1)
		}
	},
}

func runDta(path string) error {
	base := regexp.MustCompile(`.\w+$`).ReplaceAllString(path, "")
	base = strings.TrimPrefix(base, `.\`)
	base = strings.TrimPrefix(base, `./`)

	if err := splitData(path, base); err != nil {
		return err
	}
	if err := makePlt(base); err != nil {
		return err
	}

	return nil
}

func splitData(path string, base string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	num := 0
	var steps [][]int
	list := [][]string{{"#Temp.", "DTA"}}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		num++

		// read steps
		if strings.HasPrefix(line, "#HD") {
			if strings.Contains(line, "HOLD") {
				continue
			}
			var step []int
			for _, v := range strings.Split(line, "\t")[2:4] {
				idx, _ := strconv.Atoi(v)
				step = append(step, idx)
			}
			// skip first HEATING
			if step[0] == 0 {
				continue
			}
			steps = append(steps, step)
			continue
		}

		// read data row
		if strings.HasPrefix(line, "#GD") {
			idx := -1
			end := false
			for i, v := range steps {
				if num == v[1] {
					end = true
					idx = i + 1
				}
			}

			data := strings.Split(line, "\t")
			list = append(list, []string{data[2], data[4]})

			// if end of step, create csv file
			if end {
				if err := createCSV(&list, fmt.Sprint(base, "_", idx, ".csv")); err != nil {
					return err
				}
				list = [][]string{{"#Temp.", "DTA"}}
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

const PLT_TEMP = `PATH = "%PATH%"

set encoding utf8
set terminal pdfcairo enhanced font "Source Han Sans JP, 12" size 12cm, 9cm
set output PATH.".pdf"

# use csv file
set datafile separator ","

set xlabel "{/:Italic T} (℃)"
set ylabel "DTA (μV)"
set xrange [400:950]

plot PATH."_1.csv" using 1:($2 -50) title "1st Heating" with l lw 4 lc "#93003a", \
     PATH."_2.csv"                  title "1st Cooling" with l lw 4 lc "#00429d", \
     PATH."_3.csv" using 1:($2 -50) title "2nd Heating" with l lw 4 lc "#f4777f", \
     PATH."_4.csv"                  title "2nd Cooling" with l lw 4 lc "#73a2c6"



set terminal pngcairo enhanced font "Source Han Sans JP, 22" size 1000, 750
set output PATH.".png"
set key left

replot
`

var rng []float64

func makePlt(base string) error {

	f, err := os.Create(base + ".plt")
	if err != nil {
		return err
	}

	r := strings.NewReplacer("%PATH%", base)
	w := bufio.NewWriter(f)

	if _, err := w.WriteString(r.Replace(PLT_TEMP)); err != nil {
		return err
	}
	if err := w.Flush(); err != nil {
		return err
	}

	if err := f.Close(); err != nil {
		return err
	}

	return nil
}

func init() {
	rootCmd.AddCommand(dtaCmd)
	dtaCmd.Flags().Float64SliceVarP(&rng, "range", "r", []float64{400., 950.}, "range of Temp.")
}
