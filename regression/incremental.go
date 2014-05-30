package main

import (
	"encoding/json"
	"fmt"
	"github.com/GaryBoone/GoStats/stats"
	"io/ioutil"
	"math/rand"
	"time"
)

type Point struct {
	X float64
	Y float64
}

type Graph struct {
	Name             string
	DataPoints       []Point
	RegressionPoints []Point
	RSquared         float64
}

type HasGraphs struct {
	Graphs []Graph
}

func main() {
	rand.Seed(int64(time.Now().Nanosecond()))

	points := make([]Point, 0, 100)

	xmin := 0.0
	xmax := 0.0

	var r stats.Regression

	for i := 0; i < 25; i++ {
		x := float64(i) * 3.0
		y := rand.Float64()*100.0 - 25.0 // uniform samples in {-25, 75}
		r.Update(x, y)

		points = append(points, Point{x, y})

		if x < xmin {
			xmin = x
		}
		if x > xmax {
			xmax = x
		}
	}

	// transform slope and intercept into 2 points defining the line between min-x and max-x
	// x1 = min-x
	// y1 = (slope * x1) + intercept
	// x2 = max-x
	// y2 = (slope * x2) + intercept
	regression := []Point{
		Point{xmin, (r.Slope() * xmin) + r.Intercept()},
		Point{xmax, (r.Slope() * xmax) + r.Intercept()},
	}

	hg := HasGraphs{
		Graphs: []Graph{
			Graph{
				Name:             "test points",
				DataPoints:       points,
				RegressionPoints: regression,
				RSquared:         r.RSquared(),
			},
		},
		// then add another graph!
	}

	b, err := json.Marshal(hg)
	if err != nil {
		panic(err)
	}

	if err = ioutil.WriteFile("static/regression.json", b, 0644); err != nil {
		panic(err)
	}

	fmt.Println("json=", string(b))
	fmt.Println("count", r.Count())
	fmt.Println("slope", r.Slope())
	fmt.Println("intercept", r.Intercept())
	fmt.Println("r-squared", r.RSquared())
	fmt.Println("slope standard error", r.SlopeStandardError())
	fmt.Println("intercept standard error", r.InterceptStandardError())

}
