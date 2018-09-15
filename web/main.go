package main

import (
	"net/http"
	// "html/template"
	"github.com/julienschmidt/httprouter"
)

func RegisterHandler() *httprouter.Router {
	router:=httprouter.New()
	router.GET("/", homeHandler)
	router.POST("/", homeHandler)
	router.GET("/userhome", userHomeHandler)
	router.POST("/userhome", userHomeHandler)
	router.POST("/api", apiHandler)
	router.GET("/videos/:vid-id", proxyVideoHandler)
	router.GET("/videos1/:vid-id", proxyVideoHandler)
	router.POST("/upload/:vid-id", proxyUploadHandler)
	router.POST("/upload1/:vid-id", proxyUploadHandler)
	router.ServeFiles("/statics/*filepath", http.Dir("./templates"))
	return router
}


func main() {
	r:=RegisterHandler()
	http.ListenAndServe(":9094", r)
}
