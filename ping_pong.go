// http://talks.golang.org/2013/advconc.slide

package main

import (
	"fmt"
	"time"
)

type Ball struct{ hits int }

func player(name string, table chan *Ball) {
	for {
		ball := <-table
		ball.hits++
		fmt.Println(name, ball.hits)
		time.Sleep(100 * time.Millisecond)
		table <- ball
	}
}

func main() {
	table := make(chan *Ball)

	go player("first", table)
	go player("second", table)
	go player("third", table)
	go player("forth", table)

	fmt.Println("ready")
	time.Sleep(1 * time.Second)
	fmt.Println("go!")
	table <- new(Ball) // game on; toss the ball

	time.Sleep(1 * time.Second)
	fmt.Println("Game Over")
	<-table // game over; grab the ball
	panic("show me the stacks")
}
