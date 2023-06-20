package cmd

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

type xpeakOptions struct {
	lc   float64
	wl   float64
	isAC bool
}

var opt = &xpeakOptions{}

var xpeak = &cobra.Command{
	Use:   "xpeak",
	Short: "XRD peak search",
	Long:  "XRD peak search",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := runXpeak(args[0]); err != nil {
			cmd.Println(err)
			os.Exit(1)
		}
	},
}

func runXpeak(filepath string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	var calc string
	if opt.isAC {
		calc = "AC11"
	} else {
		calc = "QC"
	}
	f, err := os.Open(path.Join(home, "qc", calc+".csv"))
	if err != nil {
		return err
	}
	defer f.Close()

	r := csv.NewReader(f)

	var peaks *[][]string
	if opt.isAC {
		p, err := searchAC(r)
		if err != nil {
			return err
		}
		peaks = p
	} else {
		p, err := searchQC(r)
		if err != nil {
			return err
		}
		peaks = p
	}

	if err := createCSV(peaks, filepath+"_peak.csv"); err != nil {
		return err
	}

	return nil
}

func searchAC(r *csv.Reader) (*[][]string, error) {
	peaks := [][]string{
		{"h", "k", "l", "2theta"},
	}
	scanner := bufio.NewScanner(os.Stdin)

	for {
		rec, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if strings.HasPrefix(rec[0], "#") {
			continue
		}

		h, _ := strconv.ParseFloat(rec[0], 64)
		k, _ := strconv.ParseFloat(rec[1], 64)
		l, _ := strconv.ParseFloat(rec[2], 64)
		N := squares(h, k, l)
		towTheta := calcTowTheta(N)

		fmt.Printf("(%s, %s, %s) ~%.2f: ", rec[0], rec[1], rec[2], towTheta)
		for scanner.Scan() {
			text := scanner.Text()
			if text != "" {
				th, err := strconv.ParseFloat(text, 64)
				if err != nil {
					fmt.Print("Retry: ")
					continue
				}
				peaks = append(peaks, []string{rec[0], rec[1], rec[2], text, calcNR(th), calcLatticeConstant(N, th)})
			}
			break
		}
	}

	return &peaks, nil
}

func searchQC(r *csv.Reader) (*[][]string, error) {
	peaks := [][]string{
		{"h", "k", "l", "m", "n", "o", "2theta", "NR", "lattice constant"},
	}
	scanner := bufio.NewScanner(os.Stdin)

	for {
		rec, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if strings.HasPrefix(rec[0], "#") {
			continue
		}

		h, _ := strconv.ParseFloat(rec[0], 64)
		k, _ := strconv.ParseFloat(rec[1], 64)
		l, _ := strconv.ParseFloat(rec[2], 64)
		m, _ := strconv.ParseFloat(rec[3], 64)
		n, _ := strconv.ParseFloat(rec[4], 64)
		o, _ := strconv.ParseFloat(rec[5], 64)

		r1 := math.Sqrt(5)*h + k + l + m + n + o
		r2 := h + math.Sqrt(5)*k + l - m - n + o
		r3 := h + k + math.Sqrt(5)*l + m - n - o
		r4 := h - k + l + math.Sqrt(5)*m + n - o
		r5 := h - k - l + m + math.Sqrt(5)*n + o
		r6 := h + k - l - m + n + math.Sqrt(5)*o

		N := squares(r1, r2, r3, r4, r5, r6) / 20

		towTheta := calcTowTheta(N)

		fmt.Printf("(%s, %s, %s, %s, %s, %s) ~%.2f: ", rec[0], rec[1], rec[2], rec[3], rec[4], rec[5], towTheta)
		for scanner.Scan() {
			text := scanner.Text()
			if text != "" {
				th, err := strconv.ParseFloat(text, 64)
				if err != nil {
					fmt.Print("Retry: ")
					continue
				}
				peaks = append(peaks, []string{rec[0], rec[1], rec[2], rec[3], rec[4], rec[5], text, calcNR(th), calcLatticeConstant(N, th)})
			}
			break
		}
	}

	return &peaks, nil
}

func init() {
	root.AddCommand(xpeak)
	xpeak.Flags().Float64VarP(&opt.lc, "lc", "l", -1, "lattice constant")
	xpeak.MarkFlagRequired("lc")
	xpeak.Flags().Float64VarP(&opt.wl, "wl", "w", 1.540593, "wave length")
	xpeak.Flags().BoolVarP(&opt.isAC, "ac", "a", false, "structure of peak search \"AC11\" (default QC)")
}
