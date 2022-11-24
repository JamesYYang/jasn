package main

import (
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

type scanTask struct {
	ip      string
	port    int
	timeout int
}

var results sync.Map

func beginScan(sp scanParams) {
	tasks := generateTask(sp)
	runTask(tasks, sp.concurrency)
	printResult()
}

func generateTask(sp scanParams) []scanTask {
	tasks := make([]scanTask, 0)

	for _, ip := range sp.ipList {
		for _, port := range sp.portList {
			ipPort := scanTask{ip: ip.String(), port: port, timeout: sp.timeout}
			tasks = append(tasks, ipPort)
		}
	}

	return tasks
}

func runTask(tasks []scanTask, concurrency int) {
	wg := &sync.WaitGroup{}

	taskChan := make(chan scanTask, concurrency)

	for i := 0; i < concurrency; i++ {
		go scan(taskChan, wg)
	}

	for _, task := range tasks {
		wg.Add(1)
		taskChan <- task
	}

	close(taskChan)
	wg.Wait()
}

func scan(taskChan chan scanTask, wg *sync.WaitGroup) {
	for task := range taskChan {
		err := connect(task)
		saveResult(task, err)
		wg.Done()
	}
}

func saveResult(t scanTask, err error) {
	if err != nil {
		return
	}
	// log.Printf("find ip:%v, port: %v\n", t.ip, t.port)
	if t.port > 0 {
		v, ok := results.Load(t.ip)
		if ok {
			ports, ok1 := v.([]int)
			if ok1 {
				ports = append(ports, t.port)
				results.Store(t.ip, ports)
			}
		} else {
			ports := make([]int, 0)
			ports = append(ports, t.port)
			results.Store(t.ip, ports)
		}
	}
}

func printResult() {
	results.Range(func(key, value interface{}) bool {
		log.Printf("ip:%v\n", key)
		log.Printf("ports: %v\n", value)
		log.Println(strings.Repeat("-", 100))
		return true
	})
}

func connect(t scanTask) error {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%v:%v", t.ip, t.port), time.Duration(t.timeout)*time.Second)
	defer func() {
		if conn != nil {
			_ = conn.Close()
		}
	}()

	return err
}
