package taskrunner

import (
    "errors"
    "log"
    "sync"
    "video1/scheduler/dbops"
    "video1/scheduler/ossops"
    "time"
)

func deleteVideo(vid string) error {
    //err := os.Remove(VIDEO_PATH + vid)
    //
    //if err != nil && !os.IsNotExist(err) {
    //    log.Printf("Deleting video error: %v", err)
    //    return err
    //}

    ossfn := "videos/" + vid
    bn := "avenssi-videos2"
    ok := ossops.DeleteObject(ossfn, bn)

    if !ok {
        log.Printf("Deleting video error, oss operation faile")
        return errors.New("Deleting video error")
    }

    return nil
}

func VideoClearDispatcher(dc dataChan) error {
    res, err := dbops.ReadVideoDeletionRecord(5000)
    if err != nil {
        log.Printf("Video clear dispatcher error: %v", err)
        return err
    }

    if len(res) == 0 {
        return errors.New("All tasks finished")
    }

    for _, id := range res {
        dc <- id
    }

    return nil
}

func VideoClearExecutor(dc dataChan) error {
    errMap := &sync.Map{}
    var err error

forloop:
    for {
        select {
        case vid := <-dc:
            go func(id interface{}) {
                if err := deleteVideo(id.(string)); err != nil {
                    errMap.Store(id, err)
                    return
                }
                if err := dbops.DelVideoDeletionRecord(id.(string)); err != nil {
                    errMap.Store(id, err)
                    return
                }
            }(vid)
        default:
            time.Sleep(1 * time.Second)
            break forloop
        }
    }

    errMap.Range(func(k, v interface{}) bool {
        err = v.(error)

        if err.Error() == "Deleting video error" {
            log.Printf("errMap Deleting video error id=%v",k)
            dbops.DelVideoDeletionRecord(k.(string))
        }else {
            return false
        }

        return true
    })
    return err
}
