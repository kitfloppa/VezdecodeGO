package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
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
	Duration time.Duration `json:"timeDuration"`
	Number   int32         `json:"number"`
	resChan  chan byte
}

func Execute(t *Task) {
	time.Sleep(t.Duration)

	if t.resChan != nil {
		<-t.resChan
	}
}

type Queue struct {
	sync.RWMutex
	items   []Task
	sumTime time.Duration
}

func (q *Queue) Push(item Task) {
	q.Lock()
	defer q.Unlock()
	q.items = append(q.items, item)
	q.sumTime += item.Duration
}

func (q *Queue) Pop() Task {
	q.Lock()
	defer q.Unlock()
	if len(q.items) == 0 {
		return Task{}
	}
	item := q.items[0]
	q.items = q.items[1:]
	q.sumTime -= item.Duration
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

		}
		task := q.Pop()
		if task.Number == 0 {
			time.Sleep(time.Millisecond) // waiting 1ms for new tasks
			continue
		}
		select {
		case <-terminator:
			return
		default:

		}
		log.Printf("Starting task №%d...", task.Number)
		Execute(&task)
		select {
		case <-terminator:
			return
		default:

		}
		log.Printf("Task №%d done!", task.Number)

	}
}

func timeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	fmt.Fprintf(w, q.sumTime.String())
}

func schedule(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}
	bz, err := json.Marshal(q.GetAll())

	if err != nil {
		http.Error(w, "Serialize response error.", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, string(bz))
}

func add(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Parse data error.", http.StatusInternalServerError)
		return
	}

	data := r.Form

	keys := data["timeDuration"]
	if keys == nil {
		http.Error(w, "Require timeDuration parameter", http.StatusNotFound)
		return
	}

	tds := keys[0]

	var waiting bool

	if flag := data["sync"]; flag != nil {
		waiting = true
	} else if flag := data["async"]; flag != nil {
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
		q.Push(Task{Duration: td, Number: num, resChan: wait})
		wait <- 1
	} else {
		q.Push(Task{Duration: td, Number: num, resChan: nil})
	}

	fmt.Fprintf(w, "Succes")
}

func init() {
	http.HandleFunc("/add", add)

	http.HandleFunc("/schedule", schedule)

	http.HandleFunc("/time", timeHandler)
}

func main() {
	count = 1

	log.Printf("Starting HTTP server...")
	httpServerExitDone := &sync.WaitGroup{}

	httpServerExitDone.Add(1)

	srv := startHttpServer(httpServerExitDone)

	exit := make(chan os.Signal, 1) // we need to reserve to buffer size 1, so the notifier are not blocked
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)

	<-exit
	log.Printf("Server done, exiting...")
	if err := srv.Shutdown(context.TODO()); err != nil {
		panic(err) // failure/timeout shutting down the server gracefully
	}
	httpServerExitDone.Wait()
}

func startHttpServer(wg *sync.WaitGroup) *http.Server {
	srv := &http.Server{Addr: ":8081"}

	term := make(chan byte, 1) // for termination TaskLoop ability, not implemented

	go TaskLoop(&q, term)

	go func() {
		defer func() {
			term <- 1
			wg.Done()
		}() // let main know we are done cleaning up

		// always returns error. ErrServerClosed on graceful close
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			// unexpected error. port in use?
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()

	// returning reference so caller can call Shutdown()
	return srv
}

func New(terminator <-chan byte, httpServerExitDone *sync.WaitGroup) {

	q = Queue{
		RWMutex: sync.RWMutex{},
		items:   []Task{},
		sumTime: 0,
	}

	count = 1

	log.Printf("Starting HTTP server...")
	srv := startHttpServer(httpServerExitDone)
	<-terminator

	log.Printf("Server done, exiting...")
	if err := srv.Shutdown(context.TODO()); err != nil {
		panic(err) // failure/timeout shutting down the server gracefully
	}
}
