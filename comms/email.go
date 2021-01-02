package comms

import (
        "bytes"
        "crypto/tls"
        "encoding/json"
        "net/http"
        //"net/url"
	"fmt"
        "strconv"
        //"strings"
)

type Request struct {
        Messages []MsgToSend `json:"Messages"`
}

type MsgToSend struct {
        From Entity `json:"From"`
        To []Entity `json:"To"`
        Subject string `json:Subject"`
        TextPart string `json:TextPart"`
        HTMLPart string `json:HTMLPart"`
}

type Entity struct {
        Email string `json:"Email"`
        Name string `json:"Name"`
}

type Response struct {
        Messages []Message `json:"Messages"`
}

type Message struct {
        Status string `json:"Status"`
        Errors Error `json:"Errors"`
}

type Error struct {
        ErrorIdentifier string `json:"ErrorIdentifier"`
        ErrorCode string `json:"ErrorCode"`
        StatusCode int `json:"StatusCode"`
        ErrorMessage string`json:"ErrorMessage"`
        ErrorRelatedTo []string `json:"ErrorRelatedTo"`
}

var (
        mailjet = ""
        api = ""
        secret = ""
        body = `{"Messages":[{"From": { "Email": "asdf@gmail.com", "Name": "ABC"}, "To": [{"Email": "sdfg@gmail.com", "Name": "passenger 1"}], "Subject": "Your email flight plan!", "TextPart": "Dear passenger 1, welcome to Mailjet! May the delivery force be with you!", "HTMLPart": "with you!"}]}`
        from = Entity{"asdf@gmail.com", "ABC"}
        to = []Entity{{"sdfg@gmail.com", "DEF"}}
        msg = &Request{[]MsgToSend{{from, to, "Howdy", "Mate", ""}}}
)

/*func main() {
	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: false, ServerName: "mailjet.com", CipherSuites: []uint16{tls.TLS_RSA_WITH_AES_256_CBC_SHA}}}
	client := &http.Client{Transport: tr}

        bs, err := json.Marshal(msg)
        if err != nil {
                fmt.Printf("json marshal %v \n", err)
                return
        }
        //req, err := http.NewRequest("POST", mailjet, strings.NewReader(body))
        req, err := http.NewRequest("POST", mailjet, bytes.NewReader(bs))
        if err != nil {
                fmt.Printf("req %v \n", err)
        }
        req.Header.Add("Content-Type", "application/json")
        req.SetBasicAuth(api, secret)
        //resp, err := client.PostForm(boom, url.Values{"name":{"Doofys1"}, "email":{"doof1@us.com"}, "phone":{"999"}})
        resp, err := client.Do(req)
        if err != nil {
                fmt.Printf("client do get %v \n", err)
        }
        if resp != nil {
                defer resp.Body.Close()
        } else {
                return
        }
        cl := resp.Header.Get("Content-Length")
        icl, err := strconv.Atoi(cl)
        if err != nil {
                fmt.Printf("strconv %v \n" + err.Error())
                return
        }
        ubs := make([]byte, icl) //300) //icl*3)
        _, err = resp.Body.Read(ubs)
        message := &Response{}
        err = json.Unmarshal(ubs, message)
        ubs = ubs[0:]
        if err != nil {
                fmt.Printf("json unmarshal %v \n", err)
                fmt.Printf("bytes %v \n", string(ubs))
                return
        }
        fmt.Printf("resp status %s %s \n",resp.Status, string(ubs))
}*/

func SendMail(from Entity, to []Entity, sub, msg string) error {
        msgreq := &Request{[]MsgToSend{{from, to, sub, msg, ""}}}
	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: false, ServerName: "mailjet.com", CipherSuites: []uint16{tls.TLS_RSA_WITH_AES_256_CBC_SHA}}}
	client := &http.Client{Transport: tr}

        bs, err := json.Marshal(msgreq)
        if err != nil {
                fmt.Printf("json marshal %v \n", err)
                return err
        }
        //req, err := http.NewRequest("POST", mailjet, strings.NewReader(body))
        req, err := http.NewRequest("POST", mailjet, bytes.NewReader(bs))
        if err != nil {
                fmt.Printf("req %v \n", err)
                return err
        }
        req.Header.Add("Content-Type", "application/json")
        req.SetBasicAuth(api, secret)
        //resp, err := client.PostForm(boom, url.Values{"name":{"Doofys1"}, "email":{"doof1@us.com"}, "phone":{"999"}})
        resp, err := client.Do(req)
        if err != nil {
                fmt.Printf("client do get %v \n", err)
                return err
        }
        if resp != nil {
                defer resp.Body.Close()
        } else {
                return fmt.Errorf("Nil response %s \n", sub)
        }

        cl := resp.Header.Get("Content-Length")
        icl, err := strconv.Atoi(cl)
        if err != nil {
                fmt.Printf("strconv %v \n" + err.Error())
                return err
        }
        ubs := make([]byte, icl) //300) //icl*3)
        _, err = resp.Body.Read(ubs)
        message := &Response{}
        err = json.Unmarshal(ubs, message)
        if err != nil {
                fmt.Printf("json unmarshal %v \n", err)
                fmt.Printf("bytes %v \n", string(ubs))
                return err
        }
        fmt.Printf("resp status %s %s \n",resp.Status, string(ubs))
        return nil
}

func SendSMS() error {

        return nil
}

