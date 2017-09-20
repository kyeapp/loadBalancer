package main

import "time"
import "fmt"
import "strings"

//import "math/rand"
import "os"
import "strconv"

type serverNode struct {
	id              int
	requestsServed  int
	requestsLost    int
	currentRequests float64
	maxRequests     float64
	requestTime     int
}

func (s *serverNode) assignRequest() {
	if s.currentRequests < s.maxRequests {
		s.currentRequests++
	} else {
		s.requestsLost++
	}
}

func (s *serverNode) run() {
	ticker := time.NewTicker(time.Millisecond * time.Duration(s.requestTime))
	for {
		_ = <-ticker.C
		//println(s.id, s.requestTime)
		if s.currentRequests > 0 {
			s.currentRequests--
			s.requestsServed++
		}
	}
}

//input ID, and maxRequestTime
func newServer(id int, requestTime int) serverNode {
	return serverNode{id, 0, 0, 50, 100, requestTime}
}

//=============================================================================================
type loadBalancer struct {
	totalrequests int
}

//load balancing simulation algorithm
func roundRobin(serverList []serverNode, reqPerSec int) {
	i := 0
	end := len(serverList)
	delay := 1000000000 / reqPerSec
	ticker := time.NewTicker(time.Nanosecond * time.Duration(delay))
	for {
		_ = <-ticker.C                   //wait for new request to come in
		go serverList[i].assignRequest() //assign Request to server

		i++
		if i == end {
			i = 0
		}

	}
}

func weightedRoundRobin(serverList []serverNode, reqPerSec int) {	
	total := 0
	for _, node := range(serverList) {
		total += node.requestTime
	}

	var queue []int
	for i, node := range(serverList) {
		w := int(float64(total)/float64(node.requestTime)*4)
		t := make([]int, w, w)
		for j := range(t) {
			t[j] = i
		}
		queue = append(queue, t...);
	}


	i := 0
	end := len(queue) -1
	delay := 1000000000 / reqPerSec
	ticker := time.NewTicker(time.Nanosecond * time.Duration(delay))
	for {
		_ = <-ticker.C                   //wait for new request to come in
		go serverList[queue[i]].assignRequest() //assign Request to server

		i++
		if i == end {
			i = 0
		}
	}
}

func loadBar(s serverNode, maxLen int) string {
	load := s.currentRequests / s.maxRequests
	barLen := int(load * float64(maxLen))
	bar := strings.Repeat("=", barLen)
	space := strings.Repeat(" ", maxLen-barLen)
	return fmt.Sprintf("server_%d [%s%s] %3.0f%%  served: %d   lost: %d  current: %3.0f", s.id, bar, space, load*100, s.requestsServed, s.requestsLost, s.currentRequests)
}

func updateStats(serverNodeList []serverNode) {
	ticker := time.NewTicker(time.Millisecond * 500)
	for {
		select {
		case <-ticker.C:
			fmt.Printf("\033[0;0H")
			fmt.Println(time.Now())
			for _, serverNode := range serverNodeList {
				fmt.Println(loadBar(serverNode, 20))
			}
		}
	}

}

func clearScreen() {
	fmt.Printf("\033[0;0H")
	for i := 0; i < 50; i++ {
		fmt.Println(strings.Repeat(" ", 80))
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("no requests per second in arg 1")
		os.Exit(1)
	}
	rps, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("need valid requests per second in arg 1")
		os.Exit(1)
	}
	clearScreen()

	//create servers
	serverNodeList := []serverNode{
		newServer(0, 50),
		newServer(1, 100),
		newServer(2, 150),
	}

	go updateStats(serverNodeList)

	//start all the servers
	for i, _ := range serverNodeList {
		go serverNodeList[i].run()
	}

	weightedRoundRobin(serverNodeList, rps)
}
