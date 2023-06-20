package cmd

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var dta = &cobra.Command{
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

type Step struct {
	start int
	end   int
	w     *csv.Writer
}

func splitData(path string, base string) error {
	// reader
	r, err := os.Open(path)
	if err != nil {
		return err
	}
	defer r.Close()

	s := bufio.NewScanner(r)

	num := 0
	var steps []Step
	var step Step

	for s.Scan() {
		l := s.Text()
		num++

		// read steps
		if strings.HasPrefix(l, "#HD") {
			if strings.Contains(l, "HOLD") {
				continue
			}

			r := strings.Split(l, "\t")
			start, _ := strconv.Atoi(r[2])
			end, _ := strconv.Atoi(r[3])

			f, err := os.Create(fmt.Sprint(base, "_", len(steps)+1, ".csv"))
			if err != nil {
				return err
			}
			defer f.Close()

			w := csv.NewWriter(f)
			w.Write([]string{"Temp.", ""})

			if len(steps) == 0 {
				step = Step{start, end, w}
			} else {
				steps = append(steps, Step{start, end, w})
			}

			continue
		}

		// read data row
		if strings.HasPrefix(l, "#GD") {
			if num > step.end {
				step = steps[0]
				steps = steps[1:]
			}

			data := strings.Split(l, "\t")
			step.w.Write([]string{data[2], data[4]})

		}
	}

	if err := s.Err(); err != nil {
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

plot PATH."_2.csv" using 1:($2 -50) title "1st Heating" with l lw 4 lc "#93003a", \
     PATH."_3.csv"                  title "1st Cooling" with l lw 4 lc "#00429d", \
     PATH."_4.csv" using 1:($2 -50) title "2nd Heating" with l lw 4 lc "#f4777f", \
     PATH."_5.csv"                  title "2nd Cooling" with l lw 4 lc "#73a2c6"



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
	root.AddCommand(dta)
	dta.Flags().Float64SliceVarP(&rng, "range", "r", []float64{400., 950.}, "range of Temp.")
}
