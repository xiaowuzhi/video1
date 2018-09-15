package main

import (
    "net/http"
    "github.com/julienschmidt/httprouter"
)

type middleWareHandler struct {
    r *httprouter.Router
    l *ConnLimiter
}

func NewMiddleWareHandler(r *httprouter.Router, cc int) http.Handler {
    m := middleWareHandler{}
    m.r = r
    m.l = NewConnLimiter(cc)
    return m
}

func RegisterHandlers() *httprouter.Router {
    router := httprouter.New()

    router.GET("/videos/:vid-id", streamHandler)

    router.GET("/videos1/:vid-id", streamHandler1)

    router.POST("/upload/:vid-id", uploadHandler)
    router.POST("/upload1/:vid-id", uploadHandlerqinui)
    router.POST("/upload_qinui/:vid-id", uploadHandlerqinui)

    router.GET("/testpage", testPageHandler)

    return router
}

func (m middleWareHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    if !m.l.GetConn() {
        sendErrorResponse(w, http.StatusTooManyRequests, "Too many requests")
        return
    }

    m.r.ServeHTTP(w, r)
    defer m.l.ReleaseConn()
}

func main() {
    r := RegisterHandlers()
    mh := NewMiddleWareHandler(r, 3000)
    http.ListenAndServe(":9093", mh)
}
