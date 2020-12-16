package log

import (
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

func Draw(nodeSample map[string]*NodeSample) {

	for name, node := range nodeSample {
		p, _ := plot.New()
		p.Title.Text = name
		p.X.Label.Text = "times"
		p.Y.Label.Text = "%"
		var xys []plotter.XY = make([]plotter.XY, 0)
		for index, value := range node.SampleSumPoint {
			xy := plotter.XY{float64(index), value}
			xys = append(xys, xy)
		}
		plotutil.AddLinePoints(p, plotter.XYs(xys))
		p.Save(800*vg.Inch, 4*vg.Inch, name+".png")
		//p.Save(16*vg.Inch, 4*vg.Inch, name+".png")
	}
}
