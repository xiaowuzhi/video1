package main

import (
    "io"
    "net/http"
    "html/template"
    "log"
    "github.com/julienschmidt/httprouter"
    "fmt"
    "sync"
    "os"
    "time"
)

var bufpool *sync.Pool

func init() {
    bufpool = &sync.Pool{}
    bufpool.New = func() interface{} {
        return make([]byte, 32*1024)
    }
}

func Copy(dst io.Writer, src io.Reader) (written int64, err error) {
    if wt, ok := src.(io.WriterTo); ok {
        return wt.WriteTo(dst)
    }
    if rt, ok := dst.(io.ReaderFrom); ok {
        return rt.ReadFrom(src)
    }

    buf := bufpool.Get().([]byte)
    defer bufpool.Put(buf)

    for {
        nr, er := src.Read(buf)
        if nr > 0 {
            nw, ew := dst.Write(buf[0:nr])
            if nw > 0 {
                written += int64(nw)
            }
            if ew != nil {
                err = ew
                break
            }
            if nr != nw {
                err = io.ErrShortWrite
                break
            }
        }
        if er == io.EOF {
            break
        }
        if er != nil {
            err = er
            break
        }
    }
    return written, err
}

func testPageHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    t, _ := template.ParseFiles("./videos/upload.html")

    t.Execute(w, nil)
}

func streamHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    vid := p.ByName("vid-id")
    // vl1 := "videos/" + vid

    vl := VIDEO_DIR + vid

    video, err := os.Open(vl)
    if err != nil {
        log.Printf("Error when try to open file: %v", err)
        sendErrorResponse(w, http.StatusInternalServerError, "Internal Error")
        return
    }

    w.Header().Set("Content-Type", "video/mp4")
    http.ServeContent(w, r, "", time.Now(), video)

    defer video.Close()

    //privateAccessURL := XiaoPrivateAccessURL(vl1, 0)
    //http.Redirect(w, r, privateAccessURL, 301)
}
func streamHandler1(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    vid := p.ByName("vid-id")
    vl1 := "videos/" + vid
    privateAccessURL := XiaoPrivateAccessURL(vl1, 0)
    http.Redirect(w, r, privateAccessURL, 301)
}

func uploadHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    r.Body = http.MaxBytesReader(w, r.Body, MAX_UPLOAD_SIZE)
    if err := r.ParseMultipartForm(MAX_UPLOAD_SIZE); err != nil {
        fmt.Println(err)
        sendErrorResponse(w, http.StatusBadRequest, "File is too big")
        return
    }


    file, _, err := r.FormFile("file")
    if err != nil {
        log.Printf("Error when try to get file: %v", err)
        sendErrorResponse(w, http.StatusInternalServerError, "Internal Error108")
        return
    }
    defer file.Close()
    fn := p.ByName("vid-id")
    f, err := os.Create(VIDEO_DIR + fn)
    defer f.Close()
    ioc, err := io.Copy(f, file)
    if err != nil {
        log.Printf("Error when try to get file qinui: %v -- %v", ioc, err)
        sendErrorResponse(w, http.StatusInternalServerError, "Internal Error118")
        return
    }



    w.WriteHeader(http.StatusCreated)
    io.WriteString(w, "Uploaded successfully")
}





func uploadHandlerqinui(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    r.Body = http.MaxBytesReader(w, r.Body, MAX_UPLOAD_SIZE_QINUI)
    if err := r.ParseMultipartForm(MAX_UPLOAD_SIZE_QINUI); err != nil {
        fmt.Println(err)
        sendErrorResponse(w, http.StatusBadRequest, "File is too big")
        return
    }


    file, _, err := r.FormFile("file")
    if err != nil {
        log.Printf("Error when try to get file: %v", err)
        sendErrorResponse(w, http.StatusInternalServerError, "Internal Error144")
        return
    }
    defer file.Close()
    fn := p.ByName("vid-id")
    f, err := os.Create(VIDEO_DIR + fn)
    defer f.Close()
    ioc, err := io.Copy(f, file)
    if err != nil {
        log.Printf("Error when try to get file qinui: %v -- %v", ioc, err)
        sendErrorResponse(w, http.StatusInternalServerError, "Internal Error154")
        return
    }

    ossfn := "videos/" + fn
    path :=  VIDEO_DIR + fn
    bn := "avenssi-videos2"

    ret := UploadToOssbig(ossfn, path, bn)

    if !ret {
     sendErrorResponse(w, http.StatusInternalServerError, "Internal Error165")
     return
    }

    os.Remove(path)

    w.WriteHeader(http.StatusCreated)
    io.WriteString(w, "Uploaded successfully")
}
