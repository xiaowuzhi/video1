package taskrunner

import (
    "testing"
    "log"
    "errors"
    "time"
)

func TestRunner(t *testing.T) {
    d := func(dc dataChan) error {
        for i := 0; i < 30; i++ {
            dc <- i;
            log.Printf("Dispatcher sent: %v", i)
        }
        return nil
    }

    e := func(dc dataChan) error {
    forloop:
        for {
            select {
            case d := <-dc:
                log.Printf("Executor received: %v", d)
            default:
                break forloop
            }
        }
        return errors.New("Executor")
        //return nil
    }

    runner := NewRunner(30, true, d, e)
    go runner.StartAll()

    time.Sleep(3 * time.Second)
}
