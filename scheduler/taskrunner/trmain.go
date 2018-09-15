package taskrunner

import (
    "time"

    "fmt"
)

type Worker struct {
    ticker *time.Ticker
    runner *Runner
}

func NewWorker(interval time.Duration, r *Runner) *Worker {
    return &Worker{
        ticker: time.NewTicker(interval * time.Second),
        runner: r,
    }
}

func (w *Worker) startWorker() {
    for {
        select {
        case <-w.ticker.C:
            go w.runner.StartAll()
        }
    }
}

func (w *Worker) startWorker_local() {
    for {
        select {
        case <-w.ticker.C:
            fmt.Println(555555)
            go w.runner.StartAll()
        }
    }
}

func Start() {
    // Start video file cleaning
    r := NewRunner(1000, true, VideoClearDispatcher, VideoClearExecutor)
    w := NewWorker(10, r)

    //r_local := NewRunner(1000, true, VideoClearDispatcher_Local, VideoClearExecutor_Local)
    //w_local := NewWorker(12, r_local)

    go w.startWorker()
    //go w_local.startWorker_local()
}
