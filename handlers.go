package main

import (
        "fmt"
        "github.com/golang/glog"
        "io"
	"m0v.in/finisher/data"
        "net/http"
        "strconv"
        "strings"
        "time"
)

func handleRoot(w http.ResponseWriter, r *http.Request) {
        tlsoks.Add(1)
        glog.Infof("Secure visitor %s %v \n", r.RemoteAddr, *tlsoks)
        fmt.Fprintf(w, "Secure but nothing to see here \n")
}

func handleAdmin(w http.ResponseWriter, r *http.Request) {
        switch r.Method{
        case "GET":
                toks := strings.Split(r.URL.Path, "/")
                glog.Infof("handleSubs %s %v \n", r.Method, toks)
                if len(toks) <= 3 {
                        glog.Infof("Nothing to see at %s \n", r.URL.Path)
                        render := Render {Message: "Nothing to see here", Categories: dflt_ctgrs}
                        _ = tmpl_adm_err.ExecuteTemplate(w, "base", render)
                        return
                }
                if toks[3] == "new" {
                        render := Render {Message: "Subs", Categories: dflt_ctgrs}
                        err := tmpl_adm_sbs_new.ExecuteTemplate(w, "base", render)
                        if err != nil {
                                glog.Errorf("Https %v \n", err)
                                render = Render{Message: "Render error", Categories: dflt_ctgrs}
                                _ = tmpl_adm_err.ExecuteTemplate(w, "base", render)
                                return
                        }
                }
                pk, err := data.GetSubs(10)
                if err != nil {
                        glog.Infof("Https %v \n", err)
                        render := Render {Message: "No subs", Categories: dflt_ctgrs}
                        _ = tmpl_adm_err.ExecuteTemplate(w, "base", render)
                        return
                }
                render := Render {Message: "Subs", Subs: pk, Categories: dflt_ctgrs}
                err = tmpl_adm_sbs_lst.ExecuteTemplate(w, "base", render)
                if err != nil {
                        glog.Errorf("Https %v \n", err)
                        render = Render{Message: "Render error", Categories: dflt_ctgrs}
                        _ = tmpl_adm_err.ExecuteTemplate(w, "base", render)
                        return
                }
                return
        case "POST":
                err := r.ParseForm()
                if err != nil {
                        glog.Errorf("handleSubs post  %v", err)
                        render := Render{Message: "Parse form error", Categories: dflt_ctgrs}
                        _ = tmpl_adm_err.ExecuteTemplate(w, "base", render)
                        return
                }
                v := r.Form
                n := v.Get("name")
                e := v.Get("email")
                p := v.Get("phone")
                pw := v.Get("pswd")
                if s, err := data.GetSubByEmail(e); s != nil {
                        glog.Errorf("handlesubs post putsubs %v \n", err)
                        render := Render{Message: "Sub email exists", Categories: dflt_ctgrs}
                        _ = tmpl_adm_err.ExecuteTemplate(w, "base", render)
                        return
                }
                _, err = data.PutSub(&data.Sub{Email:e, Name:n, Phone:p, Pswd:pw})
                if err != nil {
                        glog.Errorf("handlesubs post putsubs %v \n", err)
                        render := Render{Message: "Couldn't create Sub", Categories: dflt_ctgrs}
                        _ = tmpl_adm_err.ExecuteTemplate(w, "base", render)
                        return
                }
        }
}

func handleSubs(w http.ResponseWriter, r *http.Request) {
        sub := ""
        var you *data.Sub
        // check cookie
        if cookie, err := r.Cookie("sub"); err != nil {
                glog.Infof("handleSubs no cookie %s \n", r.Method)
        } else {
                if cookie.Value != "gst" {
                        glog.Infof("handleSubs cookie %s %v \n", cookie.Value, cookie.Expires)
                        sub = cookie.Value
                        you, err = data.GetSubByEmail(sub)
                        if err != nil {
                                glog.Errorf("new %v \n", err)
                                render := Render{Message: "Cookie error", Categories: dflt_ctgrs}
                                _ = tmpl_adm_err.ExecuteTemplate(w, "base", render)
                                return
                        }
                }
        }
        switch r.Method{
        case "GET":
                toks := strings.Split(r.URL.Path, "/")
                glog.Infof("handleSubs %s %v \n", r.Method, toks)
                //forward /subs to relevant /subs/xxx path
                if len(toks) < 2 {
                        if sub == "" {
                                toks = append(toks, "new")
                        } else {
                                toks = append(toks, "you")
                        }
                }
                switch toks[2] {
                case "new":
                        if sub == "" {
                                glog.Infof("handlesubs get new no cookie \n")
                                render := Render {Message: "New", Categories: dflt_ctgrs}
                                err := tmpl_adm_sbs_new.ExecuteTemplate(w, "base", render)
                                if err != nil {
                                        glog.Errorf("Https %v \n", err)
                                        render = Render{Message: "Render error", Categories: dflt_ctgrs}
                                        _ = tmpl_adm_err.ExecuteTemplate(w, "base", render)
                                        return
                                }
                        } else {
                                glog.Infof("handlesubs get new %s \n", sub)
                                render := RenderOne {Message: "You", Sub: you, Categories: dflt_ctgrs}
                                err := tmpl_adm_sbs_you.ExecuteTemplate(w, "base", render)
                                if err != nil {
                                        glog.Errorf("Https %v \n", err)
                                        rendere := Render{Message: "Render error", Categories: dflt_ctgrs}
                                        _ = tmpl_adm_err.ExecuteTemplate(w, "base", rendere)
                                        return
                                }
                        }
                case "you":
                        if sub == "" {
                                glog.Infof("handlesubs get you no cookie \n")
                                render := Render {Message: "Login", Categories: dflt_ctgrs}
                                err := tmpl_adm_sbs_lin.ExecuteTemplate(w, "base", render)
                                if err != nil {
                                        glog.Errorf("Https %v \n", err)
                                        render = Render{Message: "Render error", Categories: dflt_ctgrs}
                                        _ = tmpl_adm_err.ExecuteTemplate(w, "base", render)
                                        return
                                }
                                return
                        } else {
                                glog.Infof("handlesubs get you %s \n", sub)
                                render := RenderOne {Message: "You", Sub: you, Categories: dflt_ctgrs}
                                err := tmpl_adm_sbs_you.ExecuteTemplate(w, "base", render)
                                if err != nil {
                                        glog.Errorf("Https %v \n", err)
                                        rendere := Render{Message: "Render error", Categories: dflt_ctgrs}
                                        _ = tmpl_adm_err.ExecuteTemplate(w, "base", rendere)
                                        return
                                }
                                return
                        }
                case "pubs":
                        if sub == "" {
                                glog.Infof("handlesubs get pubs no cookie \n")
                                render := Render {Message: "Login", Categories: dflt_ctgrs}
                                err := tmpl_adm_sbs_lin.ExecuteTemplate(w, "base", render)
                                if err != nil {
                                        glog.Errorf("Https %v \n", err)
                                        render = Render{Message: "Render error", Categories: dflt_ctgrs}
                                        _ = tmpl_adm_err.ExecuteTemplate(w, "base", render)
                                        return
                                }
                                return
                        } else {
                                glog.Infof("handlesubs get pubs %s \n", sub)
                                pubs, err := data.GetPubsForSub(sub)
                                if err != nil {
                                        glog.Errorf("Https %v \n", err)
                                        render := Render{Message: "No pubs for you", Categories: dflt_ctgrs}
                                        _ = tmpl_adm_err.ExecuteTemplate(w, "base", render)
                                        return
                                }
                                render := Render {Message: "You", Pubs: pubs, Categories: dflt_ctgrs}
                                err = tmpl_adm_pbs_lst.ExecuteTemplate(w, "base", render)
                                if err != nil {
                                        glog.Errorf("Https %v \n", err)
                                        rendere := Render{Message: "Render error", Categories: dflt_ctgrs}
                                        _ = tmpl_adm_err.ExecuteTemplate(w, "base", rendere)
                                        return
                                }
                                return
                        }
                case "forgot":
                        render := Render {Message: "Forgot", Categories: dflt_ctgrs}
                        err := tmpl_adm_sbs_new.ExecuteTemplate(w, "base", render)
                        if err != nil {
                                glog.Errorf("Https %v \n", err)
                                render = Render{Message: "Render error", Categories: dflt_ctgrs}
                                _ = tmpl_adm_err.ExecuteTemplate(w, "base", render)
                                return
                        }
                        return
                case "login":
                        if sub == "" {
                                glog.Infof("handlesubs get login no cookie \n")
                                render := Render {Message: "Login", Categories: dflt_ctgrs}
                                err := tmpl_adm_sbs_lin.ExecuteTemplate(w, "base", render)
                                if err != nil {
                                        glog.Errorf("Https %v \n", err)
                                        render = Render{Message: "Render error", Categories: dflt_ctgrs}
                                        _ = tmpl_adm_err.ExecuteTemplate(w, "base", render)
                                        return
                                }
                                return
                        } else {
                                glog.Infof("handlesubs get login %s \n", sub)
                                render := RenderOne {Message: "You", Sub: you, Categories: dflt_ctgrs}
                                err := tmpl_adm_sbs_you.ExecuteTemplate(w, "base", render)
                                if err != nil {
                                        glog.Errorf("Https %v \n", err)
                                        rendere := Render{Message: "Render error", Categories: dflt_ctgrs}
                                        _ = tmpl_adm_err.ExecuteTemplate(w, "base", rendere)
                                        return
                                }
                                return
                        }
                case "logout":
                        if sub == "" {
                                glog.Infof("handlesubs get logout no cookie \n")
                                render := Render {Message: "Login", Categories: dflt_ctgrs}
                                err := tmpl_adm_sbs_lin.ExecuteTemplate(w, "base", render)
                                if err != nil {
                                        glog.Errorf("Https %v \n", err)
                                        render = Render{Message: "Render error", Categories: dflt_ctgrs}
                                        _ = tmpl_adm_err.ExecuteTemplate(w, "base", render)
                                        return
                                }
                                return
                        } else {
                                http.SetCookie(w, &http.Cookie{Name: "sub", Value: "gst", Domain:"b00m.in", Path: "/", MaxAge: -300, HttpOnly: true, Expires: time.Now().Add(time.Second * -120)})
                                glog.Infof("handlesubs get logout %s \n", sub)
                                render := Render {Message: "Login", Categories: dflt_ctgrs}
                                err := tmpl_adm_sbs_lin.ExecuteTemplate(w, "base", render)
                                if err != nil {
                                        glog.Errorf("Https %v \n", err)
                                        render = Render{Message: "Render error", Categories: dflt_ctgrs}
                                        _ = tmpl_adm_err.ExecuteTemplate(w, "base", render)
                                        return
                                }
                                return
                        }
                case "list":
                        pk, err := data.GetSubs(10)
                        if err != nil {
                                glog.Infof("Https %v \n", err)
                                render := Render {Message: "No subs", Categories: dflt_ctgrs}
                                _ = tmpl_adm_err.ExecuteTemplate(w, "base", render)
                                return
                        }
                        render := Render {Message: "Subs", Subs: pk, Categories: dflt_ctgrs}
                        err = tmpl_adm_sbs_lst.ExecuteTemplate(w, "base", render)
                        if err != nil {
                                glog.Errorf("Https %v \n", err)
                                render = Render{Message: "Render error", Categories: dflt_ctgrs}
                                _ = tmpl_adm_err.ExecuteTemplate(w, "base", render)
                                return
                        }
                        return
                default:
                        glog.Infof("Nothing to see at %s \n", r.URL.Path)
                        render := Render {Message: "Nothing here", Categories: dflt_ctgrs}
                        _ = tmpl_adm_err.ExecuteTemplate(w, "base", render)
                        return
                }
        case "POST":
                toks := strings.Split(r.URL.Path, "/")
                glog.Infof("handleSubs %s %v \n", r.Method, toks)
                if len(toks) < 3 {
                        glog.Infof("No posting at %s \n", r.URL.Path)
                        render := Render {Message: "No posting here", Categories: dflt_ctgrs}
                        _ = tmpl_adm_err.ExecuteTemplate(w, "base", render)
                        return
                }
                err := r.ParseForm()
                if err != nil {
                        glog.Errorf("handleSubs post  %v", err)
                        render := Render{Message: "Parse form error", Categories: dflt_ctgrs}
                        _ = tmpl_adm_err.ExecuteTemplate(w, "base", render)
                        return
                }
                switch toks[2] {
                case "new":
                        if sub == "" {
                                v := r.Form
                                n := v.Get("name")
                                e := v.Get("email")
                                p := v.Get("phone")
                                pw := v.Get("pswd")
                                if s, err := data.GetSubByEmail(e); s != nil {
                                        glog.Errorf("handlesubs post putsubs %v \n", err)
                                        render := Render{Message: "Sub email exists", Categories: dflt_ctgrs}
                                        _ = tmpl_adm_err.ExecuteTemplate(w, "base", render)
                                        return
                                }
                                id, err := data.PutSub(&data.Sub{Email:e, Name:n, Phone:p, Pswd:pw})
                                if err != nil {
                                        glog.Errorf("handlesubs post putsubs %v \n", err)
                                        render := Render{Message: "Couldn't create Sub", Categories: dflt_ctgrs}
                                        _ = tmpl_adm_err.ExecuteTemplate(w, "base", render)
                                        return
                                }
                                // success
                                glog.Infof("handleSubs set cookie %s \n", e)
                                http.SetCookie(w, &http.Cookie{Name: "sub", Value: e, Domain:"b00m.in", Path: "/", MaxAge: 300, HttpOnly: true, Expires: time.Now().Add(time.Second * 120)})
                                glog.Infof("handlesubs post putsubs %s %d \n", n, id)
                                render := Render{Message: "Welcome " + n, Categories: dflt_ctgrs}
                                _ = tmpl_adm_err.ExecuteTemplate(w, "base", render)
                                return
                        } else {
                                glog.Infof("handlesubs post sub/new with %s \n", sub)
                                render := Render{Message: "Already " + sub, Categories: dflt_ctgrs}
                                _ = tmpl_adm_err.ExecuteTemplate(w, "base", render)
                                return
                        }
                case "login":
                        if sub == "" {
                                glog.Infof("handlesubs post sub/login w/o sub \n")
                                v := r.Form
                                e := v.Get("email")
                                pw := v.Get("pswd")
                                s, err := data.GetSubByEmail(e)
                                if err != nil {
                                        glog.Errorf("handlesubs post login %v \n", err)
                                        render := Render{Message: "Email doesn't exist", Categories: dflt_ctgrs}
                                        _ = tmpl_adm_err.ExecuteTemplate(w, "base", render)
                                        return
                                }
                                if !data.CheckPswd(e, pw) {
                                        glog.Errorf("handlesubs post putsubs %s == %s \n", s.Pswd, pw)
                                        render := Render{Message: "Password Doesn't match", Categories: dflt_ctgrs}
                                        _ = tmpl_adm_err.ExecuteTemplate(w, "base", render)
                                        return
                                }
                                // success
                                glog.Infof("handleSubs post login set cookie %s \n", e)
                                http.SetCookie(w, &http.Cookie{Name: "sub", Value: e, Domain:"b00m.in", Path: "/", MaxAge: 300, HttpOnly: true, Expires: time.Now().Add(time.Second * 120)})
                                glog.Infof("handlesubs post login %s %s \n", s.Name, e)
                                render := Render{Message: "Welcome " + s.Name, Categories: dflt_ctgrs}
                                _ = tmpl_adm_err.ExecuteTemplate(w, "base", render)
                                return
                        } else {
                                glog.Infof("handlesubs post sub/lin w/ %s \n", sub)
                                render := Render{Message: "Already " + sub, Categories: dflt_ctgrs}
                                _ = tmpl_adm_err.ExecuteTemplate(w, "base", render)
                                return
                        }
                default:
                        glog.Infof("No posting at %s \n", r.URL.Path)
                        render := Render {Message: "No posting here", Categories: dflt_ctgrs}
                        _ = tmpl_adm_err.ExecuteTemplate(w, "base", render)
                        return
                }
        }
}

func handlePubs(w http.ResponseWriter, r *http.Request) {
        switch r.Method{
        case "GET":
                pbs, err := data.GetPubs(10)
                if err != nil {
                        glog.Errorf("Https %v \n", err)
                        epbs := make([]*data.Pub, 3)
                        render := Render {Message: "Pubs", Pubs: epbs, Categories: dflt_ctgrs}
                        _ = tmpl_adm_pbs_lst.ExecuteTemplate(w, "base", render)
                        return
                }
                render := Render {Message: "Pubs", Pubs: pbs, Categories: dflt_ctgrs}
                err = tmpl_adm_pbs_lst.ExecuteTemplate(w, "base", render)
                if err != nil {
                        glog.Errorf("Https %v \n", err)
                        return
                }
                return
        case "POST":
        }
}

func handleAPI(w http.ResponseWriter, r *http.Request) {
        toks := strings.Split(r.URL.Path, "/")
        glog.Infof("handleApi %s %v \n", r.Method, toks)
        if len(toks) < 3 {
                glog.Infof("No get/post at %s \n", r.URL.Path)
                w.WriteHeader(http.StatusNoContent)
                return
        }
        switch r.Method{
        case "GET":
                switch toks[2] {
                case "confo":
                        devicename := toks[3] // "Movprov"
                        ssid := toks[4] // "M0V"
                        confo, err := data.GetLastConfo(devicename, ssid)
                        if err != nil {
                                w.WriteHeader(http.StatusOK)
                                w.Write([]byte("Not found"))
                                return
                        }
                        w.WriteHeader(http.StatusOK)
                        w.Write([]byte(strconv.FormatInt(confo.Hash, 10)))
                        return
                }
        case "POST":
                err := r.ParseForm()
                if err != nil {
                        glog.Errorf("handleApi post  %v", err)
                        w.WriteHeader(http.StatusBadRequest)
                        return
                }
                switch toks[2] {
                case "login":
                        //if sub == "" {
                                glog.Infof("handleApi post api/login w/o sub \n")
                                v := r.Form
                                e := v.Get("email")
                                pw := v.Get("pswd")
                                s, err := data.GetSubByEmail(e)
                                if err != nil {
                                        glog.Errorf("handleApi post login %v \n", err)
                                        w.WriteHeader(http.StatusUnauthorized)
                                        return
                                }
                                if !data.CheckPswd(e, pw) {
                                        glog.Errorf("handleApi post putsubs %s == %s \n", s.Pswd, pw)
                                        w.WriteHeader(http.StatusUnauthorized)
                                        return
                                }
                                // success
                                //glog.Infof("handleApi post login success %s \n", e)
                                //http.SetCookie(w, &http.Cookie{Name: "sub", Value: e, Domain:"b00m.in", Path: "/", MaxAge: 300, HttpOnly: true, Expires: time.Now().Add(time.Second * 120)})
                                glog.Infof("handleapi post login %s %s \n", s.Name, e)
                                w.WriteHeader(http.StatusOK)
                                io.WriteString(w, "Logged in")
                                return
                        /*} else {
                                glog.Infof("handlesubs post sub/lin w/ %s \n", sub)
                                render := Render{Message: "Already " + sub, Categories: dflt_ctgrs}
                                _ = tmpl_adm_err.ExecuteTemplate(w, "base", render)
                                return
                        }*/
                default:
                        glog.Infof("No posting at %s \n", r.URL.Path)
                        w.WriteHeader(http.StatusNotFound)
                        return
                }
        }

}

/*

/root -> (get)
/subs -> /subs/login (get/post)
/subs/new (get/post)
/subs/you (get/post)
/subs/forgot (get/post)
/subs/pubs/yours (get)
/subs/pubs/faults(get)
/subs/pubs/others(get)
/pubs (get)
/pubs/search (get/post)
/pubs/id (get)
/pubs/id/packets (get)
/admin (get)
/admin/pubs (get)
/admin/subs (get)







*/
