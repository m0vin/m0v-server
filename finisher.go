package main

import (
        //"context"
	"crypto"
	//"crypto/rand"
	"crypto/x509"
	"crypto/ecdsa"
	"crypto/rsa"
        "expvar"
        "encoding/pem"
        "errors"
	"flag"
	"fmt"
        "github.com/golang/glog"
        //"golang.org/x/crypto/acme"
        //"golang.org/x/crypto/acme/autocert"
	"html/template"
        "io/ioutil"
	"b00m.in/finisher/data"
	"b00m.in/finisher/comms"
	"b00m.in/finisher/subs"
        "net/http"
        "os"
        //"reflect"
        "regexp"
	//"rsc.io/quote"
        "strconv"
        "strings"
        "time"
)

type Render struct { //for most purposes
        Message string `json:"message"`
        Subs []*data.Sub `json:"subs,string"`
        Pubs []*data.Pub `json:"pubs,string"`
        //Categories []Category `json:"categories,string"`
        Categories subs.Category `json:"categories,string"`
        User string
        Carts []*data.Cart
}

type RenderOne struct { //for most purposes
        Message string `json:"message"`
        Sub *data.Sub `json:"sub,string"`
        Pub *data.Pub `json:"pub,string"`
        PubConfig *data.PubConfig `json:"pubconfig,string"`
        //Categories []Category `json:"categories,string"`
        Categories subs.Category `json:"categories,string"`
        User string
}

type Render1 struct { //for a packet
        Message string `json:"message"`
        //Sub sub `json:"sub"`
        Packet *data.Packet `json:"packet,string"`
        //Categories []Category `json:"categories,string"`
        Categories subs.Category `json:"categories,string"`
        User string
}

type Rendern struct { //for packets
        Message string `json:"message"`
        //Sub sub `json:"sub"`
        Packets []*data.Packet `json:"packet,string"`
        //Categories []Category `json:"categories,string"`
        Categories subs.Category `json:"categories,string"`
        User string
}

/*type Category struct {
        Name string `json:"name"`
}

func (c *Category) ToLower() string {
        return strings.ToLower(c.Name)
}*/

var (
        cfg = flag.String("c", "config/b00m.config", "path to config file")
        oks *expvar.Int
        tlsoks *expvar.Int
        httpRto int
        httpWto int
        httpPort int
        httpsPort int
        redirectHttp bool
        /*keyName = "./pv.b00m.in/b00m-key.pem" // "/home/sridhar/prod/b00m_tls/b00m-key.pem"
        keyGen = false
        privKey *rsa.PrivateKey
        man *autocert.Manager
        c *acme.Client
        prod1 = "https://acme-v01.api.letsencrypt.org/directory"
        prod2 = "https://acme-v02.api.letsencrypt.org/directory"
        reg1 = "https://acme-v01.api.letsencrypt.org/acme/reg"
        stag2 = "https://acme-staging-v02.api.letsencrypt.org/directory"
        stag1 = "https://acme-staging.api.letsencrypt.org/directory"*/
        cssFile = regexp.MustCompile("\\.css$")
        jsFile = regexp.MustCompile("\\.js$")
        staticfileserver = http.FileServer(http.Dir("static"))
        newregs chan comms.Entity
        // templates
        funcMap = template.FuncMap{
                "mult": Multiply,
        }
        tmpl_index = template.Must(template.ParseFiles("static/index.html"))
	tmpl_root = template.Must(template.ParseFiles("templates/root"))
	tmpl_adm_err = template.Must(template.ParseFiles("templates/adm/error", "templates/adm/cmn/body", "templates/adm/cmn/right", "templates/adm/cmn/center_errs", "templates/adm/cmn/search", "templates/cmn/base", "templates/cmn/head_2back", "templates/cmn/menu", "templates/cmn/footer"))
	tmpl_adm_pbs_lst = template.Must(template.ParseFiles("templates/adm/pubs_list", "templates/adm/cmn/body", "templates/adm/cmn/right", "templates/adm/cmn/center_pubs", "templates/adm/cmn/search", "templates/cmn/base", "templates/cmn/head", "templates/cmn/menu", "templates/cmn/footer"))
	tmpl_adm_sbs_lst = template.Must(template.ParseFiles("templates/adm/subs_list", "templates/adm/cmn/body", "templates/adm/cmn/right", "templates/adm/cmn/center_subs", "templates/adm/cmn/search", "templates/cmn/base", "templates/cmn/head_2back", "templates/cmn/menu", "templates/cmn/footer"))
	tmpl_adm_sbs_new = template.Must(template.ParseFiles("templates/adm/subs_new", "templates/adm/cmn/body", "templates/adm/cmn/right", "templates/adm/cmn/center_subs", "templates/adm/cmn/search", "templates/cmn/base", "templates/cmn/head_2back", "templates/cmn/menu", "templates/cmn/footer"))
	tmpl_adm_sbs_you = template.Must(template.ParseFiles("templates/adm/subs_you", "templates/adm/cmn/body", "templates/adm/cmn/right", "templates/adm/cmn/center_subs", "templates/adm/cmn/search", "templates/cmn/base", "templates/cmn/head_2back", "templates/cmn/menu", "templates/cmn/footer"))
	tmpl_adm_sbs_pmt = template.Must(template.ParseFiles("templates/adm/pmt/paypal", "templates/adm/cmn/body", "templates/adm/cmn/right", "templates/adm/cmn/center_subs", "templates/adm/cmn/search", "templates/cmn/base", "templates/cmn/head_2back", "templates/cmn/menu", "templates/cmn/footer"))
	tmpl_adm_pbs_dee = template.Must(template.ParseFiles("templates/adm/pub_deet", "templates/adm/cmn/body", "templates/adm/cmn/right", "templates/adm/cmn/center_pubs", "templates/adm/cmn/search", "templates/cmn/base", "templates/cmn/head_2back", "templates/cmn/menu", "templates/cmn/footer"))
	tmpl_adm_sbs_lin = template.Must(template.ParseFiles("templates/adm/subs_login", "templates/adm/cmn/body", "templates/adm/cmn/right", "templates/adm/cmn/center_subs", "templates/adm/cmn/search", "templates/cmn/base", "templates/cmn/head_2back", "templates/cmn/menu", "templates/cmn/footer"))
	tmpl_adm_pck_lst = template.Must(template.ParseFiles("templates/adm/pcks_list", "templates/adm/cmn/body", "templates/adm/cmn/right", "templates/adm/cmn/center", "templates/adm/cmn/search", "templates/cmn/base", "templates/cmn/head_2back", "templates/cmn/menu", "templates/cmn/footer"))
	tmpl_adm_pck_one = template.Must(template.ParseFiles("templates/adm/pcks_one", "templates/adm/cmn/body", "templates/adm/cmn/right", "templates/adm/cmn/center", "templates/adm/cmn/search", "templates/cmn/base", "templates/cmn/head_2back", "templates/cmn/menu", "templates/cmn/footer"))
	tmpl_adm_sbs_prv = template.Must(template.ParseFiles("templates/adm/privacy", "templates/adm/cmn/body-vha", "templates/adm/cmn/right", "templates/adm/cmn/center_subs", "templates/adm/cmn/search", "templates/cmn/base", "templates/cmn/head_2back", "templates/cmn/menu", "templates/cmn/footer"))
	tmpl_adm_sbs_cnt = template.Must(template.ParseFiles("templates/adm/contact", "templates/adm/cmn/body", "templates/adm/cmn/right", "templates/adm/cmn/center_subs", "templates/adm/cmn/search", "templates/cmn/base", "templates/cmn/head_2back", "templates/cmn/menu", "templates/cmn/footer"))
	tmpl_adm_sbs_trm = template.Must(template.ParseFiles("templates/adm/terms", "templates/adm/cmn/body", "templates/adm/cmn/right", "templates/adm/cmn/center_subs", "templates/adm/cmn/search", "templates/cmn/base", "templates/cmn/head_2back", "templates/cmn/menu", "templates/cmn/footer"))
        tmpl_gds_ods = template.Must(template.New("tmpl_gds_orders").Funcs(funcMap).ParseFiles("templates/adm/orders", "templates/adm/cmn/body", "templates/adm/cmn/right", "templates/adm/cmn/center_orders", "templates/adm/cmn/search", "templates/cmn/base", "templates/cmn/head", "templates/cmn/menu", "templates/cmn/footer"))
        //dflt_ctgrs = []Category{Category{Name: "Docs", }, Category{Name: "News", }, Category{Name: "Gridwatch", }, Category{Name: "Leaderboard"}, Category{Name: "Community"}, Category{Name: "Github"}}
        dflt_ctgrs = subs.Category{}
	tmpl_grw = template.Must(template.ParseFiles("templates/adm/cmn/body1", "templates/adm/cmn/right", "templates/adm/cmn/center_grw", "templates/adm/cmn/search", "templates/cmn/base", "templates/cmn/head_2back", "templates/cmn/menu", "templates/cmn/footer"))
)

func Multiply(a float64, b int) (float64){
        f64 := a * float64(b)
        return f64
}

func init() {
        /*data, err := ioutil.ReadFile(keyName)
        if err != nil {
                fmt.Printf("%s %v \n", "Generate rsa key", err)
                keyGen = true
        }
        if keyGen {
                privKey, err = rsa.GenerateKey(rand.Reader, 2048)
                if err != nil {
                        fmt.Printf("%s \n", "Generating rsa key")
                        os.Exit(1) // no other option but to exit
                }
        } else { // using key from file
                // private key
                priv, _ := pem.Decode(data) // ignore public key
                if priv == nil || !strings.Contains(priv.Type, "PRIVATE") {
                        fmt.Printf("%s \n", "Nil rsa key")
                        os.Exit(1) // no other option but to exit
                        //if key == nil {
                        //        key, err = rsa.GenerateKey(rand.Reader, 2048)
                        //        if err != nil {
                        //                fmt.Printf("%s \n", "Generating rsa key")
                        //                os.Exit(1) // no other option but to exit
                        //        }
                        //}
                }
                signer, err := parsePrivateKey(priv.Bytes)
                if err != nil {
                        fmt.Printf("%s \n", "Parsing rsa key")
                        os.Exit(1) // no other option but to exit
                }
                privKey = signer.(*rsa.PrivateKey)
        }
        c = &acme.Client{DirectoryURL: prod2, Key: privKey}
        man = &autocert.Manager{
                Client: c,
                Email: "rcs@m0v.in",
                Prompt: autocert.AcceptTOS,
                Cache: autocert.DirCache("./pv.b00m.in"), // ("/b00m.in/b00m_tls"), //"/home/sridhar/prod/b00m_tls"),
                //HostPolicy: autocert.HostWhitelist("m0v.in", "www.m0v.in"),
        }*/
        oks = expvar.NewInt("oks")
        tlsoks = expvar.NewInt("tlsoks")
        flag.IntVar(&httpRto, "wto", 10, "Read timeout")
        flag.IntVar(&httpWto, "rto", 10, "Write timeout")
        //flag.IntVar(&httpPort, "http_port", 80, "Http server port")
        //flag.IntVar(&httpsPort, "https_port", 443, "Http server port")
        //flag.BoolVar(&redirectHttp, "redirect_http", false, "Redirect http to https")
}

func main() {
	glog.Infof("%s \n", "subs server starting...")
        config, err := subs.ConfigureConfig(os.Args[1:])
        if err != nil {
                glog.Errorf("configureconfig %v \n", err)
        }
	flag.Parse()
        s, err := subs.NewServer(config)
        if err != nil {
                glog.Errorf("newserver %v \n", err)
        }
        s.Start()
        dflt_ctgrs = s.Cfg.Categories
        p, err := strconv.Atoi(s.Cfg.HTTPPort)
        if err != nil {
                glog.Errorf("strconv httpport exiting %v \n", err)
                return
        }
        httpPort = p
        redirectHttp = false
        b := make(chan bool, 1)
        newregs = make(chan comms.Entity, 3)
	//data.CacheGeoJSON()
	//data.LoadStations()
	//data.LoadPubdeetsLocal()
	//if err := data.LoadPubdeets(); err != nil {
	if err := data.LoadDummyPubdeets(); err != nil {
                glog.Infof("data.LoadPubdeets %v \n", err)
        }
        go loadPubdeets()
        go startHttp()
        if httpsPort > 0 {
                go startHttps()
        }
	/*id, dup, err := data.PutToStore("1234", 5678)
	if dup {
		fmt.Printf("a request is already active for controller %s", "1234")
	}
	if err != nil {
		fmt.Printf("Failed to put to store:%s\n", err)
	}
	fmt.Println(id)
	if err = data.UpdateRequest(1, "ACCEPTED"); err != nil {

		fmt.Printf("Failed to update store:%s\n", err)
	}*/

        // Discover queries letsencrypt for service urls
        /*dir, err := c.Discover(context.Background())
        if err != nil {
                fmt.Printf("discover %v \n", err)
        }
        fmt.Printf("Dir %v \n", dir)
        */


        // GetKey loads key from file into client
        /*k, err := GetKey("/home/sridhar/dev/m0v/gcp/acme_account+key")
        if err != nil || k == nil {
                fmt.Printf("account %v \n", err)
                return
        }
        c.Key = k

        // GetReg requires the complete uri of the account including account id - needs the key loaded in client
        acc, err := c.GetReg(context.Background(), reg1 + "/52032422")
        //acc, err := c.Register(context.Background(), &acme.Account{AgreeTerms: "intentionally_failing"}, autocert.AcceptTOS)
        if err != nil {
                fmt.Printf("account %v \n", err)
        }*/

        // Can use existing client key to read the existing registration's account number from the response header
        /*var acc *acme.Account
        if acc, err = c.Register(context.Background(), &acme.Account{AgreedTerms: "intentionally_failing"}, func(tos string) bool {
                return false
        }); err != nil {
                switch aErr := err.(type) {
                case *acme.Error:
                        if aErr.StatusCode == 409 {
                                fmt.Printf("Found existing registration: %s\n", aErr.Header.Get("Boulder-Requester"))
                        }
                        fmt.Printf("Couldn't find key, probably not registered: (%v)\n", aErr.Detail)
                default:
                        fmt.Printf("Failed to create registration: %v (%v)\n", err, reflect.TypeOf(err))
                }
        }*/
        //fmt.Printf("Account %v \n", acc)
        <-b
}

func startHttp() {
        mux := http.NewServeMux()
        //mux.Handle("/debug/vars", expvar.Handler())
        /*mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                oks.Add(1)
                fmt.Printf("Visitor %v \n", *oks)
                fmt.Fprintf(w, "Nothing to see here \n")
        }))*/
        if redirectHttp {
                mux.Handle("/", http.HandlerFunc(RedirectHttp))
        } else {
                //mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
                mux.Handle("/static/", http.StripPrefix("/static/", http.HandlerFunc(fileHandler)))
                mux.Handle("/api/", http.HandlerFunc(handleAPI))
                mux.Handle("/subs/", http.HandlerFunc(handleSubs))
                mux.Handle("/pubs/", http.HandlerFunc(handlePubs))
                mux.Handle("/data/subway-stations", http.HandlerFunc(data.SubwayStationsHandler))
                mux.Handle("/data/subway-lines", http.HandlerFunc(subwayLinesHandler))
                mux.Handle("/data/publishers", http.HandlerFunc(data.PubDeetsHandler))
                mux.Handle("/gridwatch/", http.HandlerFunc(indexHandler))
                mux.Handle("/", http.HandlerFunc(handleRoot))
        }
        hs := http.Server{
                ReadTimeout: time.Duration(httpRto) * time.Second,
                WriteTimeout: time.Duration(httpWto) * time.Second,
                Addr: fmt.Sprintf(":%d", httpPort),
                Handler: mux,
        }
        err := hs.ListenAndServe()
        if err != nil {
                fmt.Printf("Oops: %v \n", err)
        }
}

func startHttps() {
        mux := http.NewServeMux()
        mux.Handle("/debug/vars", expvar.Handler())
        mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
        mux.Handle("/admin/packets/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                toks := strings.Split(r.URL.Path, "/")
                if len(toks) <= 3 {
                        glog.Infof("Nothing to see at %s \n", r.URL.Path)
                        epbs := make([]*data.Pub, 3)
                        render := Render {Message: "Nothing to see here", Pubs: epbs, Categories: dflt_ctgrs}
                        _ = tmpl_adm_pbs_lst.ExecuteTemplate(w, "base", render)
                        return
                }
                id, err := strconv.ParseInt(toks[3], 10, 64)
                if err != nil {
                        glog.Infof("strconv: %v \n", err)
                        render := Render1 {Message: "Nothing to see here", Packet: &data.Packet{}, Categories: dflt_ctgrs}
                        _ = tmpl_adm_pck_lst.ExecuteTemplate(w, "base", render)
                        return
                }
                pk, err := data.GetLastPacket(id)
                if err != nil {
                        glog.Infof("Https %v \n", err)
                        render := Render1 {Message: "Packets", Packet: &data.Packet{}, Categories: dflt_ctgrs}
                        _ = tmpl_adm_pck_one.ExecuteTemplate(w, "base", render)
                        return
                }
                render := Render1 {Message: "Packets", Packet: pk, Categories: dflt_ctgrs}
                err = tmpl_adm_pck_one.ExecuteTemplate(w, "base", render)
                if err != nil {
                        fmt.Printf("Https %v \n", err)
                        return
                }
                return
        }))
        mux.Handle("/api/", http.HandlerFunc(handleAPI))
        mux.Handle("/pubs/", http.HandlerFunc(handlePubs))
        mux.Handle("/subs/", http.HandlerFunc(handleSubs))
        mux.Handle("/admin/subs/", http.HandlerFunc(handleAdmin))
        mux.Handle("/admin/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                //pbs := make([]*data.Pub, 0)
                //render := Render {"Pubs", pbs, dflt_ctgrs}
                //err = tmpl_adm_gds_lst.ExecuteTemplate(w, "admin", s0)
                pbs, err := data.GetPubs(10)
                if err != nil {
                        fmt.Printf("Https %v \n", err)
                        epbs := make([]*data.Pub, 3)
                        render := Render {Message: "Pubs", Pubs: epbs, Categories: dflt_ctgrs}
                        _ = tmpl_adm_pbs_lst.ExecuteTemplate(w, "base", render)
                        return
                }
                render := Render {Message: "Pubs", Pubs: pbs, Categories: dflt_ctgrs}
                err = tmpl_adm_pbs_lst.ExecuteTemplate(w, "base", render)
                if err != nil {
                        fmt.Printf("Https %v \n", err)
                        return
                }
                return
        }))
        mux.Handle("/", http.HandlerFunc(handleRoot))
        hs := http.Server{
                ReadTimeout: time.Duration(httpRto) * time.Second,
                WriteTimeout: time.Duration(httpWto) * time.Second,
                Addr: fmt.Sprintf(":%d", httpsPort), //":https", //
                //TLSConfig: man.TLSConfig(),
                Handler: mux,
        }
        err := hs.ListenAndServeTLS("", "")
        if err != nil {
                fmt.Printf("Https %v \n", err)
                return
        }
}

func RedirectHttp(w http.ResponseWriter, r *http.Request) {
        if r.TLS != nil || r.Host == "" {
                http.Error(w, "Not Found", 404)
        }
        u := r.URL
        u.Host = r.Host
        u.Scheme = "https"
        switch r.Method {
        case "GET":
                http.Redirect(w, r, u.String(), 302)
        case "POST":
                http.Redirect(w, r, u.String(), 307)
        }
}

func GetKey(path string) (*ecdsa.PrivateKey, error) {
        keybs, err := ioutil.ReadFile(path)
        if err != nil {
                return nil, err
        }
        d, _ := pem.Decode(keybs)
        if d == nil {
                return nil, fmt.Errorf("pem no block found")
        }
        k, err := x509.ParseECPrivateKey(d.Bytes)
        if err != nil {
                return nil, err
        }
        return k, nil
}

// Attempt to parse the given private key DER block. OpenSSL 0.9.8 generates
// PKCS#1 private keys by default, while OpenSSL 1.0.0 generates PKCS#8 keys.
// OpenSSL ecparam generates SEC1 EC private keys for ECDSA. We try all three.
//
// Inspired by parsePrivateKey in crypto/tls/tls.go.
func parsePrivateKey(der []byte) (crypto.Signer, error) {
	if key, err := x509.ParsePKCS1PrivateKey(der); err == nil {
		return key, nil
	}
	if key, err := x509.ParsePKCS8PrivateKey(der); err == nil {
		switch key := key.(type) {
		case *rsa.PrivateKey:
			return key, nil
		case *ecdsa.PrivateKey:
			return key, nil
		default:
			return nil, errors.New("acme/autocert: unknown private key type in PKCS#8 wrapping")
		}
	}
	if key, err := x509.ParseECPrivateKey(der); err == nil {
		return key, nil
	}

	return nil, errors.New("acme/autocert: failed to parse private key")
}

// loadPubDdeets also listens on channel to recieve data of new registrants to send email
func loadPubdeets() {
        from := comms.Entity{"rcs@b00m.in", "RCS"}
        to := make([]comms.Entity, 0)
        for {
                select {
                case newreg:= <-newregs:
                        to = append(to, newreg)
                        sha1str := data.Sha1Str(newreg.Email)
                        err := comms.SendMail(from, to, "Welcome to B00M","Welcome. Please verify your email https://pv.b00m.in/subs/verify/" + sha1str)
                        if err != nil {
                                glog.Infof("comms.SendEmail %v \n", err)
                        }
                        to = to[:0]
                case <-time.After(time.Duration(600 * time.Second)):
                        if err := data.LoadDummyPubdeets(); err != nil {
                                glog.Infof("data.LoadStations %v \n", err)
                        }
                }
        }
}
