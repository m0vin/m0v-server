package tmpl

import (
        "errors"
        "html/template"
        "net/http"
        //"b00m-in/summarise/data"
        "b00m.in/finisher/data"
        "github.com/golang/glog"
)

type Render struct { //for most purposes
        Message string `json:"message"`
        Error string `json:"message-error"`
        Pubs []*data.Pub `json:"pubs,string"`
        User string `json:"user"`
        Categories []string `json:"categories,string"`
}

type Renderm struct { //for most purposes with categories map
        Message string `json:"message"`
        Error string `json:"message-error"`
        Pubs []*data.Pub `json:"pubs,string"`
        User string `json:"user"`
        Categories map[string]string `json:"categories,map"`
}

type RenderSummary struct { //for summary
        Message string `json:"message"`
        Error string `json:"message-error"`
        Total int64
        Protected int64
        Online int64
        P1 float64
        P2 float64
        P3 float64
        P4 float64
        User string `json:"user"`
        Categories []string `json:"categories,string"`
}

type Rendern struct { //for packets
        Message string `json:"message"`
        Packets []*data.Packet `json:"pubs,string"`
        Id int64
        Start string
        End string
        Freq string
        User string `json:"user"`
        Categories []string `json:"categories,string"`
}

type RenderOne struct { //for most purposes
        Message string `json:"message"`
        Sub *data.Sub `json:"sub,string"`
        Pub *data.Pub `json:"pub,string"`
        PubConfig *data.PubConfig `json:"pubconfig,string"`
        //Categories []Category `json:"categories,string"`
        Categories []string `json:"categories,string"`
        User string
}

var (
        FuncMap = template.FuncMap{
                "div": Divide,
                "divf": Dividef,
                "mul": Multiply,
                "sub": Subtract,
                "subf": Subtractf,
        }
        dflt_ctgrs = []string{"News", "Docs", "Gridwatch", "Leaderboard", "Community", "Github"}
)

func Divide(a int64, b int64) (float64){
        var f64 float64
        f64 = float64(a) / float64(b)
        return f64
}

func Dividef(a float64, b float64) (float64){
        var f64 float64
        f64 = a / b
        return f64
}

func Multiply(a, b float64) (float64){
        f64 := a * b
        return f64
}

func Subtract(a, b int64) (int64){
        i64 := a - b
        return i64
}

func Subtractf(a, b float64) (float64){
        f64 := a - b
        return f64
}

func ExecuteTemplate(w http.ResponseWriter, t *template.Template, e *template.Template, stuff ...interface{}) error {
        if t == nil || e == nil {
                return errors.New("desired templates missing")
        }
        switch len(stuff) {
        case 0: // login template
                render := Render {Message: "Login", Categories: dflt_ctgrs}
                err := t.ExecuteTemplate(w, "base", render)
                if err != nil {
                        glog.Errorf("tmpl.executetemplate %v \n", err)
                        render.Message = "Something went wrong"
                        _ = e.ExecuteTemplate(w, "base", render)
                        return err
                }
                return nil
        case 1: // subs template
                sub, ok := stuff[0].(*data.Sub)
                if ok {
                        render := RenderOne {Message: "Your Details", Categories: dflt_ctgrs, Sub: sub, User: sub.Name}
                        err := t.ExecuteTemplate(w, "base", render)
                        if err != nil {
                                glog.Errorf("tmpl.executetemplate %v \n", err)
                                render.Message = "Something went wrong"
                                _ = e.ExecuteTemplate(w, "base", render)
                                return err
                        }
                }
                return nil
        }
        return nil
}

func GetSummaryRender(you *data.Sub) RenderSummary {
        // get device summary (total, online, protected, power1, .., power4)
        ps, err := data.GetPubStatsForSub(you.Id)
        if err != nil {
                glog.Errorf("GetPubStatsForSub %v \n", err)
        }
        if ps != nil {
                return RenderSummary {Message: "Welcome " + you.Name + " - Summary", Categories: dflt_ctgrs, User: you.Name, Total: ps.T, Online: ps.O, Protected: ps.P, P1: ps.P1, P2: ps.P2, P3: ps.P3, P4: ps.P4}
        }
        return RenderSummary {Message: "Welcome " + you.Name + " - Summary", Categories: dflt_ctgrs, User: you.Name, Total: int64(0), Online: int64(0), Protected: int64(0), P1: float64(0.0), P2: float64(0.0), P3: float64(0.0), P4: float64(0.0)}
}

func GetConfigRender(you *data.Sub, hash int64) RenderOne {
        var renderone RenderOne
        pub, err := data.GetPubByHash(hash)
        if err != nil {
                glog.Errorf("getconfigrender update pubconfig get pub %v \n", err)
        }
        if pub != nil {
                renderone = RenderOne{Message: "tmpl.ConfigRender error", Categories: dflt_ctgrs, User: you.Name}
                //renderone = RenderOne{Message: pc.Nickname, Pub: pub, PubConfig: pc, Categories: dflt_ctgrs, User: you.Name}
        } else {
                renderone = RenderOne{Message: "tmpl.ConfigRender error", Categories: dflt_ctgrs, User: you.Name}
        }
        return renderone
}
