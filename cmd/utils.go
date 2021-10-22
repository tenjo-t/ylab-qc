package cmd

import (
	"encoding/csv"
	"math"
	"os"
	"strconv"
	"strings"
)

func createCSV(data *[][]string, path string) error {
	output, err := os.Create(strings.Replace(path, ".ras", ".csv", 1))
	if err != nil {
		return err
	}
	defer output.Close()

	w := csv.NewWriter(output)
	if err := w.WriteAll(*data); err != nil {
		return err
	}
	w.Flush()

	return nil
}

func squares(n ...float64) float64 {
	sum := 0.
	for _, i := range n {
		sum += math.Pow(i, 2)
	}

	return sum
}

func calcTowTheta(n float64) float64 {
	return math.Asin(opt.wl*math.Sqrt(n)/2/opt.lc) * 180 / math.Pi
}

func calcLatticeConstant(n, th float64) string {
	return strconv.FormatFloat(opt.wl*math.Sqrt(n)/2/math.Sin(th*math.Pi/180), 'f', -1, 64)
}

func calcNR(th float64) string {
	r := th * math.Pi / 180
	i := math.Pow(math.Cos(r), 2)
	return strconv.FormatFloat(i/math.Sin(r)+i/r, 'f', -1, 64)
}
