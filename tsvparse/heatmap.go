package main

import (
	"fmt"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/palette/moreland"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

var _ plotter.GridXYZ = &HeatMap{}

var minimumTime, minimumGage uint64

type HeatMap struct {
	Values [][]float32
}

func (hm *HeatMap) Dims() (c, r int) {
	return len(hm.Values), len(hm.Values[0])
}

func (hm *HeatMap) Z(c, r int) float64 {
	f := float64(hm.Values[c][r])

	return f
}

func (hm *HeatMap) X(c int) float64 {
	if optimizeHeatmap {
		return float64(c*10) + float64(minimumTime)
	}
	return float64(c) + float64(minimumTime)
}

func (hm *HeatMap) Y(r int) float64 {
	return float64(r) + float64(minimumGage)
}

func plotHeatmap(res *BinaryResult, minTime, maxTime, minGage, maxGage uint64) {
	minimumGage = minGage
	minimumTime = minTime

	timeToValueMap := make(map[int64][]float32)
	if maxGage == uint64(res.CollCount) {
		maxGage--
	}

	startUnix := res.Timestamps[0]

	lineData := make([][]float32, 0)

	for idx, v := range res.Values {
		timeToValueMap[res.Timestamps[idx]-startUnix] = v
	}

	for i := minTime; i < maxTime; i++ {
		if optimizeHeatmap && i%10 != 0 {
			continue
		}

		gages, ok := timeToValueMap[int64(i)]

		if !ok {
			gages = make([]float32, res.CollCount) // Array of zeroes
		}
		lineData = append(lineData, gages[minGage:maxGage])
	}

	p := plot.New()
	hm := plotter.NewHeatMap(&HeatMap{lineData}, moreland.SmoothBlueRed().Palette(255))
	p.Add(hm)

	p.X.Label.Text = "Time (Seconds)"
	// if optimizeHeatmap {
	// 	p.X.Max = float64(maxTime / 10)
	// } else {
	// 	p.X.Max = float64(maxTime)
	// }
	p.X.Max = float64(maxTime)
	p.X.Min = float64(minTime)
	p.X.DashOffs = vg.Centimeter

	p.Y.Label.Text = "Gage"
	p.Y.Max = float64(maxGage)
	p.Y.Min = float64(minGage)
	p.Y.DashOffs = vg.Centimeter

	// fmt.Println(a, b)
	// hm.Max, hm.Min = 10, -10

	err := p.Save(8*vg.Inch, 8*vg.Inch, GRAPH_SAVE_LOCATION)

	if err != nil {
		fmt.Println("error:", err)
	}
}

func _plotHeatmap(res *BinaryResult, minTime, maxTime, minGage, maxGage uint64) {
	fmt.Println(len(res.Values), len(res.Timestamps))
	startUnix := res.Timestamps[0]

	lineData := make([][]float32, 0)
	for i := 0; i < res.RowCount; i++ {
		if optimizeHeatmap && i%10 != 0 {
			continue
		}

		seconds := res.Timestamps[i] - startUnix

		if minTime > uint64(seconds) || uint64(seconds) > maxTime {
			continue
		}

		arr := make([]float32, 0)
		for j := 0; j < len(res.Values[i]); j++ {
			if optimizeHeatmap && j%10 != 0 {
				continue
			}

			if minGage > uint64(j) || uint64(j) > maxGage {
				continue
			}
			arr = append(arr, res.Values[i][j])
		}
		lineData = append(lineData, arr)

		// lineData = append(lineData, plotter.XY{
		// 	X: float64(res.Timestamps[i] - startUnix),
		// 	Y: float64(res.Values[i][gage]),
		// })
	}

	// var startIdx, endIdx = 0, 0
	// for idx, t := range res.Timestamps {
	// 	if minTime == uint64(t-startUnix) {
	// 		startIdx = idx
	// 	}
	// 	if maxTime == uint64(t-startUnix) {
	// 		endIdx = idx
	// 		break
	// 	}
	// }

	// lineData := res.Values[startIdx:endIdx]
	// fmt.Println("Indicies", startIdx, endIdx)

	if len(lineData) == 0 || len(lineData[0]) == 0 {
		fmt.Println("LineData is empty")
		return
	}

	data := &HeatMap{lineData}
	// data := &HeatMap{res.Values}

	p := plot.New()
	// p.X.Max, p.Y.Max = data.Dims()
	hm := plotter.NewHeatMap(data, moreland.SmoothBlueRed().Palette(255))
	// fmt.Println(hm.Min, hm.Max)
	p.Add(hm)

	p.X.Label.Text = "Time (Seconds)"
	p.X.Max = float64(maxTime)
	p.X.Min = float64(minTime)

	p.Y.Label.Text = "Gage"
	p.Y.Max = float64(maxGage)
	p.Y.Min = float64(minGage)

	fmt.Println(p.Save(8*vg.Inch, 8*vg.Inch, GRAPH_SAVE_LOCATION))
}
