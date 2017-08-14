package main

import "time"
import "fmt"
import "strings"
import "math/rand"

type serverNode struct {
	id int
	serverLoad int
	maxRequestTime int // max amount of time it takes to complete a request
	c chan int
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


	maxRequestDelay int //delay before next request is sent out.
	
}

func loadBar(s serverNode) string  {
	maxLen := 20
	scale := 100/maxLen
	barLen := s.serverLoad/scale
	bar := strings.Repeat("=", barLen)
	space := strings.Repeat(" ", maxLen-barLen)
	return fmt.Sprintf("server_%d [%s%s] %d      ", s.id, bar, space, s.serverLoad)
}

func updateStats(serverNodeList []serverNode) {
	fmt.Printf("\033[0;0H")
	for _, serverNode := range serverNodeList {
		fmt.Println(loadBar(serverNode))
	}
}

func mockCpu(serverNodeList []serverNode) {
	for i := 0; i < len(serverNodeList); i++ {
		serverNodeList[i].serverLoad = rand.Intn(100)
	}
}

func clearScreen() {
	fmt.Printf("\033[0;0H")
	for i := 0; i < 50; i ++ {
		fmt.Println(strings.Repeat(" ", 80))
	}
}

func runServer(s serverNode) {
	for {
		_ = <-s.c
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(s.maxRequestTime)))
		s.c <- 1
	}
}

func newServer(id int, maxTime int) serverNode {
	ch := make(chan int)
	return serverNode{id, 0, maxTime, ch}
}

func main() {
	//setup
	clearScreen()
	
	serverNodeList := []serverNode{
		newServer(0, 50),
		newServer(1, 100),
		newServer(2, 150),
	}

	for _, server := range(serverNodeList) {
		go runServer(server)
	}

	for i := 0; i < 10; i++ {
		updateStats(serverNodeList)
		time.Sleep(500 * time.Millisecond)
		mockCpu(serverNodeList)
	}
}
