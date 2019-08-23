package main

import (
        "github.com/golang/glog"
	"m0v.in/finisher/data"
        "net/http"
        "strings"
)

func handleSubs(w http.ResponseWriter, r *http.Request) {
        glog.Infof("handleSubs %s \n", r.Method)
        switch r.Method{
        case "GET":
                toks := strings.Split(r.URL.Path, "/")
                if len(toks) <= 3 {
                        glog.Infof("Nothing to see at %s \n", r.URL.Path)
                        render := Render {Message: "Nothing to see here", Categories: dflt_ctgrs}
                        _ = tmpl_adm_err.ExecuteTemplate(w, "base", render)
                        return
                }
                if toks[2] == "new" {
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
