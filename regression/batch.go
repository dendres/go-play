package main

import (
	"fmt"
	"github.com/GaryBoone/GoStats/stats"
	"math/rand"
	"time"
)

type Point struct {
	x float64
	y float64
}

const NUM_SAMPLES = 5

func main() {

	rand.Seed(int64(time.Now().Nanosecond()))

	xData := make([]float64, NUM_SAMPLES)
	for i := 0; i < NUM_SAMPLES; i++ {
		xData[i] = float64(i) * 3.0
	}
	fmt.Println("xData", xData)

	yData := make([]float64, NUM_SAMPLES)
	for i := 0; i < NUM_SAMPLES; i++ {
		x := rand.Float64()*100 - 25 // uniform samples in {-25, 75}
		yData[i] = x
	}

	fmt.Println("yData", yData)

	var slope, intercept, rsquared, count, slopeStdErr, intcptStdErr = stats.LinearRegression(xData, yData)

	fmt.Println("slope", slope)
	fmt.Println("intercept", intercept)
	fmt.Println("r-squared", rsquared)
	fmt.Println("count", count)
	fmt.Println("slope standard error", slopeStdErr)
	fmt.Println("intercept standard error", intcptStdErr)

	// output the data and regression in a form suitable for d3.js graph drawing

	// start with the original data

}
