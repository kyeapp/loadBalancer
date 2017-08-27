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
	min int
	max int
	diff int
	requestTime  int // max amount of time it takes to complete a request +- 10%
}

func (s *serverNode) processRequest() {
	if s.currentRequests < s.maxRequests {
		s.currentRequests++
		// if you don't contrain the random bounds tightly it could take a while to spool up the average load
		//time.Sleep(time.Millisecond * time.Duration(rand.Intn(s.requestTime)))
		time.Sleep(time.Millisecond * time.Duration(s.requestTime))
		s.currentRequests--
		s.requestsServed++
	} else {
		//request is lost if the server is full
		s.requestsLost++
	}
}

//input ID, and maxRequestTime
func newServer(id int, requestTime int) serverNode {
	return serverNode{id, 0, 0, 0, 32, 0, 0, 0, requestTime-5}
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
		_ = <-ticker.C
		go serverList[i].processRequest()
		i++
		if i == end {
			i = 0
		}

	}
}

func loadBar(s serverNode, maxLen int) string {
	load := s.currentRequests/s.maxRequests
	barLen := int(load * float64(maxLen))
	bar := strings.Repeat("=", barLen)
	space := strings.Repeat(" ", maxLen-barLen)
	return fmt.Sprintf("server_%d [%s%s] %0.2f  served: %d   lost: %d  current: %3.0f", s.id, bar, space, load, s.requestsServed, s.requestsLost, s.currentRequests)
}

func updateStats(serverNodeList []serverNode) {
	ticker := time.NewTicker(time.Millisecond * 500)
	go func() {
		for {
			select {
			case <-ticker.C:
				go func() {
					
					fmt.Printf("\033[0;0H")
					fmt.Println(time.Now())
					for _, serverNode := range serverNodeList {
						fmt.Println(loadBar(serverNode, 20))
					}
				}()
			}
		}
	}()

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

	updateStats(serverNodeList)
	roundRobin(serverNodeList, rps)
}
