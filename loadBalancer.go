package main

import "fmt"
import "time"
import "strings"

type serverNode struct {
	id int
	serverLoad float32
	maxRequestTime float32 // max amount of time it takes to complete a request

}

type loadBalancer struct {
	//stats
	requestsMade int
	requestOk int
	requestError int

	// FIFO stack that tracks available servers
	readyStack []int

	//
}

type users struct {


	maxRequestDelay int //delay before next request is sen[M Ct out.
	
}

func main() {
	for i := 0; i < 101; i++ {
		fmt.Printf("\r[%s%s] %v", strings.Repeat("=", i/5), strings.Repeat(" ", 20-(i/5)), i)
		time.Sleep(200 * time.Millisecond)
	}
}
