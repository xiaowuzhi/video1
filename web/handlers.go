package main

import (
    "html/template"
    "net/http"
    "log"
    "github.com/julienschmidt/httprouter"
    "encoding/json"
    "io"
    "io/ioutil"
    "net/url"
    "net/http/httputil"
    "time"
    "strconv"
    "avenssi/config"
    "fmt"
)

type HomePage struct {
    Name string
}

type UserPage struct {
    Name string
}

func homeHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    cname, err1 := r.Cookie("username")
    sid, err2 := r.Cookie("session")

    tt := time.Now()
    tts1 := strconv.FormatInt(tt.Unix(), 10)
    w.Header().Set("X-Xiao-T1", tts1)

    if err1 != nil || err2 != nil {
        p := &HomePage{
            Name: "xiaoxiao",
        }
        t, e := template.ParseFiles("./templates/home.html")
        if e != nil {
            log.Printf("Parsing template home.html error: %s", e)
            return
        }

        t.Execute(w, p)
        return
    }

    if len(cname.Value) != 0 && len(sid.Value) != 0 {
        http.Redirect(w, r, "/userhome", http.StatusFound)
        return
    }

}

func userHomeHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    cname, err1 := r.Cookie("username")
    _, err2 := r.Cookie("session")
    if err1 != nil || err2 != nil {
        http.Redirect(w, r, "/", http.StatusFound)
        return
    }

    fname := r.FormValue("username")
    var p *UserPage
    if len(cname.Value) != 0 {
        p = &UserPage{Name: cname.Value}
    } else if len(fname) != 0 {
        p = &UserPage{Name: fname}
    }
    t, e := template.ParseFiles("./templates/userhome.html")
    if e != nil {
        log.Printf("Parsing userhome.html error: %s", e)
        return
    }
    t.Execute(w, p)
}

func apiHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    if r.Method != http.MethodPost {
        re, _ := json.Marshal(ErrorRequestNotRecognized)
        io.WriteString(w, string(re))
        return
    }

    res, _ := ioutil.ReadAll(r.Body)
    apibody := &ApiBody{}
    if err := json.Unmarshal(res, apibody); err != nil {
        re, _ := json.Marshal(ErrorRequestBodyParseFailed)
        io.WriteString(w, string(re))
        return
    }
    request(apibody, w, r)
    defer r.Body.Close()
}

func proxyVideoHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    p := new(Proxy)
    //host := "www.google.com" // WORKS AS EXPECTED
    host := config.GetLbAddr() // GIVES AN ERROR
    u, err := url.Parse(fmt.Sprintf("http://%v:9093/", host))
    if err != nil {
        log.Printf("Error parsing URL")
    }
    p.proxy = httputil.NewSingleHostReverseProxy(u)
    p.proxy.ServeHTTP(w, r)
}
//func proxyUploadHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
//    u, _ := url.Parse("http://" + config.GetLbAddr() + ":9093/")
//    proxy := httputil.NewSingleHostReverseProxy(u)
//    proxy.ServeHTTP(w, r)
//}




func proxyUploadHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    p := new(Proxy)
    //host := "www.google.com" // WORKS AS EXPECTED
    host := config.GetLbAddr() // GIVES AN ERROR
    u, err := url.Parse(fmt.Sprintf("http://%v:9093/", host))
    if err != nil {
        log.Printf("Error parsing URL")
    }
    p.proxy = httputil.NewSingleHostReverseProxy(u)
    p.proxy.ServeHTTP(w, r)
}

type Proxy struct {
    proxy *httputil.ReverseProxy
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    p.proxy.ServeHTTP(w, r)
}
