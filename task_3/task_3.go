package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type Task struct {
	duration time.Duration
	number   int
}

func Execute(t *Task, ch chan byte, wg *sync.WaitGroup) {
	log.Printf("Starting task №%d...", t.number)
	time.Sleep(t.duration)
	log.Printf("Task №%d done!", t.number)

	<-ch
	wg.Done()
}

func main() {

	args := os.Args

	if len(args) != 2 {
		fmt.Printf("Usage: ./%s [input-file]\n", filepath.Base(args[0]))
		os.Exit(-1)
	}

	fmt.Println("Input count of processors:")

	var maxProcessorCount int
	_, err := fmt.Scan(&maxProcessorCount)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	timefile, err := ioutil.ReadFile(args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	lines := strings.Split(string(timefile), "\n")

	var tasks []Task

	for i, str := range lines {
		str = strings.Trim(str, "\r") // For windows OS
		t, err := time.ParseDuration(str)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		i := i
		tasks = append(tasks, Task{duration: t, number: i})
	}

	ch := make(chan byte, maxProcessorCount)

	var wg sync.WaitGroup

	for i := 0; i < len(tasks); i++ {
		wg.Add(1)
		ch <- 1
		go Execute(&tasks[i], ch, &wg)
	}

	wg.Wait()
}
