package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Task struct {
	duration time.Duration
	number   int
}

func (t *Task) Execute() {
	log.Printf("Starting task №%d...", t.number)
	time.Sleep(t.duration)
	log.Printf("Task №%d done!", t.number)
}

func main() {

	args := os.Args

	if len(args) != 2 {
		fmt.Printf("Usage: ./%s [input-file]\n", filepath.Base(args[0]))
		os.Exit(-1)
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
		tasks = append(tasks, Task{duration: t, number: i})
	}

	for _, task := range tasks {
		task.Execute()
	}

}
