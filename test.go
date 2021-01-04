package main

import (
	"errors"
	"fmt"
	"log"
	"sync/atomic"
	"time"
)

type Task struct {
	Handler func(v ...interface{})
	Params  []interface{}
}

var ErrInvalidPoolCap = errors.New("invalid pool cap")

const (
	RUNNING = 1
	STOPPED = 0
)

func NewPool(capacity uint64) (*Pool, error) {
	if capacity <= 0 {
		return nil, ErrInvalidPoolCap
	}
	return &Pool{
		capacity: capacity,
		state:    RUNNING,
		taskC:    make(chan *Task, capacity),
		closeC:   make(chan bool),
	}, nil
}

func (p *Pool) incRunning() {
	atomic.AddUint64(&p.runningWorkers, 1)

}

func (p *Pool) decRunning() {
	atomic.AddUint64(&p.runningWorkers, ^uint64(0))

}

func (p *Pool) GetRunningWorkers() uint64 {
	return atomic.LoadUint64(&p.runningWorkers)

}

func (p *Pool) GetCap() uint64 {
	return atomic.LoadUint64(&p.capacity)
}

type Pool struct {
	capacity       uint64
	runningWorkers uint64
	state          int64
	taskC          chan *Task
	closeC         chan bool
	PanicHandler   func(interface{})
}

var ErrPoolAlreadyClosed = errors.New("pool already closed")

func (p *Pool) Put(task *Task) error {

	if p.state == STOPPED { // 如果任务池处于关闭状态, 再 put 任务会返回 ErrPoolAlreadyClosed 错误
		return ErrPoolAlreadyClosed
	}

	if p.GetRunningWorkers() < p.GetCap() {
		p.run()
	}

	p.taskC <- task

	return nil
}

func (p *Pool) run() {
	p.incRunning()

	go func() {
		defer func() {
			p.decRunning()
			if r := recover(); r != nil { // 恢复 panic
				if p.PanicHandler != nil { // 如果设置了 PanicHandler, 调用
					p.PanicHandler(r)
				} else { // 默认处理
					log.Printf("Worker panic: %s\n", r)
				}
			}
		}()

		for {
			select {
			case task, ok := <-p.taskC:
				if !ok {
					return
				}
				task.Handler(task.Params...)
			case <-p.closeC:
				return
			}
		}
	}()
}

func main() {
	// 创建任务池.
	pool, err := NewPool(10)
	if err != nil {
		panic(err)
	}

	for i := 0; i < 10; i++ {
		// 任务放入池中
		pool.Put(&Task{
			Handler: func(v ...interface{}) {
				fmt.Println(v)
			},
			Params: []interface{}{i},
		})
	}

	time.Sleep(1e9) // 等待执行
}
