package main

import (
	"fmt"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

func plotGraph(res *BinaryResult, gage, minTime, maxTime int) {
	fmt.Println("Plotting with", minTime, maxTime)

	p := plot.New()
	p.Add(plotter.NewGrid())

	// p.X.Max = float64(res.RowCount)
	p.X.Label.Text = "Time"
	p.Y.Label.Text = "Value"

	// lineData := make(plotter.XYs, res.RowCount)
	lineData := make(plotter.XYs, 0)

	startUnix := res.Timestamps[0]
	for i := 0; i < res.RowCount; i++ {
		seconds := res.Timestamps[i] - startUnix

		if minTime > int(seconds) || int(seconds) > maxTime {
			// fmt.Println("skipping")
			continue
		}

		lineData = append(lineData, plotter.XY{
			X: float64(res.Timestamps[i] - startUnix),
			Y: float64(res.Values[i][gage]),
		})
	}

	plotutil.AddLinePoints(p, lineData)
	p.Save(8*vg.Inch, 8*vg.Inch, GRAPH_SAVE_LOCATION)

	// return lineData, nil
}
