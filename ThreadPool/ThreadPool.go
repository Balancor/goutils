package ThreadPool

import (
	"container/list"
	"sync"
	"sync/atomic"
)

const (
	RESTART  = 1
	SCHEDULE = 2
	PAUSE    = 3
	STOP     = 4
	QUIT     = 5
)

type Job struct {
	Run        func() error
	OnFinished func(bool)
}

type ThreadPool struct {
	MaxWorkers     int32 //最大的线程数目
	WorkerNum      int32
	Mutex          sync.Mutex
	ControlChannel chan int
	JobChannel     chan *Job
	JobQueue       *list.List
}

func NewThreadPool(max_workers int32) *ThreadPool {
	return &ThreadPool{
		MaxWorkers:     max_workers,
		WorkerNum:      0,
		ControlChannel: make(chan int),
		JobChannel:     make(chan *Job),
		JobQueue:       list.New(),
	}
}

func (pool *ThreadPool) dispatch() {
	if pool.JobQueue.Len() == 0 {
		return
	}

	if pool.WorkerNum >= pool.MaxWorkers {
		return
	}

	pool.Mutex.Lock()
	jobElement := pool.JobQueue.Front()
	job := jobElement.Value.(*Job)
	go func() {
		atomic.AddInt32(&pool.WorkerNum, 1)
		err := job.Run()
		if err == nil {
			job.OnFinished(true)
		} else {
			job.OnFinished(false)
		}
		atomic.AddInt32(&pool.WorkerNum, -1)
		pool.ControlChannel <- SCHEDULE
	}()
	pool.JobQueue.Remove(jobElement)
	pool.Mutex.Unlock()

	pool.ControlChannel <- SCHEDULE
}

func (pool *ThreadPool) master() {
	go func() {
		for {
			select {
			case job := <-pool.JobChannel:
				pool.JobQueue.PushBack(job)
				pool.ControlChannel <- SCHEDULE
				break
			case control := <-pool.ControlChannel:
				switch control {
				case RESTART:
					break
				case SCHEDULE:
					pool.dispatch()
					break
				case PAUSE:
					break
				case STOP:
					break
				case QUIT:
					break
				}
				// we have received a signal to stop
				return
			}
		}
	}()
}

func (pool *ThreadPool) AddJob(job *Job) {
	pool.JobChannel <- job
}

func (pool *ThreadPool) Pause() {
	pool.ControlChannel <- PAUSE
}

func (pool *ThreadPool) Stop() {
	pool.ControlChannel <- STOP
}

func (pool *ThreadPool) Quit() {
	pool.ControlChannel <- QUIT
}
