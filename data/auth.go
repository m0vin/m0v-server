package data

import (
        "context"
	"net/http"
)

// middleware is a definition of  what a middleware is, 
// take in one handlerfunc and wrap it within another handlerfunc
type Middleware func(http.HandlerFunc) http.HandlerFunc

// buildChain builds the middlware chain recursively, functions are first class
func BuildChain(f http.HandlerFunc, m ...Middleware) http.HandlerFunc {
        // if our chain is done, use the original handlerfunc
        if len(m) == 0 {
                return f
        }
        // otherwise nest the handlerfuncs
        return m[0](BuildChain(f, m[1:cap(m)]...))
}


// AuthMiddleware - takes in a http.HandlerFunc, and returns a http.HandlerFunc
var AuthMiddleware = func(f http.HandlerFunc) http.HandlerFunc {
        // one time scope setup area for middleware
        return func(w http.ResponseWriter, r *http.Request) {
                // ... pre handler functionality
                sub := ""
                var you *Sub
                // check cookie
                if cookie, err := r.Cookie("sub"); err != nil {
                        //glog.Infof("handleNew no cookie %s \n", r.Method)
                        you = &Sub{Name: "new"}
                } else {
                        if cookie.Value != "gst" {
                                //glog.Infof("handleSubs cookie %s %v \n", cookie.Value, cookie.Expires)
                                sub = cookie.Value
                                you, err = GetSubByEmail(sub)
                                if err != nil {
                                        //glog.Errorf("new %v \n", err)
                                        //rendere := Render{Message: "Cookie error", Categories: dflt_ctgrs, User: "new"}
                                        //_ = tmpl_adm_err.ExecuteTemplate(w, "base", rendere)
                                        return
                                }
                        } else {
                                you = &Sub{Name: "new"}
                        }
                }
                r = r.Clone(context.WithValue(context.Background(), "user", you))
                f(w, r)
                // ... post handler functionality
        }
}
