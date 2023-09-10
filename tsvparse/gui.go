package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"

	"fyne.io/fyne/v2"
	fyneapp "fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var app fyne.App

const GRAPH_SAVE_LOCATION = "./temp_generated_graph.png"

func StartGUI() {
	app = fyneapp.New()
	app.Settings().SetTheme(theme.DarkTheme())
	defer app.Quit()

	w := app.NewWindow(fmt.Sprintf("Graph Plotter (%s)", VERSION))
	defer w.Close()

	w.SetCloseIntercept(func() {
		fmt.Println(os.Remove(GRAPH_SAVE_LOCATION))
		w.Close()
	})

	w.Resize(fyne.NewSize(600, 480))
	renderUploadWindow(w)
	w.ShowAndRun()

}

func renderUploadWindow(w fyne.Window) {
	w.SetContent(widget.NewButton("Upload", func() {
		rc := make(chan io.Reader, 1)
		dialog.ShowFileOpen(func(uc fyne.URIReadCloser, err error) {
			rc <- uc
		}, w)

		go func() {
			renderLoadingWindow(w, bufio.NewReader(<-rc))
			close(rc)
		}()
	}))
}

func renderLoadingWindow(w fyne.Window, r *bufio.Reader) {
	loadingBar := widget.NewProgressBar()

	w.SetContent(loadingBar)
	// w.Resize(fyne.NewSize(200, 100))

	ProgressCallback = func(progress float32, done bool) {
		loadingBar.SetValue(float64(progress))
	}

	res, err := parseBinary(r)

	if err != nil {
		dialog.ShowError(err, w)
		return
	}

	renderPlotterWindow(w, &res)
}

var contextSpecificCanvas fyne.CanvasObject = layout.NewSpacer()
var currentSelection string

const (
	GRAPH   = "Graph"
	HEATMAP = "Heatmap"
)

func renderPlotterWindow(w fyne.Window, res *BinaryResult) {
	plotButton := widget.NewButton("Plot", func() {
		plotFunc()
		// renderPlotterWindow(w, res)
	})
	plotButton.Importance = widget.HighImportance

	topLabel := widget.NewLabel(fmt.Sprintf("Processed %d rows and %d collumns.", res.RowCount, res.CollCount))

	var selection, selectionLabel fyne.CanvasObject

	selectionLabel = widget.NewLabel("Select an option to plot.")
	selection = widget.NewSelect([]string{GRAPH, HEATMAP}, func(s string) {
		currentSelection = s
		updateSelection(w, res)
		// w.Canvas().Refresh(contextSpecificCanvas)

		vbucks := container.NewVBox(topLabel, selectionLabel, selection, contextSpecificCanvas, plotButton)

		w.SetContent(vbucks)
	})

	// plotButton := widget.NewButton("Plot", func() { execPlot(w, res) })

	vbucks := container.NewVBox(topLabel, selectionLabel, selection, contextSpecificCanvas, plotButton)

	w.SetContent(vbucks)
}

var label = widget.NewLabel

var plotFunc func() = func() { /* Empty so that it dont SIGSEGV*/ }

var optimizeHeatmap bool = false

func updateSelection(w fyne.Window, res *BinaryResult) {
	// var err error

	showDiag := func(s string) {
		dialog.ShowCustom("Invalid input data.", "Ok", label(s), w)
	}

	if currentSelection == GRAPH {
		gageSelection := widget.NewEntry()
		gageSelection.SetPlaceHolder(fmt.Sprintf("Select a gage from 0-%d", res.CollCount))

		minTime := widget.NewEntry()
		minTime.SetPlaceHolder(fmt.Sprint(0))

		maxTime := widget.NewEntry()
		ts := res.Timestamps
		maxTime.SetPlaceHolder(fmt.Sprint(ts[len(ts)-1] - ts[0]))
		plotFunc = func() {
			// fmt.Println("Plotfunc")
			var gage, minT, maxT uint64 = 0, 0, uint64(res.CollCount)
			var err error
			gage, err = strconv.ParseUint(gageSelection.Text, 10, 32)

			if err != nil {
				showDiag("Gage is not a valid number")
				return
			}

			minT, err = strconv.ParseUint(minTime.Text, 10, 32)

			skip1 := false
			if len(minTime.Text) == 0 {
				minT = 0
				skip1 = true
			}

			if err != nil && !skip1 {
				showDiag("Min time is not a valid number")
				return
			}

			maxT, err = strconv.ParseUint(maxTime.Text, 10, 32)

			skip2 := false
			if len(maxTime.Text) == 0 {
				maxT = uint64(ts[len(ts)-1] - ts[0])
				skip2 = true
			}

			if err != nil && !skip2 {
				showDiag("Max time is not a valid number")
				return
			}

			if minT > maxT {
				showDiag("Min time is greater than Max time")
				return
			}

			plotGraph(res, int(gage), int(minT), int(maxT))
			renderResultWindow(w, res)
		}

		contextSpecificCanvas = container.NewVBox(gageSelection,
			label("Min time:"), minTime,
			label("Max time:"), maxTime,
		)
	} else if currentSelection == HEATMAP {
		mingage := widget.NewEntry()
		mingage.SetPlaceHolder(fmt.Sprint(0))

		maxgage := widget.NewEntry()
		maxgage.SetPlaceHolder(fmt.Sprint(res.CollCount))

		minTime := widget.NewEntry()
		minTime.SetPlaceHolder(fmt.Sprint(0))

		maxTime := widget.NewEntry()
		ts := res.Timestamps
		maxTime.SetPlaceHolder(fmt.Sprint(ts[len(ts)-1] - ts[0]))

		optCheckbox := widget.NewCheck("Optimize", func(b bool) { optimizeHeatmap = b })
		optCheckbox.SetChecked(optimizeHeatmap)

		plotFunc = func() {
			// Validate input, later

			var minT, maxT, minG, maxG uint64 = 0, 0, 0, uint64(res.CollCount)
			var err error

			// Time clamping
			{
				minT, err = strconv.ParseUint(minTime.Text, 10, 32)

				skip1 := false
				if len(minTime.Text) == 0 {
					minT = 0
					skip1 = true
				}

				if err != nil && !skip1 {
					showDiag("Min time is not a valid number")
					return
				}

				maxT, err = strconv.ParseUint(maxTime.Text, 10, 32)

				skip2 := false
				if len(maxTime.Text) == 0 {
					maxT = uint64(ts[len(ts)-1] - ts[0])
					skip2 = true
				}

				if err != nil && !skip2 {
					showDiag("Max time is not a valid number")
					return
				}

				if minT > maxT {
					showDiag("Min time is greater than Max time")
					return
				}
			}

			// Gage clamping
			{
				minG, err = strconv.ParseUint(mingage.Text, 10, 32)

				skip1 := false
				if len(mingage.Text) == 0 {
					minG = 0
					skip1 = true
				}

				if err != nil && !skip1 {
					showDiag("Min gage is not a valid gage")
					return
				}

				maxG, err = strconv.ParseUint(maxgage.Text, 10, 32)

				skip2 := false
				if len(maxgage.Text) == 0 {
					maxG = uint64(res.CollCount)
					skip2 = true
				}

				if err != nil && !skip2 {
					showDiag("Max gage is not a valid gage")
					return
				}

				if minG > maxG {
					showDiag("Min gage is greater than Max gage")
					return
				}
			}

			plotHeatmap(res, minT, maxT, minG, maxG)
			renderResultWindow(w, res)
		}

		contextSpecificCanvas = container.NewVBox(
			label("Min gage:"), mingage,
			label("Max gage:"), maxgage,
			label("Min time:"), minTime,
			label("Max time:"), maxTime,
			optCheckbox,
		)
	}
}

func renderResultWindow(w fyne.Window, res *BinaryResult) {
	backBtn := widget.NewButton("Back", func() { renderPlotterWindow(w, res) })
	backBtn.Importance = widget.DangerImportance

	saveBtn := widget.NewButton("Save", func() {
		dialog.ShowFileSave(func(uc fyne.URIWriteCloser, err error) {
			file, err := os.ReadFile(GRAPH_SAVE_LOCATION)

			if err != nil {
				dialog.ShowError(err, w)
				return
			}

			uc.Write(file)
		}, w)
	})
	saveBtn.Importance = widget.HighImportance

	img := canvas.NewImageFromFile(GRAPH_SAVE_LOCATION)
	img.FillMode = canvas.ImageFillOriginal

	w.SetContent(container.NewVBox(img, container.NewHBox(backBtn, layout.NewSpacer(), saveBtn)))
}
