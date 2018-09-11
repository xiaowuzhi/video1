package main

import (
    "net/http"
    "github.com/julienschmidt/httprouter"
    "log"
    "video1/api/session"
)

type middleWareHandler struct {
    r *httprouter.Router
}

func NewMiddleWareHandler(r *httprouter.Router) http.Handler {
    m := middleWareHandler{}
    m.r = r
    return m
}

func (m middleWareHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    //check session
    validateUserSession(r)
    m.r.ServeHTTP(w, r)
}

func RegisterHandlers() *httprouter.Router {
    log.Printf("preparing to post request\n")
    router := httprouter.New()
    router.POST("/user", CreateUser)
    router.POST("/user/:username", Login)
    router.GET("/user/:username", GetUserInfo)
    router.POST("/user/:username/videos", AddNewVideo)
    router.GET("/user/:username/videos", ListAllVideos)
    router.DELETE("/user/:username/videos/:vid-id", DeleteVideo)
    router.POST("/videos/:vid-id/comments", PostComment)
    router.GET("/videos/:vid-id/comments", ShowComments)
    return router
}

func Prepare() {
    session.LoadSessionsFromDB()
}

func main() {
    Prepare()
    r := RegisterHandlers()
    mh := NewMiddleWareHandler(r)
    http.ListenAndServe(":9091", mh)
}

//handler->validation{1.request, 2.user}->business logic->reponse.
//1.data model
//2.error handling
//session 13701370121
