package main

import (
        "encoding/json"
        "fmt"
        "github.com/golang/glog"
        "io"
	"b00m.in/finisher/data"
	"b00m.in/finisher/comms"
        "net/http"
        "strconv"
        "strings"
        "time"
)

func handleRoot(w http.ResponseWriter, r *http.Request) {
        tlsoks.Add(1)
        glog.Infof("Secure visitor %s %v \n", r.RemoteAddr, *tlsoks)
        //fmt.Fprintf(w, "Secure but nothing to see here \n")
        switch r.Method{
        case "GET":
                render := Render {Message: "Coming soon ...", Categories: dflt_ctgrs}
                err := tmpl_root.ExecuteTemplate(w, "base", render)
                if err != nil {
                        glog.Errorf("Https %v \n", err)
                        render = Render{Message: "Render error", Categories: dflt_ctgrs}
                        _ = tmpl_adm_err.ExecuteTemplate(w, "base", render)
                        return
                }
                return
        case "POST":
                render := Render {Message: "Coming soon ...", Categories: dflt_ctgrs}
                err := tmpl_root.ExecuteTemplate(w, "base", render)
                if err != nil {
                        glog.Errorf("Https %v \n", err)
                        render = Render{Message: "Render error", Categories: dflt_ctgrs}
                        _ = tmpl_adm_err.ExecuteTemplate(w, "base", render)
                        return
                }
                return
        }
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
                case "pubs":
                        //pubs, err := data.GetPubs(100)
                        pubs, err := data.GetPubDummies(50)
                        if err != nil {
                                str := fmt.Sprintf("Couldn't query results: %s", err)
                                http.Error(w, str, 500)
                                return
                        }
                        glog.Infof("Dummies: %v \n", pubs[0])
                        err = json.NewEncoder(w).Encode(pubs)
                        if err != nil {
                                str := fmt.Sprintf("Couldn't encode results: %s", err)
                                http.Error(w, str, 500)
                                return
                        }
                        return
                case "subs":
                        return
                case "packets":
                        return
                case "summary": // /api/summary/hourly/2021-Jan-01/2021-Jan-31
                        var err error
                        freq := "hourly"
                        from := time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC)
                        to := time.Date(2021, time.January, 31, 0, 0, 0, 0, time.UTC)
                        const shortform = "2006-Jan-02"
                        if len(toks) > 3 {
                                freq = toks[3]
                                if len(toks) > 4 {
                                        from, err = time.Parse(shortform, toks[4])
                                        if err != nil {
                                                glog.Errorf("time.Parse %v \n", err)
                                        }
                                        if len(toks) > 5 {
                                                to, err = time.Parse(shortform, toks[5])
                                                if err != nil {
                                                        glog.Errorf("time.Parse %v \n", err)
                                                }
                                        }
                                }
                        }
                        ss, err := data.GetSummaries(from, to, freq)
                        if err != nil {
                                glog.Errorf("GetSummaries %v \n", err)
                        }
                        err = json.NewEncoder(w).Encode(ss)
                        if err != nil {
                                glog.Errorf("json.Encode %v \n", err)
                                str := fmt.Sprintf("Couldn't encode results: %s", err)
                                http.Error(w, str, 500)
                                return
                        }
                }
        case "POST":
                err := r.ParseForm()
                if err != nil {
                        glog.Errorf("handleApi post  %v", err)
                        w.WriteHeader(http.StatusBadRequest)
                        return
                }
                switch toks[2] {
                case "register":
                        //if sub == "" {
                                glog.Infof("handleApi post api/register w/o sub \n")
                                v := r.Form
                                fn := v.Get("fullname")
                                e := v.Get("email")
                                pw := v.Get("pswd")
                                ph := v.Get("phone")
                                cc := v.Get("confirm")
                                s, err := data.GetSubByEmail(e)
                                if s != nil || err == nil {
                                        glog.Errorf("handleApi post register s: %s err: %v \n", e, err)
                                        w.WriteHeader(http.StatusConflict)
                                        io.WriteString(w, "Email already registered")
                                        return
                                }
                                if pw != cc {
                                        glog.Errorf("handleApi post register %s == %s \n", cc, pw)
                                        w.WriteHeader(http.StatusConflict)
                                        io.WriteString(w, "Password mismatch")
                                        return
                                }
                                // success
                                //glog.Infof("handleApi post login success %s \n", e)
                                //http.SetCookie(w, &http.Cookie{Name: "sub", Value: e, Domain:"b00m.in", Path: "/", MaxAge: 300, HttpOnly: true, Expires: time.Now().Add(time.Second * 120)})
                                glog.Infof("handleapi post register %s \n", e)
                                sha1str := data.Sha1Str(e) // 16 characters of sha1 hash
                                i, err := data.PutSub(&data.Sub{Email:e, Name:fn, Phone:ph, Pswd:pw, Verification: sha1str})
                                if err != nil {
                                        glog.Infof("handleapi post register %s %v \n", e, err)
                                        w.WriteHeader(http.StatusServiceUnavailable)
                                        io.WriteString(w, "Sorry couldn't register")
                                        return
                                }
                                newregs <- comms.Entity{e, fn}
                                w.WriteHeader(http.StatusOK)
                                io.WriteString(w, "Registered " + strconv.Itoa(int(i)))
                                return
                        /*} else {
                                glog.Infof("handlesubs post sub/lin w/ %s \n", sub)
                                render := Render{Message: "Already " + sub, Categories: dflt_ctgrs}
                                _ = tmpl_adm_err.ExecuteTemplate(w, "base", render)
                                return
                        }*/
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

func handlePubs(w http.ResponseWriter, r *http.Request) {
        sub := ""
        var you *data.Sub
        // check cookie
        if cookie, err := r.Cookie("sub"); err != nil {
                glog.Infof("handleSubs no cookie %s \n", r.Method)
                you = &data.Sub{Name: "new"}
        } else {
                if cookie.Value != "gst" {
                        glog.Infof("handleSubs cookie %s %v \n", cookie.Value, cookie.Expires)
                        sub = cookie.Value
                        you, err = data.GetSubByEmail(sub)
                        if err != nil {
                                glog.Errorf("new %v \n", err)
                                rendere := Render{Message: "Cookie error", Categories: dflt_ctgrs, User: "new"}
                                _ = tmpl_adm_err.ExecuteTemplate(w, "base", rendere)
                                return
                        }
                } else {
                        you = &data.Sub{Name: "new"}
                }
        }
        switch r.Method{
        case "GET":
                toks := strings.Split(r.URL.Path, "/")
                glog.Infof("handlePubs %s %v %d \n", r.Method, toks, len(toks))
                //forward /pubs to relevant /pubs/xxx path
                if len(toks) < 3 {
                        toks = append(toks, "top")
                }
                if len(toks) >= 3 && toks[2] == "" {
                        toks[2] = "top"
                }
                switch toks[2] {
                case "all":
                        pbs, err := data.GetPubs(10)
                        if err != nil {
                                glog.Errorf("Https %v \n", err)
                                epbs := make([]*data.Pub, 3)
                                rendere := Render {Message: "Pubs", Pubs: epbs, Categories: dflt_ctgrs, User: you.Name}
                                _ = tmpl_adm_pbs_lst.ExecuteTemplate(w, "base", rendere)
                                return
                        }
                        render := Render {Message: "Pubs", Pubs: pbs, Categories: dflt_ctgrs, User: you.Name}
                        err = tmpl_adm_pbs_lst.ExecuteTemplate(w, "base", render)
                        if err != nil {
                                glog.Errorf("Https %v \n", err)
                                return
                        }
                        return
                case "top":
                        pbs, err := data.GetPubs(10)
                        if err != nil {
                                glog.Errorf("Https %v \n", err)
                                epbs := make([]*data.Pub, 3)
                                rendere := Render {Message: "Pubs", Pubs: epbs, Categories: dflt_ctgrs, User: you.Name}
                                _ = tmpl_adm_pbs_lst.ExecuteTemplate(w, "base", rendere)
                                return
                        }
                        render := Render {Message: "Pubs", Pubs: pbs, Categories: dflt_ctgrs, User: you.Name}
                        err = tmpl_adm_pbs_lst.ExecuteTemplate(w, "base", render)
                        if err != nil {
                                glog.Errorf("Https %v \n", err)
                                return
                        }
                        return
                case "packets":
                        if len(toks) <= 3 {
                                glog.Infof("Nothing to see at %s \n", r.URL.Path)
                                render := Render {Message: "Nothing to see here", Categories: dflt_ctgrs, User: you.Name}
                                _ = tmpl_adm_err.ExecuteTemplate(w, "base", render)
                                return
                        } else {
                                id, err := strconv.ParseInt(toks[3], 10, 64)
                                if err != nil {
                                        glog.Infof("strconv: %v \n", err)
                                        rendere := Render1 {Message: "Nothing to see here", Packet: &data.Packet{}, Categories: dflt_ctgrs, User: you.Name}
                                        _ = tmpl_adm_pck_lst.ExecuteTemplate(w, "base", rendere)
                                        return
                                }
                                pks, err := data.GetLastPackets(id, 100)
                                if err != nil {
                                        glog.Infof("Https %v \n", err)
                                        rendere := Render {Message: "Packets error", Categories: dflt_ctgrs, User: you.Name}
                                        _ = tmpl_adm_err.ExecuteTemplate(w, "base", rendere)
                                        return
                                }
                                render := Rendern {Message: "Packets", Packets: pks, Categories: dflt_ctgrs, User: you.Name}
                                err = tmpl_adm_pck_lst.ExecuteTemplate(w, "base", render)
                                if err != nil {
                                        fmt.Printf("Https %v \n", err)
                                        return
                                }
                                return
                        }
                case "default":
                        glog.Infof("Nothing to see at %s \n", r.URL.Path)
                        render := Render {Message: "Nothing here", Categories: dflt_ctgrs, User: you.Name}
                        _ = tmpl_adm_err.ExecuteTemplate(w, "base", render)
                        return
                }
        case "POST":
        }
}

func handleSubs(w http.ResponseWriter, r *http.Request) {
        var err error
        sub := ""
        var you *data.Sub
        // check cookie
        if cookie, err := r.Cookie("sub"); err != nil {
                glog.Infof("handleSubs no cookie %s \n", r.Method)
                you = &data.Sub{Name: "new"}
        } else {
                if cookie.Value != "gst" {
                        glog.Infof("handleSubs cookie %s %v \n", cookie.Value, cookie.Expires)
                        sub = cookie.Value
                        you, err = data.GetSubByEmail(sub)
                        if err != nil {
                                glog.Errorf("new %v \n", err)
                                rendere := Render{Message: "Cookie error", Categories: dflt_ctgrs, User: "new"}
                                _ = tmpl_adm_err.ExecuteTemplate(w, "base", rendere)
                                return
                        }
                } else {
                        you = &data.Sub{Name: "new"}
                }
        }
        switch r.Method{
        case "GET":
                toks := strings.Split(r.URL.Path, "/")
                glog.Infof("handleSubs %s %v %d \n", r.Method, toks, len(toks))
                //forward /subs to relevant /subs/xxx path
                if len(toks) < 3 {
                        if sub == "" {
                                toks = append(toks, "new")
                        } else {
                                toks = append(toks, "you")
                        }
                }
                if len(toks) >= 3 && toks[2] == "" {
                        if sub == "" {
                                toks[2] = "new"
                        } else {
                                toks[2] = "you"
                        }
                }
                switch toks[2] {
                case "new": // /subs/new
                        if sub == "" {
                                glog.Infof("handlesubs get new no cookie \n")
                                render := Render {Message: "New", Categories: dflt_ctgrs, User: "new"}
                                err := tmpl_adm_sbs_new.ExecuteTemplate(w, "base", render)
                                if err != nil {
                                        glog.Errorf("Https %v \n", err)
                                        rendere := Render{Message: "Render error", Categories: dflt_ctgrs, User: "new"}
                                        _ = tmpl_adm_err.ExecuteTemplate(w, "base", rendere)
                                        return
                                }
                        } else {
                                glog.Infof("handlesubs get new %s \n", sub)
                                render := RenderOne {Message: "You", Sub: you, Categories: dflt_ctgrs, User: you.Name}
                                err := tmpl_adm_sbs_you.ExecuteTemplate(w, "base", render)
                                if err != nil {
                                        glog.Errorf("Https %v \n", err)
                                        rendere := Render{Message: "Render error", Categories: dflt_ctgrs, User: "new"}
                                        _ = tmpl_adm_err.ExecuteTemplate(w, "base", rendere)
                                        return
                                }
                        }
                case "you": // /subs/you
                        if sub == "" {
                                glog.Infof("handlesubs get you no cookie \n")
                                render := Render {Message: "Login", Categories: dflt_ctgrs, User: "new"}
                                err := tmpl_adm_sbs_lin.ExecuteTemplate(w, "base", render)
                                if err != nil {
                                        glog.Errorf("Https %v \n", err)
                                        rendere := Render{Message: "Render error", Categories: dflt_ctgrs, User: "new"}
                                        _ = tmpl_adm_err.ExecuteTemplate(w, "base", rendere)
                                        return
                                }
                                return
                        } else {
                                glog.Infof("handlesubs get you %s \n", sub)
                                render := RenderOne {Message: "You", Sub: you, Categories: dflt_ctgrs, User: you.Name}
                                err := tmpl_adm_sbs_you.ExecuteTemplate(w, "base", render)
                                if err != nil {
                                        glog.Errorf("Https %v \n", err)
                                        rendere := Render{Message: "Render error", Categories: dflt_ctgrs, User: "new"}
                                        _ = tmpl_adm_err.ExecuteTemplate(w, "base", rendere)
                                        return
                                }
                                return
                        }
                case "forgot":
                        render := Render {Message: "Forgot", Categories: dflt_ctgrs, User: "new"}
                        err := tmpl_adm_sbs_new.ExecuteTemplate(w, "base", render)
                        if err != nil {
                                glog.Errorf("Https %v \n", err)
                                rendere := Render{Message: "Render error", Categories: dflt_ctgrs, User: "new"}
                                _ = tmpl_adm_err.ExecuteTemplate(w, "base", rendere)
                                return
                        }
                        return
                case "faults":
                        if sub == "" {
                                glog.Infof("handlesubs get pubs no cookie \n")
                                render := Render {Message: "Login", Categories: dflt_ctgrs, User: "new"}
                                err := tmpl_adm_sbs_lin.ExecuteTemplate(w, "base", render)
                                if err != nil {
                                        glog.Errorf("Https %v \n", err)
                                        rendere := Render{Message: "Render error", Categories: dflt_ctgrs, User: "new"}
                                        _ = tmpl_adm_err.ExecuteTemplate(w, "base", rendere)
                                        return
                                }
                                return
                        } else {
                                glog.Infof("handlesubs get pubs %s : %d \n", sub, you.Id)
                                pubs, err := data.GetPubFaultsForSub(you.Id)
                                if err != nil {
                                        glog.Errorf("Https %v \n", err)
                                        rendere := Render{Message: "No pubs for you", Categories: dflt_ctgrs, User: you.Name}
                                        _ = tmpl_adm_err.ExecuteTemplate(w, "base", rendere)
                                        return
                                }
                                render := Render {Message: "Faults", Pubs: pubs, Categories: dflt_ctgrs, User: you.Name}
                                err = tmpl_adm_pbs_lst.ExecuteTemplate(w, "base", render)
                                if err != nil {
                                        glog.Errorf("Https %v \n", err)
                                        rendere := Render{Message: "Render error", Categories: dflt_ctgrs, User: you.Name}
                                        _ = tmpl_adm_err.ExecuteTemplate(w, "base", rendere)
                                        return
                                }
                                return
                        }
                case "checkout":
                        if len(toks) == 3 { // /subs/checkout
                                if sub == "" {
                                        glog.Infof("handlesubs get pubs no cookie \n")
                                        render := Render {Message: "Login", Categories: dflt_ctgrs, User: "new"}
                                        err := tmpl_adm_sbs_lin.ExecuteTemplate(w, "base", render)
                                        if err != nil {
                                                glog.Errorf("Https %v \n", err)
                                                rendere := Render{Message: "Render error", Categories: dflt_ctgrs, User: "new"}
                                                _ = tmpl_adm_err.ExecuteTemplate(w, "base", rendere)
                                                return
                                        }
                                        return
                                } else {
                                        render := Render {Message: "Checkout", Categories: dflt_ctgrs, User: you.Name}
                                        err := tmpl_adm_sbs_pmt.ExecuteTemplate(w, "base", render)
                                        if err != nil {
                                                glog.Errorf("Https %v \n", err)
                                                rendere := Render{Message: "Render error", Categories: dflt_ctgrs, User: you.Name}
                                                _ = tmpl_adm_err.ExecuteTemplate(w, "base", rendere)
                                                return
                                        }
                                        return
                                }
                        } else { // /subs/checkout/<etc> - 404 
                                glog.Infof("Nothing to see at %s \n", r.URL.Path)
                                render := Render {Message: "Nothing here", Categories: dflt_ctgrs, User: you.Name}
                                _ = tmpl_adm_err.ExecuteTemplate(w, "base", render)
                                return
                        }
                case "contact":
                        render := Render {Message: "Contact",  Categories: dflt_ctgrs, User: you.Name}
                        err = tmpl_adm_sbs_cnt.ExecuteTemplate(w, "base", render)
                        if err != nil {
                                glog.Errorf("privacy %v \n", err)
                        }
                        return
                case "list":
                        if sub == "" {
                                glog.Infof("handlesubs get list no cookie \n")
                                render := Render {Message: "Login", Categories: dflt_ctgrs, User: "new"}
                                err := tmpl_adm_sbs_lin.ExecuteTemplate(w, "base", render)
                                if err != nil {
                                        glog.Errorf("Https %v \n", err)
                                        rendere := Render{Message: "Render error", Categories: dflt_ctgrs, User: "new"}
                                        _ = tmpl_adm_err.ExecuteTemplate(w, "base", rendere)
                                        return
                                }
                                return
                        } else {
                                pk, err := data.GetSubs(10)
                                if err != nil {
                                        glog.Infof("Https %v \n", err)
                                        rendere := Render {Message: "No subs", Categories: dflt_ctgrs}
                                        _ = tmpl_adm_err.ExecuteTemplate(w, "base", rendere)
                                        return
                                }
                                render := Render {Message: "Subs", Subs: pk, Categories: dflt_ctgrs, User: you.Name}
                                err = tmpl_adm_sbs_lst.ExecuteTemplate(w, "base", render)
                                if err != nil {
                                        glog.Errorf("Https %v \n", err)
                                        rendere := Render{Message: "Render error", Categories: dflt_ctgrs, User: you.Name}
                                        _ = tmpl_adm_err.ExecuteTemplate(w, "base", rendere)
                                        return
                                }
                                return
                        }
                case "login":
                        if sub == "" {
                                glog.Infof("handlesubs get login no cookie \n")
                                render := Render {Message: "Login", Categories: dflt_ctgrs, User: "new"}
                                err := tmpl_adm_sbs_lin.ExecuteTemplate(w, "base", render)
                                if err != nil {
                                        glog.Errorf("Https %v \n", err)
                                        rendere := Render{Message: "Render error", Categories: dflt_ctgrs, User: "new"}
                                        _ = tmpl_adm_err.ExecuteTemplate(w, "base", rendere)
                                        return
                                }
                                return
                        } else {
                                glog.Infof("handlesubs get login %s \n", sub)
                                render := RenderOne {Message: "You", Sub: you, Categories: dflt_ctgrs, User: you.Name}
                                err := tmpl_adm_sbs_you.ExecuteTemplate(w, "base", render)
                                if err != nil {
                                        glog.Errorf("Https %v \n", err)
                                        rendere := Render{Message: "Render error", Categories: dflt_ctgrs, User: you.Name}
                                        _ = tmpl_adm_err.ExecuteTemplate(w, "base", rendere)
                                        return
                                }
                                return
                        }
                case "logout":
                        if sub == "" {
                                glog.Infof("handlesubs get logout no cookie \n")
                                render := Render {Message: "Login", Categories: dflt_ctgrs, User: "new"}
                                err := tmpl_adm_sbs_lin.ExecuteTemplate(w, "base", render)
                                if err != nil {
                                        glog.Errorf("Https %v \n", err)
                                        rendere := Render{Message: "Render error", Categories: dflt_ctgrs, User: "new"}
                                        _ = tmpl_adm_err.ExecuteTemplate(w, "base", rendere)
                                        return
                                }
                                return
                        } else {
                                http.SetCookie(w, &http.Cookie{Name: "sub", Value: "gst", Domain:"b00m.in", Path: "/", MaxAge: -600, HttpOnly: true, Expires: time.Now().Add(time.Second * -120)})
                                glog.Infof("handlesubs get logout %s \n", sub)
                                render := Render {Message: "Login", Categories: dflt_ctgrs, User: "new"}
                                err := tmpl_adm_sbs_lin.ExecuteTemplate(w, "base", render)
                                if err != nil {
                                        glog.Errorf("Https %v \n", err)
                                        rendere := Render{Message: "Render error", Categories: dflt_ctgrs}
                                        _ = tmpl_adm_err.ExecuteTemplate(w, "base", rendere)
                                        return
                                }
                                return
                        }
                case "orders": // /subs/orders
                        if sub == "" {
                                w.WriteHeader(http.StatusOK)
                                render := Render {Message: "Please login!", Categories: dflt_ctgrs, User: you.Name}
                                err = tmpl_adm_sbs_lin.ExecuteTemplate(w, "base", render)
                                if err != nil {
                                        glog.Errorf("Https %v \n", err)
                                        return
                                }
                                return

                        } else {
                                carts, err := data.GetCarts(sub)
                                if err != nil {
                                        glog.Errorf("getcarts %v \n", err)
                                        if carts == nil {
                                                carts = make([]*data.Cart, 0)
                                        }
                                }
                                w.WriteHeader(http.StatusOK)
                                render := Render {Message: "Your Orders", Categories: dflt_ctgrs, User: you.Name,  Carts: carts}
                                err = tmpl_gds_ods.ExecuteTemplate(w, "base", render)
                                if err != nil {
                                        glog.Errorf("Https %v \n", err)
                                        rendere := Render{Message: "Render error", Categories: dflt_ctgrs, User: you.Name}
                                        _ = tmpl_adm_err.ExecuteTemplate(w, "base", rendere)
                                        return
                                }
                                return
                        }
                case "pubs":
                        if len(toks) == 3 { // /subs/pubs
                                if sub == "" {
                                        glog.Infof("handlesubs get pubs no cookie \n")
                                        render := Render {Message: "Login", Categories: dflt_ctgrs, User: "new"}
                                        err := tmpl_adm_sbs_lin.ExecuteTemplate(w, "base", render)
                                        if err != nil {
                                                glog.Errorf("Https %v \n", err)
                                                rendere := Render{Message: "Render error", Categories: dflt_ctgrs, User: "new"}
                                                _ = tmpl_adm_err.ExecuteTemplate(w, "base", rendere)
                                                return
                                        }
                                        return
                                } else {
                                        glog.Infof("handlesubs get pubs %s : %d \n", sub, you.Id)
                                        pubs, err := data.GetPubsForSub(you.Id)
                                        if err != nil {
                                                glog.Errorf("Https %v \n", err)
                                                rendere := Render{Message: "No pubs for you", Categories: dflt_ctgrs, User: you.Name}
                                                _ = tmpl_adm_err.ExecuteTemplate(w, "base", rendere)
                                                return
                                        }
                                        render := Render {Message: "Yours", Pubs: pubs, Categories: dflt_ctgrs, User: you.Name}
                                        err = tmpl_adm_pbs_lst.ExecuteTemplate(w, "base", render)
                                        if err != nil {
                                                glog.Errorf("Https %v \n", err)
                                                rendere := Render{Message: "Render error", Categories: dflt_ctgrs, User: you.Name}
                                                _ = tmpl_adm_err.ExecuteTemplate(w, "base", rendere)
                                                return
                                        }
                                        return
                                }
                        } else { // /subs/pubs/<id> 
                                if sub == "" {
                                        glog.Infof("handlesubs /subs/pubs/<id> no cookie \n")
                                        render := Render {Message: "Login", Categories: dflt_ctgrs, User: "new"}
                                        err := tmpl_adm_sbs_lin.ExecuteTemplate(w, "base", render)
                                        if err != nil {
                                                glog.Errorf("Https %v \n", err)
                                                rendere := Render{Message: "Render error", Categories: dflt_ctgrs, User: "new"}
                                                _ = tmpl_adm_err.ExecuteTemplate(w, "base", rendere)
                                                return
                                        }
                                        return
                                } else {
                                        id, err := strconv.ParseInt(toks[3], 10, 64)
                                        if err != nil {
                                                glog.Errorf("Https %v \n", err)
                                                rendere := Render{Message: "Render error", Categories: dflt_ctgrs, User: you.Name}
                                                _ = tmpl_adm_err.ExecuteTemplate(w, "base", rendere)
                                                return
                                        }
                                        pub, err := data.GetPubById(id)
                                        if err != nil {
                                                glog.Errorf("Https %v \n", err)
                                                rendere := Render{Message: "No pubs for you", Categories: dflt_ctgrs, User: you.Name}
                                                _ = tmpl_adm_err.ExecuteTemplate(w, "base", rendere)
                                                return
                                        }
                                        var render RenderOne
                                        devicename,err := data.GetPubDeviceName(pub.Hash)
                                        if err != nil {
                                                glog.Errorf("Https handle subs/pubs/%d : %v \n", id, err)
                                                render = RenderOne{Message: "Unknown", Pub: pub, Categories: dflt_ctgrs, User: you.Name}
                                        } else {
                                                pc, err := data.GetPubConfigByHash(pub.Hash)
                                                if err != nil {
                                                        pc = &data.PubConfig{Kwp: 0, Kwpmake: "unknown", Kwr: 0, Kwrmake: "unknown"}
                                                }
                                                render = RenderOne{Message: devicename, Pub: pub, PubConfig: pc, Categories: dflt_ctgrs, User: you.Name}
                                        }
                                        err = tmpl_adm_pbs_dee.ExecuteTemplate(w, "base", render)
                                        if err != nil {
                                                glog.Errorf("Https %v \n", err)
                                                rendere := Render{Message: "Render error", Categories: dflt_ctgrs, User: you.Name}
                                                _ = tmpl_adm_err.ExecuteTemplate(w, "base", rendere)
                                                return
                                        }
                                        return
                                }
                        }
                case "verify":
                        //check if toks[3] exists
                        if len(toks) != 4 {
                                glog.Infof("Nothing to see at %s \n", r.URL.Path)
                                render := Render {Message: "Nothing here", Categories: dflt_ctgrs, User: you.Name}
                                _ = tmpl_adm_err.ExecuteTemplate(w, "base", render)
                                return
                        }
                        //check if verification exists in table sub
                        s1, err := data.CheckVerification(toks[3])
                        if err != nil {
                                glog.Infof("Nothing to see at %s \n", r.URL.Path)
                                rendere := Render {Message: "Verification doesn't exist", Categories: dflt_ctgrs, User: you.Name}
                                _ = tmpl_adm_err.ExecuteTemplate(w, "base", rendere)
                                return
                        }
                        glog.Infof("Verified %s \n", s1.Email)
                        render := Render {Message: "Thanks for verifying, please login", Categories: dflt_ctgrs, User: you.Name}
                        _ = tmpl_adm_err.ExecuteTemplate(w, "base", render)
                        return
                case "privacy":
                        render := Render {Message: "Privacy Policy", Categories: dflt_ctgrs, User: you.Name}
                        err = tmpl_adm_sbs_prv.ExecuteTemplate(w, "base", render)
                        if err != nil {
                                glog.Errorf("privacy %v \n", err)
                        }
                        return
                case "terms":
                        render := Render {Message: "Terms & Conditions", Categories: dflt_ctgrs, User: you.Name}
                        _ = tmpl_adm_sbs_trm.ExecuteTemplate(w, "base", render)
                        return
                default:
                        glog.Infof("Nothing to see at %s \n", r.URL.Path)
                        render := Render {Message: "Nothing here", Categories: dflt_ctgrs, User: you.Name}
                        _ = tmpl_adm_err.ExecuteTemplate(w, "base", render)
                        return
                }
        case "POST":
                toks := strings.Split(r.URL.Path, "/")
                glog.Infof("handleSubs %s %v \n", r.Method, toks)
                if len(toks) < 3 {
                        glog.Infof("No posting at %s \n", r.URL.Path)
                        render := Render {Message: "No posting here", Categories: dflt_ctgrs, User: you.Name}
                        _ = tmpl_adm_err.ExecuteTemplate(w, "base", render)
                        return
                }
                err := r.ParseForm()
                if err != nil {
                        glog.Errorf("handleSubs post  %v", err)
                        render := Render{Message: "Parse form error", Categories: dflt_ctgrs, User: you.Name}
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
                                        render := Render{Message: "Sub email exists", Categories: dflt_ctgrs, User: you.Name}
                                        _ = tmpl_adm_err.ExecuteTemplate(w, "base", render)
                                        return
                                }
                                sha1str := data.Sha1Str(e) // 16 characters of sha1 hash
                                you = &data.Sub{Email:e, Name:n, Phone:p, Pswd:pw, Verification: sha1str}
                                id, err := data.PutSub(you)
                                if err != nil {
                                        glog.Errorf("handlesubs post putsubs %v \n", err)
                                        rendere := Render{Message: "Couldn't create Sub", Categories: dflt_ctgrs, User: you.Name}
                                        _ = tmpl_adm_err.ExecuteTemplate(w, "base", rendere)
                                        return
                                }
                                // success
                                newregs <- comms.Entity{e, n} // put in channel to send email
                                glog.Infof("handleSubs set cookie %s \n", e)
                                http.SetCookie(w, &http.Cookie{Name: "sub", Value: e, Domain:"b00m.in", Path: "/", MaxAge: 600, HttpOnly: true, Expires: time.Now().Add(time.Second * 600)})
                                glog.Infof("handlesubs post putsubs %s %d \n", n, id)
                                render := Render{Message: "Welcome " + n, Categories: dflt_ctgrs, User: you.Name}
                                _ = tmpl_adm_err.ExecuteTemplate(w, "base", render)
                                return
                        } else {
                                glog.Infof("handlesubs post sub/new with %s \n", sub)
                                render := Render{Message: "Already " + sub, Categories: dflt_ctgrs, User: you.Name}
                                _ = tmpl_adm_err.ExecuteTemplate(w, "base", render)
                                return
                        }
                case "login":
                        if sub == "" {
                                glog.Infof("handlesubs post sub/login w/o sub \n")
                                v := r.Form
                                agr := v.Get("agree")
                                if agr != "1" {
                                        glog.Infof("agreed: %v \n", agr)
                                        render := Render {Message: "Login - please agree to Privacy Policy", Categories: dflt_ctgrs, User: "new"}
                                        err := tmpl_adm_sbs_lin.ExecuteTemplate(w, "base", render)
                                        if err != nil {
                                                glog.Errorf("Https %v \n", err)
                                                rendere := Render{Message: "Render error", Categories: dflt_ctgrs, User: "new"}
                                                _ = tmpl_adm_err.ExecuteTemplate(w, "base", rendere)
                                                return
                                        }
                                        return
                                }
                                e := v.Get("email")
                                pw := v.Get("pswd")
                                s, err := data.GetSubByEmail(e)
                                if err != nil {
                                        glog.Errorf("handlesubs post login %v \n", err)
                                        rendere := Render{Message: "Email doesn't exist", Categories: dflt_ctgrs, User: you.Name}
                                        _ = tmpl_adm_err.ExecuteTemplate(w, "base", rendere)
                                        return
                                }
                                if !data.CheckPswd(e, pw) {
                                        glog.Errorf("handlesubs post putsubs %s == %s \n", s.Pswd, pw)
                                        render := Render{Message: "Password Doesn't match", Categories: dflt_ctgrs, User: you.Name}
                                        _ = tmpl_adm_err.ExecuteTemplate(w, "base", render)
                                        return
                                }
                                // success
                                glog.Infof("handleSubs post login set cookie %s \n", e)
                                http.SetCookie(w, &http.Cookie{Name: "sub", Value: e, Domain:"b00m.in", Path: "/", MaxAge: 600, HttpOnly: true, Expires: time.Now().Add(time.Second * 600)})
                                you = s
                                glog.Infof("handlesubs post login %s %s \n", s.Name, e)
                                render := Render{Message: "Welcome " + s.Name, Categories: dflt_ctgrs, User: you.Name}
                                _ = tmpl_adm_err.ExecuteTemplate(w, "base", render)
                                return
                        } else {
                                glog.Infof("handlesubs post sub/lin w/ %s \n", sub)
                                render := Render{Message: "Already " + sub, Categories: dflt_ctgrs, User: you.Name}
                                _ = tmpl_adm_err.ExecuteTemplate(w, "base", render)
                                return
                        }
                case "pubs":
                        if sub == "" {
                                glog.Infof("handlesubs post sub/pubs w/o cookie \n")
                                render := Render {Message: "Login", Categories: dflt_ctgrs, User: you.Name}
                                err := tmpl_adm_sbs_lin.ExecuteTemplate(w, "base", render)
                                if err != nil {
                                        glog.Errorf("Https %v \n", err)
                                        rendere := Render{Message: "Render error", Categories: dflt_ctgrs, User: you.Name}
                                        _ = tmpl_adm_err.ExecuteTemplate(w, "base", rendere)
                                        return
                                }
                                return
                        } else {
                                glog.Infof("handlesubs post sub/pubs w cookie %s \n", sub)
                                hash := 0
                                if len(toks) < 4 {
                                        glog.Errorf("Https %v \n", err)
                                        render := Render{Message: "No pubid provided", Categories: dflt_ctgrs, User: you.Name}
                                        _ = tmpl_adm_err.ExecuteTemplate(w, "base", render)
                                        return
                                } else {
                                        hash, err = strconv.Atoi(toks[3])
                                        if err != nil {
                                                glog.Errorf("Https %v \n", err)
                                                rendere := Render{Message: "No pubid provided", Categories: dflt_ctgrs, User: you.Name}
                                                _ = tmpl_adm_err.ExecuteTemplate(w, "base", rendere)
                                                return
                                        }
                                }
                                v := r.Form
                                //nickname can only be set during app provisioning
                                //nn := v.Get("nickname")                                
                                kwps := v.Get("kwp")
                                kwpm := v.Get("kwpm")
                                kwrs := v.Get("kwr")
                                kwrm := v.Get("kwrm")
                                kwp, err := strconv.ParseFloat(kwps, 32)
                                if err != nil {
                                        glog.Errorf("Https handlesubs post pubconfig %v \n", err)
                                        rendere := Render{Message: "Couldn't update pubconfig", Categories: dflt_ctgrs, User: you.Name}
                                        _ = tmpl_adm_err.ExecuteTemplate(w, "base", rendere)
                                        return
                                }
                                kwr, err := strconv.ParseFloat(kwrs, 32)
                                if err != nil {
                                        glog.Errorf("Https handlesubs post pubconfig %v \n", err)
                                        rendere := Render{Message: "Couldn't update pubconfig", Categories: dflt_ctgrs, User: you.Name}
                                        _ = tmpl_adm_err.ExecuteTemplate(w, "base", rendere)
                                        return
                                }
                                notify := v.Get("notify")
                                notific := false
                                if notify == "on" {
                                        notific = true
                                }
                                pc := &data.PubConfig{Hash:int64(hash), Kwp: float32(kwp), Kwpmake: kwpm, Kwr: float32(kwr), Kwrmake: kwrm, Notify: notific}
                                if err := data.UpdatePubConfig(pc); err != nil {
                                        glog.Errorf("Https handlesubs post putpubconfig %v \n", err)
                                        rendere := Render{Message: "Couldn't update pub", Categories: dflt_ctgrs, User: you.Name}
                                        _ = tmpl_adm_err.ExecuteTemplate(w, "base", rendere)
                                        return
                                }
                                pub, err := data.GetPubByHash(pc.Hash)
                                if err != nil {
                                        glog.Errorf("Https handlesubs post update pubconfig get pub %v \n", err)
                                        rendere := Render{Message: "Couldn't update pub", Categories: dflt_ctgrs, User: you.Name}
                                        _ = tmpl_adm_err.ExecuteTemplate(w, "base", rendere)
                                        return
                                }
                                render := RenderOne{Message: pc.Nickname, Pub: pub, PubConfig: pc, Categories: dflt_ctgrs, User: you.Name}
                                err = tmpl_adm_pbs_dee.ExecuteTemplate(w, "base", render)
                                if err != nil {
                                        glog.Errorf("Https %v \n", err)
                                        rendere := Render{Message: "Render error", Categories: dflt_ctgrs, User: you.Name}
                                        _ = tmpl_adm_err.ExecuteTemplate(w, "base", rendere)
                                        return
                                }
                        }
                case "you": // /subs/you
                        if sub == "" {
                                render := Render {Message: "Login - please try again", Categories: dflt_ctgrs, User: "new"}
                                err := tmpl_adm_sbs_lin.ExecuteTemplate(w, "base", render)
                                if err != nil {
                                        glog.Errorf("Https %v \n", err)
                                        rendere := Render{Message: "Render error", Categories: dflt_ctgrs, User: "new"}
                                        _ = tmpl_adm_err.ExecuteTemplate(w, "base", rendere)
                                        return
                                }
                                return
                        }
                        v := r.Form
                        oldp := v.Get("old_password")
                        newp := v.Get("password")
                        conf := v.Get("confirm")
                        if conf != newp {
                                render := RenderOne {Message: "You - password not changed - try again", Sub: you, Categories: dflt_ctgrs, User: you.Name}
                                err := tmpl_adm_sbs_you.ExecuteTemplate(w, "base", render)
                                if err != nil {
                                        glog.Errorf("Https %v \n", err)
                                        rendere := Render{Message: "Render error", Categories: dflt_ctgrs, User: "new"}
                                        _ = tmpl_adm_err.ExecuteTemplate(w, "base", rendere)
                                        return
                                }
                                return
                        } else {
                                check := data.CheckPswd(sub, oldp)
                                if !check {
                                        glog.Errorf("/subs/you post check pswd %v \n", oldp)
                                        render := RenderOne {Message: "You - password not changed - try again", Sub: you, Categories: dflt_ctgrs, User: you.Name}
                                        err := tmpl_adm_sbs_you.ExecuteTemplate(w, "base", render)
                                        if err != nil {
                                                glog.Errorf("Https %v \n", err)
                                                rendere := Render{Message: "Render error", Categories: dflt_ctgrs, User: "new"}
                                                _ = tmpl_adm_err.ExecuteTemplate(w, "base", rendere)
                                                return
                                        }
                                        return
                                }
                                err := data.UpdateSub(&data.Sub{Email: sub, Pswd: newp})
                                if err != nil {
                                        render := RenderOne {Message: "You - password not changed - try again", Sub: you, Categories: dflt_ctgrs, User: you.Name}
                                        err := tmpl_adm_sbs_you.ExecuteTemplate(w, "base", render)
                                        if err != nil {
                                                glog.Errorf("Https %v \n", err)
                                                rendere := Render{Message: "Render error", Categories: dflt_ctgrs, User: "new"}
                                                _ = tmpl_adm_err.ExecuteTemplate(w, "base", rendere)
                                                return
                                        }
                                        return
                                }
                                http.SetCookie(w, &http.Cookie{Name: "sub", Value: sub, Domain:"b00m.in", Path: "/", MaxAge: 600, HttpOnly: true, Expires: time.Now().Add(time.Second * 600)})
                                render := RenderOne {Message: "You - password changed", Sub: you, Categories: dflt_ctgrs, User: you.Name}
                                err = tmpl_adm_sbs_you.ExecuteTemplate(w, "base", render)
                                if err != nil {
                                        glog.Errorf("Https %v \n", err)
                                        rendere := Render{Message: "Render error", Categories: dflt_ctgrs, User: "new"}
                                        _ = tmpl_adm_err.ExecuteTemplate(w, "base", rendere)
                                        return
                                }
                                return
                        }
                default:
                        glog.Infof("No posting at %s \n", r.URL.Path)
                        render := Render {Message: "No posting here", Categories: dflt_ctgrs, User: you.Name}
                        _ = tmpl_adm_err.ExecuteTemplate(w, "base", render)
                        return
                }
        }
}

func fileHandler(w http.ResponseWriter, r *http.Request) {
        ruri := r.RequestURI
        //glog.Infof("content-type %s \n", ruri)
        if cssFile.MatchString(ruri) {
        //        glog.Infof("Setting content-type css %s \n", ruri)
                w.Header().Set("Content-Type", "text/css")
        }
        if jsFile.MatchString(ruri) {
        //        glog.Infof("Setting content-type js %s \n", ruri)
                w.Header().Set("Content-Type", "text/javascript")
        }
        staticfileserver.ServeHTTP(w, r)

}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html")
        render := Render {Message: "Gridwatch", Categories: dflt_ctgrs}
	err := tmpl_grw.ExecuteTemplate(w, "base", render)
        if err !=nil {
                fmt.Printf("Error: %v \n", err)
                return
        }
        return
}

func subwayLinesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	w.Write(data.GeoJSON["subway-lines.geojson"])
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
/subs/checkout (get)
/pubs -> /pubs/top (get)
/pubs/search (get/post)
/pubs/hash (get)
/pubs/packets/hash (get)
/admin (get)
/admin/pubs (get)
/admin/subs (get)


*/
