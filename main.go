package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

var (
	count count32
	q     Queue
)

type count32 int32

func (c *count32) inc() int32 {
	return atomic.AddInt32((*int32)(c), 1)
}

func (c *count32) get() int32 {
	return atomic.LoadInt32((*int32)(c))
}

type Task struct {
	duration time.Duration
	number   int32
	resChan  chan byte
}

func Execute(t *Task) {
	log.Printf("Starting task №%d...", t.number)
	time.Sleep(t.duration)
	log.Printf("Task №%d done!", t.number)

	if t.resChan != nil {
		<-t.resChan
	}
}

type Queue struct {
	sync.RWMutex
	items []Task
}

func (q *Queue) Push(item Task) {
	q.Lock()
	defer q.Unlock()
	q.items = append(q.items, item)
}

func (q *Queue) Pop() Task {
	q.Lock()
	defer q.Unlock()
	if len(q.items) == 0 {
		return Task{}
	}
	item := q.items[0]
	q.items = q.items[1:]
	return item
}

func (q *Queue) GetAll() []Task {
	q.RLock()
	defer q.RUnlock()
	return q.items
}

func TaskLoop(q *Queue, terminator <-chan byte) {
	for {
		select {
		case <-terminator:
			return
		default:
			task := q.Pop()
			if task.number == 0 {
				time.Sleep(time.Millisecond) // waiting 1s for new tasks
				continue
			}
			Execute(&task)
		}
	}

}

func add(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	keys := r.URL.Query()["timeDuration"]
	if keys == nil {
		http.Error(w, "Require timeDuration parameter", http.StatusNotFound)
		return
	}

	tds := keys[0]

	var waiting bool

	if flag := r.URL.Query()["sync"]; flag != nil {
		waiting = true
	} else if flag := r.URL.Query()["async"]; flag != nil {
		waiting = false
	} else {
		http.Error(w, "Require flag sync/async", http.StatusNotFound)
		return
	}

	td, err := time.ParseDuration(tds)
	if err != nil {
		http.Error(w, "Require timeDuration parameter", http.StatusNotFound)
		return
	}

	num := count.get()
	count.inc()
	if waiting {
		wait := make(chan byte, 1)
		wait <- 1
		q.Push(Task{duration: td, number: num, resChan: wait})
		wait <- 1
	} else {
		q.Push(Task{duration: td, number: num, resChan: nil})
	}

	fmt.Fprintf(w, "Succes")
}

func main() {
	count.inc()

	term := make(chan byte)

	go TaskLoop(&q, term)

	http.HandleFunc("/add", add)

	log.Printf("Starting http server ...\n")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatal(err)
	}
}
