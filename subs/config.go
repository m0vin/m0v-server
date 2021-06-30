package subs

import (
	"crypto/tls"
	"crypto/x509"
        "encoding/json"
	"errors"
	"flag"
	"fmt"
        "io/ioutil"
	"os"
)

type Config struct {
	Worker   int       `json:"workerNum"`
	HTTPPort string    `json:"httpPort"`
	Host     string    `json:"host"`
	Port     string    `json:"port"`
	Cluster  RouteInfo `json:"cluster"`
	Router   string    `json:"router"`
	TlsHost  string    `json:"tlsHost"`
	TlsPort  string    `json:"tlsPort"`
	TlsConfoPort  string    `json:"tlsConfoPort"`
	WsPath   string    `json:"wsPath"`
	WsPort   string    `json:"wsPort"`
	WsTLS    bool      `json:"wsTLS"`
	TlsInfo  TLSInfo   `json:"tlsInfo"`
	Debug    bool      `json:"debug"`
	//Plugin   Plugins   `json:"plugins"`
        Categories Category `json:"categories"`
        Db       Store     `json:"store"`
        // new fields
        ClientId string `json:"clientId"`
        Secret   string `json:"secret"`
        WebhookId string `json:"webhookid"`
        Sk string `json:"sk"`
        Pk string `json:"pk"`
}

type Category map[string]string
/*type Category struct {
        Name string `json:"name"`
        Link string `json:"link"`
}

func (c *Categories) ToLower(key string) string {
        if c.Link != "" {
                return c.Link
        }
        return strings.ToLower(c.Name)
}*/

type Store struct {
        Name string `json:"name"`
        User string `json:"user"`
        Password string `json:"password"`
}

type RouteInfo struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

type TLSInfo struct {
	Verify   bool   `json:"verify"`
	CaFile   string `json:"caFile"`
	CertFile string `json:"certFile"`
	KeyFile  string `json:"keyFile"`
}

var DefaultConfig *Config = &Config{
	Worker: 4096,
	Host:   "0.0.0.0",
	Port:   "1883",
}


func showHelp() {
	fmt.Printf("%s\n", usageStr)
	os.Exit(0)
}

func ConfigureConfig(args []string) (*Config, error) {
	config := &Config{}
	var (
		help       bool
		configFile string
	)
	fs := flag.NewFlagSet("subs", flag.ExitOnError)
	fs.Usage = showHelp

	fs.BoolVar(&help, "h", false, "Show this message.")
	fs.BoolVar(&help, "help", false, "Show this message.")
	fs.IntVar(&config.Worker, "w", 1024, "worker num to process message, perfer (client num)/10.")
	fs.IntVar(&config.Worker, "worker", 1024, "worker num to process message, perfer (client num)/10.")
	fs.StringVar(&config.HTTPPort, "httpport", "8080", "Port to listen on.")
	fs.StringVar(&config.HTTPPort, "hp", "8080", "Port to listen on.")
	fs.StringVar(&config.Port, "port", "1883", "Port to listen on.")
	fs.StringVar(&config.Port, "p", "1883", "Port to listen on.")
	fs.StringVar(&config.Host, "host", "0.0.0.0", "Network host to listen on")
	fs.StringVar(&config.Cluster.Port, "cp", "", "Cluster port from which members can connect.")
	fs.StringVar(&config.Cluster.Port, "clusterport", "", "Cluster port from which members can connect.")
	fs.StringVar(&config.Router, "r", "", "Router who maintenance cluster info")
	fs.StringVar(&config.Router, "router", "", "Router who maintenance cluster info")
	fs.StringVar(&config.WsPort, "ws", "", "port for ws to listen on")
	fs.StringVar(&config.WsPort, "wsport", "", "port for ws to listen on")
	fs.StringVar(&config.WsPath, "wsp", "", "path for ws to listen on")
	fs.StringVar(&config.WsPath, "wspath", "", "path for ws to listen on")
	fs.StringVar(&configFile, "config", "", "config file for hmq")
	fs.StringVar(&configFile, "c", "", "config file for hmq")
	fs.BoolVar(&config.Debug, "debug", false, "enable Debug logging.")
	fs.BoolVar(&config.Debug, "d", false, "enable Debug logging.")

        fs.String("stderrthreshold", "INFO", "for glog")
        fs.StringVar(&config.Db.Name, "dbdb", "", "database name")
        fs.StringVar(&config.Db.User, "dbuser", "", "database user")
        fs.StringVar(&config.Db.Password, "dbpw", "", "database user password")

	fs.Bool("D", true, "enable Debug logging.")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	if help {
		showHelp()
		return nil, nil
	}

	fs.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "D":
			config.Debug = true
		}
	})

	if configFile != "" {
		tmpConfig, e := LoadConfig(configFile)
		if e != nil {
			return nil, e
		} else {
			config = tmpConfig
		}
	}

	if config.Debug {
	}

        if config.Db.Name != "" {
                flag.Set("dbdb", config.Db.Name)
                flag.Set("dbuser", config.Db.User)
                flag.Set("dbpw", config.Db.Password)
        }

	if err := config.check(); err != nil {
		return nil, err
	}

	return config, nil

}

func LoadConfig(filename string) (*Config, error) {
        cfgs, err := ioutil.ReadFile(filename)
        if err != nil {
                return nil, err
        }
        var config Config
        err = json.Unmarshal(cfgs, &config)
        if err != nil {
                return nil, err
        }
        return &config, nil
}

func (config *Config) check() error {
	if config.Worker == 0 {
		config.Worker = 1024
	}

	if config.Port != "" {
		if config.Host == "" {
			config.Host = "0.0.0.0"
		}
	}

	if config.Cluster.Port != "" {
		if config.Cluster.Host == "" {
			config.Cluster.Host = "0.0.0.0"
		}
	}
	if config.Router != "" {
		if config.Cluster.Port == "" {
			return errors.New("cluster port is null")
		}
	}

	if config.TlsPort != "" {
		if config.TlsInfo.CertFile == "" || config.TlsInfo.KeyFile == "" {
			fmt.Println("tls config error, no cert or key file.")
			return errors.New("tls config error, no cert or key file.")
		}
		if config.TlsHost == "" {
			config.TlsHost = "0.0.0.0"
		}
	}
	return nil
}

func NewTLSConfig(tlsInfo TLSInfo) (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(tlsInfo.CertFile, tlsInfo.KeyFile)
	if err != nil {
		return nil, fmt.Errorf("error parsing X509 certificate/key pair: %v", err)
	}
	cert.Leaf, err = x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		return nil, fmt.Errorf("error parsing certificate: %v", err)
	}

	// Create TLSConfig
	// We will determine the cipher suites that we prefer.
	config := tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
	}

	// Require client certificates as needed
	if tlsInfo.Verify {
		config.ClientAuth = tls.RequireAndVerifyClientCert
	}
	// Add in CAs if applicable.
	if tlsInfo.CaFile != "" {
		rootPEM, err := ioutil.ReadFile(tlsInfo.CaFile)
		if err != nil || rootPEM == nil {
			return nil, err
		}
		pool := x509.NewCertPool()
		ok := pool.AppendCertsFromPEM([]byte(rootPEM))
		if !ok {
			return nil, fmt.Errorf("failed to parse root ca certificate")
		}
		config.ClientCAs = pool
	}

	return &config, nil
}

func getClientCert(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
        fmt.Println(hello.CipherSuites)
        fmt.Println(hello.SupportedVersions)
        return nil, nil
}
