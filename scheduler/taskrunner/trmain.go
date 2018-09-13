package taskrunner

import (
    "time"

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

func Start() {
    // Start video file cleaning
    r := NewRunner(1000, true, VideoClearDispatcher, VideoClearExecutor)
    w := NewWorker(5, r)
    go w.startWorker()
}
