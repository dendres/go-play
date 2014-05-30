//
// regression_demo.go
//
// Can be compiled and run from the main GoStats directory using:
// $ (cd demos/regression/ && make clean && make)
// $ demos/regression/regression_demo
//
// Author:    Gary Boone
//

package main

import (
	"fmt"
	"github.com/GaryBoone/GoStats/stats"
	"math/rand"
	"time"
)

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

	fmt.Printf("then request regression results, for example:\n")
	fmt.Printf("      var slope, intercept, _, _, _, _ = stats.LinearRegression(xData, yData)\n")
	fmt.Printf("** Linear Regression:\n")
	fmt.Printf("slope = %v\n", slope)
	fmt.Printf("intercept = %v\n", intercept)
	fmt.Printf("r-squared = %v\n", rsquared)
	fmt.Printf("count = %v\n", count)
	fmt.Printf("slope standard error = %v\n", slopeStdErr)
	fmt.Printf("intercept standard error = %v\n", intcptStdErr)

}
